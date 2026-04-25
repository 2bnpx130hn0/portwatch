package history

import (
	"fmt"
	"os"
	"time"
)

// HotspotCmd provides a runnable command for the hotspot feature.
type HotspotCmd struct {
	HistoryPath string
	Format      string
	TopN        int
	Action      string
	Protocol    string
	SinceHours  int
}

// NewHotspotCmd returns a HotspotCmd with sensible defaults.
func NewHotspotCmd(historyPath string) *HotspotCmd {
	return &HotspotCmd{
		HistoryPath: historyPath,
		Format:      "text",
		TopN:        10,
	}
}

// Run loads history, computes hotspots, and renders the results.
func (c *HotspotCmd) Run() error {
	h, err := New(c.HistoryPath)
	if err != nil {
		return fmt.Errorf("hotspot: load history: %w", err)
	}

	entries := h.Entries()

	opts := HotspotOptions{
		TopN:     c.TopN,
		Action:   c.Action,
		Protocol: c.Protocol,
	}
	if c.SinceHours > 0 {
		opts.Since = time.Now().Add(-time.Duration(c.SinceHours) * time.Hour)
	}

	hotspots := Hotspot(entries, opts)
	RenderHotspot(hotspots, c.Format, os.Stdout)
	return nil
}
