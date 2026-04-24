package history

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func baseWatchdogEntries() []Entry {
	now := time.Now()
	return []Entry{
		{Port: 80, Protocol: "tcp", Action: "alert", Timestamp: now.Add(-1 * time.Minute)},
		{Port: 80, Protocol: "tcp", Action: "alert", Timestamp: now.Add(-2 * time.Minute)},
		{Port: 80, Protocol: "tcp", Action: "alert", Timestamp: now.Add(-3 * time.Minute)},
		{Port: 443, Protocol: "tcp", Action: "allow", Timestamp: now.Add(-1 * time.Minute)},
		{Port: 22, Protocol: "tcp", Action: "alert", Timestamp: now.Add(-90 * time.Minute)},
	}
}

func TestEvalWatchdog_DetectsViolation(t *testing.T) {
	entries := baseWatchdogEntries()
	rules := []WatchdogRule{
		{Port: 80, Protocol: "tcp", Action: "alert", MaxCount: 2, Window: 10 * time.Minute},
	}
	violations := EvalWatchdog(entries, rules)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Count != 3 {
		t.Errorf("expected count 3, got %d", violations[0].Count)
	}
}

func TestEvalWatchdog_NoViolationBelowThreshold(t *testing.T) {
	entries := baseWatchdogEntries()
	rules := []WatchdogRule{
		{Port: 80, Protocol: "tcp", Action: "alert", MaxCount: 5, Window: 10 * time.Minute},
	}
	violations := EvalWatchdog(entries, rules)
	if len(violations) != 0 {
		t.Errorf("expected no violations, got %d", len(violations))
	}
}

func TestEvalWatchdog_WindowExcludesOldEntries(t *testing.T) {
	entries := baseWatchdogEntries()
	// port 22 alert is 90 minutes ago, window is 60 minutes
	rules := []WatchdogRule{
		{Port: 22, Protocol: "tcp", Action: "alert", MaxCount: 0, Window: 60 * time.Minute},
	}
	violations := EvalWatchdog(entries, rules)
	if len(violations) != 0 {
		t.Errorf("expected no violations (entry outside window), got %d", len(violations))
	}
}

func TestEvalWatchdog_MultipleRules(t *testing.T) {
	entries := baseWatchdogEntries()
	rules := []WatchdogRule{
		{Port: 80, Protocol: "tcp", Action: "alert", MaxCount: 2, Window: 10 * time.Minute},
		{Port: 443, Protocol: "tcp", Action: "allow", MaxCount: 0, Window: 10 * time.Minute},
	}
	violations := EvalWatchdog(entries, rules)
	if len(violations) != 2 {
		t.Fatalf("expected 2 violations, got %d", len(violations))
	}
}

func TestRenderWatchdog_Text(t *testing.T) {
	violations := []WatchdogViolation{
		{Rule: WatchdogRule{Port: 80, Protocol: "tcp", Action: "alert", MaxCount: 2}, Count: 5},
	}
	var buf bytes.Buffer
	RenderWatchdog(violations, "text", &buf)
	out := buf.String()
	if !strings.Contains(out, "VIOLATION") {
		t.Errorf("expected VIOLATION in output, got: %s", out)
	}
	if !strings.Contains(out, "port=80") {
		t.Errorf("expected port=80 in output, got: %s", out)
	}
}

func TestRenderWatchdog_NoViolations(t *testing.T) {
	var buf bytes.Buffer
	RenderWatchdog(nil, "text", &buf)
	if !strings.Contains(buf.String(), "no watchdog violations") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}
