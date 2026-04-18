package history

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

// RenderTimeline writes timeline entries to w in the given format ("text" or "json").
func RenderTimeline(entries []TimelineEntry, format string, w io.Writer) error {
	if w == nil {
		w = os.Stdout
	}
	switch strings.ToLower(format) {
	case "json":
		return renderTimelineJSON(entries, w)
	default:
		return renderTimelineText(entries, w)
	}
}

func renderTimelineText(entries []TimelineEntry, w io.Writer) error {
	if len(entries) == 0 {
		_, err := fmt.Fprintln(w, "No timeline data.")
		return err
	}
	for _, e := range entries {
		line := fmt.Sprintf("%s  total=%-4d", e.Bucket.Format("2006-01-02 15:04"), e.Total)
		for action, count := range e.ByAction {
			line += fmt.Sprintf("  %s=%d", action, count)
		}
		if _, err := fmt.Fprintln(w, line); err != nil {
			return err
		}
	}
	return nil
}

func renderTimelineJSON(entries []TimelineEntry, w io.Writer) error {
	return json.NewEncoder(w).Encode(entries)
}
