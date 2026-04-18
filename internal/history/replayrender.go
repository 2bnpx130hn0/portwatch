package history

import (
	"encoding/json"
	"fmt"
	"io"
)

func renderReplayText(r ReplayResult, w io.Writer) {
	if len(r.Entries) == 0 {
		fmt.Fprintln(w, "no entries to replay")
		return
	}
	fmt.Fprintf(w, "replaying %d entries\n", r.Total)
	fmt.Fprintf(w, "%-30s %-6s %-6s %s\n", "timestamp", "proto", "port", "action")
	for _, e := range r.Entries {
		fmt.Fprintf(w, "%-30s %-6s %-6d %s\n",
			e.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
			e.Protocol, e.Port, e.Action)
	}
}

func renderReplayJSON(r ReplayResult, w io.Writer) {
	type payload struct {
		Total   int     `json:"total"`
		Entries []Entry `json:"entries"`
	}
	_ = json.NewEncoder(w).Encode(payload{Total: r.Total, Entries: r.Entries})
}
