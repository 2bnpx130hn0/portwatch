package history

import (
	"testing"
	"time"
)

var baseDedupEntries = func() []Entry {
	now := time.Now().Truncate(time.Second)
	return []Entry{
		{Port: 80, Protocol: "tcp", Action: "allow", Timestamp: now},
		{Port: 80, Protocol: "tcp", Action: "allow", Timestamp: now},
		{Port: 443, Protocol: "tcp", Action: "alert", Timestamp: now},
		{Port: 80, Protocol: "tcp", Action: "allow", Timestamp: now.Add(2 * time.Minute)},
		{Port: 80, Protocol: "UDP", Action: "allow", Timestamp: now},
	}
}

func TestDedupe_ExactTimestampDuplicates(t *testing.T) {
	entries := baseDedupEntries()
	result := Dedupe(entries, DedupeOptions{})
	// entries[0] and entries[1] are exact duplicates; one should be removed
	if len(result) != 4 {
		t.Fatalf("expected 4 entries, got %d", len(result))
	}
}

func TestDedupe_WithinWindow(t *testing.T) {
	entries := baseDedupEntries()
	// 5-minute window: entries[0], entries[1] (exact dup), and entries[3]
	// (2 min later, same key) should all collapse to one
	result := Dedupe(entries, DedupeOptions{Window: 5 * time.Minute})
	// port80/tcp/allow collapses to 1; port443/tcp/alert = 1; port80/udp/allow = 1
	if len(result) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(result))
	}
}

func TestDedupe_OutsideWindowKept(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	entries := []Entry{
		{Port: 22, Protocol: "tcp", Action: "alert", Timestamp: now},
		{Port: 22, Protocol: "tcp", Action: "alert", Timestamp: now.Add(10 * time.Minute)},
	}
	result := Dedupe(entries, DedupeOptions{Window: 5 * time.Minute})
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
}

func TestDedupe_CaseInsensitiveProtocol(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	entries := []Entry{
		{Port: 53, Protocol: "UDP", Action: "allow", Timestamp: now},
		{Port: 53, Protocol: "udp", Action: "allow", Timestamp: now},
	}
	result := Dedupe(entries, DedupeOptions{})
	if len(result) != 1 {
		t.Fatalf("expected 1 entry after dedup, got %d", len(result))
	}
}

func TestDedupe_EmptyInput(t *testing.T) {
	result := Dedupe([]Entry{}, DedupeOptions{})
	if len(result) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(result))
	}
}

func TestDedupe_PartialFields_PortOnly(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	entries := []Entry{
		{Port: 8080, Protocol: "tcp", Action: "allow", Timestamp: now},
		{Port: 8080, Protocol: "udp", Action: "alert", Timestamp: now},
	}
	// Only matching on port: both share port 8080, so second is a dup
	result := Dedupe(entries, DedupeOptions{Fields: []string{"port"}})
	if len(result) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(result))
	}
}
