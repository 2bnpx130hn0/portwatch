package history

import "time"

// CompareOptions controls how two slices of entries are compared.
type CompareOptions struct {
	Since  time.Time
	Before time.Time
	Action string
}

// CompareResult holds the output of a Compare operation.
type CompareResult struct {
	OnlyInA []Entry
	OnlyInB []Entry
	InBoth  []Entry
}

// Compare returns entries unique to a, unique to b, and present in both,
// matched by port+protocol. Optional filters are applied before comparison.
func Compare(a, b []Entry, opts CompareOptions) CompareResult {
	a = applyCompareFilter(a, opts)
	b = applyCompareFilter(b, opts)

	keyOf := func(e Entry) string {
		return itoa(e.Port) + "/" + e.Protocol
	}

	aMap := make(map[string]Entry, len(a))
	for _, e := range a {
		aMap[keyOf(e)] = e
	}

	bMap := make(map[string]Entry, len(b))
	for _, e := range b {
		bMap[keyOf(e)] = e
	}

	var result CompareResult
	for k, e := range aMap {
		if _, ok := bMap[k]; ok {
			result.InBoth = append(result.InBoth, e)
		} else {
			result.OnlyInA = append(result.OnlyInA, e)
		}
	}
	for k, e := range bMap {
		if _, ok := aMap[k]; !ok {
			result.OnlyInB = append(result.OnlyInB, e)
		}
	}
	return result
}

func applyCompareFilter(entries []Entry, opts CompareOptions) []Entry {
	var out []Entry
	for _, e := range entries {
		if !opts.Since.IsZero() && e.Timestamp.Before(opts.Since) {
			continue
		}
		if !opts.Before.IsZero() && !e.Timestamp.Before(opts.Before) {
			continue
		}
		if opts.Action != "" && !equalFold(e.Action, opts.Action) {
			continue
		}
		out = append(out, e)
	}
	return out
}
