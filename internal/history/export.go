package history

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// Exporter writes history entries to an io.Writer in a specified format.
type Exporter struct {
	format string
	w      io.Writer
}

// NewExporter creates an Exporter writing to w in the given format ("csv" or "json").
func NewExporter(w io.Writer, format string) *Exporter {
	if format == "" {
		format = "csv"
	}
	return &Exporter{format: format, w: w}
}

// Export writes entries to the underlying writer.
func (e *Exporter) Export(entries []Entry) error {
	switch e.format {
	case "json":
		return e.exportJSON(entries)
	default:
		return e.exportCSV(entries)
	}
}

func (e *Exporter) exportCSV(entries []Entry) error {
	w := csv.NewWriter(e.w)
	if err := w.Write([]string{"timestamp", "protocol", "port", "action", "rule"}); err != nil {
		return err
	}
	for _, en := range entries {
		row := []string{
			en.Timestamp.Format(time.RFC3339),
			en.Protocol,
			fmt.Sprintf("%d", en.Port),
			en.Action,
			en.Rule,
		}
		if err := w.Write(row); err != nil {
			return err
		}
	}
	w.Flush()
	return w.Error()
}

func (e *Exporter) exportJSON(entries []Entry) error {
	enc := json.NewEncoder(e.w)
	enc.SetIndent("", "  ")
	return enc.Encode(entries)
}
