package history

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func sampleWatchSummary() WatchSummary {
	base := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	return NewWatchSummary([]Entry{
		{Port: 8080, Protocol: "tcp", Action: "alert", Timestamp: base.Add(time.Second)},
		{Port: 443, Protocol: "tcp", Action: "allow", Timestamp: base.Add(2 * time.Second)},
		{Port: 22, Protocol: "tcp", Action: "removed", Timestamp: base.Add(3 * time.Second)},
	}, base)
}

func TestRenderWatchSummary_Text(t *testing.T) {
	var buf bytes.Buffer
	if err := RenderWatchSummary(sampleWatchSummary(), "text", &buf); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "Cycle:") {
		t.Error("expected Cycle label")
	}
	if !strings.Contains(out, "alerted: 1") {
		t.Error("expected alerted count")
	}
	if !strings.Contains(out, "tcp:22") {
		t.Error("expected removed port 22")
	}
}

func TestRenderWatchSummary_JSON(t *testing.T) {
	var buf bytes.Buffer
	if err := RenderWatchSummary(sampleWatchSummary(), "json", &buf); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, `"alerted"`) {
		t.Error("expected alerted field in JSON")
	}
	if !strings.Contains(out, `"removed"`) {
		t.Error("expected removed field in JSON")
	}
}

func TestRenderWatchSummary_DefaultsToText(t *testing.T) {
	var buf bytes.Buffer
	if err := RenderWatchSummary(sampleWatchSummary(), "", &buf); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "Cycle:") {
		t.Error("expected text output as default")
	}
}

func TestRenderWatchSummary_NilWriter(t *testing.T) {
	// should not panic
	if err := RenderWatchSummary(sampleWatchSummary(), "text", nil); err != nil {
		t.Fatal(err)
	}
}
