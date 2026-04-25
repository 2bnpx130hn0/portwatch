package history

import (
	"testing"
	"time"
)

func baseDecayEntries() []Entry {
	now := time.Now()
	return []Entry{
		{Port: 80, Protocol: "tcp", Action: "allow", Timestamp: now.Add(-1 * time.Hour)},
		{Port: 443, Protocol: "tcp", Action: "alert", Timestamp: now.Add(-2 * time.Hour)},
		{Port: 8080, Protocol: "tcp", Action: "warn", Timestamp: now.Add(-30 * time.Minute)},
		{Port: 9000, Protocol: "tcp", Action: "alert", Timestamp: now.Add(-10 * time.Minute)},
	}
}

func TestDecay_AlertScoresHigherThanAllow(t *testing.T) {
	entries := baseDecayEntries()
	results := Decay(entries, DecayOptions{})

	if len(results) != 4 {
		t.Fatalf("expected 4 results, got %d", len(results))
	}
	// The recent alert (9000) should rank first.
	if results[0].Entry.Port != 9000 {
		t.Errorf("expected port 9000 first, got %d", results[0].Entry.Port)
	}
}

func TestDecay_OlderEntriesScoreLower(t *testing.T) {
	now := time.Now()
	entries := []Entry{
		{Port: 1, Protocol: "tcp", Action: "alert", Timestamp: now.Add(-48 * time.Hour)},
		{Port: 2, Protocol: "tcp", Action: "alert", Timestamp: now.Add(-1 * time.Hour)},
	}
	results := Decay(entries, DecayOptions{HalfLife: 24 * time.Hour})

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0].Entry.Port != 2 {
		t.Errorf("expected newer entry first, got port %d", results[0].Entry.Port)
	}
	if results[0].Decayed <= results[1].Decayed {
		t.Errorf("expected newer entry to have higher decayed score")
	}
}

func TestDecay_SinceFilter(t *testing.T) {
	now := time.Now()
	entries := []Entry{
		{Port: 80, Protocol: "tcp", Action: "alert", Timestamp: now.Add(-5 * time.Hour)},
		{Port: 443, Protocol: "tcp", Action: "alert", Timestamp: now.Add(-1 * time.Hour)},
	}
	results := Decay(entries, DecayOptions{Since: now.Add(-2 * time.Hour)})

	if len(results) != 1 {
		t.Fatalf("expected 1 result after since filter, got %d", len(results))
	}
	if results[0].Entry.Port != 443 {
		t.Errorf("expected port 443, got %d", results[0].Entry.Port)
	}
}

func TestDecay_ActionFilter(t *testing.T) {
	entries := baseDecayEntries()
	results := Decay(entries, DecayOptions{Action: "alert"})

	for _, r := range results {
		if r.Entry.Action != "alert" {
			t.Errorf("expected only alert entries, got %q", r.Entry.Action)
		}
	}
	if len(results) != 2 {
		t.Errorf("expected 2 alert entries, got %d", len(results))
	}
}

func TestDecay_EmptyEntries(t *testing.T) {
	results := Decay([]Entry{}, DecayOptions{})
	if len(results) != 0 {
		t.Errorf("expected empty results, got %d", len(results))
	}
}

func TestDecay_DefaultHalfLife(t *testing.T) {
	now := time.Now()
	entries := []Entry{
		{Port: 22, Protocol: "tcp", Action: "warn", Timestamp: now.Add(-1 * time.Minute)},
	}
	// Zero HalfLife should default to 24h without panicking.
	results := Decay(entries, DecayOptions{HalfLife: 0})
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Decayed <= 0 {
		t.Errorf("expected positive decayed score, got %f", results[0].Decayed)
	}
}
