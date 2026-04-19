package history

import (
	"bytes"
	"testing"
)

func makeFlagHistory(t *testing.T) *History {
	t.Helper()
	h, err := New(tempPath(t))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	for _, e := range baseFlagEntries() {
		if err := h.Record(e); err != nil {
			t.Fatalf("Record: %v", err)
		}
	}
	return h
}

func TestFlagCmd_Set(t *testing.T) {
	h := makeFlagHistory(t)
	var buf bytes.Buffer
	cmd := NewFlagCmd(h.path, &buf)
	if err := cmd.Set(80, "tcp"); err != nil {
		t.Fatalf("Set: %v", err)
	}
	if buf.String() == "" {
		t.Error("expected output from Set")
	}
	h2, _ := New(h.path)
	flagged := FilterFlagged(h2.entries)
	if len(flagged) != 1 || flagged[0].Port != 80 {
		t.Errorf("expected 1 flagged entry for port 80, got %v", flagged)
	}
}

func TestFlagCmd_Remove(t *testing.T) {
	h := makeFlagHistory(t)
	cmd := NewFlagCmd(h.path, &bytes.Buffer{})
	_ = cmd.Set(443, "tcp")
	var buf bytes.Buffer
	cmd2 := NewFlagCmd(h.path, &buf)
	if err := cmd2.Remove(443, "tcp"); err != nil {
		t.Fatalf("Remove: %v", err)
	}
	h2, _ := New(h.path)
	for _, e := range h2.entries {
		if e.Port == 443 && e.Labels["flagged"] == "true" {
			t.Error("expected port 443 to be unflagged")
		}
	}
}

func TestFlagCmd_List_NoPanic(t *testing.T) {
	h := makeFlagHistory(t)
	var buf bytes.Buffer
	cmd := NewFlagCmd(h.path, &buf)
	if err := cmd.List(); err != nil {
		t.Fatalf("List: %v", err)
	}
}
