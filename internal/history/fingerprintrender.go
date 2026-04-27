package history

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

// RenderFingerprint writes fingerprint results to w in the given format.
func RenderFingerprint(results []FingerprintResult, format string, w io.Writer) {
	if w == nil {
		w = os.Stdout
	}
	switch strings.ToLower(format) {
	case "json":
		renderFingerprintJSON(results, w)
	default:
		renderFingerprintText(results, w)
	}
}

func renderFingerprintText(results []FingerprintResult, w io.Writer) {
	if len(results) == 0 {
		fmt.Fprintln(w, "no fingerprint data")
		return
	}
	fmt.Fprintf(w, "%-8s %-8s %-10s %-6s %s\n", "PORT", "PROTO", "EVENTS", "SPAN", "FINGERPRINT")
	fmt.Fprintln(w, strings.Repeat("-", 72))
	for _, r := range results {
		span := r.LastSeen.Sub(r.FirstSeen).Round(1e9)
		fmt.Fprintf(w, "%-8d %-8s %-10d %-6s %s\n",
			r.Port, r.Protocol, r.EventCount, span.String(), r.Fingerprint)
		// print action breakdown
		keys := make([]string, 0, len(r.Actions))
		for k := range r.Actions {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fmt.Fprintf(w, "         %-8s %d\n", k, r.Actions[k])
		}
	}
}

func renderFingerprintJSON(results []FingerprintResult, w io.Writer) {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(results)
}
