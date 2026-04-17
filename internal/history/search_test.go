package history

import (
	"testing"
	"time"
)

func baseSearchEntries() []Entry {
	now := time.Now()
	return []Entry{
		{Port: 80, Protocol: "tcp", Action: "allow", Host: "localhost", Timestamp: now},
		{Port: 443, Protocol: "tcp", Action: "alert", Host: "localhost", Timestamp: now},
		{Port: 53, Protocol: "udp", Action: "allow", Host: "remotehost", Timestamp: now},
		{Port: 8080, Protocol: "tcp", Action: "warn", Host: "remotehost", Timestamp: now},
	}
}

func TestSearch_ByPort(t *testing.T) {
	res := Search(baseSearchEntries(), SearchOptions{Port: 80})
	if len(res) != 1 || res[0].Port != 80 {
		t.Fatalf("expected 1 entry with port 80, got %v", res)
	}
}

func TestSearch_ByProtocol(t *testing.T) {
	res := Search(baseSearchEntries(), SearchOptions{Protocol: "udp"})
	if len(res) != 1 || res[0].Protocol != "udp" {
		t.Fatalf("expected 1 udp entry, got %v", res)
	}
}

func TestSearch_ByAction(t *testing.T) {
	res := Search(baseSearchEntries(), SearchOptions{Action: "allow"})
	if len(res) != 2 {
		t.Fatalf("expected 2 allow entries, got %d", len(res))
	}
}

func TestSearch_MultipleFields(t *testing.T) {
	res := Search(baseSearchEntries(), SearchOptions{Protocol: "tcp", Action: "allow"})
	if len(res) != 1 || res[0].Port != 80 {
		t.Fatalf("expected 1 entry, got %v", res)
	}
}

func TestSearch_NoMatch(t *testing.T) {
	res := Search(baseSearchEntries(), SearchOptions{Port: 9999})
	if len(res) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(res))
	}
}

func TestSearchAny_MatchesMultiple(t *testing.T) {
	res := SearchAny(baseSearchEntries(), SearchOptions{Port: 80, Protocol: "udp"})
	if len(res) != 2 {
		t.Fatalf("expected 2 entries (OR match), got %d", len(res))
	}
}

func TestSearch_CaseInsensitiveProtocol(t *testing.T) {
	res := Search(baseSearchEntries(), SearchOptions{Protocol: "TCP"})
	if len(res) != 3 {
		t.Fatalf("expected 3 tcp entries, got %d", len(res))
	}
}
