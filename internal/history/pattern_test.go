package history

import (
	"testing"
	"time"
)

func basePatternEntries() []Entry {
	now := time.Now()
	return []Entry{
		{Port: 80, Protocol: "tcp", Action: "allow", Timestamp: now.Add(-5 * time.Hour)},
		{Port: 80, Protocol: "tcp", Action: "allow", Timestamp: now.Add(-3 * time.Hour)},
		{Port: 80, Protocol: "tcp", Action: "allow", Timestamp: now.Add(-1 * time.Hour)},
		{Port: 443, Protocol: "tcp", Action: "alert", Timestamp: now.Add(-4 * time.Hour)},
		{Port: 443, Protocol: "tcp", Action: "alert", Timestamp: now.Add(-2 * time.Hour)},
		{Port: 9999, Protocol: "udp", Action: "warn", Timestamp: now.Add(-6 * time.Hour)},
	}
}

func TestDetectPatterns_FindsRecurring(t *testing.T) {
	entries := basePatternEntries()
	results := DetectPatterns(entries, PatternOptions{MinOccurrences: 2})
	if len(results) != 2 {
		t.Fatalf("expected 2 patterns, got %d", len(results))
	}
	// highest count first
	if results[0].Port != 80 || results[0].Count != 3 {
		t.Errorf("expected port 80 count 3, got port %d count %d", results[0].Port, results[0].Count)
	}
	if results[1].Port != 443 || results[1].Count != 2 {
		t.Errorf("expected port 443 count 2, got port %d count %d", results[1].Port, results[1].Count)
	}
}

func TestDetectPatterns_MinOccurrencesFilters(t *testing.T) {
	entries := basePatternEntries()
	results := DetectPatterns(entries, PatternOptions{MinOccurrences: 3})
	if len(results) != 1 {
		t.Fatalf("expected 1 pattern, got %d", len(results))
	}
	if results[0].Port != 80 {
		t.Errorf("expected port 80, got %d", results[0].Port)
	}
}

func TestDetectPatterns_FilterByAction(t *testing.T) {
	entries := basePatternEntries()
	results := DetectPatterns(entries, PatternOptions{MinOccurrences: 2, Action: "alert"})
	if len(results) != 1 {
		t.Fatalf("expected 1 pattern, got %d", len(results))
	}
	if results[0].Action != "alert" {
		t.Errorf("expected action alert, got %s", results[0].Action)
	}
}

func TestDetectPatterns_FilterByProtocol(t *testing.T) {
	entries := basePatternEntries()
	// add a second udp entry so it meets threshold
	entries = append(entries, Entry{Port: 9999, Protocol: "udp", Action: "warn", Timestamp: time.Now().Add(-1 * time.Hour)})
	results := DetectPatterns(entries, PatternOptions{MinOccurrences: 2, Protocol: "udp"})
	if len(results) != 1 {
		t.Fatalf("expected 1 pattern, got %d", len(results))
	}
	if results[0].Protocol != "udp" {
		t.Errorf("expected protocol udp, got %s", results[0].Protocol)
	}
}

func TestDetectPatterns_SinceFilter(t *testing.T) {
	entries := basePatternEntries()
	// only entries within last 2h: port 80 once, port 443 once — neither hits threshold
	results := DetectPatterns(entries, PatternOptions{
		MinOccurrences: 2,
		Since:          time.Now().Add(-2 * time.Hour),
	})
	if len(results) != 0 {
		t.Errorf("expected 0 patterns after since filter, got %d", len(results))
	}
}

func TestDetectPatterns_AvgInterval(t *testing.T) {
	now := time.Now()
	entries := []Entry{
		{Port: 22, Protocol: "tcp", Action: "alert", Timestamp: now.Add(-4 * time.Hour)},
		{Port: 22, Protocol: "tcp", Action: "alert", Timestamp: now.Add(-2 * time.Hour)},
		{Port: 22, Protocol: "tcp", Action: "alert", Timestamp: now},
	}
	results := DetectPatterns(entries, PatternOptions{MinOccurrences: 2})
	if len(results) == 0 {
		t.Fatal("expected at least one result")
	}
	expected := 2 * 3600.0
	if diff := results[0].AvgIntervalSecs - expected; diff > 1 || diff < -1 {
		t.Errorf("expected avg interval ~%.0f, got %.2f", expected, results[0].AvgIntervalSecs)
	}
}
