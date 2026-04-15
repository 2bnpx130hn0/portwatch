package notifier

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Format controls how notifications are rendered.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Event represents a port change event to be notified.
type Event struct {
	Timestamp time.Time
	Protocol  string
	Port      uint16
	Action    string // "added" | "removed"
	Level     string // "info" | "warn" | "alert"
}

// Notifier writes formatted event notifications to a writer.
type Notifier struct {
	format Format
	out    io.Writer
}

// New creates a Notifier writing to out with the given format.
// If out is nil, os.Stdout is used.
func New(format Format, out io.Writer) *Notifier {
	if out == nil {
		out = os.Stdout
	}
	if format == "" {
		format = FormatText
	}
	return &Notifier{format: format, out: out}
}

// Send writes a notification for the given event.
func (n *Notifier) Send(e Event) error {
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now()
	}
	switch n.format {
	case FormatJSON:
		return n.writeJSON(e)
	default:
		return n.writeText(e)
	}
}

func (n *Notifier) writeText(e Event) error {
	_, err := fmt.Fprintf(
		n.out,
		"%s [%s] port %s/%d %s\n",
		e.Timestamp.Format(time.RFC3339),
		e.Level,
		e.Protocol,
		e.Port,
		e.Action,
	)
	return err
}

func (n *Notifier) writeJSON(e Event) error {
	_, err := fmt.Fprintf(
		n.out,
		`{"time":%q,"level":%q,"protocol":%q,"port":%d,"action":%q}\n`,
		e.Timestamp.Format(time.RFC3339),
		e.Level,
		e.Protocol,
		e.Port,
		e.Action,
	)
	return err
}
