package history

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// RenderFlow writes flow edges to w in the given format ("text" or "json").
// If w is nil, os.Stdout is used.
func RenderFlow(edges []FlowEdge, format string, w io.Writer) error {
	if w == nil {
		w = os.Stdout
	}
	switch format {
	case "json":
		return renderFlowJSON(edges, w)
	default:
		return renderFlowText(edges, w)
	}
}

func renderFlowText(edges []FlowEdge, w io.Writer) error {
	if len(edges) == 0 {
		_, err := fmt.Fprintln(w, "no flow edges found")
		return err
	}
	_, err := fmt.Fprintf(w, "%-10s %-10s %-8s %s\n", "FROM", "TO", "PROTO", "COUNT")
	if err != nil {
		return err
	}
	for _, e := range edges {
		_, err = fmt.Fprintf(w, "%-10d %-10d %-8s %d\n",
			e.FromPort, e.ToPort, e.Protocol, e.Count)
		if err != nil {
			return err
		}
	}
	return nil
}

func renderFlowJSON(edges []FlowEdge, w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(edges)
}
