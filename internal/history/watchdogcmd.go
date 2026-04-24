package history

import (
	"fmt"
	"os"
	"time"
)

// WatchdogCmd evaluates stored history against configured watchdog rules.
type WatchdogCmd struct {
	HistoryPath string
	Rules       []WatchdogRule
	Format      string
	Window      time.Duration
}

// NewWatchdogCmd returns a WatchdogCmd with sensible defaults.
func NewWatchdogCmd(historyPath string) *WatchdogCmd {
	return &WatchdogCmd{
		HistoryPath: historyPath,
		Format:      "text",
		Window:      60 * time.Minute,
	}
}

// AddRule appends a rule to the watchdog command.
func (c *WatchdogCmd) AddRule(port int, protocol, action string, maxCount int) {
	c.Rules = append(c.Rules, WatchdogRule{
		Port:     port,
		Protocol: protocol,
		Action:   action,
		MaxCount: maxCount,
		Window:   c.Window,
	})
}

// Run loads history, evaluates rules, and renders any violations.
func (c *WatchdogCmd) Run() error {
	if len(c.Rules) == 0 {
		fmt.Fprintln(os.Stderr, "watchdog: no rules configured")
		return nil
	}

	h, err := New(c.HistoryPath)
	if err != nil {
		return fmt.Errorf("watchdog: load history: %w", err)
	}

	entries := h.Entries()
	violations := EvalWatchdog(entries, c.Rules)
	RenderWatchdog(violations, c.Format, os.Stdout)

	if len(violations) > 0 {
		return fmt.Errorf("watchdog: %d rule(s) violated", len(violations))
	}
	return nil
}
