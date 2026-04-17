package history

import (
	"testing"
	"time"
)

func baseDiffEntries() ([]Entry, []Entry) {
	now := time.Now()
	before := []Entry{
		{Port: 80, Protocol: "tcp", Action: "allow", Timestamp: now},
		{Port: 443, Protocol: "tcp", Action: "allow", Timestamp: now},
	}
	after := []Entry{
		{Port: 443, Protocol: "tcp", Action: "allow", Timestamp: now},
		{Port: 8080, Protocol: "tcp", Action: "alert", Timestamp: now},
	}
	return before, after
}

func TestDiff_DetectsAdded(t *testing.T) {
	before, after := baseDiffEntries()
	diffs := Diff(before, after)
	var added []DiffEntry
	for _, d := range diffs {
		if d.Added {
			added = append(added, d)
		}
	}
	if len(added) != 1 || added[0].Port != 8080 {
		t.Errorf("expected added port 8080, got %+v", added)
	}
}

func TestDiff_DetectsRemoved(t *testing.T) {
	before, after := baseDiffEntries()
	diffs := Diff(before, after)
	var removed []DiffEntry
	for _, d := range diffs {
		if d.Removed {
			removed = append(removed, d)
		}
	}
	if len(removed) != 1 || removed[0].Port != 80 {
		t.Errorf("expected removed port 80, got %+v", removed)
	}
}

func TestDiff_NoChanges(t *testing.T) {
	now := time.Now()
	entries := []Entry{
		{Port: 80, Protocol: "tcp", Action: "allow", Timestamp: now},
	}
	diffs := Diff(entries, entries)
	if len(diffs) != 0 {
		t.Errorf("expected no diffs, got %d", len(diffs))
	}
}

func TestDiff_EmptyBefore(t *testing.T) {
	now := time.Now()
	after := []Entry{
		{Port: 22, Protocol: "tcp", Action: "alert", Timestamp: now},
	}
	diffs := Diff(nil, after)
	if len(diffs) != 1 || !diffs[0].Added {
		t.Errorf("expected one added entry, got %+v", diffs)
	}
}

func TestDiff_EmptyAfter(t *testing.T) {
	now := time.Now()
	before := []Entry{
		{Port: 22, Protocol: "tcp", Action: "alert", Timestamp: now},
	}
	diffs := Diff(before, nil)
	if len(diffs) != 1 || !diffs[0].Removed {
		t.Errorf("expected one removed entry, got %+v", diffs)
	}
}
