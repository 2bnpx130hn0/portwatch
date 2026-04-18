package history

import (
	"testing"
	"time"
)

var t0 = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func baseMergeEntries() ([]Entry, []Entry) {
	a := []Entry{
		{Port: 80, Protocol: "tcp", Action: "allow", Timestamp: t0},
		{Port: 443, Protocol: "tcp", Action: "allow", Timestamp: t0.Add(2 * time.Second)},
	}
	b := []Entry{
		{Port: 22, Protocol: "tcp", Action: "alert", Timestamp: t0.Add(1 * time.Second)},
		{Port: 8080, Protocol: "tcp", Action: "warn", Timestamp: t0.Add(3 * time.Second)},
	}
	return a, b
}

func TestMerge_CombinesAndSorts(t *testing.T) {
	a, b := baseMergeEntries()
	result := Merge(a, b, MergeOptions{})
	if len(result) != 4 {
		t.Fatalf("expected 4 entries, got %d", len(result))
	}
	if result[0].Port != 80 || result[1].Port != 22 || result[2].Port != 443 || result[3].Port != 8080 {
		t.Errorf("unexpected order: %+v", result)
	}
}

func TestMerge_NoDuplicateWindowKeepsAll(t *testing.T) {
	e := Entry{Port: 80, Protocol: "tcp", Action: "allow", Timestamp: t0}
	result := Merge([]Entry{e}, []Entry{e}, MergeOptions{})
	if len(result) != 2 {
		t.Fatalf("expected 2, got %d", len(result))
	}
}

func TestMerge_DeduplicatesWithinWindow(t *testing.T) {
	e1 := Entry{Port: 80, Protocol: "tcp", Action: "allow", Timestamp: t0}
	e2 := Entry{Port: 80, Protocol: "tcp", Action: "allow", Timestamp: t0.Add(3 * time.Second)}
	e3 := Entry{Port: 80, Protocol: "tcp", Action: "allow", Timestamp: t0.Add(10 * time.Second)}

	result := Merge([]Entry{e1, e3}, []Entry{e2}, MergeOptions{DeduplicateWindow: 5 * time.Second})
	if len(result) != 2 {
		t.Fatalf("expected 2 after dedup, got %d", len(result))
	}
}

func TestMerge_DifferentActionNotDeduped(t *testing.T) {
	e1 := Entry{Port: 80, Protocol: "tcp", Action: "allow", Timestamp: t0}
	e2 := Entry{Port: 80, Protocol: "tcp", Action: "alert", Timestamp: t0.Add(1 * time.Second)}

	result := Merge([]Entry{e1}, []Entry{e2}, MergeOptions{DeduplicateWindow: 10 * time.Second})
	if len(result) != 2 {
		t.Fatalf("expected 2, got %d", len(result))
	}
}

func TestMerge_EmptyInputs(t *testing.T) {
	result := Merge(nil, nil, MergeOptions{})
	if len(result) != 0 {
		t.Fatalf("expected 0, got %d", len(result))
	}
}
