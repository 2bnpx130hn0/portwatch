package history

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestReplayCmd_Run_NoFile(t *testing.T) {
	cmd := NewReplayCmd("/tmp/portwatch_replay_nofile_test.json")
	cmd.Out = &bytes.Buffer{}
	// New() on missing file should succeed with empty history
	if err := cmd.Run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestReplayCmd_Run_WithEntries(t *testing.T) {
	p := tempPath(t)
	h, _ := New(p)
	now := time.Now()
	_ = h.Record(Entry{Port: 80, Protocol: "tcp", Action: "allow", Timestamp: now.Add(-1 * time.Hour)})
	_ = h.Record(Entry{Port: 443, Protocol: "tcp", Action: "alert", Timestamp: now})

	var buf bytes.Buffer
	cmd := NewReplayCmd(p)
	cmd.Out = &buf
	if err := cmd.Run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "replaying 2") {
		t.Errorf("unexpected output: %s", buf.String())
	}
}

func TestReplayCmd_WithLimit(t *testing.T) {
	p := tempPath(t)
	h, _ := New(p)
	for i := 0; i < 5; i++ {
		_ = h.Record(Entry{Port: 80 + i, Protocol: "tcp", Action: "allow", Timestamp: time.Now()})
	}
	var buf bytes.Buffer
	cmd := NewReplayCmd(p)
	cmd.WithLimit(3)
	cmd.Out = &buf
	_ = cmd.Run()
	if !strings.Contains(buf.String(), "replaying 3") {
		t.Errorf("unexpected output: %s", buf.String())
	}
}
