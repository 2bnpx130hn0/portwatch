package history_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourorg/portwatch/internal/history"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "history.json")
}

func TestNew_NoFile(t *testing.T) {
	h, err := history.New(tempPath(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := len(h.Entries()); got != 0 {
		t.Errorf("expected 0 entries, got %d", got)
	}
}

func TestRecord_PersistsEntry(t *testing.T) {
	p := tempPath(t)
	h, _ := history.New(p)

	e := history.Entry{
		Proto:   "tcp",
		Port:    8080,
		Action:  "alert",
		Message: "new port opened",
	}
	if err := h.Record(e); err != nil {
		t.Fatalf("Record failed: %v", err)
	}

	if _, err := os.Stat(p); err != nil {
		t.Fatalf("history file not created: %v", err)
	}

	entries := h.Entries()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Port != 8080 {
		t.Errorf("expected port 8080, got %d", entries[0].Port)
	}
}

func TestRecord_TimestampAutoSet(t *testing.T) {
	h, _ := history.New(tempPath(t))
	before := time.Now().UTC()
	_ = h.Record(history.Entry{Proto: "udp", Port: 53, Action: "allow"})
	after := time.Now().UTC()

	e := h.Entries()[0]
	if e.Timestamp.Before(before) || e.Timestamp.After(after) {
		t.Errorf("timestamp %v not in expected range [%v, %v]", e.Timestamp, before, after)
	}
}

func TestNew_LoadsExisting(t *testing.T) {
	p := tempPath(t)
	h1, _ := history.New(p)
	_ = h1.Record(history.Entry{Proto: "tcp", Port: 443, Action: "warn", Message: "unexpected"})
	_ = h1.Record(history.Entry{Proto: "tcp", Port: 80, Action: "allow", Message: "ok"})

	h2, err := history.New(p)
	if err != nil {
		t.Fatalf("reload failed: %v", err)
	}
	if got := len(h2.Entries()); got != 2 {
		t.Errorf("expected 2 entries after reload, got %d", got)
	}
}
