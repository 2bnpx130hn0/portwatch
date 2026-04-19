package history

import (
	"testing"
)

func baseFlagEntries() []Entry {
	return []Entry{
		{Port: 80, Protocol: "tcp", Action: "allow"},
		{Port: 443, Protocol: "tcp", Action: "allow"},
		{Port: 8080, Protocol: "tcp", Action: "alert"},
	}
}

func TestFlag_AddsFlaggedLabel(t *testing.T) {
	entries := baseFlagEntries()
	result := Flag(entries, 80, "tcp")
	if result[0].Labels["flagged"] != "true" {
		t.Errorf("expected flagged=true, got %v", result[0].Labels["flagged"])
	}
	if result[0].Labels["flagged_at"] == "" {
		t.Error("expected flagged_at to be set")
	}
}

func TestFlag_CaseInsensitiveProtocol(t *testing.T) {
	entries := baseFlagEntries()
	result := Flag(entries, 443, "TCP")
	if result[1].Labels["flagged"] != "true" {
		t.Error("expected flag on case-insensitive match")
	}
}

func TestFlag_NoMatchUnchanged(t *testing.T) {
	entries := baseFlagEntries()
	result := Flag(entries, 9999, "tcp")
	for _, e := range result {
		if e.Labels != nil && e.Labels["flagged"] == "true" {
			t.Error("expected no entries to be flagged")
		}
	}
}

func TestUnflag_RemovesLabel(t *testing.T) {
	entries := Flag(baseFlagEntries(), 80, "tcp")
	result := Unflag(entries, 80, "tcp")
	if result[0].Labels["flagged"] == "true" {
		t.Error("expected flagged label to be removed")
	}
	if result[0].Labels["flagged_at"] != "" {
		t.Error("expected flagged_at to be removed")
	}
}

func TestFilterFlagged_ReturnsOnlyFlagged(t *testing.T) {
	entries := Flag(baseFlagEntries(), 8080, "tcp")
	result := FilterFlagged(entries)
	if len(result) != 1 {
		t.Fatalf("expected 1 flagged entry, got %d", len(result))
	}
	if result[0].Port != 8080 {
		t.Errorf("expected port 8080, got %d", result[0].Port)
	}
}

func TestFilterFlagged_EmptyWhenNoneFlagged(t *testing.T) {
	result := FilterFlagged(baseFlagEntries())
	if len(result) != 0 {
		t.Errorf("expected 0 results, got %d", len(result))
	}
}
