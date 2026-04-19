package history

import "time"

// ScoreOptions controls how risk scores are calculated.
type ScoreOptions struct {
	AlertWeight   float64
	WarnWeight    float64
	AllowWeight   float64
	RecencyHours  float64 // entries within this window get a recency boost
	RecencyBoost  float64
}

// ScoredEntry pairs an entry with its computed risk score.
type ScoredEntry struct {
	Entry Entry
	Score float64
}

// Score computes a risk score for each entry based on action weights and recency.
func Score(entries []Entry, opts ScoreOptions) []ScoredEntry {
	if opts.AlertWeight == 0 {
		opts.AlertWeight = 3.0
	}
	if opts.WarnWeight == 0 {
		opts.WarnWeight = 1.5
	}
	if opts.AllowWeight == 0 {
		opts.AllowWeight = 0.5
	}
	if opts.RecencyHours == 0 {
		opts.RecencyHours = 24
	}
	if opts.RecencyBoost == 0 {
		opts.RecencyBoost = 2.0
	}

	now := time.Now()
	result := make([]ScoredEntry, 0, len(entries))
	for _, e := range entries {
		var base float64
		switch e.Action {
		case "alert":
			base = opts.AlertWeight
		case "warn":
			base = opts.WarnWeight
		default:
			base = opts.AllowWeight
		}
		if !e.Timestamp.IsZero() && now.Sub(e.Timestamp).Hours() <= opts.RecencyHours {
			base *= opts.RecencyBoost
		}
		result = append(result, ScoredEntry{Entry: e, Score: base})
	}
	return result
}

// TopScored returns the top n entries by score (descending).
func TopScored(scored []ScoredEntry, n int) []ScoredEntry {
	sorted := make([]ScoredEntry, len(scored))
	copy(sorted, scored)
	for i := 0; i < len(sorted)-1; i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[j].Score > sorted[i].Score {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}
	if n > 0 && n < len(sorted) {
		return sorted[:n]
	}
	return sorted
}
