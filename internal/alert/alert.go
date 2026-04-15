package alert

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/rules"
)

// Level represents the severity of an alert.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelAlert Level = "ALERT"
)

// Event describes a port change that triggered an alert.
type Event struct {
	Timestamp time.Time
	Port      int
	Protocol  string
	Action    rules.Action
	Level     Level
	Message   string
}

// Notifier sends alert events to a destination.
type Notifier struct {
	out io.Writer
}

// New creates a Notifier that writes to the given writer.
// If w is nil, os.Stderr is used.
func New(w io.Writer) *Notifier {
	if w == nil {
		w = os.Stderr
	}
	return &Notifier{out: w}
}

// Notify formats and writes an Event to the configured writer.
func (n *Notifier) Notify(e Event) error {
	level := e.Level
	if level == "" {
		level = levelForAction(e.Action)
	}
	_, err := fmt.Fprintf(
		n.out,
		"[%s] %s port=%d protocol=%s action=%s message=%q\n",
		e.Timestamp.UTC().Format(time.RFC3339),
		level,
		e.Port,
		e.Protocol,
		e.Action,
		e.Message,
	)
	return err
}

// levelForAction maps a rules.Action to an alert Level.
func levelForAction(a rules.Action) Level {
	switch a {
	case rules.ActionAllow:
		return LevelInfo
	case rules.ActionWarn:
		return LevelWarn
	default:
		return LevelAlert
	}
}
