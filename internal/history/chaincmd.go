package history

import (
	"fmt"
	"io"
	"os"
	"time"
)

// ChainCmd builds and renders an activity chain from a history file.
type ChainCmd struct {
	Path     string
	Port     int
	Protocol string
	Action   string
	Since    time.Duration
	MaxGap   time.Duration
	Format   string
	Output   io.Writer
}

// Run executes the chain command.
func (c *ChainCmd) Run() error {
	h, err := New(c.Path)
	if err != nil {
		return fmt.Errorf("chain: load history: %w", err)
	}

	if c.Port == 0 {
		return fmt.Errorf("chain: --port is required")
	}

	opts := ChainOptions{
		Port:     c.Port,
		Protocol: c.Protocol,
		Action:   c.Action,
		MaxGap:   c.MaxGap,
	}
	if c.Since > 0 {
		opts.Since = time.Now().Add(-c.Since)
	}

	entries := h.All()
	chain := BuildChain(entries, opts)

	w := c.Output
	if w == nil {
		w = os.Stdout
	}
	RenderChain(chain, c.Format, w)
	return nil
}
