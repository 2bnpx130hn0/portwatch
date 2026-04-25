package history

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

// RenderHotspot writes hotspot results to w in the given format ("text" or "json").
// If w is nil, os.Stdout is used.
func RenderHotspot(hotspots []HotspotEntry, format string, w io.Writer) {
	if w == nil {
		w = os.Stdout
	}
	switch strings.ToLower(format) {
	case "json":
		renderHotspotJSON(hotspots, w)
	default:
		renderHotspotText(hotspots, w)
	}
}

func renderHotspotText(hotspots []HotspotEntry, w io.Writer) {
	if len(hotspots) == 0 {
		fmt.Fprintln(w, "no hotspots found")
		return
	}
	fmt.Fprintf(w, "%-8s %-8s %-8s %s\n", "PORT", "PROTO", "COUNT", "LAST SEEN")
	fmt.Fprintln(w, strings.Repeat("-", 44))
	for _, h := range hotspots {
		fmt.Fprintf(w, "%-8d %-8s %-8d %s\n",
			h.Port,
			h.Protocol,
			h.Count,
			h.LastSeen.Format("2006-01-02 15:04:05"),
		)
	}
}

func renderHotspotJSON(hotspots []HotspotEntry, w io.Writer) {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if hotspots == nil {
		hotspots = []HotspotEntry{}
	}
	_ = enc.Encode(hotspots)
}
