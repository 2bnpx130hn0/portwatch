package history

import (
	"testing"
)

func basePinEntries() []Entry {
	return []Entry{
		{Port: 80, Protocol: "tcp", Action: "allow"},
		{Port: 443, Protocol: "tcp", Action: "allow"},
		{Port: 22, Protocol: "tcp", Action: "warn"},
	}
}

func TestPin_AddsPinnedLabel(t *testing.T) {
	entries := basePinEntries()
	result := Pin(entries, 80, "tcp")
	if result[0].Labels["pinned"] != "true" {
		t.Error("expected port 80 to be pinned")
	}
	if result[0].Labels["pinned_at"] == "" {
		t.Error("expected pinned_at to be set")
	}
}

func TestPin_CaseInsensitiveProtocol(t *testing.T) {
	entries := basePinEntries()
	result := Pin(entries, 443, "TCP")
	if result[1].Labels["pinned"] != "true" {
		t.Error("expected port 443 to be pinned with case-insensitive match")
	}
}

func TestPin_NoMatchUnchanged(t *testing.T) {
	entries := basePinEntries()
	result := Pin(entries, 9999, "tcp")
	for _, e := range result {
		if e.Labels != nil && e.Labels["pinned"] == "true" {
			t.Error("expected no entries to be pinned")
		}
	}
}

func TestUnpin_RemovesPinnedLabel(t *testing.T) {
	entries := Pin(basePinEntries(), 80, "tcp")
	result := Unpin(entries, 80, "tcp")
	if result[0].Labels != nil && result[0].Labels["pinned"] == "true" {
		t.Error("expected port 80 to be unpinned")
	}
}

func TestFilterPinned_ReturnsOnlyPinned(t *testing.T) {
	entries := Pin(basePinEntries(), 22, "tcp")
	result := FilterPinned(entries)
	if len(result) != 1 {
		t.Fatalf("expected 1 pinned entry, got %d", len(result))
	}
	if result[0].Port != 22 {
		t.Errorf("expected port 22, got %d", result[0].Port)
	}
}

func TestFilterPinned_EmptyWhenNonePinned(t *testing.T) {
	entries := basePinEntries()
	result := FilterPinned(entries)
	if len(result) != 0 {
		t.Errorf("expected 0 pinned entries, got %d", len(result))
	}
}
