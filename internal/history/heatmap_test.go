package history

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func baseHeatmapEntries() []Entry {
	base := time.Date(2024, 1, 8, 10, 0, 0, 0, time.UTC) // Monday 10:00
	return []Entry{
		{Port: 80, Protocol: "tcp", Action: "alert", Timestamp: base},
		{Port: 80, Protocol: "tcp", Action: "alert", Timestamp: base.Add(time.Minute)},
		{Port: 443, Protocol: "tcp", Action: "allow", Timestamp: base.Add(2 * time.Hour)}, // Monday 12:00
		{Port: 22, Protocol: "tcp", Action: "alert", Timestamp: base.Add(24 * time.Hour)}, // Tuesday 10:00
	}
}

func TestHeatmap_CountsCorrectCell(t *testing.T) {
	cells := Heatmap(baseHeatmapEntries(), HeatmapOptions{})
	if len(cells) != 7*24 {
		t.Fatalf("expected 168 cells, got %d", len(cells))
	}
	// Monday = 1, hour 10 → 2 events
	for _, c := range cells {
		if c.DayOfWeek == time.Monday && c.Hour == 10 {
			if c.Count != 2 {
				t.Errorf("expected 2 at Mon/10, got %d", c.Count)
			}
		}
	}
}

func TestHeatmap_FilterByAction(t *testing.T) {
	cells := Heatmap(baseHeatmapEntries(), HeatmapOptions{Action: "allow"})
	total := 0
	for _, c := range cells {
		total += c.Count
	}
	if total != 1 {
		t.Errorf("expected 1 allow event, got %d", total)
	}
}

func TestHeatmap_SinceFilter(t *testing.T) {
	base := time.Date(2024, 1, 8, 10, 0, 0, 0, time.UTC)
	cells := Heatmap(baseHeatmapEntries(), HeatmapOptions{Since: base.Add(23 * time.Hour)})
	total := 0
	for _, c := range cells {
		total += c.Count
	}
	if total != 1 {
		t.Errorf("expected 1 entry after since filter, got %d", total)
	}
}

func TestRenderHeatmap_Text(t *testing.T) {
	cells := Heatmap(baseHeatmapEntries(), HeatmapOptions{})
	var buf bytes.Buffer
	RenderHeatmap(cells, "text", &buf)
	out := buf.String()
	if !strings.Contains(out, "Mon") {
		t.Error("expected 'Mon' in text output")
	}
	if !strings.Contains(out, "Tue") {
		t.Error("expected 'Tue' in text output")
	}
}

func TestRenderHeatmap_JSON(t *testing.T) {
	cells := Heatmap(baseHeatmapEntries(), HeatmapOptions{})
	var buf bytes.Buffer
	RenderHeatmap(cells, "json", &buf)
	var out []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(out) != 7*24 {
		t.Errorf("expected 168 JSON cells, got %d", len(out))
	}
}

func TestRenderHeatmap_DefaultsToText(t *testing.T) {
	cells := Heatmap(baseHeatmapEntries(), HeatmapOptions{})
	var buf bytes.Buffer
	RenderHeatmap(cells, "", &buf)
	if !strings.Contains(buf.String(), "Sun") {
		t.Error("default format should be text and contain day abbreviations")
	}
}
