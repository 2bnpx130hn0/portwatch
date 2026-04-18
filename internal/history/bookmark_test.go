package history

import (
	"testing"
	"time"
)

func baseBookmarkEntries() []Entry {
	now := time.Now()
	return []Entry{
		{Port: 80, Protocol: "tcp", Action: "allow", Timestamp: now},
		{Port: 443, Protocol: "tcp", Action: "allow", Timestamp: now},
		{Port: 22, Protocol: "tcp", Action: "alert", Timestamp: now},
	}
}

func TestBookmark_AddsTag(t *testing.T) {
	entries := baseBookmarkEntries()
	out := Bookmark(entries, 80, "tcp", "web")
	if !containsTag(out[0].Tags, "bookmark:web") {
		t.Error("expected bookmark:web tag on port 80")
	}
	if containsTag(out[1].Tags, "bookmark:web") {
		t.Error("unexpected bookmark:web on port 443")
	}
}

func TestBookmark_NoDuplicates(t *testing.T) {
	entries := baseBookmarkEntries()
	out := Bookmark(entries, 80, "tcp", "web")
	out = Bookmark(out, 80, "tcp", "web")
	count := 0
	for _, t2 := range out[0].Tags {
		if t2 == "bookmark:web" {
			count++
		}
	}
	if count != 1 {
		t.Errorf("expected 1 bookmark tag, got %d", count)
	}
}

func TestRemoveBookmark_RemovesTag(t *testing.T) {
	entries := baseBookmarkEntries()
	out := Bookmark(entries, 80, "tcp", "web")
	out = RemoveBookmark(out, "web")
	if containsTag(out[0].Tags, "bookmark:web") {
		t.Error("expected bookmark:web to be removed")
	}
}

func TestFilterByBookmark_ReturnsMatching(t *testing.T) {
	entries := baseBookmarkEntries()
	out := Bookmark(entries, 443, "tcp", "secure")
	filtered := FilterByBookmark(out, "secure")
	if len(filtered) != 1 || filtered[0].Port != 443 {
		t.Errorf("expected 1 entry with port 443, got %v", filtered)
	}
}

func TestListBookmarks_ReturnsNames(t *testing.T) {
	entries := baseBookmarkEntries()
	out := Bookmark(entries, 80, "tcp", "web")
	out = Bookmark(out, 22, "tcp", "ssh")
	names := ListBookmarks(out)
	if len(names) != 2 {
		t.Errorf("expected 2 bookmark names, got %d", len(names))
	}
}

func TestBookmark_CaseInsensitiveProtocol(t *testing.T) {
	entries := baseBookmarkEntries()
	out := Bookmark(entries, 80, "TCP", "web")
	if !containsTag(out[0].Tags, "bookmark:web") {
		t.Error("expected bookmark with uppercase protocol to match")
	}
}
