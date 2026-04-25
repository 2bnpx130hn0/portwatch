package history

import (
	"flag"
	"fmt"
	"os"
	"time"
)

// SessionCmd is a CLI command that prints port activity sessions.
type SessionCmd struct {
	HistoryFile string
	Format      string
	Action      string
	Protocol    string
	Gap         time.Duration
	Since       time.Duration
}

// NewSessionCmd returns a SessionCmd with defaults.
func NewSessionCmd() *SessionCmd {
	return &SessionCmd{
		Format: "text",
		Gap:    5 * time.Minute,
	}
}

// RegisterFlags binds CLI flags to the command.
func (c *SessionCmd) RegisterFlags(fs *flag.FlagSet) {
	fs.StringVar(&c.HistoryFile, "history", "portwatch_history.json", "path to history file")
	fs.StringVar(&c.Format, "format", "text", "output format: text|json")
	fs.StringVar(&c.Action, "action", "", "filter by action (allow|warn|alert)")
	fs.StringVar(&c.Protocol, "protocol", "", "filter by protocol")
	fs.DurationVar(&c.Gap, "gap", 5*time.Minute, "max gap between events to merge into one session")
	fs.DurationVar(&c.Since, "since", 0, "only include entries newer than this duration ago")
}

// Run executes the session command.
func (c *SessionCmd) Run() error {
	h, err := New(c.HistoryFile)
	if err != nil {
		return fmt.Errorf("load history: %w", err)
	}

	opts := SessionOptions{
		Gap:      c.Gap,
		Action:   c.Action,
		Protocol: c.Protocol,
	}
	if c.Since > 0 {
		opts.Since = time.Now().UTC().Add(-c.Since)
	}

	sessions := BuildSessions(h.Entries(), opts)
	RenderSessions(sessions, c.Format, os.Stdout)
	return nil
}
