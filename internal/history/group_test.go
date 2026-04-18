package history

import (
	"testing"
	"time"
)

func baseGroupEntries() []Entry {
	now := time.Now()
	return []Entry{
		{Port: 80, Protocol: "tcp", Action: "allow", Timestamp: now},
		{Port: 80, Protocol: "tcp", Action: "alert", Timestamp: now},
		{Port: 443, Protocol: "tcp", Action: "allow", Timestamp: now},
		{Port: 53, Protocol: "udp", Action: "warn", Timestamp: now},
		{Port: 53, Protocol: "udp", Action: "allow", Timestamp: now},
	}
}

func TestGroupBy_Port(t *testing.T) {
	groups := GroupBy(baseGroupEntries(), GroupByPort)
	if len(groups) != 3 {
		t.Fatalf("expected 3 groups, got %d", len(groups))
	}
	// sorted: 443, 53, 80
	if groups[0].Key != "443" || groups[1].Key != "53" || groups[2].Key != "80" {
		t.Errorf("unexpected group keys: %v", groups)
	}
	if len(groups[2].Entries) != 2 {
		t.Errorf("expected 2 entries for port 80, got %d", len(groups[2].Entries))
	}
}

func TestGroupBy_Protocol(t *testing.T) {
	groups := GroupBy(baseGroupEntries(), GroupByProtocol)
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
	if groups[0].Key != "tcp" || groups[1].Key != "udp" {
		t.Errorf("unexpected keys: %v", groups)
	}
}

func TestGroupBy_Action(t *testing.T) {
	groups := GroupBy(baseGroupEntries(), GroupByAction)
	if len(groups) != 3 {
		t.Fatalf("expected 3 groups, got %d", len(groups))
	}
}

func TestGroupBy_Empty(t *testing.T) {
	groups := GroupBy([]Entry{}, GroupByPort)
	if len(groups) != 0 {
		t.Errorf("expected empty groups")
	}
}

func TestGroupCounts(t *testing.T) {
	counts := GroupCounts(baseGroupEntries(), GroupByProtocol)
	if counts["tcp"] != 3 {
		t.Errorf("expected tcp=3, got %d", counts["tcp"])
	}
	if counts["udp"] != 2 {
		t.Errorf("expected udp=2, got %d", counts["udp"])
	}
}

func TestGroupBy_SingleEntry(t *testing.T) {
	now := time.Now()
	entries := []Entry{
		{Port: 8080, Protocol: "tcp", Action: "allow", Timestamp: now},
	}
	groups := GroupBy(entries, GroupByPort)
	if len(groups) != 1 {
		t.Fatalf("expected 1 group, got %d", len(groups))
	}
	if groups[0].Key != "8080" {
		t.Errorf("expected key 8080, got %s", groups[0].Key)
	}
	if len(groups[0].Entries) != 1 {
		t.Errorf("expected 1 entry, got %d", len(groups[0].Entries))
	}
}
