package history

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// RenderSnapshot writes a human-readable or JSON diff of two snapshots.
func RenderSnapshot(baseline, current Snapshot, format string, w io.Writer) {
	if w == nil {
		w = os.Stdout
	}
	added, removed := SnapshotDiff(baseline, current)
	switch format {
	case "json":
		renderSnapshotJSON(added, removed, baseline.TakenAt.String(), current.TakenAt.String(), w)
	default:
		renderSnapshotText(added, removed, w)
	}
}

func renderSnapshotText(added, removed []Entry, w io.Writer) {
	if len(added) == 0 && len(removed) == 0 {
		fmt.Fprintln(w, "No changes between snapshots.")
		return
	}
	for _, e := range added {
		fmt.Fprintf(w, "+ %s/%d (%s)\n", e.Protocol, e.Port, e.Action)
	}
	for _, e := range removed {
		fmt.Fprintf(w, "- %s/%d (%s)\n", e.Protocol, e.Port, e.Action)
	}
}

func renderSnapshotJSON(added, removed []Entry, baseline, current string, w io.Writer) {
	out := map[string]interface{}{
		"baseline": baseline,
		"current":  current,
		"added":    added,
		"removed":  removed,
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(out)
}
