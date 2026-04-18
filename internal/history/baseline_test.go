package history

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func baselineEntries() []Entry {
	return []Entry{
		{Port: 80, Protocol: "tcp", Action: "allow", Timestamp: time.Now()},
		{Port: 443, Protocol: "tcp", Action: "allow", Timestamp: time.Now()},
	}
}

func TestSaveAndLoadBaseline(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")

	entries := baselineEntries()
	if err := SaveBaseline(entries, path); err != nil {
		t.Fatalf("SaveBaseline: %v", err)
	}

	b, err := LoadBaseline(path)
	if err != nil {
		t.Fatalf("LoadBaseline: %v", err)
	}
	if len(b.Entries) != len(entries) {
		t.Errorf("expected %d entries, got %d", len(entries), len(b.Entries))
	}
	if b.CreatedAt.IsZero() {
		t.Error("expected non-zero CreatedAt")
	}
}

func TestLoadBaseline_MissingFile(t *testing.T) {
	_, err := LoadBaseline("/nonexistent/baseline.json")
	if !os.IsNotExist(err) {
		t.Errorf("expected not-exist error, got %v", err)
	}
}

func TestBaselineDiff_AddedAndRemoved(t *testing.T) {
	baseline := baselineEntries()
	current := []Entry{
		{Port: 80, Protocol: "tcp", Action: "allow"},
		{Port: 8080, Protocol: "tcp", Action: "alert"},
	}

	added, removed := BaselineDiff(baseline, current)

	if len(added) != 1 || added[0].Port != 8080 {
		t.Errorf("expected added port 8080, got %+v", added)
	}
	if len(removed) != 1 || removed[0].Port != 443 {
		t.Errorf("expected removed port 443, got %+v", removed)
	}
}

func TestBaselineDiff_NoChanges(t *testing.T) {
	entries := baselineEntries()
	added, removed := BaselineDiff(entries, entries)
	if len(added) != 0 || len(removed) != 0 {
		t.Errorf("expected no diff, got added=%v removed=%v", added, removed)
	}
}

func TestBaselineDiff_EmptyBaseline(t *testing.T) {
	current := baselineEntries()
	added, removed := BaselineDiff(nil, current)
	if len(added) != len(current) {
		t.Errorf("expected all current as added")
	}
	if len(removed) != 0 {
		t.Errorf("expected no removed")
	}
}
