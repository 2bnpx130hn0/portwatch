package history

import (
	"math"
	"time"
)

// SpikeOptions configures spike detection behaviour.
type SpikeOptions struct {
	Window    time.Duration // lookback window for baseline
	Multiple  float64      // how many times the baseline to consider a spike
	Action    string       // optional action filter ("alert", "warn", "allow")
	Protocol  string       // optional protocol filter
	MinEvents int          // minimum events required to evaluate
}

// SpikeResult describes a detected spike for a port/protocol pair.
type SpikeResult struct {
	Port      int
	Protocol  string
	Baseline  float64 // average events per window period before the latest period
	Actual    int     // events observed in the latest period
	Multiple  float64 // actual / baseline
}

// DetectSpikes identifies port/protocol pairs whose event count in the most
// recent window period is significantly higher than the historical average.
func DetectSpikes(entries []Entry, opts SpikeOptions) []SpikeResult {
	if opts.Window <= 0 {
		opts.Window = time.Hour
	}
	if opts.Multiple <= 0 {
		opts.Multiple = 3.0
	}
	if opts.MinEvents <= 0 {
		opts.MinEvents = 2
	}

	now := time.Now()
	cutoff := now.Add(-opts.Window)

	type key struct{ port int; proto string }

	// Split entries into "current" (last window) and "historical" (prior).
	current := map[key]int{}
	historical := map[key][]int{} // counts per window period

	for _, e := range entries {
		if opts.Action != "" && !strings.EqualFold(e.Action, opts.Action) {
			continue
		}
		if opts.Protocol != "" && !strings.EqualFold(e.Protocol, opts.Protocol) {
			continue
		}
		k := key{e.Port, strings.ToLower(e.Protocol)}
		if e.Timestamp.After(cutoff) {
			current[k]++
		} else {
			// bucket into window-sized periods for averaging
			periodIdx := int(now.Sub(e.Timestamp) / opts.Window)
			for len(historical[k]) <= periodIdx {
				historical[k] = append(historical[k], 0)
			}
			historical[k][periodIdx]++
		}
	}

	var results []SpikeResult
	for k, cnt := range current {
		hist := historical[k]
		if len(hist) == 0 {
			continue
		}
		if cnt < opts.MinEvents {
			continue
		}
		baseline := mean(hist)
		if baseline == 0 {
			continue
		}
		multiple := float64(cnt) / baseline
		if multiple >= opts.Multiple {
			results = append(results, SpikeResult{
				Port:     k.port,
				Protocol: k.proto,
				Baseline: math.Round(baseline*100) / 100,
				Actual:   cnt,
				Multiple: math.Round(multiple*100) / 100,
			})
		}
	}
	return results
}

func mean(vals []int) float64 {
	if len(vals) == 0 {
		return 0
	}
	sum := 0
	for _, v := range vals {
		sum += v
	}
	return float64(sum) / float64(len(vals))
}
