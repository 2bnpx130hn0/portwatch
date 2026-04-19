package history

import (
	"testing"
	"time"
)

func baseCompareEntries() ([]Entry, []Entry) {
	now := time.Now()
	a := []Entry{
		{Port: 80, Protocol: "tcp", Action: "allow", Timestamp: now},
		{Port: 443, Protocol: "tcp", Action: "allow", Timestamp: now},
		{Port: 8080, Protocol: "tcp", Action: "alert", Timestamp: now},
	}
	b := []Entry{
		{Port: 80, Protocol: "tcp", Action: "allow", Timestamp: now},
		{Port: 443, Protocol: "tcp", Action: "allow", Timestamp: now},
		{Port: 9090, Protocol: "tcp", Action: "alert", Timestamp: now},
	}
	return a, b
}

func TestCompare_OnlyInA(t *testing.T) {
	a, b := baseCompareEntries()
	r := Compare(a, b, CompareOptions{})
	if len(r.OnlyInA) != 1 || r.OnlyInA[0].Port != 8080 {
		t.Errorf("expected port 8080 only in A, got %+v", r.OnlyInA)
	}
}

func TestCompare_OnlyInB(t *testing.T) {
	a, b := baseCompareEntries()
	r := Compare(a, b, CompareOptions{})
	if len(r.OnlyInB) != 1 || r.OnlyInB[0].Port != 9090 {
		t.Errorf("expected port 9090 only in B, got %+v", r.OnlyInB)
	}
}

func TestCompare_InBoth(t *testing.T) {
	a, b := baseCompareEntries()
	r := Compare(a, b, CompareOptions{})
	if len(r.InBoth) != 2 {
		t.Errorf("expected 2 in both, got %d", len(r.InBoth))
	}
}

func TestCompare_FilterByAction(t *testing.T) {
	a, b := baseCompareEntries()
	r := Compare(a, b, CompareOptions{Action: "alert"})
	if len(r.OnlyInA) != 1 || r.OnlyInA[0].Port != 8080 {
		t.Errorf("expected only port 8080 alert in A")
	}
	if len(r.InBoth) != 0 {
		t.Errorf("expected no shared alert entries")
	}
}

func TestCompare_FilterBySince(t *testing.T) {
	now := time.Now()
	a := []Entry{
		{Port: 80, Protocol: "tcp", Action: "allow", Timestamp: now.Add(-2 * time.Hour)},
		{Port: 443, Protocol: "tcp", Action: "allow", Timestamp: now},
	}
	b := []Entry{
		{Port: 80, Protocol: "tcp", Action: "allow", Timestamp: now},
	}
	r := Compare(a, b, CompareOptions{Since: now.Add(-time.Hour)})
	if len(r.InBoth) != 1 || r.InBoth[0].Port != 443 {
		t.Errorf("expected only port 443 in both after since filter")
	}
}
