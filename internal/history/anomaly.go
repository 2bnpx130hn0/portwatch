package history

import (
	"math"
	"time"
)

// AnomalyResult holds a detected anomaly for a port/protocol pair.
type AnomalyResult struct {
	Port     int
	Protocol string
	Action   string
	ZScore   float64
	Mean     float64
	StdDev   float64
	Count    int
}

// AnomalyOptions configures anomaly detection.
type AnomalyOptions struct {
	Since     time.Time
	Threshold float64 // z-score threshold; defaults to 2.0
	MinSamples int    // minimum bucket count to consider; defaults to 3
}

// DetectAnomalies identifies port/protocol pairs whose activity count
// deviates significantly from the mean across all observed pairs.
func DetectAnomalies(entries []Entry, opts AnomalyOptions) []AnomalyResult {
	if opts.Threshold == 0 {
		opts.Threshold = 2.0
	}
	if opts.MinSamples == 0 {
		opts.MinSamples = 3
	}

	filtered := entries
	if !opts.Since.IsZero() {
		filtered = filterSince(entries, opts.Since)
	}

	type key struct {
		port     int
		protocol string
		action   string
	}

	counts := make(map[key]int)
	for _, e := range filtered {
		k := key{e.Port, normProto(e.Protocol), e.Action}
		counts[k]++
	}

	if len(counts) < opts.MinSamples {
		return nil
	}

	vals := make([]float64, 0, len(counts))
	for _, c := range counts {
		vals = append(vals, float64(c))
	}

	mean, stddev := meanStddev(vals)
	if stddev == 0 {
		return nil
	}

	var results []AnomalyResult
	for k, c := range counts {
		z := math.Abs((float64(c) - mean) / stddev)
		if z >= opts.Threshold {
			results = append(results, AnomalyResult{
				Port:     k.port,
				Protocol: k.protocol,
				Action:   k.action,
				ZScore:   math.Round(z*100) / 100,
				Mean:     math.Round(mean*100) / 100,
				StdDev:   math.Round(stddev*100) / 100,
				Count:    c,
			})
		}
	}
	return results
}

func filterSince(entries []Entry, since time.Time) []Entry {
	out := make([]Entry, 0, len(entries))
	for _, e := range entries {
		if !e.Timestamp.Before(since) {
			out = append(out, e)
		}
	}
	return out
}

func normProto(p string) string {
	if p == "" {
		return "tcp"
	}
	return p
}

func meanStddev(vals []float64) (float64, float64) {
	if len(vals) == 0 {
		return 0, 0
	}
	var sum float64
	for _, v := range vals {
		sum += v
	}
	mean := sum / float64(len(vals))
	var variance float64
	for _, v := range vals {
		d := v - mean
		variance += d * d
	}
	variance /= float64(len(vals))
	return mean, math.Sqrt(variance)
}
