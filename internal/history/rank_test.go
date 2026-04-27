package history

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func baseRankEntries() []Entry {
	now := time.Now()
	return []Entry{
		{Port: 80, Protocol: "tcp", Action: "alert", Timestamp: now.Add(-1 * time.Hour)},
		{Port: 80, Protocol: "tcp", Action: "alert", Timestamp: now.Add(-2 * time.Hour)},
		{Port: 80, Protocol: "tcp", Action: "alert", Timestamp: now.Add(-3 * time.Hour)},
		{Port: 443, Protocol: "tcp", Action: "allow", Timestamp: now.Add(-30 * time.Minute)},
		{Port: 443, Protocol: "tcp", Action: "allow", Timestamp: now.Add(-45 * time.Minute)},
		{Port: 22, Protocol: "tcp", Action: "alert", Timestamp: now.Add(-72 * time.Hour)},
	}
}

func TestRank_OrderedByScore(t *testing.T) {
	entries := baseRankEntries()
	results := Rank(entries, RankOptions{})
	if len(results) < 2 {
		t.Fatalf("expected at least 2 results, got %d", len(results))
	}
	if results[0].Score < results[1].Score {
		t.Errorf("expected results ordered descending by score")
	}
}

func TestRank_TopN(t *testing.T) {
	entries := baseRankEntries()
	results := Rank(entries, RankOptions{TopN: 2})
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}
}

func TestRank_FilterByAction(t *testing.T) {
	entries := baseRankEntries()
	results := Rank(entries, RankOptions{Action: "allow"})
	for _, r := range results {
		if r.Port != 443 {
			t.Errorf("expected only port 443 for action=allow, got %d", r.Port)
		}
	}
}

func TestRank_SinceFilter(t *testing.T) {
	entries := baseRankEntries()
	cutoff := time.Now().Add(-4 * time.Hour)
	results := Rank(entries, RankOptions{Since: cutoff})
	for _, r := range results {
		if r.Port == 22 {
			t.Errorf("port 22 should be excluded by since filter")
		}
	}
}

func TestRenderRank_Text(t *testing.T) {
	results := []RankResult{
		{Port: 80, Protocol: "tcp", Count: 3, Score: 4.5},
	}
	var buf bytes.Buffer
	RenderRank(results, "text", &buf)
	out := buf.String()
	if !strings.Contains(out, "80") {
		t.Errorf("expected port 80 in text output")
	}
	if !strings.Contains(out, "tcp") {
		t.Errorf("expected protocol tcp in text output")
	}
}

func TestRenderRank_JSON(t *testing.T) {
	results := []RankResult{
		{Port: 443, Protocol: "tcp", Count: 2, Score: 3.2},
	}
	var buf bytes.Buffer
	RenderRank(results, "json", &buf)
	out := buf.String()
	if !strings.Contains(out, `"Port"`) && !strings.Contains(out, `"port"`) {
		t.Errorf("expected JSON output to contain port field")
	}
}

func TestRenderRank_DefaultsToText(t *testing.T) {
	var buf bytes.Buffer
	RenderRank([]RankResult{}, "", &buf)
	if buf.Len() == 0 {
		t.Errorf("expected non-empty output for default format")
	}
}
