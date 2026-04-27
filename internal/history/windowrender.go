package history

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
)

// RenderWindow writes window buckets in the specified format.
func RenderWindow(buckets []WindowBucket, format string, w io.Writer) {
	if w == nil {
		w = os.Stdout
	}
	if format == "json" {
		renderWindowJSON(buckets, w)
	} else {
		renderWindowText(buckets, w)
	}
}

func renderWindowText(buckets []WindowBucket, w io.Writer) {
	if len(buckets) == 0 {
		fmt.Fprintln(w, "no window data")
		return
	}
	for _, b := range buckets {
		fmt.Fprintf(w, "[%s - %s] count=%d\n",
			b.Start.Format("15:04"),
			b.End.Format("15:04"),
			b.Count,
		)
		actions := make([]string, 0, len(b.Actions))
		for a := range b.Actions {
			actions = append(actions, a)
		}
		sort.Strings(actions)
		for _, a := range actions {
			fmt.Fprintf(w, "  %s: %d\n", a, b.Actions[a])
		}
	}
}

func renderWindowJSON(buckets []WindowBucket, w io.Writer) {
	type jsonBucket struct {
		Start   string         `json:"start"`
		End     string         `json:"end"`
		Count   int            `json:"count"`
		Ports   map[int]int    `json:"ports"`
		Actions map[string]int `json:"actions"`
	}
	out := make([]jsonBucket, len(buckets))
	for i, b := range buckets {
		out[i] = jsonBucket{
			Start:   b.Start.Format("2006-01-02T15:04:05Z07:00"),
			End:     b.End.Format("2006-01-02T15:04:05Z07:00"),
			Count:   b.Count,
			Ports:   b.Ports,
			Actions: b.Actions,
		}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(out)
}
