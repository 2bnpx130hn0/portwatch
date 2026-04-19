package history

import (
	"testing"
	"time"
)

func baseScoreEntries() []Entry {
	now := time.Now()
	old := now.Add(-48 * time.Hour)
	return []Entry{
		{Port: 22, Protocol: "tcp", Action: "alert", Timestamp: now},
		{Port: 80, Protocol: "tcp", Action: "allow", Timestamp: now},
		{Port: 443, Protocol: "tcp", Action: "warn", Timestamp: old},
		{Port: 8080, Protocol: "tcp", Action: "alert", Timestamp: old},
	}
}

func TestScore_AlertHigherThanAllow(t *testing.T) {
	entries := baseScoreEntries()
	scored := Score(entries, ScoreOptions{})
	var alertScore, allowScore float64
	for _, s := range scored {
		if s.Entry.Port == 22 {
			alertScore = s.Score
		}
		if s.Entry.Port == 80 {
			allowScore = s.Score
		}
	}
	if alertScore <= allowScore {
		t.Errorf("expected alert score %.2f > allow score %.2f", alertScore, allowScore)
	}
}

func TestScore_RecencyBoostApplied(t *testing.T) {
	now := time.Now()
	old := now.Add(-48 * time.Hour)
	entries := []Entry{
		{Port: 1, Action: "alert", Timestamp: now},
		{Port: 2, Action: "alert", Timestamp: old},
	}
	scored := Score(entries, ScoreOptions{})
	var recent, stale float64
	for _, s := range scored {
		if s.Entry.Port == 1 {
			recent = s.Score
		} else {
			stale = s.Score
		}
	}
	if recent <= stale {
		t.Errorf("expected recent score %.2f > stale score %.2f", recent, stale)
	}
}

func TestTopScored_LimitsResults(t *testing.T) {
	scored := Score(baseScoreEntries(), ScoreOptions{})
	top := TopScored(scored, 2)
	if len(top) != 2 {
		t.Fatalf("expected 2 results, got %d", len(top))
	}
}

func TestTopScored_OrderedDescending(t *testing.T) {
	scored := Score(baseScoreEntries(), ScoreOptions{})
	top := TopScored(scored, 0)
	for i := 1; i < len(top); i++ {
		if top[i].Score > top[i-1].Score {
			t.Errorf("scores not descending at index %d", i)
		}
	}
}

func TestScore_DefaultWeights(t *testing.T) {
	entries := []Entry{
		{Port: 1, Action: "alert", Timestamp: time.Time{}},
		{Port: 2, Action: "warn", Timestamp: time.Time{}},
		{Port: 3, Action: "allow", Timestamp: time.Time{}},
	}
	scored := Score(entries, ScoreOptions{})
	if len(scored) != 3 {
		t.Fatalf("expected 3 scored entries")
	}
}
