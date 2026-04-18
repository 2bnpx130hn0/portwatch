package history

import (
	"path/filepath"
	"testing"
	"time"
)

func makeLabelHistory(t *testing.T) *History {
	t.Helper()
	p := filepath.Join(t.TempDir(), "history.json")
	h, err := New(p)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	for _, e := range []Entry{
		{Port: 80, Protocol: "tcp", Action: "allow", Timestamp: time.Now()},
		{Port: 443, Protocol: "tcp", Action: "allow", Timestamp: time.Now()},
	} {
		if err := h.Record(e); err != nil {
			t.Fatalf("Record: %v", err)
		}
	}
	return h
}

func TestLabelCmd_Set(t *testing.T) {
	h := makeLabelHistory(t)
	cmd := NewLabelCmd(h)
	if err := cmd.Set(80, "tcp", "env", "prod"); err != nil {
		t.Fatalf("Set: %v", err)
	}
	entries := h.All()
	if entries[0].Labels["env"] != "prod" {
		t.Fatalf("expected label env=prod")
	}
}

func TestLabelCmd_Remove(t *testing.T) {
	h := makeLabelHistory(t)
	cmd := NewLabelCmd(h)
	_ = cmd.Set(80, "tcp", "env", "prod")
	if err := cmd.Remove(80, "tcp", "env"); err != nil {
		t.Fatalf("Remove: %v", err)
	}
	entries := h.All()
	if v := entries[0].Labels["env"]; v != "" {
		t.Fatalf("expected label removed, got %q", v)
	}
}

func TestLabelCmd_List(t *testing.T) {
	h := makeLabelHistory(t)
	cmd := NewLabelCmd(h)
	_ = cmd.Set(80, "tcp", "env", "prod")
	_ = cmd.Set(443, "tcp", "tier", "web")
	list := cmd.List()
	if len(list) != 2 {
		t.Fatalf("expected 2 labels, got %d: %v", len(list), list)
	}
}
