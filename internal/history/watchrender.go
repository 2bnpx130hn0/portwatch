package history

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

// RenderWatchSummary writes a WatchSummary to w in the given format ("text" or "json").
// If w is nil, os.Stdout is used.
func RenderWatchSummary(s WatchSummary, format string, w io.Writer) error {
	if w == nil {
		w = os.Stdout
	}
	switch strings.ToLower(format) {
	case "json":
		return renderWatchSummaryJSON(s, w)
	default:
		return renderWatchSummaryText(s, w)
	}
}

func renderWatchSummaryText(s WatchSummary, w io.Writer) error {
	fmt.Fprintf(w, "Cycle: %s\n", s.CycleAt.Format("2006-01-02 15:04:05"))
	fmt.Fprintf(w, "Added: %d (alerted: %d, allowed: %d)\n", len(s.Added), s.Alerted, s.Allowed)
	for _, e := range s.Added {
		fmt.Fprintf(w, "  + [%s] %s:%d\n", e.Action, e.Protocol, e.Port)
	}
	fmt.Fprintf(w, "Removed: %d\n", len(s.Removed))
	for _, e := range s.Removed {
		fmt.Fprintf(w, "  - %s:%d\n", e.Protocol, e.Port)
	}
	return nil
}

func renderWatchSummaryJSON(s WatchSummary, w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(s)
}
