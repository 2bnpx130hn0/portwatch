package history

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// RenderSessions writes session data to w in the given format.
func RenderSessions(sessions []Session, format string, w io.Writer) {
	if w == nil {
		w = os.Stdout
	}
	switch format {
	case "json":
		renderSessionsJSON(sessions, w)
	default:
		renderSessionsText(sessions, w)
	}
}

func renderSessionsText(sessions []Session, w io.Writer) {
	if len(sessions) == 0 {
		fmt.Fprintln(w, "no sessions found")
		return
	}
	fmt.Fprintf(w, "%-6s %-8s %-8s %-26s %-26s %-12s %s\n",
		"PORT", "PROTO", "ACTION", "START", "END", "DURATION", "COUNT")
	for _, s := range sessions {
		fmt.Fprintf(w, "%-6d %-8s %-8s %-26s %-26s %-12s %d\n",
			s.Port,
			s.Protocol,
			s.Action,
			s.Start.Format("2006-01-02 15:04:05"),
			s.End.Format("2006-01-02 15:04:05"),
			s.Duration.Round(1e6).String(),
			s.Count,
		)
	}
}

func renderSessionsJSON(sessions []Session, w io.Writer) {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(sessions)
}
