package history

import (
	"testing"
	"time"
)

func seedHistory(t *testing.T) *History {
	t.Helper()
	h := &History{}
	now := time.Now()
	h.entries = []Entry{
		{Port: 80, Protocol: "tcp", Action: "allow", Timestamp: now.Add(-3 * time.Hour)},
		{Port: 443, Protocol: "tcp", Action: "allow", Timestamp: now.Add(-2 * time.Hour)},
		{Port: 8080, Protocol: "tcp", Action: "alert", Timestamp: now.Add(-1 * time.Hour)},
		{Port: 53, Protocol: "udp", Action: "warn", Timestamp: now},
	}
	return h
}

func TestQuery_ByProtocol(t *testing.T) {
	h := seedHistory(t)
	results := h.Query(Filter{Protocol: "udp"})
	if len(results) != 1 || results[0].Port != 53 {
		t.Fatalf("expected 1 udp entry (port 53), got %v", results)
	}
}

func TestQuery_ByPort(t *testing.T) {
	h := seedHistory(t)
	results := h.Query(Filter{Port: 443})
	if len(results) != 1 || results[0].Port != 443 {
		t.Fatalf("expected 1 entry for port 443, got %v", results)
	}
}

func TestQuery_ByAction(t *testing.T) {
	h := seedHistory(t)
	results := h.Query(Filter{Action: "alert"})
	if len(results) != 1 || results[0].Port != 8080 {
		t.Fatalf("expected 1 alert entry, got %v", results)
	}
}

func TestQuery_Since(t *testing.T) {
	h := seedHistory(t)
	cutoff := time.Now().Add(-90 * time.Minute)
	results := h.Query(Filter{Since: cutoff})
	if len(results) != 2 {
		t.Fatalf("expected 2 entries after cutoff, got %d", len(results))
	}
}

func TestQuery_Limit(t *testing.T) {
	h := seedHistory(t)
	results := h.Query(Filter{Limit: 2})
	if len(results) != 2 {
		t.Fatalf("expected 2 entries with limit, got %d", len(results))
	}
}

func TestQuery_NoFilter(t *testing.T) {
	h := seedHistory(t)
	results := h.Query(Filter{})
	if len(results) != 4 {
		t.Fatalf("expected all 4 entries, got %d", len(results))
	}
}

func TestLatest_LessThanTotal(t *testing.T) {
	h := seedHistory(t)
	results := h.Latest(2)
	if len(results) != 2 {
		t.Fatalf("expected 2 latest entries, got %d", len(results))
	}
	if results[0].Port != 8080 || results[1].Port != 53 {
		t.Fatalf("unexpected latest entries: %v", results)
	}
}

func TestLatest_ZeroReturnsAll(t *testing.T) {
	h := seedHistory(t)
	results := h.Latest(0)
	if len(results) != 4 {
		t.Fatalf("expected all entries, got %d", len(results))
	}
}
