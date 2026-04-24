package history

import (
	"testing"
	"time"
)

func baseCorrelateEntries() []Entry {
	now := time.Now()
	return []Entry{
		{Port: 80, Protocol: "tcp", Action: "alert", Timestamp: now},
		{Port: 80, Protocol: "tcp", Action: "alert", Timestamp: now.Add(1 * time.Minute)},
		{Port: 80, Protocol: "tcp", Action: "alert", Timestamp: now.Add(2 * time.Minute)},
		{Port: 443, Protocol: "tcp", Action: "allow", Timestamp: now},
		{Port: 443, Protocol: "tcp", Action: "allow", Timestamp: now.Add(30 * time.Second)},
		{Port: 9000, Protocol: "udp", Action: "alert", Timestamp: now},
		{Port: 9000, Protocol: "udp", Action: "alert", Timestamp: now.Add(10 * time.Hour)}, // outside window
	}
}

func TestCorrelate_GroupsByPortProtocol(t *testing.T) {
	entries := baseCorrelateEntries()
	result := Correlate(entries, CorrelateOptions{Window: 5 * time.Minute, MinEntries: 2})

	if len(result) != 2 {
		t.Fatalf("expected 2 correlations, got %d", len(result))
	}
}

func TestCorrelate_OrderedByScoreDescending(t *testing.T) {
	entries := baseCorrelateEntries()
	result := Correlate(entries, CorrelateOptions{Window: 5 * time.Minute, MinEntries: 2})

	if result[0].Score < result[1].Score {
		t.Errorf("expected descending score order, got %.0f then %.0f", result[0].Score, result[1].Score)
	}
	if result[0].Port != 80 {
		t.Errorf("expected port 80 first (3 entries), got %d", result[0].Port)
	}
}

func TestCorrelate_FilterByAction(t *testing.T) {
	entries := baseCorrelateEntries()
	result := Correlate(entries, CorrelateOptions{
		Window:     5 * time.Minute,
		MinEntries: 2,
		Action:     "allow",
	})

	if len(result) != 1 {
		t.Fatalf("expected 1 correlation for action=allow, got %d", len(result))
	}
	if result[0].Port != 443 {
		t.Errorf("expected port 443, got %d", result[0].Port)
	}
}

func TestCorrelate_RespectsWindow(t *testing.T) {
	entries := baseCorrelateEntries()
	// port 9000 udp entries are 10 hours apart — should not correlate
	result := Correlate(entries, CorrelateOptions{Window: 5 * time.Minute, MinEntries: 2})

	for _, c := range result {
		if c.Port == 9000 {
			t.Errorf("port 9000 should not correlate within 5m window")
		}
	}
}

func TestCorrelate_MinEntriesNotMet(t *testing.T) {
	now := time.Now()
	entries := []Entry{
		{Port: 22, Protocol: "tcp", Action: "alert", Timestamp: now},
	}
	result := Correlate(entries, CorrelateOptions{Window: 5 * time.Minute, MinEntries: 2})
	if len(result) != 0 {
		t.Errorf("expected no correlations for single entry, got %d", len(result))
	}
}

func TestCorrelate_DefaultOptions(t *testing.T) {
	now := time.Now()
	entries := []Entry{
		{Port: 8080, Protocol: "tcp", Action: "alert", Timestamp: now},
		{Port: 8080, Protocol: "tcp", Action: "alert", Timestamp: now.Add(1 * time.Minute)},
	}
	// zero-value options should use defaults
	result := Correlate(entries, CorrelateOptions{})
	if len(result) != 1 {
		t.Errorf("expected 1 correlation with default options, got %d", len(result))
	}
}
