package history

import (
	"testing"
	"time"
)

func baseSpikeEntries() []Entry {
	now := time.Now()
	var entries []Entry

	// Port 8080/tcp: 1 event per historical hour, then 10 in current window → spike
	for i := 1; i <= 5; i++ {
		entries = append(entries, Entry{
			Port:      8080,
			Protocol:  "tcp",
			Action:    "alert",
			Timestamp: now.Add(-time.Duration(i) * time.Hour),
		})
	}
	for i := 0; i < 10; i++ {
		entries = append(entries, Entry{
			Port:      8080,
			Protocol:  "tcp",
			Action:    "alert",
			Timestamp: now.Add(-time.Minute * time.Duration(i+1)),
		})
	}

	// Port 443/tcp: steady 2 events per hour — no spike
	for i := 1; i <= 4; i++ {
		for j := 0; j < 2; j++ {
			entries = append(entries, Entry{
				Port:      443,
				Protocol:  "tcp",
				Action:    "allow",
				Timestamp: now.Add(-time.Duration(i)*time.Hour - time.Duration(j)*time.Minute),
			})
		}
	}
	for j := 0; j < 2; j++ {
		entries = append(entries, Entry{
			Port:      443,
			Protocol:  "tcp",
			Action:    "allow",
			Timestamp: now.Add(-time.Duration(j+1) * time.Minute),
		})
	}

	return entries
}

func TestDetectSpikes_FindsSpike(t *testing.T) {
	entries := baseSpikeEntries()
	results := DetectSpikes(entries, SpikeOptions{Window: time.Hour, Multiple: 3.0, MinEvents: 2})
	if len(results) != 1 {
		t.Fatalf("expected 1 spike, got %d", len(results))
	}
	if results[0].Port != 8080 {
		t.Errorf("expected port 8080, got %d", results[0].Port)
	}
	if results[0].Multiple < 3.0 {
		t.Errorf("expected multiple >= 3.0, got %f", results[0].Multiple)
	}
}

func TestDetectSpikes_NoSpikeForSteadyPort(t *testing.T) {
	entries := baseSpikeEntries()
	results := DetectSpikes(entries, SpikeOptions{Window: time.Hour, Multiple: 3.0, MinEvents: 2})
	for _, r := range results {
		if r.Port == 443 {
			t.Errorf("port 443 should not be flagged as a spike")
		}
	}
}

func TestDetectSpikes_MinEventsFilters(t *testing.T) {
	entries := baseSpikeEntries()
	// Require at least 20 events in current window — nothing should match
	results := DetectSpikes(entries, SpikeOptions{Window: time.Hour, Multiple: 2.0, MinEvents: 20})
	if len(results) != 0 {
		t.Errorf("expected 0 results with high MinEvents, got %d", len(results))
	}
}

func TestDetectSpikes_ActionFilter(t *testing.T) {
	entries := baseSpikeEntries()
	// Filter to "allow" — 8080 is alert so should be excluded
	results := DetectSpikes(entries, SpikeOptions{Window: time.Hour, Multiple: 2.0, MinEvents: 1, Action: "allow"})
	for _, r := range results {
		if r.Port == 8080 {
			t.Errorf("port 8080 (alert) should be excluded when filtering by allow")
		}
	}
}

func TestDetectSpikes_EmptyEntries(t *testing.T) {
	results := DetectSpikes([]Entry{}, SpikeOptions{Window: time.Hour, Multiple: 3.0})
	if len(results) != 0 {
		t.Errorf("expected empty results for empty entries, got %d", len(results))
	}
}
