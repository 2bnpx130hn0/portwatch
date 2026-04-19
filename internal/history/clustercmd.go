package history

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"
)

// ClusterCmd runs the cluster command against a history file.
type ClusterCmd struct {
	Path     string
	Action   string
	MinCount int
	Format   string
	Writer   io.Writer
}

// NewClusterCmd creates a ClusterCmd with defaults.
func NewClusterCmd(path string) *ClusterCmd {
	return &ClusterCmd{Path: path, Format: "text", Writer: os.Stdout}
}

func (c *ClusterCmd) Run() error {
	h, err := New(c.Path)
	if err != nil {
		return fmt.Errorf("cluster: load history: %w", err)
	}

	entries := h.All()
	results := Cluster(entries, ClusterOptions{
		Action:   c.Action,
		MinCount: c.MinCount,
	})

	w := c.Writer
	if w == nil {
		w = os.Stdout
	}

	if c.Format == "json" {
		return c.writeJSON(w, results)
	}
	return c.writeText(w, results)
}

func (c *ClusterCmd) writeText(w io.Writer, results []ClusterResult) error {
	if len(results) == 0 {
		fmt.Fprintln(w, "no clusters found")
		return nil
	}
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "PORT\tPROTOCOL\tCOUNT")
	for _, r := range results {
		fmt.Fprintf(tw, "%d\t%s\t%d\n", r.Port, r.Protocol, r.Count)
	}
	return tw.Flush()
}

func (c *ClusterCmd) writeJSON(w io.Writer, results []ClusterResult) error {
	fmt.Fprintln(w, "[")
	for i, r := range results {
		comma := ","
		if i == len(results)-1 {
			comma = ""
		}
		fmt.Fprintf(w, `  {"port":%d,"protocol":%q,"count":%d}%s`+"\n",
			r.Port, r.Protocol, r.Count, comma)
	}
	fmt.Fprintln(w, "]")
	return nil
}
