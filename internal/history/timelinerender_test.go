package history

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func sampleTimeline() []TimelineEntry {
	now := time.Now().Truncate(time.Hour)
	return []TimelineEntry{
		{Bucket: now, Total: 3, ByAction: map[string]int{"alert": 2, "allow": 1}},
		{Bucket: now.Add(time.Hour), Total: 1, ByAction: map[string]int{"warn": 1}},
	}
}

func TestRenderTimeline_Text(t *testing.T) {
	var buf bytes.Buffer
	if err := RenderTimeline(sampleTimeline(), "text", &buf); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "total=3") {
		t.Errorf("expected total=3 in output, got: %s", out)
	}
}

func TestRenderTimeline_JSON(t *testing.T) {
	var buf bytes.Buffer
	if err := RenderTimeline(sampleTimeline(), "json", &buf); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "Total") {
		t.Errorf("expected JSON output with Total field")
	}
}

func TestRenderTimeline_DefaultsToText(t *testing.T) {
	var buf bytes.Buffer
	if err := RenderTimeline(sampleTimeline(), "", &buf); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "total=") {
		t.Errorf("expected text output as default")
	}
}

func TestRenderTimeline_Empty(t *testing.T) {
	var buf bytes.Buffer
	if err := RenderTimeline([]TimelineEntry{}, "text", &buf); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "No timeline data") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}
