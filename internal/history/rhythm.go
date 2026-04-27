package history

import (
	"math"
	"sort"
	"time"
)

// RhythmResult describes the periodicity detected for a port/protocol pair.
type RhythmResult struct {
	Port     int
	Protocol string
	PeriodAvg time.Duration
	PeriodStddev time.Duration
	Occurrences int
	Regular     bool // true when stddev/avg < 0.25
}

// RhythmOptions controls Rhythm detection.
type RhythmOptions struct {
	MinOccurrences int
	Since          time.Time
	Action         string
	RegularOnly    bool
}

// Rhythm analyses inter-arrival times for each port/protocol pair and
// reports whether events arrive with a regular cadence.
func Rhythm(entries []Entry, opts RhythmOptions) []RhythmResult {
	if opts.MinOccurrences <= 0 {
		opts.MinOccurrences = 3
	}

	type key struct {
		port  int
		proto string
	}

	buckets := map[key][]time.Time{}
	for _, e := range entries {
		if !opts.Since.IsZero() && e.Timestamp.Before(opts.Since) {
			continue
		}
		if opts.Action != "" && !strings.EqualFold(e.Action, opts.Action) {
			continue
		}
		k := key{port: e.Port, proto: normProto(e.Protocol)}
		buckets[k] = append(buckets[k], e.Timestamp)
	}

	var results []RhythmResult
	for k, times := range buckets {
		if len(times) < opts.MinOccurrences {
			continue
		}
		sort.Slice(times, func(i, j int) bool { return times[i].Before(times[j]) })

		gaps := make([]float64, len(times)-1)
		for i := 1; i < len(times); i++ {
			gaps[i-1] = float64(times[i].Sub(times[i-1]))
		}

		avg, stddev := meanStddev(gaps)
		periodAvg := time.Duration(avg)
		periodStd := time.Duration(stddev)

		regular := avg > 0 && (stddev/avg) < 0.25
		if opts.RegularOnly && !regular {
			continue
		}

		results = append(results, RhythmResult{
			Port:         k.port,
			Protocol:     k.proto,
			PeriodAvg:    periodAvg,
			PeriodStddev: periodStd,
			Occurrences:  len(times),
			Regular:      regular,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].Regular != results[j].Regular {
			return results[i].Regular
		}
		return results[i].Port < results[j].Port
	})
	return results
}

// rhythmCV returns coefficient of variation (stddev/mean), or 0 if mean==0.
func rhythmCV(avg, stddev float64) float64 {
	if avg == 0 {
		return 0
	}
	return math.Abs(stddev / avg)
}
