package history

import "sort"

// GroupKey identifies a grouping dimension.
type GroupKey string

const (
	GroupByPort     GroupKey = "port"
	GroupByProtocol GroupKey = "protocol"
	GroupByAction   GroupKey = "action"
)

// Group holds entries sharing the same key value.
type Group struct {
	Key     string
	Entries []Entry
}

// GroupBy partitions entries by the given key dimension.
// Groups are returned sorted by key ascending.
func GroupBy(entries []Entry, key GroupKey) []Group {
	m := make(map[string][]Entry)
	for _, e := range entries {
		var k string
		switch key {
		case GroupByPort:
			k = itoa(e.Port)
		case GroupByProtocol:
			k = e.Protocol
		case GroupByAction:
			k = e.Action
		default:
			k = "unknown"
		}
		m[k] = append(m[k], e)
	}

	groups := make([]Group, 0, len(m))
	for k, es := range m {
		groups = append(groups, Group{Key: k, Entries: es})
	}
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Key < groups[j].Key
	})
	return groups
}

// GroupCounts returns a map of key -> count for quick aggregation.
func GroupCounts(entries []Entry, key GroupKey) map[string]int {
	counts := make(map[string]int)
	for _, g := range GroupBy(entries, key) {
		counts[g.Key] = len(g.Entries)
	}
	return counts
}
