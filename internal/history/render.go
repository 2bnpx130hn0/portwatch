package history

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"text/tabwriter"
)

// RenderStats writes PortStat results to w in the given format ("text" or "json").
func RenderStats(w io.Writer, stats []PortStat, format string) error {
	sort.Slice(stats, func(i, j int) bool {
		if stats[i].Seen != stats[j].Seen {
			return stats[i].Seen > stats[j].Seen
		}
		return stats[i].Port < stats[j].Port
	})

	switch format {
	case "json":
		return json.NewEncoder(w).Encode(stats)
	default:
		return renderStatsText(w, stats)
	}
}

func renderStatsText(w io.Writer, stats []PortStat) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "PORT\tPROTOCOL\tSEEN\tLAST SEEN\tACTIONS")
	for _, s := range stats {
		actStr := ""
		for a, c := range s.Actions {
			if actStr != "" {
				actStr += " "
			}
			actStr += fmt.Sprintf("%s:%d", a, c)
		}
		fmt.Fprintf(tw, "%d\t%s\t%d\t%s\t%s\n",
			s.Port, s.Protocol, s.Seen,
			s.LastSeen.Format("2006-01-02 15:04:05"),
			actStr,
		)
	}
	return tw.Flush()
}
