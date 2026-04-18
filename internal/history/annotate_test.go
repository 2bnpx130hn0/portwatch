package history

import (
	"strings"
	"testing"
)

func baseAnnotateEntries() []Entry {
	return []Entry{
		{Port: 80, Protocol: "tcp", Action: "allow"},
		{Port: 443, Protocol: "tcp", Action: "alert"},
		{Port: 80, Protocol: "udp", Action: "warn"},
	}
}

func TestAnnotate_AddsNote(t *testing.T) {
	entries := baseAnnotateEntries()
	result, n := Annotate(entries, 80, "tcp", "known web server")
	if n != 1 {
		t.Fatalf("expected 1 updated, got %d", n)
	}
	if result[0].Meta["annotation"] != "known web server" {
		t.Errorf("unexpected annotation: %s", result[0].Meta["annotation"])
	}
}

func TestAnnotate_CaseInsensitiveProtocol(t *testing.T) {
	entries := baseAnnotateEntries()
	_, n := Annotate(entries, 443, "TCP", "secure")
	if n != 1 {
		t.Fatalf("expected 1 updated, got %d", n)
	}
}

func TestAnnotate_SetsTimestamp(t *testing.T) {
	entries := baseAnnotateEntries()
	result, _ := Annotate(entries, 80, "tcp", "note")
	if _, ok := result[0].Meta["annotation_at"]; !ok {
		t.Error("expected annotation_at to be set")
	}
}

func TestAnnotate_NoMatch(t *testing.T) {
	entries := baseAnnotateEntries()
	_, n := Annotate(entries, 9999, "tcp", "ghost")
	if n != 0 {
		t.Errorf("expected 0 updated, got %d", n)
	}
}

func TestClearAnnotation_RemovesNote(t *testing.T) {
	entries := baseAnnotateEntries()
	entries, _ = Annotate(entries, 80, "tcp", "temp note")
	entries, n := ClearAnnotation(entries, 80, "tcp")
	if n != 1 {
		t.Fatalf("expected 1 cleared, got %d", n)
	}
	if _, ok := entries[0].Meta["annotation"]; ok {
		t.Error("annotation should have been removed")
	}
}

func TestFilterAnnotated_ReturnsOnlyAnnotated(t *testing.T) {
	entries := baseAnnotateEntries()
	entries, _ = Annotate(entries, 443, "tcp", "important")
	result := FilterAnnotated(entries)
	if len(result) != 1 {
		t.Fatalf("expected 1 annotated entry, got %d", len(result))
	}
	if result[0].Port != 443 {
		t.Errorf("unexpected port: %d", result[0].Port)
	}
}

func TestFilterAnnotated_EmptyWhenNone(t *testing.T) {
	entries := baseAnnotateEntries()
	result := FilterAnnotated(entries)
	if len(result) != 0 {
		t.Errorf("expected 0, got %d", len(result))
	}
}

// ensure strings import used in annotate.go compiles
var _ = strings.EqualFold
