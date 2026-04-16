package history

import (
	"testing"
	"time"
)

func seedForPrune(t *testing.T) *History {
	t.Helper()
	h := New(tempPath(t))
	now := time.Now()
	for i := 0; i < 10; i++ {
		e := Entry{
			Port:      i + 1,
			Protocol:  "tcp",
			Action:    "allow",
			Timestamp: now.Add(-time.Duration(10-i) * time.Hour),
		}
		h.entries = append(h.entries, e)
	}
	return h
}

func TestPrune_ByMaxAge(t *testing.T) {
	h := seedForPrune(t)
	removed, err := h.Prune(PruneOptions{MaxAge: 5 * time.Hour})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if removed != 5 {
		t.Errorf("expected 5 removed, got %d", removed)
	}
	if len(h.entries) != 5 {
		t.Errorf("expected 5 remaining, got %d", len(h.entries))
	}
}

func TestPrune_ByMaxEntries(t *testing.T) {
	h := seedForPrune(t)
	removed, err := h.Prune(PruneOptions{MaxEntries: 4})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if removed != 6 {
		t.Errorf("expected 6 removed, got %d", removed)
	}
	if len(h.entries) != 4 {
		t.Errorf("expected 4 remaining, got %d", len(h.entries))
	}
}

func TestPrune_Combined(t *testing.T) {
	h := seedForPrune(t)
	// age prune leaves 5, then max entries trims to 3
	removed, err := h.Prune(PruneOptions{MaxAge: 5 * time.Hour, MaxEntries: 3})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if removed != 7 {
		t.Errorf("expected 7 removed, got %d", removed)
	}
}

func TestPrune_NoOptions(t *testing.T) {
	h := seedForPrune(t)
	removed, err := h.Prune(PruneOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if removed != 0 {
		t.Errorf("expected 0 removed, got %d", removed)
	}
}
