package history

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func makeChainHistory(t *testing.T) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "history.json")
	h, err := New(p)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now()
	_ = h.Record(Entry{Port: 8080, Protocol: "tcp", Action: "alert", Timestamp: now.Add(-2 * time.Hour)})
	_ = h.Record(Entry{Port: 8080, Protocol: "tcp", Action: "allow", Timestamp: now.Add(-1 * time.Hour)})
	_ = h.Record(Entry{Port: 443, Protocol: "tcp", Action: "allow", Timestamp: now.Add(-30 * time.Minute)})
	return p
}

func TestChainCmd_Run_ReturnsChain(t *testing.T) {
	p := makeChainHistory(t)
	var buf bytes.Buffer
	cmd := &ChainCmd{
		Path:     p,
		Port:     8080,
		Protocol: "tcp",
		Format:   "text",
		Output:   &buf,
	}
	if err := cmd.Run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "8080/tcp") {
		t.Errorf("expected 8080/tcp in output")
	}
}

func TestChainCmd_Run_MissingPort(t *testing.T) {
	p := makeChainHistory(t)
	cmd := &ChainCmd{Path: p, Output: os.Discard}
	if err := cmd.Run(); err == nil {
		t.Fatal("expected error for missing port")
	}
}

func TestChainCmd_Run_NoFile(t *testing.T) {
	cmd := &ChainCmd{
		Path:   filepath.Join(t.TempDir(), "missing.json"),
		Port:   80,
		Output: os.Discard,
	}
	// New() creates the file if missing, so Run should succeed with empty chain
	if err := cmd.Run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
