package history

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// RenderChain writes a Chain to w in the given format ("text" or "json").
func RenderChain(c Chain, format string, w io.Writer) {
	if w == nil {
		w = os.Stdout
	}
	switch format {
	case "json":
		renderChainJSON(c, w)
	default:
		renderChainText(c, w)
	}
}

func renderChainText(c Chain, w io.Writer) {
	if len(c.Links) == 0 {
		fmt.Fprintln(w, "no chain entries")
		return
	}
	fmt.Fprintf(w, "chain: %s (%d events)\n", c.Key, len(c.Links))
	for i, link := range c.Links {
		gapStr := ""
		if i > 0 {
			gapStr = fmt.Sprintf(" (+%s)", link.GapSince.Round(1e6).String())
		}
		fmt.Fprintf(w, "  [%d] %s  %-8s  %d/%s%s\n",
			i+1,
			link.Entry.Timestamp.Format("2006-01-02 15:04:05"),
			link.Entry.Action,
			link.Entry.Port,
			link.Entry.Protocol,
			gapStr,
		)
	}
}

func renderChainJSON(c Chain, w io.Writer) {
	type jsonLink struct {
		Timestamp string  `json:"timestamp"`
		Port      int     `json:"port"`
		Protocol  string  `json:"protocol"`
		Action    string  `json:"action"`
		GapMs     float64 `json:"gap_ms,omitempty"`
	}
	type out struct {
		Key   string     `json:"key"`
		Links []jsonLink `json:"links"`
	}
	o := out{Key: c.Key}
	for i, l := range c.Links {
		jl := jsonLink{
			Timestamp: l.Entry.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
			Port:      l.Entry.Port,
			Protocol:  l.Entry.Protocol,
			Action:    l.Entry.Action,
		}
		if i > 0 {
			jl.GapMs = float64(l.GapSince.Milliseconds())
		}
		o.Links = append(o.Links, jl)
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(o)
}
