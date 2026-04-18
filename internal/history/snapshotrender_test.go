package history

import (
	"bytes"
	"strings"
	"testing"
)

func makeSnapshots() (Snapshot, Snapshot) {
	baseline := TakeSnapshot([]Entry{
		{Port: 22, Protocol: "tcp", Action: "allow"},
		{Port: 80, Protocol: "tcp", Action: "allow"},
	})
	current := TakeSnapshot([]Entry{
		{Port: 80, Protocol: "tcp", Action: "allow"},
		{Port: 443, Protocol: "tcp", Action: "allow"},
	})
	return baseline, current
}

func TestRenderSnapshot_Text(t *testing.T) {
	baseline, current := makeSnapshots()
	var buf bytes.Buffer
	RenderSnapshot(baseline, current, "text", &buf)
	out := buf.String()
	if !strings.Contains(out, "+ tcp/443") {
		t.Errorf("expected added port 443, got: %s", out)
	}
	if !strings.Contains(out, "- tcp/22") {
		t.Errorf("expected removed port 22, got: %s", out)
	}
}

func TestRenderSnapshot_JSON(t *testing.T) {
	baseline, current := makeSnapshots()
	var buf bytes.Buffer
	RenderSnapshot(baseline, current, "json", &buf)
	out := buf.String()
	if !strings.Contains(out, "added") || !strings.Contains(out, "removed") {
		t.Errorf("expected JSON keys, got: %s", out)
	}
}

func TestRenderSnapshot_NoChanges(t *testing.T) {
	snap := TakeSnapshot([]Entry{{Port: 80, Protocol: "tcp", Action: "allow"}})
	var buf bytes.Buffer
	RenderSnapshot(snap, snap, "text", &buf)
	if !strings.Contains(buf.String(), "No changes") {
		t.Errorf("expected no-change message, got: %s", buf.String())
	}
}

func TestRenderSnapshot_DefaultsToText(t *testing.T) {
	baseline, current := makeSnapshots()
	var buf bytes.Buffer
	RenderSnapshot(baseline, current, "", &buf)
	if !strings.Contains(buf.String(), "+") {
		t.Errorf("expected text output, got: %s", buf.String())
	}
}

func TestRenderSnapshot_NilWriter(t *testing.T) {
	baseline, current := makeSnapshots()
	// Should not panic
	RenderSnapshot(baseline, current, "text", nil)
}
