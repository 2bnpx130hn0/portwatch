package history

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func sampleStats() []PortStat {
	return []PortStat{
		{
			Port: 80, Protocol: "tcp", Seen: 5,
			LastSeen: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			Actions: map[string]int{"allow": 4, "alert": 1},
		},
		{
			Port: 443, Protocol: "tcp", Seen: 2,
			LastSeen: time.Date(2024, 1, 15, 11, 0, 0, 0, time.UTC),
			Actions: map[string]int{"allow": 2},
		},
	}
}

func TestRenderStats_TextFormat(t *testing.T) {
	var buf bytes.Buffer
	if err := RenderStats(&buf, sampleStats(), "text"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "PORT") {
		t.Error("expected header row")
	}
	if !strings.Contains(out, "80") {
		t.Error("expected port 80 in output")
	}
}

func TestRenderStats_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	if err := RenderStats(&buf, sampleStats(), "json"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out []PortStat
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if len(out) != 2 {
		t.Errorf("expected 2 entries, got %d", len(out))
	}
}

func TestRenderStats_DefaultsToText(t *testing.T) {
	var buf bytes.Buffer
	if err := RenderStats(&buf, sampleStats(), ""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "PROTOCOL") {
		t.Error("expected text output")
	}
}

func TestRenderStats_SortedBySeen(t *testing.T) {
	var buf bytes.Buffer
	_ = RenderStats(&buf, sampleStats(), "text")
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	// header + 2 data rows
	if len(lines) < 3 {
		t.Fatalf("expected at least 3 lines")
	}
	if !strings.Contains(lines[1], "80") {
		t.Error("expected port 80 first (higher seen count)")
	}
}
