package history

import (
	"testing"
	"time"
)

func baseWindowEntries() []Entry {
	now := time.Now().Truncate(time.Hour)
	return []Entry{
		{Port: 80, Protocol: "tcp", Action: "alert", Timestamp: now},
		{Port: 443, Protocol: "tcp", Action: "allow", Timestamp: now.Add(10 * time.Minute)},
		{Port: 80, Protocol: "tcp", Action: "alert", Timestamp: now.Add(20 * time.Minute)},
		{Port: 8080, Protocol: "tcp", Action: "warn", Timestamp: now.Add(70 * time.Minute)},
		{Port: 22, Protocol: "tcp", Action: "alert", Timestamp: now.Add(80 * time.Minute)},
	}
}

func TestWindow_BucketsBySize(t *testing.T) {
	entries := baseWindowEntries()
	buckets := Window(entries, WindowOptions{Size: time.Hour})
	if len(buckets) < 2 {
		t.Fatalf("expected at least 2 buckets, got %d", len(buckets))
	}
	if buckets[0].Count != 3 {
		t.Errorf("expected 3 in first bucket, got %d", buckets[0].Count)
	}
}

func TestWindow_FilterByAction(t *testing.T) {
	entries := baseWindowEntries()
	buckets := Window(entries, WindowOptions{Size: time.Hour, Action: "alert"})
	for _, b := range buckets {
		for action := range b.Actions {
			if action != "alert" {
				t.Errorf("unexpected action %q in bucket", action)
			}
		}
	}
}

func TestWindow_SlidingStep(t *testing.T) {
	entries := baseWindowEntries()
	buckets := Window(entries, WindowOptions{Size: time.Hour, Step: 30 * time.Minute})
	if len(buckets) < 3 {
		t.Errorf("expected at least 3 sliding buckets, got %d", len(buckets))
	}
}

func TestWindow_SinceFilter(t *testing.T) {
	entries := baseWindowEntries()
	now := time.Now().Truncate(time.Hour)
	buckets := Window(entries, WindowOptions{Size: time.Hour, Since: now.Add(60 * time.Minute)})
	for _, b := range buckets {
		if b.Count > 2 {
			t.Errorf("since filter should exclude early entries, got count %d", b.Count)
		}
	}
}

func TestWindow_EmptyEntries(t *testing.T) {
	buckets := Window(nil, WindowOptions{Size: time.Hour})
	if len(buckets) != 0 {
		t.Errorf("expected no buckets for empty input")
	}
}

func TestWindow_ZeroSize(t *testing.T) {
	entries := baseWindowEntries()
	buckets := Window(entries, WindowOptions{Size: 0})
	if len(buckets) != 0 {
		t.Errorf("expected no buckets for zero size")
	}
}

func TestWindow_PortCounts(t *testing.T) {
	entries := baseWindowEntries()
	buckets := Window(entries, WindowOptions{Size: time.Hour})
	if buckets[0].Ports[80] != 2 {
		t.Errorf("expected port 80 count=2 in first bucket, got %d", buckets[0].Ports[80])
	}
}
