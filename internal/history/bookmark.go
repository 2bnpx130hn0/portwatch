package history

import (
	"strings"
	"time"
)

// Bookmark attaches a named bookmark to entries matching port/protocol.
func Bookmark(entries []Entry, port int, protocol, name string) []Entry {
	result := make([]Entry, len(entries))
	copy(result, entries)
	for i, e := range result {
		if e.Port == port && strings.EqualFold(e.Protocol, protocol) {
			if result[i].Tags == nil {
				result[i].Tags = []string{}
			}
			bookmarkTag := "bookmark:" + name
			if !containsTag(result[i].Tags, bookmarkTag) {
				result[i].Tags = append(result[i].Tags, bookmarkTag)
			}
		}
	}
	return result
}

// RemoveBookmark removes a named bookmark from all matching entries.
func RemoveBookmark(entries []Entry, name string) []Entry {
	bookmarkTag := "bookmark:" + name
	result := make([]Entry, len(entries))
	copy(result, entries)
	for i, e := range result {
		filtered := e.Tags[:0:0]
		for _, t := range e.Tags {
			if !strings.EqualFold(t, bookmarkTag) {
				filtered = append(filtered, t)
			}
		}
		result[i].Tags = filtered
	}
	return result
}

// FilterByBookmark returns entries that carry the named bookmark.
func FilterByBookmark(entries []Entry, name string) []Entry {
	bookmarkTag := "bookmark:" + name
	var out []Entry
	for _, e := range entries {
		if containsTag(e.Tags, bookmarkTag) {
			out = append(out, e)
		}
	}
	return out
}

// ListBookmarks returns distinct bookmark names found across entries.
func ListBookmarks(entries []Entry) []string {
	seen := map[string]bool{}
	var names []string
	for _, e := range entries {
		for _, t := range e.Tags {
			if strings.HasPrefix(t, "bookmark:") {
				n := strings.TrimPrefix(t, "bookmark:")
				if !seen[n] {
					seen[n] = true
					names = append(names, n)
				}
			}
		}
	}
	return names
}

func containsTag(tags []string, tag string) bool {
	for _, t := range tags {
		if strings.EqualFold(t, tag) {
			return true
		}
	}
	return false
}

var _ = time.Now // suppress unused import
