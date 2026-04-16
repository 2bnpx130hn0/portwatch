package history

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func sampleExportEntries() []Entry {
	ts := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	return []Entry{
		{Timestamp: ts, Protocol: "tcp", Port: 8080, Action: "alert", Rule: "default"},
		{Timestamp: ts, Protocol: "udp", Port: 53, Action: "allow", Rule: "dns-rule"},
	}
}

func TestExport_CSVFormat(t *testing.T) {
	var buf bytes.Buffer
	ex := NewExporter(&buf, "csv")
	if err := ex.Export(sampleExportEntries()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "timestamp,protocol,port,action,rule") {
		t.Error("missing CSV header")
	}
	if !strings.Contains(out, "8080") {
		t.Error("missing port 8080 in CSV output")
	}
	if !strings.Contains(out, "dns-rule") {
		t.Error("missing rule name in CSV output")
	}
}

func TestExport_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	ex := NewExporter(&buf, "json")
	if err := ex.Export(sampleExportEntries()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var entries []Entry
	if err := json.Unmarshal(buf.Bytes(), &entries); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if len(entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Port != 8080 {
		t.Errorf("expected port 8080, got %d", entries[0].Port)
	}
}

func TestExport_DefaultsToCSV(t *testing.T) {
	var buf bytes.Buffer
	ex := NewExporter(&buf, "")
	if err := ex.Export(sampleExportEntries()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "timestamp,protocol,port,action,rule") {
		t.Error("expected CSV output when format is empty")
	}
}

func TestExport_EmptyEntries(t *testing.T) {
	var buf bytes.Buffer
	ex := NewExporter(&buf, "csv")
	if err := ex.Export([]Entry{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "timestamp") {
		t.Error("expected header even with no entries")
	}
}
