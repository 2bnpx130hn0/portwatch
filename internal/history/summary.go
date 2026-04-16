package history

// Summary holds aggregated statistics over a slice of history entries.
type Summary struct {
	Total      int            `json:"total"`
	ByAction   map[string]int `json:"by_action"`
	ByProtocol map[string]int `json:"by_protocol"`
	TopPorts   []PortCount    `json:"top_ports"`
}

// PortCount pairs a port number with its occurrence count.
type PortCount struct {
	Port  int `json:"port"`
	Count int `json:"count"`
}

// Summarize computes a Summary from the provided entries.
func Summarize(entries []Entry) Summary {
	s := Summary{
		Total:      len(entries),
		ByAction:   make(map[string]int),
		ByProtocol: make(map[string]int),
	}
	portCounts := make(map[int]int)
	for _, e := range entries {
		s.ByAction[e.Action]++
		s.ByProtocol[e.Protocol]++
		portCounts[e.Port]++
	}
	s.TopPorts = topN(portCounts, 5)
	return s
}

func topN(counts map[int]int, n int) []PortCount {
	result := make([]PortCount, 0, len(counts))
	for port, count := range counts {
		result = append(result, PortCount{Port: port, Count: count})
	}
	// simple insertion sort for small slices
	for i := 1; i < len(result); i++ {
		for j := i; j > 0 && result[j].Count > result[j-1].Count; j-- {
			result[j], result[j-1] = result[j-1], result[j]
		}
	}
	if len(result) > n {
		return result[:n]
	}
	return result
}
