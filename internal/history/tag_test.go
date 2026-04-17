package history

import (
	"testing"
)

func baseTagEntries() []Entry {
	return []Entry{
		{Port: 80, Protocol: "tcp", Action: "allow", Tags: []string{"web"}},
		{Port: 443, Protocol: "tcp", Action: "allow", Tags: nil},
		{Port: 22, Protocol: "tcp", Action: "alert", Tags: []string{"ssh"}},
	}
}

func TestTag_AddsTagToMatching(t *testing.T) {
	entries := baseTagEntries()
	result := Tag(entries, []string{"secure"}, func(e Entry) bool {
		return e.Port == 443
	})
	if !containsTag(result[1].Tags, "secure") {
		t.Errorf("expected 'secure' tag on port 443")
	}
	if containsTag(result[0].Tags, "secure") {
		t.Errorf("port 80 should not have 'secure' tag")
	}
}

func TestTag_PreservesExistingTags(t *testing.T) {
	entries := baseTagEntries()
	result := Tag(entries, []string{"monitored"}, func(e Entry) bool {
		return e.Port == 80
	})
	if !containsTag(result[0].Tags, "web") {
		t.Errorf("expected existing 'web' tag to be preserved")
	}
	if !containsTag(result[0].Tags, "monitored") {
		t.Errorf("expected new 'monitored' tag")
	}
}

func TestTag_NoDuplicates(t *testing.T) {
	entries := baseTagEntries()
	result := Tag(entries, []string{"web"}, func(e Entry) bool {
		return e.Port == 80
	})
	count := 0
	for _, t2 := range result[0].Tags {
		if t2 == "web" {
			count++
		}
	}
	if count != 1 {
		t.Errorf("expected exactly 1 'web' tag, got %d", count)
	}
}

func TestUntag_RemovesTags(t *testing.T) {
	entries := baseTagEntries()
	result := Untag(entries, []string{"web", "ssh"})
	for _, e := range result {
		if containsTag(e.Tags, "web") || containsTag(e.Tags, "ssh") {
			t.Errorf("expected tags to be removed from port %d", e.Port)
		}
	}
}

func TestFilterByTag_ReturnsMatching(t *testing.T) {
	entries := baseTagEntries()
	result := FilterByTag(entries, []string{"web"})
	if len(result) != 1 || result[0].Port != 80 {
		t.Errorf("expected only port 80, got %v", result)
	}
}

func TestFilterByTag_NoMatch(t *testing.T) {
	entries := baseTagEntries()
	result := FilterByTag(entries, []string{"nonexistent"})
	if len(result) != 0 {
		t.Errorf("expected no results, got %d", len(result))
	}
}

func containsTag(tags []string, tag string) bool {
	for _, t := range tags {
		if t == tag {
			return true
		}
	}
	return false
}
