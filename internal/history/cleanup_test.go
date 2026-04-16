package history

import (
	"testing"
	"time"
)

func seedForCleanup(t *testing.T, path string) *History {
	t.Helper()
	h, err := New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	now := time.Now().UTC()
	entries := []Entry{
		{Port: 80, Protocol: "tcp", Action: "allow", Timestamp: now.Add(-72 * time.Hour)},
		{Port: 443, Protocol: "tcp", Action: "allow", Timestamp: now.Add(-48 * time.Hour)},
		{Port: 8080, Protocol: "tcp", Action: "alert", Timestamp: now.Add(-1 * time.Hour)},
		{Port: 9090, Protocol: "tcp", Action: "warn", Timestamp: now.Add(-10 * time.Minute)},
	}
	for _, e := range entries {
		if err := h.Record(e); err != nil {
			t.Fatalf("Record: %v", err)
		}
	}
	return h
}

func TestCleanup_ByMaxAge(t *testing.T) {
	h := seedForCleanup(t, tempPath(t))
	removed, err := h.Cleanup(CleanupOptions{MaxAge: 36 * time.Hour})
	if err != nil {
		t.Fatalf("Cleanup: %v", err)
	}
	if removed != 2 {
		t.Errorf("expected 2 removed, got %d", removed)
	}
	if len(h.entries) != 2 {
		t.Errorf("expected 2 remaining, got %d", len(h.entries))
	}
}

func TestCleanup_ByMaxEntries(t *testing.T) {
	h := seedForCleanup(t, tempPath(t))
	removed, err := h.Cleanup(CleanupOptions{MaxEntries: 2})
	if err != nil {
		t.Fatalf("Cleanup: %v", err)
	}
	if removed != 2 {
		t.Errorf("expected 2 removed, got %d", removed)
	}
	if len(h.entries) != 2 {
		t.Errorf("expected 2 remaining, got %d", len(h.entries))
	}
}

func TestCleanup_NoLimits(t *testing.T) {
	h := seedForCleanup(t, tempPath(t))
	removed, err := h.Cleanup(CleanupOptions{})
	if err != nil {
		t.Fatalf("Cleanup: %v", err)
	}
	if removed != 0 {
		t.Errorf("expected 0 removed, got %d", removed)
	}
}

func TestCleanup_PersistsToDisk(t *testing.T) {
	path := tempPath(t)
	h := seedForCleanup(t, path)
	if _, err := h.Cleanup(CleanupOptions{MaxEntries: 1}); err != nil {
		t.Fatalf("Cleanup: %v", err)
	}
	h2, err := New(path)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if len(h2.entries) != 1 {
		t.Errorf("expected 1 persisted entry, got %d", len(h2.entries))
	}
}
