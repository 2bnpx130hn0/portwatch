package history

import (
	"math"
	"time"
)

// DecayOptions controls how score decay is applied.
type DecayOptions struct {
	// HalfLife is the duration after which a score is halved.
	// Defaults to 24 hours.
	HalfLife time.Duration

	// Since excludes entries older than this time.
	Since time.Time

	// Action filters entries to a specific action (e.g. "alert").
	Action string
}

// DecayResult holds an entry paired with its time-decayed score.
type DecayResult struct {
	Entry   Entry
	RawScore float64
	Decayed  float64
}

// Decay applies exponential time-decay to entry scores and returns results
// sorted by decayed score descending.
func Decay(entries []Entry, opts DecayOptions) []DecayResult {
	if opts.HalfLife <= 0 {
		opts.HalfLife = 24 * time.Hour
	}

	now := time.Now()
	lambda := math.Log(2) / opts.HalfLife.Seconds()

	var results []DecayResult
	for _, e := range entries {
		if !opts.Since.IsZero() && e.Timestamp.Before(opts.Since) {
			continue
		}
		if opts.Action != "" && !strings.EqualFold(e.Action, opts.Action) {
			continue
		}

		age := now.Sub(e.Timestamp).Seconds()
		if age < 0 {
			age = 0
		}

		raw := baseDecayScore(e)
		decayed := raw * math.Exp(-lambda*age)

		results = append(results, DecayResult{
			Entry:    e,
			RawScore: raw,
			Decayed:  decayed,
		})
	}

	// Sort descending by decayed score.
	for i := 1; i < len(results); i++ {
		for j := i; j > 0 && results[j].Decayed > results[j-1].Decayed; j-- {
			results[j], results[j-1] = results[j-1], results[j]
		}
	}

	return results
}

func baseDecayScore(e Entry) float64 {
	switch strings.ToLower(e.Action) {
	case "alert":
		return 10.0
	case "warn":
		return 5.0
	default:
		return 1.0
	}
}
