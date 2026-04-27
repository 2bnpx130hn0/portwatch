package history

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// RenderRhythm writes rhythm results to w (defaults to stdout) in the given format.
func RenderRhythm(results []RhythmResult, format string, w io.Writer) {
	if w == nil {
		w = os.Stdout
	}
	switch format {
	case "json":
		renderRhythmJSON(results, w)
	default:
		renderRhythmText(results, w)
	}
}

func renderRhythmText(results []RhythmResult, w io.Writer) {
	if len(results) == 0 {
		fmt.Fprintln(w, "no rhythm data")
		return
	}
	fmt.Fprintf(w, "%-8s %-6s %-14s %-14s %-6s %s\n",
		"PORT", "PROTO", "AVG_PERIOD", "STDDEV", "COUNT", "REGULAR")
	for _, r := range results {
		reg := "no"
		if r.Regular {
			reg = "yes"
		}
		fmt.Fprintf(w, "%-8d %-6s %-14s %-14s %-6d %s\n",
			r.Port,
			r.Protocol,
			r.PeriodAvg.Round(1000000).String(),
			r.PeriodStddev.Round(1000000).String(),
			r.Occurrences,
			reg,
		)
	}
}

func renderRhythmJSON(results []RhythmResult, w io.Writer) {
	type row struct {
		Port         int    `json:"port"`
		Protocol     string `json:"protocol"`
		PeriodAvgMs  int64  `json:"period_avg_ms"`
		PeriodStdMs  int64  `json:"period_stddev_ms"`
		Occurrences  int    `json:"occurrences"`
		Regular      bool   `json:"regular"`
	}
	rows := make([]row, len(results))
	for i, r := range results {
		rows[i] = row{
			Port:        r.Port,
			Protocol:    r.Protocol,
			PeriodAvgMs: r.PeriodAvg.Milliseconds(),
			PeriodStdMs: r.PeriodStddev.Milliseconds(),
			Occurrences: r.Occurrences,
			Regular:     r.Regular,
		}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(rows)
}
