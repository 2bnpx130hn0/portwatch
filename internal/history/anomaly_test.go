package history

import (
	"testing"
	"time"
)

func baseAnomalyEntries() []Entry {
	now := time.Now()
	entries := []Entry{}
	// Port 80/tcp appears 10 times — anomalously high
	for i := 0; i < 10; i++ {
		entries = append(entries, Entry{Port: 80, Protocol: "tcp", Action: "allow", Timestamp: now})
	}
	// Other ports appear once each
	for _, p := range []int{22, 443, 8080, 9090} {
		entries = append(entries, Entry{Port: p, Protocol: "tcp", Action: "allow", Timestamp: now})
	}
	return entries
}

func TestDetectAnomalies_FindsHighCount(t *testing.T) {
	entries := baseAnomalyEntries()
	results := DetectAnomalies(entries, AnomalyOptions{Threshold: 1.5, MinSamples: 3})
	if len(results) == 0 {
		t.Fatal("expected at least one anomaly, got none")
	}
	found := false
	for _, r := range results {
		if r.Port == 80 && r.Protocol == "tcp" {
			found = true
			if r.Count != 10 {
				t.Errorf("expected count 10, got %d", r.Count)
			}
			if r.ZScore <= 0 {
				t.Errorf("expected positive z-score, got %f", r.ZScore)
			}
		}
	}
	if !found {
		t.Error("expected port 80/tcp to be flagged as anomaly")
	}
}

func TestDetectAnomalies_BelowThreshold(t *testing.T) {
	now := time.Now()
	// All ports appear exactly twice — no anomaly
	entries := []Entry{}
	for _, p := range []int{80, 443, 8080, 9090, 3000} {
		entries = append(entries, Entry{Port: p, Protocol: "tcp", Action: "allow", Timestamp: now})
		entries = append(entries, Entry{Port: p, Protocol: "tcp", Action: "allow", Timestamp: now})
	}
	results := DetectAnomalies(entries, AnomalyOptions{Threshold: 2.0, MinSamples: 3})
	if len(results) != 0 {
		t.Errorf("expected no anomalies, got %d", len(results))
	}
}

func TestDetectAnomalies_SinceFilter(t *testing.T) {
	old := time.Now().Add(-2 * time.Hour)
	recent := time.Now()
	entries := []Entry{
		{Port: 80, Protocol: "tcp", Action: "alert", Timestamp: old},
		{Port: 80, Protocol: "tcp", Action: "alert", Timestamp: old},
		{Port: 80, Protocol: "tcp", Action: "alert", Timestamp: old},
		{Port: 22, Protocol: "tcp", Action: "allow", Timestamp: recent},
		{Port: 443, Protocol: "tcp", Action: "allow", Timestamp: recent},
		{Port: 8080, Protocol: "tcp", Action: "allow", Timestamp: recent},
		{Port: 9090, Protocol: "tcp", Action: "allow", Timestamp: recent},
	}
	// With since=now-1h, old port 80 entries should be excluded
	results := DetectAnomalies(entries, AnomalyOptions{
		Since:     time.Now().Add(-1 * time.Hour),
		Threshold: 1.5,
		MinSamples: 3,
	})
	for _, r := range results {
		if r.Port == 80 {
			t.Error("port 80 should not appear after since filter excludes old entries")
		}
	}
}

func TestDetectAnomalies_TooFewSamples(t *testing.T) {
	now := time.Now()
	entries := []Entry{
		{Port: 80, Protocol: "tcp", Action: "allow", Timestamp: now},
		{Port: 443, Protocol: "tcp", Action: "allow", Timestamp: now},
	}
	results := DetectAnomalies(entries, AnomalyOptions{MinSamples: 5})
	if len(results) != 0 {
		t.Errorf("expected no results with too few samples, got %d", len(results))
	}
}

func TestDetectAnomalies_DefaultOptions(t *testing.T) {
	now := time.Now()
	entries := []Entry{}
	for i := 0; i < 20; i++ {
		entries = append(entries, Entry{Port: 80, Protocol: "tcp", Action: "alert", Timestamp: now})
	}
	for _, p := range []int{22, 443, 8080} {
		entries = append(entries, Entry{Port: p, Protocol: "tcp", Action: "allow", Timestamp: now})
	}
	// Zero-value options should use defaults (threshold=2.0, minSamples=3)
	results := DetectAnomalies(entries, AnomalyOptions{})
	if len(results) == 0 {
		t.Error("expected anomaly with default options")
	}
}
