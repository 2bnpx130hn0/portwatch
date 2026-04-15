package history_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/portwatch/internal/history"
)

func sampleEntries() []history.Entry {
	return []history.Entry{
		{
			Timestamp: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			Proto:     "tcp",
			Port:      8080,
			Action:    "alert",
			Message:   "unexpected port",
		},
		{
			Timestamp: time.Date(2024, 1, 15, 11, 0, 0, 0, time.UTC),
			Proto:     "udp",
			Port:      53,
			Action:    "allow",
			Message:   "dns ok",
		},
	}
}

func TestPrint_TextFormat(t *testing.T) {
	var buf bytes.Buffer
	r := history.NewReporter(&buf, history.FormatText)
	if err := r.Print(sampleEntries()); err != nil {
		t.Fatalf("Print failed: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "TIMESTAMP") {
		t.Error("expected header row")
	}
	if !strings.Contains(out, "8080") {
		t.Error("expected port 8080 in output")
	}
	if !strings.Contains(out, "alert") {
		t.Error("expected action 'alert' in output")
	}
}

func TestPrint_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	r := history.NewReporter(&buf, history.FormatJSON)
	if err := r.Print(sampleEntries()); err != nil {
		t.Fatalf("Print failed: %v", err)
	}
	var entries []history.Entry
	if err := json.Unmarshal(buf.Bytes(), &entries); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if len(entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(entries))
	}
}

func TestPrint_DefaultsToText(t *testing.T) {
	var buf bytes.Buffer
	r := history.NewReporter(&buf, "")
	if err := r.Print(sampleEntries()); err != nil {
		t.Fatalf("Print failed: %v", err)
	}
	if !strings.Contains(buf.String(), "TIMESTAMP") {
		t.Error("expected text format as default")
	}
}

func TestPrint_EmptyEntries(t *testing.T) {
	var buf bytes.Buffer
	r := history.NewReporter(&buf, history.FormatText)
	if err := r.Print([]history.Entry{}); err != nil {
		t.Fatalf("Print failed on empty entries: %v", err)
	}
}
