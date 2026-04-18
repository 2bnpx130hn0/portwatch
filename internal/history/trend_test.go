package history

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func baseTrendEntries() []Entry {
	now := time.Now().Truncate(time.Hour)
	return []Entry{
		{Port: 80, Protocol: "tcp", Action: "alert", Timestamp: now},
		{Port: 443, Protocol: "tcp", Action: "allow", Timestamp: now},
		{Port: 22, Protocol: "tcp", Action: "alert", Timestamp: now.Add(-2 * time.Hour)},
		{Port: 8080, Protocol: "udp", Action: "alert", Timestamp: now.Add(-2 * time.Hour)},
	}
}

func TestTrend_BucketsByHour(t *testing.T) {
	points := Trend(baseTrendEntries(), TrendOptions{})
	if len(points) != 2 {
		t.Fatalf("expected 2 buckets, got %d", len(points))
	}
	if points[0].Count != 2 || points[1].Count != 2 {
		t.Errorf("unexpected counts: %+v", points)
	}
}

func TestTrend_FilterByAction(t *testing.T) {
	points := Trend(baseTrendEntries(), TrendOptions{Action: "alert"})
	total := 0
	for _, p := range points {
		total += p.Count
	}
	if total != 3 {
		t.Errorf("expected 3 alert entries, got %d", total)
	}
}

func TestTrend_FilterByProtocol(t *testing.T) {
	points := Trend(baseTrendEntries(), TrendOptions{Protocol: "udp"})
	if len(points) != 1 || points[0].Count != 1 {
		t.Errorf("unexpected udp trend: %+v", points)
	}
}

func TestTrend_SinceFilter(t *testing.T) {
	now := time.Now().Truncate(time.Hour)
	points := Trend(baseTrendEntries(), TrendOptions{Since: now.Add(-time.Minute)})
	if len(points) != 1 {
		t.Fatalf("expected 1 bucket after since filter, got %d", len(points))
	}
}

func TestRenderTrend_Text(t *testing.T) {
	points := []TrendPoint{
		{Window: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC), Count: 5},
	}
	var buf bytes.Buffer
	if err := RenderTrend(points, "text", &buf); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "5") {
		t.Errorf("expected count in output: %s", buf.String())
	}
}

func TestRenderTrend_JSON(t *testing.T) {
	points := []TrendPoint{
		{Window: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC), Count: 3},
	}
	var buf bytes.Buffer
	if err := RenderTrend(points, "json", &buf); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "\"count\"") {
		t.Errorf("expected json output: %s", buf.String())
	}
}

func TestRenderTrend_Empty(t *testing.T) {
	var buf bytes.Buffer
	if err := RenderTrend(nil, "text", &buf); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "no trend") {
		t.Errorf("expected empty message: %s", buf.String())
	}
}
