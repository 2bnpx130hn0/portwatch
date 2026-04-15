package alert_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/rules"
)

func makeEvent(action rules.Action, port int, protocol string) alert.Event {
	return alert.Event{
		Timestamp: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		Port:      port,
		Protocol:  protocol,
		Action:    action,
		Message:   "test message",
	}
}

func TestNotify_AlertAction(t *testing.T) {
	var buf bytes.Buffer
	n := alert.New(&buf)

	err := n.Notify(makeEvent(rules.ActionAlert, 8080, "tcp"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "ALERT") {
		t.Errorf("expected ALERT level in output, got: %s", out)
	}
	if !strings.Contains(out, "port=8080") {
		t.Errorf("expected port=8080 in output, got: %s", out)
	}
}

func TestNotify_WarnAction(t *testing.T) {
	var buf bytes.Buffer
	n := alert.New(&buf)

	err := n.Notify(makeEvent(rules.ActionWarn, 443, "tcp"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "WARN") {
		t.Errorf("expected WARN level in output, got: %s", out)
	}
}

func TestNotify_AllowAction(t *testing.T) {
	var buf bytes.Buffer
	n := alert.New(&buf)

	err := n.Notify(makeEvent(rules.ActionAllow, 22, "tcp"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "INFO") {
		t.Errorf("expected INFO level in output, got: %s", out)
	}
}

func TestNotify_ExplicitLevel(t *testing.T) {
	var buf bytes.Buffer
	n := alert.New(&buf)

	e := makeEvent(rules.ActionAlert, 9090, "udp")
	e.Level = alert.LevelInfo // override

	_ = n.Notify(e)
	out := buf.String()
	if !strings.Contains(out, "INFO") {
		t.Errorf("expected explicit INFO level, got: %s", out)
	}
}

func TestNew_NilWriterUsesStderr(t *testing.T) {
	// Should not panic when w is nil
	n := alert.New(nil)
	if n == nil {
		t.Fatal("expected non-nil Notifier")
	}
}
