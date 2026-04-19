package history

import (
	"fmt"
	"io"
	"os"
)

// ScoreCmd prints the top-scored entries from a history file.
type ScoreCmd struct {
	Path   string
	Top    int
	Format string
	Writer io.Writer
}

// NewScoreCmd returns a ScoreCmd with defaults.
func NewScoreCmd(path string) *ScoreCmd {
	return &ScoreCmd{
		Path:   path,
		Top:    10,
		Format: "text",
		Writer: os.Stdout,
	}
}

// Run loads history, scores entries, and prints the top results.
func (c *ScoreCmd) Run() error {
	h := New(c.Path)
	if err := h.Load(); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("load history: %w", err)
	}

	entries := h.All()
	scored := Score(entries, ScoreOptions{})
	top := TopScored(scored, c.Top)

	w := c.Writer
	if w == nil {
		w = os.Stdout
	}

	if c.Format == "json" {
		return c.printJSON(w, top)
	}
	return c.printText(w, top)
}

func (c *ScoreCmd) printText(w io.Writer, scored []ScoredEntry) error {
	if len(scored) == 0 {
		_, err := fmt.Fprintln(w, "no entries")
		return err
	}
	for _, s := range scored {
		_, err := fmt.Fprintf(w, "score=%.2f port=%d proto=%s action=%s\n",
			s.Score, s.Entry.Port, s.Entry.Protocol, s.Entry.Action)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *ScoreCmd) printJSON(w io.Writer, scored []ScoredEntry) error {
	_, err := fmt.Fprintln(w, "[")
	if err != nil {
		return err
	}
	for i, s := range scored {
		comma := ","
		if i == len(scored)-1 {
			comma = ""
		}
		_, err = fmt.Fprintf(w, `  {"score":%.2f,"port":%d,"protocol":%q,"action":%q}%s\n`,
			s.Score, s.Entry.Port, s.Entry.Protocol, s.Entry.Action, comma)
		if err != nil {
			return err
		}
	}
	_, err = fmt.Fprintln(w, "]")
	return err
}
