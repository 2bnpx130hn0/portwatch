package history

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

var baseReplayEntries = func() []Entry {
	now := time.Now()
	return []Entry{
		{Port: 80, Protocol: "tcp", Action: "allow", Timestamp: now.Add(-3 * time.Hour)},
		{Port: 443, Protocol: "tcp", Action: "alert", Timestamp: now.Add(-2 * time.Hour)},
		{Port: 22, Protocol: "tcp", Action: "allow", Timestamp: now.Add(-1 * time.Hour)},
		{Port: 9000, Protocol: "udp", Action: "alert", Timestamp: now},
	}
}()

func TestReplay_AllEntries(t *testing.T) {
	r := Replay(baseReplayEntries, ReplayOptions{})
	if r.Total != 4 {
		t.Fatalf("expected 4, got %d", r.Total)
	}
}

func TestReplay_ByAction(t *testing.T) {
	r := Replay(baseReplayEntries, ReplayOptions{Action: "alert"})
	if r.Total != 2 {
		t.Fatalf("expected 2, got %d", r.Total)
	}
}

func TestReplay_WithLimit(t *testing.T) {
	r := Replay(baseReplayEntries, ReplayOptions{Limit: 2})
	if r.Total != 2 {
		t.Fatalf("expected 2, got %d", r.Total)
	}
}

func TestReplay_SinceFilter(t *testing.T) {
	cutoff := time.Now().Add(-90 * time.Minute)
	r := Replay(baseReplayEntries, ReplayOptions{Since: cutoff})
	if r.Total != 2 {
		t.Fatalf("expected 2, got %d", r.Total)
	}
}

func TestRenderReplay_Text(t *testing.T) {
	var buf bytes.Buffer
	r := Replay(baseReplayEntries, ReplayOptions{})
	RenderReplay(r, "text", &buf)
	if !strings.Contains(buf.String(), "replaying 4") {
		t.Errorf("unexpected output: %s", buf.String())
	}
}

func TestRenderReplay_JSON(t *testing.T) {
	var buf bytes.Buffer
	r := Replay(baseReplayEntries, ReplayOptions{Limit: 1})
	RenderReplay(r, "json", &buf)
	if !strings.Contains(buf.String(), `"total"`) {
		t.Errorf("expected json output, got: %s", buf.String())
	}
}

func TestRenderReplay_Empty(t *testing.T) {
	var buf bytes.Buffer
	RenderReplay(ReplayResult{}, "text", &buf)
	if !strings.Contains(buf.String(), "no entries") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}
