package history

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestArchive_MovesOldEntries(t *testing.T) {
	h := newTestHistory(t)
	seedOldAndNew(t, h)

	dir := t.TempDir()
	a := NewArchiver(h, dir)

	n, err := a.Archive(24 * time.Hour)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 1 {
		t.Fatalf("expected 1 archived entry, got %d", n)
	}
	if len(h.entries) != 1 {
		t.Fatalf("expected 1 remaining entry, got %d", len(h.entries))
	}
}

func TestArchive_CreatesFile(t *testing.T) {
	h := newTestHistory(t)
	seedOldAndNew(t, h)

	dir := t.TempDir()
	a := NewArchiver(h, dir)

	_, err := a.Archive(24 * time.Hour)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entries, _ := os.ReadDir(dir)
	if len(entries) != 1 {
		t.Fatalf("expected 1 archive file, got %d", len(entries))
	}

	path := filepath.Join(dir, entries[0].Name())
	f, _ := os.Open(path)
	defer f.Close()

	var archived []Entry
	if err := json.NewDecoder(f).Decode(&archived); err != nil {
		t.Fatalf("decode archive: %v", err)
	}
	if len(archived) != 1 {
		t.Fatalf("expected 1 entry in archive file, got %d", len(archived))
	}
}

func TestArchive_NothingToArchive(t *testing.T) {
	h := newTestHistory(t)
	h.entries = []Entry{{Port: 80, Protocol: "tcp", Action: "allow", Timestamp: time.Now()}}

	dir := t.TempDir()
	a := NewArchiver(h, dir)

	n, err := a.Archive(24 * time.Hour)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 0 {
		t.Fatalf("expected 0 archived, got %d", n)
	}
}

func newTestHistory(t *testing.T) *History {
	t.Helper()
	h, err := New(tempPath(t))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return h
}

func seedOldAndNew(t *testing.T, h *History) {
	t.Helper()
	h.entries = []Entry{
		{Port: 22, Protocol: "tcp", Action: "alert", Timestamp: time.Now().Add(-48 * time.Hour)},
		{Port: 80, Protocol: "tcp", Action: "allow", Timestamp: time.Now()},
	}
}
