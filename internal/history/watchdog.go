package history

import (
	"fmt"
	"io"
	"os"
	"time"
)

// WatchdogRule defines a threshold-based alert rule for history entries.
type WatchdogRule struct {
	Port     int
	Protocol string
	Action   string
	MaxCount int
	Window   time.Duration
}

// WatchdogViolation describes a rule that has been triggered.
type WatchdogViolation struct {
	Rule    WatchdogRule
	Count   int
	Entries []Entry
}

// EvalWatchdog evaluates history entries against a set of watchdog rules and
// returns any violations where the observed count exceeds the rule threshold.
func EvalWatchdog(entries []Entry, rules []WatchdogRule) []WatchdogViolation {
	now := time.Now()
	var violations []WatchdogViolation

	for _, rule := range rules {
		since := now.Add(-rule.Window)
		var matched []Entry
		for _, e := range entries {
			if rule.Port > 0 && e.Port != rule.Port {
				continue
			}
			if rule.Protocol != "" && !equalFold(e.Protocol, rule.Protocol) {
				continue
			}
			if rule.Action != "" && !equalFold(e.Action, rule.Action) {
				continue
			}
			if !e.Timestamp.IsZero() && e.Timestamp.Before(since) {
				continue
			}
			matched = append(matched, e)
		}
		if len(matched) > rule.MaxCount {
			violations = append(violations, WatchdogViolation{
				Rule:    rule,
				Count:   len(matched),
				Entries: matched,
			})
		}
	}
	return violations
}

// RenderWatchdog writes watchdog violations to w in the given format.
func RenderWatchdog(violations []WatchdogViolation, format string, w io.Writer) {
	if w == nil {
		w = os.Stdout
	}
	if format == "json" {
		fmt.Fprintln(w, "[")
		for i, v := range violations {
			comma := ","
			if i == len(violations)-1 {
				comma = ""
			}
			fmt.Fprintf(w, `  {"port":%d,"protocol":%q,"action":%q,"max_count":%d,"observed":%d}%s\n`,
				v.Rule.Port, v.Rule.Protocol, v.Rule.Action, v.Rule.MaxCount, v.Count, comma)
		}
		fmt.Fprintln(w, "]")
		return
	}
	if len(violations) == 0 {
		fmt.Fprintln(w, "no watchdog violations")
		return
	}
	for _, v := range violations {
		fmt.Fprintf(w, "VIOLATION port=%d proto=%s action=%s threshold=%d observed=%d\n",
			v.Rule.Port, v.Rule.Protocol, v.Rule.Action, v.Rule.MaxCount, v.Count)
	}
}
