package state_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/portwatch/internal/state"
)

func tempFile(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "state.json")
}

func TestLoad_NoFile(t *testing.T) {
	store := state.New(tempFile(t))
	snap, err := store.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if snap != nil {
		t.Fatal("expected nil snapshot when no file exists")
	}
}

func TestSaveAndLoad(t *testing.T) {
	path := tempFile(t)
	store := state.New(path)

	snap := &state.Snapshot{
		Timestamp: time.Now().UTC().Truncate(time.Second),
		Ports:     map[string][]int{"tcp": {80, 443}},
	}
	if err := store.Save(snap); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := store.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if len(loaded.Ports["tcp"]) != 2 {
		t.Errorf("expected 2 tcp ports, got %d", len(loaded.Ports["tcp"]))
	}
}

func TestCurrent_ReturnsInMemory(t *testing.T) {
	store := state.New(tempFile(t))
	if store.Current() != nil {
		t.Fatal("expected nil before any save")
	}
	snap := &state.Snapshot{Ports: map[string][]int{"udp": {53}}}
	_ = store.Save(snap)
	if store.Current() == nil {
		t.Fatal("expected non-nil after save")
	}
}

func TestDiff_AddedAndRemoved(t *testing.T) {
	prev := &state.Snapshot{Ports: map[string][]int{"tcp": {80, 22}}}
	next := &state.Snapshot{Ports: map[string][]int{"tcp": {80, 8080}}}

	added, removed := state.Diff(prev, next)

	if len(added["tcp"]) != 1 || added["tcp"][0] != 8080 {
		t.Errorf("expected added tcp:[8080], got %v", added)
	}
	if len(removed["tcp"]) != 1 || removed["tcp"][0] != 22 {
		t.Errorf("expected removed tcp:[22], got %v", removed)
	}
}

func TestDiff_NoDifference(t *testing.T) {
	snap := &state.Snapshot{Ports: map[string][]int{"tcp": {443}}}
	added, removed := state.Diff(snap, snap)
	if len(added) != 0 || len(removed) != 0 {
		t.Errorf("expected no diff, got added=%v removed=%v", added, removed)
	}
}

func TestSave_CreatesFile(t *testing.T) {
	path := tempFile(t)
	_ = os.Remove(path) // ensure it doesn't exist
	store := state.New(path)
	if err := store.Save(&state.Snapshot{Ports: map[string][]int{}}); err != nil {
		t.Fatalf("Save failed: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("file not created: %v", err)
	}
}
