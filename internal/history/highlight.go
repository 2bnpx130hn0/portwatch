package history

import "time"

// HighlightOptions controls which entries are considered notable.
type HighlightOptions struct {
	MinScore     float64
	Actions      []string
	Since        time.Time
	OnlyFlagged  bool
	OnlyBookmark bool
}

// Highlight returns entries that meet any of the highlight criteria.
func Highlight(entries []Entry, opts HighlightOptions) []Entry {
	var out []Entry
	for _, e := range entries {
		if !opts.Since.IsZero() && e.Timestamp.Before(opts.Since) {
			continue
		}
		if opts.OnlyFlagged && !hasLabel(e, "flagged") {
			continue
		}
		if opts.OnlyBookmark && !hasLabel(e, "bookmarked") {
			continue
		}
		if len(opts.Actions) > 0 && !actionMatch(e.Action, opts.Actions) {
			continue
		}
		if opts.MinScore > 0 {
			scored := Score([]Entry{e})
			if len(scored) == 0 || scored[0].Score < opts.MinScore {
				continue
			}
		}
		out = append(out, e)
	}
	return out
}

func hasLabel(e Entry, key string) bool {
	_, ok := e.Labels[key]
	return ok
}

func actionMatch(action string, actions []string) bool {
	for _, a := range actions {
		if equalFold(action, a) {
			return true
		}
	}
	return false
}
