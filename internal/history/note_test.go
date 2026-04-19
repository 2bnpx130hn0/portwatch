package history

import (
	"testing"
)

func baseNoteEntries() []Entry {
	return []Entry{
		{Port: 80, Protocol: "tcp", Action: "allow"},
		{Port: 443, Protocol: "tcp", Action: "allow"},
		{Port: 22, Protocol: "tcp", Action: "alert"},
	}
}

func TestNoteEntry_AddsNote(t *testing.T) {
	entries := baseNoteEntries()
	result := NoteEntry(entries, 80, "tcp", "web traffic")
	if result[0].Labels["note"] != "web traffic" {
		t.Errorf("expected note 'web traffic', got %q", result[0].Labels["note"])
	}
	if result[0].Labels["note_at"] == "" {
		t.Error("expected note_at to be set")
	}
}

func TestNoteEntry_CaseInsensitiveProtocol(t *testing.T) {
	entries := baseNoteEntries()
	result := NoteEntry(entries, 443, "TCP", "https")
	if result[1].Labels["note"] != "https" {
		t.Errorf("expected note 'https', got %q", result[1].Labels["note"])
	}
}

func TestNoteEntry_NoMatchUnchanged(t *testing.T) {
	entries := baseNoteEntries()
	result := NoteEntry(entries, 9999, "tcp", "nope")
	for _, e := range result {
		if e.Labels != nil && e.Labels["note"] == "nope" {
			t.Error("unexpected note set on non-matching entry")
		}
	}
}

func TestRemoveNote_RemovesLabel(t *testing.T) {
	entries := baseNoteEntries()
	withNote := NoteEntry(entries, 22, "tcp", "suspicious")
	result := RemoveNote(withNote, 22, "tcp")
	if result[2].Labels != nil {
		if _, ok := result[2].Labels["note"]; ok {
			t.Error("expected note to be removed")
		}
	}
}

func TestFilterNoted_ReturnsOnlyNoted(t *testing.T) {
	entries := baseNoteEntries()
	withNote := NoteEntry(entries, 80, "tcp", "web")
	result := FilterNoted(withNote)
	if len(result) != 1 {
		t.Fatalf("expected 1 noted entry, got %d", len(result))
	}
	if result[0].Port != 80 {
		t.Errorf("expected port 80, got %d", result[0].Port)
	}
}

func TestFilterNoted_EmptyReturnsNone(t *testing.T) {
	result := FilterNoted(baseNoteEntries())
	if len(result) != 0 {
		t.Errorf("expected 0, got %d", len(result))
	}
}
