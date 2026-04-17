package history

import (
	"testing"
	"time"
	"strings"
)

func baseEntries() []Entry {
	now := time.Now()
	return []Entry{
		{Port: 80, Protocol: "tcp", Action: "allow", Timestamp: now.Add(-2 * time.Hour)},
		{Port: 443, Protocol: "tcp", Action: "alert", Timestamp: now.Add(-1 * time.Hour)},
		{Port: 53, Protocol: "udp", Action: "allow", Timestamp: now},
	}
}

func TestFilter_ByProtocol(t *testing.T) {
	res := Filter{Protocol: "udp"}.Apply(baseEntries())
	if len(res) != 1 || res[0].Port != 53 {
		t.Fatalf("expected 1 udp entry, got %+v", res)
	}
}

func TestFilter_ByPort(t *testing.T) {
	res := Filter{Port: 443}.Apply(baseEntries())
	if len(res) != 1 || res[0].Port != 443 {
		t.Fatalf("expected port 443, got %+v", res)
	}
}

func TestFilter_ByAction(t *testing.T) {
	res := Filter{Action: "allow"}.Apply(baseEntries())
	if len(res) != 2 {
		t.Fatalf("expected 2 allow entries, got %d", len(res))
	}
}

func TestFilter_Since(t *testing.T) {
	res := Filter{Since: time.Now().Add(-90 * time.Minute)}.Apply(baseEntries())
	if len(res) != 2 {
		t.Fatalf("expected 2 recent entries, got %d", len(res))
	}
}

func TestFilter_Limit(t *testing.T) {
	res := Filter{Limit: 2}.Apply(baseEntries())
	if len(res) != 2 {
		t.Fatalf("expected limit 2, got %d", len(res))
	}
}

func TestFilter_CaseInsensitiveAction(t *testing.T) {
	res := Filter{Action: "ALERT"}.Apply(baseEntries())
	if len(res) != 1 || !strings.EqualFold(res[0].Action, "alert") {
		t.Fatalf("expected case-insensitive match, got %+v", res)
	}
}

func TestFilter_NoMatch(t *testing.T) {
	res := Filter{Port: 9999}.Apply(baseEntries())
	if len(res) != 0 {
		t.Fatalf("expected no results, got %+v", res)
	}
}

func TestFilter_LimitExceedsEntries(t *testing.T) {
	res := Filter{Limit: 100}.Apply(baseEntries())
	if len(res) != 3 {
		t.Fatalf("expected all 3 entries when limit exceeds total, got %d", len(res))
	}
}
