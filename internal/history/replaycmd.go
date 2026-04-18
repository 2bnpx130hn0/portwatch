package history

import (
	"fmt"
	"io"
	"os"
	"time"
)

// ReplayCmd is a high-level command that loads history from path,
// applies opts, and renders the result.
type ReplayCmd struct {
	Path   string
	Opts   ReplayOptions
	Format string
	Out    io.Writer
}

// NewReplayCmd creates a ReplayCmd with sensible defaults.
func NewReplayCmd(path string) *ReplayCmd {
	return &ReplayCmd{
		Path:   path,
		Format: "text",
		Out:    os.Stdout,
	}
}

// Run executes the replay command.
func (c *ReplayCmd) Run() error {
	h, err := New(c.Path)
	if err != nil {
		return fmt.Errorf("replay: load history: %w", err)
	}
	entries := h.All()
	result := Replay(entries, c.Opts)
	RenderReplay(result, c.Format, c.Out)
	return nil
}

// WithSince sets the Since filter.
func (c *ReplayCmd) WithSince(t time.Time) *ReplayCmd {
	c.Opts.Since = t
	return c
}

// WithLimit caps the number of replayed entries.
func (c *ReplayCmd) WithLimit(n int) *ReplayCmd {
	c.Opts.Limit = n
	return c
}
