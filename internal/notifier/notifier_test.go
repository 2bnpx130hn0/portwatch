package notifier_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/portwatch/internal/notifier"
)

var fixedTime = time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

func makeEvent(protocol string, port uint16, action, level string) notifier.Event {
	return notifier.Event{
		Timestamp: fixedTime,
		Protocol:  protocol,
		Port:      port,
		Action:    action,
		Level:     level,
	}
}

func TestSend_TextFormat(t *testing.T) {
	var buf bytes.Buffer
	n := notifier.New(notifier.FormatText, &buf)

	err := n.Send(makeEvent("tcp", 8080, "added", "alert"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "[alert]") {
		t.Errorf("expected level in output, got: %s", out)
	}
	if !strings.Contains(out, "tcp/8080") {
		t.Errorf("expected port in output, got: %s", out)
	}
	if !strings.Contains(out, "added") {
		t.Errorf("expected action in output, got: %s", out)
	}
}

func TestSend_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	n := notifier.New(notifier.FormatJSON, &buf)

	err := n.Send(makeEvent("udp", 53, "removed", "warn"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, `"level":"warn"`) {
		t.Errorf("expected level field, got: %s", out)
	}
	if !strings.Contains(out, `"port":53`) {
		t.Errorf("expected port field, got: %s", out)
	}
	if !strings.Contains(out, `"action":"removed"`) {
		t.Errorf("expected action field, got: %s", out)
	}
}

func TestSend_DefaultsToText(t *testing.T) {
	var buf bytes.Buffer
	n := notifier.New("", &buf)

	err := n.Send(makeEvent("tcp", 443, "added", "info"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	// text format uses RFC3339, not JSON braces
	if strings.HasPrefix(out, "{") {
		t.Errorf("expected text format, got JSON-like output: %s", out)
	}
}

func TestSend_NilWriterUsesStdout(t *testing.T) {
	// Just ensure New does not panic with nil writer
	n := notifier.New(notifier.FormatText, nil)
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}
