package history

import (
	"fmt"
	"io"
	"os"
)

// FlagCmd provides CLI operations for flagging history entries.
type FlagCmd struct {
	Path   string
	Writer io.Writer
}

// NewFlagCmd creates a FlagCmd writing to w (defaults to stdout).
func NewFlagCmd(path string, w io.Writer) *FlagCmd {
	if w == nil {
		w = os.Stdout
	}
	return &FlagCmd{Path: path, Writer: w}
}

// Set flags all entries matching port/protocol.
func (c *FlagCmd) Set(port int, protocol string) error {
	h, err := New(c.Path)
	if err != nil {
		return err
	}
	h.entries = Flag(h.entries, port, protocol)
	if err := h.save(); err != nil {
		return err
	}
	fmt.Fprintf(c.Writer, "flagged port %d/%s\n", port, protocol)
	return nil
}

// Remove unflags all entries matching port/protocol.
func (c *FlagCmd) Remove(port int, protocol string) error {
	h, err := New(c.Path)
	if err != nil {
		return err
	}
	h.entries = Unflag(h.entries, port, protocol)
	if err := h.save(); err != nil {
		return err
	}
	fmt.Fprintf(c.Writer, "unflagged port %d/%s\n", port, protocol)
	return nil
}

// List prints all flagged entries.
func (c *FlagCmd) List() error {
	h, err := New(c.Path)
	if err != nil {
		return err
	}
	flagged := FilterFlagged(h.entries)
	if len(flagged) == 0 {
		fmt.Fprintln(c.Writer, "no flagged entries")
		return nil
	}
	for _, e := range flagged {
		fmt.Fprintf(c.Writer, "port=%d protocol=%s action=%s flagged_at=%s\n",
			e.Port, e.Protocol, e.Action, e.Labels["flagged_at"])
	}
	return nil
}
