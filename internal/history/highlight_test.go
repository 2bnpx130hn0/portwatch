package history

import (
	"testing"
	"time"
)

func baseHighlightEntries() []Entry {
	now := time.Now()
	return []Entry{
		{Port: 22, Protocol: "tcp", Action: "alert", Timestamp: now, Labels: map[string]string{"flagged": "true"}},
		{Port: 80, Protocol: "tcp", Action: "allow", Timestamp: now, Labels: map[string]string{}},
		{Port: 443, Protocol: "tcp", Action: "warn", Timestamp: now.Add(-48 * time.Hour), Labels: map[string]string{"bookmarked": "true"}},
		{Port: 8080, Protocol: "tcp", Action: "alert", Timestamp: now, Labels: map[string]string{}},
	}
}

func TestHighlight_ByAction(t *testing.T) {
	entries := baseHighlightEntries()
	out := Highlight(entries, HighlightOptions{Actions: []string{"alert"}})
	if len(out) != 2 {
		t.Fatalf("expected 2, got %d", len(out))
	}
}

func TestHighlight_OnlyFlagged(t *testing.T) {
	entries := baseHighlightEntries()
	out := Highlight(entries, HighlightOptions{OnlyFlagged: true})
	if len(out) != 1 || out[0].Port != 22 {
		t.Fatalf("expected port 22, got %+v", out)
	}
}

func TestHighlight_OnlyBookmarked(t *testing.T) {
	entries := baseHighlightEntries()
	out := Highlight(entries, HighlightOptions{OnlyBookmark: true})
	if len(out) != 1 || out[0].Port != 443 {
		t.Fatalf("expected port 443, got %+v", out)
	}
}

func TestHighlight_Since(t *testing.T) {
	entries := baseHighlightEntries()
	out := Highlight(entries, HighlightOptions{Since: time.Now().Add(-1 * time.Hour)})
	for _, e := range out {
		if e.Port == 443 {
			t.Fatal("expected old entry to be filtered")
		}
	}
	if len(out) != 3 {
		t.Fatalf("expected 3, got %d", len(out))
	}
}

func TestHighlight_NoOptions(t *testing.T) {
	entries := baseHighlightEntries()
	out := Highlight(entries, HighlightOptions{})
	if len(out) != len(entries) {
		t.Fatalf("expected all entries, got %d", len(out))
	}
}
