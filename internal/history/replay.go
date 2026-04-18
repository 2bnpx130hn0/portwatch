package history

import (
	"io"
	"os"
	"time"
)

// ReplayOptions controls how a replay is performed.
type ReplayOptions struct {
	Since  time.Time
	Until  time.Time
	Action string
	Limit  int
}

// ReplayResult holds the entries selected for replay.
type ReplayResult struct {
	Entries []Entry
	Total   int
}

// Replay filters history entries according to opts and returns them in
// chronological order, ready for re-processing or display.
func Replay(entries []Entry, opts ReplayOptions) ReplayResult {
	var out []Entry
	for _, e := range entries {
		if !opts.Since.IsZero() && e.Timestamp.Before(opts.Since) {
			continue
		}
		if !opts.Until.IsZero() && e.Timestamp.After(opts.Until) {
			continue
		}
		if opts.Action != "" && !equalFold(e.Action, opts.Action) {
			continue
		}
		out = append(out, e)
		if opts.Limit > 0 && len(out) >= opts.Limit {
			break
		}
	}
	return ReplayResult{Entries: out, Total: len(out)}
}

// RenderReplay writes replay results to w (defaults to stdout).
func RenderReplay(r ReplayResult, format string, w io.Writer) {
	if w == nil {
		w = os.Stdout
	}
	if format == "json" {
		renderReplayJSON(r, w)
		return
	}
	renderReplayText(r, w)
}
