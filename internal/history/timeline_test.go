package history

import (
	"testing"
	"time"
)

func baseTimelineEntries() []Entry {
	now := time.Now().Truncate(time.Hour)
	return []Entry{
		{Port: 80, Protocol: "tcp", Action: "alert", Timestamp: now},
		{Port: 443, Protocol: "tcp", Action: "allow", Timestamp: now},
		{Port: 22, Protocol: "tcp", Action: "alert", Timestamp: now.Add(time.Hour)},
		{Port: 8080, Protocol: "tcp", Action: "warn", Timestamp: now.Add(2 * time.Hour)},
	}
}

func TestTimeline_BucketsByHour(t *testing.T) {
	entries := baseTimelineEntries()
	result := Timeline(entries, TimelineOptions{BucketSize: time.Hour})
	if len(result) != 3 {
		t.Fatalf("expected 3 buckets, got %d", len(result))
	}
	if result[0].Total != 2 {
		t.Errorf("expected 2 entries in first bucket, got %d", result[0].Total)
	}
}

func TestTimeline_ByAction(t *testing.T) {
	entries := baseTimelineEntries()
	result := Timeline(entries, TimelineOptions{BucketSize: time.Hour})
	if result[0].ByAction["alert"] != 1 {
		t.Errorf("expected 1 alert in first bucket")
	}
	if result[0].ByAction["allow"] != 1 {
		t.Errorf("expected 1 allow in first bucket")
	}
}

func TestTimeline_SinceFilter(t *testing.T) {
	entries := baseTimelineEntries()
	now := time.Now().Truncate(time.Hour)
	result := Timeline(entries, TimelineOptions{
		BucketSize: time.Hour,
		Since:      now.Add(time.Hour),
	})
	if len(result) != 2 {
		t.Fatalf("expected 2 buckets after since filter, got %d", len(result))
	}
}

func TestTimeline_DefaultBucketSize(t *testing.T) {
	entries := baseTimelineEntries()
	result := Timeline(entries, TimelineOptions{})
	if len(result) == 0 {
		t.Error("expected non-empty timeline with default bucket size")
	}
}

func TestTimeline_Empty(t *testing.T) {
	result := Timeline([]Entry{}, TimelineOptions{BucketSize: time.Hour})
	if len(result) != 0 {
		t.Errorf("expected empty timeline")
	}
}
