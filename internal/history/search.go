package history

import "strings"

// SearchOptions defines criteria for full-text/field search over history entries.
type SearchOptions struct {
	Port     int
	Protocol string
	Action   string
	Host     string
}

// Search returns entries matching all non-zero fields in opts.
func Search(entries []Entry, opts SearchOptions) []Entry {
	var result []Entry
	for _, e := range entries {
		if opts.Port != 0 && e.Port != opts.Port {
			continue
		}
		if opts.Protocol != "" && !strings.EqualFold(e.Protocol, opts.Protocol) {
			continue
		}
		if opts.Action != "" && !strings.EqualFold(e.Action, opts.Action) {
			continue
		}
		if opts.Host != "" && !strings.EqualFold(e.Host, opts.Host) {
			continue
		}
		result = append(result, e)
	}
	return result
}

// SearchAny returns entries matching ANY non-zero field in opts (OR semantics).
func SearchAny(entries []Entry, opts SearchOptions) []Entry {
	var result []Entry
	for _, e := range entries {
		if opts.Port != 0 && e.Port == opts.Port {
			result = append(result, e)
			continue
		}
		if opts.Protocol != "" && strings.EqualFold(e.Protocol, opts.Protocol) {
			result = append(result, e)
			continue
		}
		if opts.Action != "" && strings.EqualFold(e.Action, opts.Action) {
			result = append(result, e)
			continue
		}
		if opts.Host != "" && strings.EqualFold(e.Host, opts.Host) {
			result = append(result, e)
			continue
		}
	}
	return result
}
