package history

import "time"

// PortStat holds aggregated statistics for a single port.
type PortStat struct {
	Port     int
	Protocol string
	Seen     int
	LastSeen time.Time
	Actions  map[string]int
}

// Stats computes per-port statistics from a slice of entries.
func Stats(entries []Entry) []PortStat {
	index := map[string]*PortStat{}

	for _, e := range entries {
		key := e.Protocol + ":" + itoa(e.Port)
		st, ok := index[key]
		if !ok {
			st = &PortStat{
				Port:     e.Port,
				Protocol: e.Protocol,
				Actions:  map[string]int{},
			}
			index[key] = st
		}
		st.Seen++
		st.Actions[e.Action]++
		if e.Timestamp.After(st.LastSeen) {
			st.LastSeen = e.Timestamp
		}
	}

	out := make([]PortStat, 0, len(index))
	for _, st := range index {
		out = append(out, *st)
	}
	return out
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	buf := [20]byte{}
	pos := len(buf)
	for n > 0 {
		pos--
		buf[pos] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[pos:])
}
