package history

import (
	"testing"
	"time"
)

func baseOutlierEntries() []Entry {
	now := time.Now()
	var entries []Entry
	// Port 80/tcp appears 10 times — high count outlier.
	for i := 0; i < 10; i++ {
		entries = append(entries, Entry{Port: 80, Protocol: "tcp", Action: "alert", Timestamp: now})
	}
	// Port 443/tcp appears 2 times — normal.
	for i := 0; i < 2; i++ {
		entries = append(entries, Entry{Port: 443, Protocol: "tcp", Action: "allow", Timestamp: now})
	}
	// Port 22/tcp appears 1 time — normal.
	entries = append(entries, Entry{Port: 22, Protocol: "tcp", Action: "allow", Timestamp: now})
	return entries
}

func TestDetectOutliers_FindsHighCount(t *testing.T) {
	entries := baseOutlierEntries()
	results := DetectOutliers(entries, OutlierOptions{})
	if len(results) == 0 {
		t.Fatal("expected at least one outlier, got none")
	}
	if results[0].Port != 80 {
		t.Errorf("expected port 80 as top outlier, got %d", results[0].Port)
	}
}

func TestDetectOutliers_BelowThreshold(t *testing.T) {
	now := time.Now()
	// All ports appear equally — no outliers.
	var entries []Entry
	for _, p := range []int{80, 443, 22} {
		for i := 0; i < 3; i++ {
			entries = append(entries, Entry{Port: p, Protocol: "tcp", Action: "allow", Timestamp: now})
		}
	}
	results := DetectOutliers(entries, OutlierOptions{MinZScore: 2.0})
	if len(results) != 0 {
		t.Errorf("expected no outliers for uniform distribution, got %d", len(results))
	}
}

func TestDetectOutliers_SinceFilter(t *testing.T) {
	now := time.Now()
	old := now.Add(-48 * time.Hour)
	// Old high-count entries that should be excluded.
	var entries []Entry
	for i := 0; i < 10; i++ {
		entries = append(entries, Entry{Port: 80, Protocol: "tcp", Action: "alert", Timestamp: old})
	}
	// Recent entries that are uniform.
	for _, p := range []int{443, 22, 8080} {
		entries = append(entries, Entry{Port: p, Protocol: "tcp", Action: "allow", Timestamp: now})
	}
	results := DetectOutliers(entries, OutlierOptions{Since: now.Add(-1 * time.Hour)})
	if len(results) != 0 {
		t.Errorf("expected no outliers after since filter, got %d", len(results))
	}
}

func TestDetectOutliers_ActionFilter(t *testing.T) {
	entries := baseOutlierEntries()
	// Only consider "allow" actions; port 80 (alert) should be excluded.
	results := DetectOutliers(entries, OutlierOptions{Action: "allow"})
	for _, r := range results {
		if r.Port == 80 {
			t.Errorf("port 80 should be excluded by action filter")
		}
	}
}

func TestDetectOutliers_EmptyInput(t *testing.T) {
	results := DetectOutliers(nil, OutlierOptions{})
	if results != nil {
		t.Errorf("expected nil for empty input, got %v", results)
	}
}

func TestDetectOutliers_OrderedByZScore(t *testing.T) {
	entries := baseOutlierEntries()
	results := DetectOutliers(entries, OutlierOptions{MinZScore: 0.1})
	for i := 1; i < len(results); i++ {
		if results[i].ZScore > results[i-1].ZScore {
			t.Errorf("results not sorted descending by z-score at index %d", i)
		}
	}
}
