package history

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

// RenderRank writes rank results to w in the given format ("text" or "json").
// If w is nil, os.Stdout is used.
func RenderRank(results []RankResult, format string, w io.Writer) {
	if w == nil {
		w = os.Stdout
	}
	switch strings.ToLower(format) {
	case "json":
		renderRankJSON(results, w)
	default:
		renderRankText(results, w)
	}
}

func renderRankText(results []RankResult, w io.Writer) {
	if len(results) == 0 {
		fmt.Fprintln(w, "no rank results")
		return
	}
	fmt.Fprintf(w, "%-8s %-10s %-8s %s\n", "PORT", "PROTOCOL", "COUNT", "SCORE")
	fmt.Fprintln(w, strings.Repeat("-", 40))
	for _, r := range results {
		fmt.Fprintf(w, "%-8d %-10s %-8d %.4f\n", r.Port, r.Protocol, r.Count, r.Score)
	}
}

func renderRankJSON(results []RankResult, w io.Writer) {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(results); err != nil {
		fmt.Fprintf(w, `{"error":%q}\n`, err.Error())
	}
}
