package history

import (
	"testing"
	"time"
)

func baseWatchEntries(base time.Time) []Entry {
	return []Entry{
		{Port: 80, Protocol: "tcp", Action: "allow", Timestamp: base.Add(1 * time.Second)},
		{Port: 443, Protocol: "tcp", Action: "alert", Timestamp: base.Add(2 * time.Second)},
		{Port: 22, Protocol: "tcp", Action: "removed", Timestamp: base.Add(3 * time.Second)},
		{Port: 9000, Protocol: "tcp", Action: "allow", Timestamp: base.Add(-1 * time.Second)}, // before cycle
	}
}

func TestWatchSummary_FiltersBeforeCycle(t *testing.T) {
	base := time.Now()
	entries := baseWatchEntries(base)
	s := NewWatchSummary(entries, base)
	if len(s.Added)+len(s.Removed) != 3 {
		t.Fatalf("expected 3 events, got %d", len(s.Added)+len(s.Removed))
	}
}

func TestWatchSummary_Counts(t *testing.T) {
	base := time.Now()
	entries := baseWatchEntries(base)
	s := NewWatchSummary(entries, base)
	if s.Alerted != 1 {
		t.Errorf("expected 1 alerted, got %d", s.Alerted)
	}
	if s.Allowed != 1 {
		t.Errorf("expected 1 allowed, got %d", s.Allowed)
	}
}

func TestWatchSummary_Removed(t *testing.T) {
	base := time.Now()
	entries := baseWatchEntries(base)
	s := NewWatchSummary(entries, base)
	if len(s.Removed) != 1 || s.Removed[0].Port != 22 {
		t.Errorf("expected removed port 22, got %+v", s.Removed)
	}
}

func TestWatchSummary_HasChanges(t *testing.T) {
	base := time.Now()
	s := NewWatchSummary(baseWatchEntries(base), base)
	if !s.HasChanges() {
		t.Error("expected HasChanges to be true")
	}
}

func TestWatchSummary_NoChanges(t *testing.T) {
	base := time.Now()
	// all entries before cycle
	entries := []Entry{
		{Port: 80, Protocol: "tcp", Action: "allow", Timestamp: base.Add(-5 * time.Second)},
	}
	s := NewWatchSummary(entries, base)
	if s.HasChanges() {
		t.Error("expected HasChanges to be false")
	}
}
