package history_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/history"
)

func makeTaggedEntries() []history.Entry {
	now := time.Now()
	return []history.Entry{
		{Timestamp: now, Port: 80, Protocol: "tcp", Action: "allow", Tags: []string{"web"}},
		{Timestamp: now, Port: 8080, Protocol: "tcp", Action: "allow", Tags: []string{"web", "dev"}},
		{Timestamp: now, Port: 22, Protocol: "tcp", Action: "alert"},
	}
}

func TestTagRoundtrip_TagThenFilter(t *testing.T) {
	entries := makeTaggedEntries()
	tagged := history.Tag(entries, []string{"critical"}, func(e history.Entry) bool {
		return e.Action == "alert"
	})
	result := history.FilterByTag(tagged, []string{"critical"})
	if len(result) != 1 || result[0].Port != 22 {
		t.Errorf("expected port 22 as critical, got %v", result)
	}
}

func TestTagRoundtrip_TagUntagFilter(t *testing.T) {
	entries := makeTaggedEntries()
	tagged := history.Tag(entries, []string{"temp"}, func(e history.Entry) bool {
		return e.Port == 8080
	})
	cleaned := history.Untag(tagged, []string{"temp"})
	result := history.FilterByTag(cleaned, []string{"temp"})
	if len(result) != 0 {
		t.Errorf("expected no entries with 'temp' tag after untag, got %d", len(result))
	}
}

func TestFilterByTag_MultipleTagsRequired(t *testing.T) {
	entries := makeTaggedEntries()
	result := history.FilterByTag(entries, []string{"web", "dev"})
	if len(result) != 1 || result[0].Port != 8080 {
		t.Errorf("expected only port 8080 with both tags, got %v", result)
	}
}
