package history

import (
	"encoding/json"
	"fmt"
	"io"
	"text/tabwriter"
	"time"
)

// Format controls the output style of the Reporter.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Reporter renders history entries to an io.Writer.
type Reporter struct {
	format Format
	w      io.Writer
}

// NewReporter creates a Reporter writing to w in the given format.
func NewReporter(w io.Writer, format Format) *Reporter {
	if format == "" {
		format = FormatText
	}
	return &Reporter{w: w, format: format}
}

// Print writes all entries to the configured writer.
func (r *Reporter) Print(entries []Entry) error {
	switch r.format {
	case FormatJSON:
		return r.printJSON(entries)
	default:
		return r.printText(entries)
	}
}

func (r *Reporter) printText(entries []Entry) error {
	tw := tabwriter.NewWriter(r.w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "TIMESTAMP\tPROTO\tPORT\tACTION\tMESSAGE")
	for _, e := range entries {
		fmt.Fprintf(tw, "%s\t%s\t%d\t%s\t%s\n",
			e.Timestamp.Format(time.RFC3339),
			e.Proto,
			e.Port,
			e.Action,
			e.Message,
		)
	}
	return tw.Flush()
}

func (r *Reporter) printJSON(entries []Entry) error {
	enc := json.NewEncoder(r.w)
	enc.SetIndent("", "  ")
	return enc.Encode(entries)
}
