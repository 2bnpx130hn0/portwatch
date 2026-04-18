package history

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func baseSnapshotEntries() []Entry {
	return []Entry{
		{Port: 80, Protocol: "tcp", Action: "allow", Timestamp: time.Now()},
		{Port: 443, Protocol: "tcp", Action: "allow", Timestamp: time.Now()},
	}
}

func TestTakeSnapshot_CopiesEntries(t *testing.T) {
	entries := baseSnapshotEntries()
	snap := TakeSnapshot(entries)
	if len(snap.Entries) != len(entries) {
		t.Fatalf("expected %d entries, got %d", len(entries), len(snap.Entries))
	}
	if snap.TakenAt.IsZero() {
		t.Fatal("expected non-zero TakenAt")
	}
}

func TestSaveAndLoadSnapshot(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")
	original := TakeSnapshot(baseSnapshotEntries())
	if err := SaveSnapshot(path, original); err != nil {
		t.Fatalf("save: %v", err)
	}
	loaded, err := LoadSnapshot(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(loaded.Entries) != len(original.Entries) {
		t.Fatalf("expected %d entries, got %d", len(original.Entries), len(loaded.Entries))
	}
}

func TestLoadSnapshot_MissingFile(t *testing.T) {
	_, err := LoadSnapshot("/nonexistent/snap.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestSnapshotDiff_AddedAndRemoved(t *testing.T) {
	baseline := TakeSnapshot([]Entry{
		{Port: 80, Protocol: "tcp"},
		{Port: 22, Protocol: "tcp"},
	})
	current := TakeSnapshot([]Entry{
		{Port: 80, Protocol: "tcp"},
		{Port: 8080, Protocol: "tcp"},
	})
	added, removed := SnapshotDiff(baseline, current)
	if len(added) != 1 || added[0].Port != 8080 {
		t.Fatalf("expected port 8080 added, got %v", added)
	}
	if len(removed) != 1 || removed[0].Port != 22 {
		t.Fatalf("expected port 22 removed, got %v", removed)
	}
}

func TestSnapshotDiff_NoChanges(t *testing.T) {
	entries := baseSnapshotEntries()
	baseline := TakeSnapshot(entries)
	current := TakeSnapshot(entries)
	added, removed := SnapshotDiff(baseline, current)
	if len(added) != 0 || len(removed) != 0 {
		t.Fatalf("expected no changes, got added=%v removed=%v", added, removed)
	}
}

func TestTakeSnapshot_IsolatesOriginal(t *testing.T) {
	entries := baseSnapshotEntries()
	snap := TakeSnapshot(entries)
	entries[0].Port = 9999
	if snap.Entries[0].Port == 9999 {
		t.Fatal("snapshot should not reflect mutation of original slice")
	}
	_ = os.Getenv("") // suppress unused import
}
