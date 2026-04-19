package history

import (
	"testing"
)

func makePinHistory(t *testing.T) *History {
	t.Helper()
	p := tempPath(t)
	h, err := New(p)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	for _, e := range basePinEntries() {
		if err := h.Record(e); err != nil {
			t.Fatalf("Record: %v", err)
		}
	}
	return h
}

func TestPinCmd_PinPort(t *testing.T) {
	h := makePinHistory(t)
	cmd := NewPinCmd(h)
	if err := cmd.PinPort(80, "tcp"); err != nil {
		t.Fatalf("PinPort: %v", err)
	}
	pinned := FilterPinned(h.entries)
	if len(pinned) != 1 || pinned[0].Port != 80 {
		t.Errorf("expected port 80 pinned, got %+v", pinned)
	}
}

func TestPinCmd_UnpinPort(t *testing.T) {
	h := makePinHistory(t)
	cmd := NewPinCmd(h)
	_ = cmd.PinPort(443, "tcp")
	if err := cmd.UnpinPort(443, "tcp"); err != nil {
		t.Fatalf("UnpinPort: %v", err)
	}
	pinned := FilterPinned(h.entries)
	if len(pinned) != 0 {
		t.Errorf("expected no pinned entries, got %d", len(pinned))
	}
}

func TestPinCmd_ListPinned_NoPanic(t *testing.T) {
	h := makePinHistory(t)
	cmd := NewPinCmd(h)
	_ = cmd.PinPort(22, "tcp")
	// should not panic
	cmd.ListPinned()
}
