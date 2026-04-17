package history

import "strings"

// Tag adds one or more tags to entries matching a predicate.
func Tag(entries []Entry, tags []string, match func(Entry) bool) []Entry {
	result := make([]Entry, len(entries))
	for i, e := range entries {
		if match(e) {
			existing := map[string]struct{}{}
			for _, t := range e.Tags {
				existing[t] = struct{}{}
			}
			for _, t := range tags {
				t = strings.TrimSpace(t)
				if t != "" {
					existing[t] = struct{}{}
				}
			}
			merged := make([]string, 0, len(existing))
			for t := range existing {
				merged = append(merged, t)
			}
			e.Tags = merged
		}
		result[i] = e
	}
	return result
}

// Untag removes specified tags from all entries.
func Untag(entries []Entry, tags []string) []Entry {
	remove := map[string]struct{}{}
	for _, t := range tags {
		remove[t] = struct{}{}
	}
	result := make([]Entry, len(entries))
	for i, e := range entries {
		kept := e.Tags[:0:0]
		for _, t := range e.Tags {
			if _, skip := remove[t]; !skip {
				kept = append(kept, t)
			}
		}
		e.Tags = kept
		result[i] = e
	}
	return result
}

// FilterByTag returns entries that have all of the specified tags.
func FilterByTag(entries []Entry, tags []string) []Entry {
	var out []Entry
	for _, e := range entries {
		set := map[string]struct{}{}
		for _, t := range e.Tags {
			set[t] = struct{}{}
		}
		match := true
		for _, t := range tags {
			if _, ok := set[t]; !ok {
				match = false
				break
			}
		}
		if match {
			out = append(out, e)
		}
	}
	return out
}
