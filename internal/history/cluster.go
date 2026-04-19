package history

import "sort"

// ClusterResult holds a group of related entries sharing port+protocol.
type ClusterResult struct {
	Port     int
	Protocol string
	Entries  []Entry
	Count    int
}

// ClusterOptions controls how clustering is performed.
type ClusterOptions struct {
	Action   string
	MinCount int
}

// Cluster groups entries by port+protocol and returns clusters sorted by count desc.
func Cluster(entries []Entry, opts ClusterOptions) []ClusterResult {
	type key struct {
		port     int
		protocol string
	}

	m := make(map[key][]Entry)
	for _, e := range entries {
		if opts.Action != "" && !equalFold(e.Action, opts.Action) {
			continue
		}
		k := key{port: e.Port, protocol: e.Protocol}
		m[k] = append(m[k], e)
	}

	var results []ClusterResult
	for k, es := range m {
		if opts.MinCount > 0 && len(es) < opts.MinCount {
			continue
		}
		results = append(results, ClusterResult{
			Port:     k.port,
			Protocol: k.protocol,
			Entries:  es,
			Count:    len(es),
		})
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].Count != results[j].Count {
			return results[i].Count > results[j].Count
		}
		return results[i].Port < results[j].Port
	})

	return results
}
