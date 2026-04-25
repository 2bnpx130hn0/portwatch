package history

import (
	"testing"
	"time"
)

func baseVelocityEntries() []Entry {
	now := time.Now()
	return []Entry{
		// current window: 3 events on port 80
		{Port: 80, Protocol: "tcp", Action: "alert", Timestamp: now.Add(-10 * time.Minute)},
		{Port: 80, Protocol: "tcp", Action: "alert", Timestamp: now.Add(-20 * time.Minute)},
		{Port: 80, Protocol: "tcp", Action: "alert", Timestamp: now.Add(-30 * time.Minute)},
		// previous window: 1 event on port 80
		{Port: 80, Protocol: "tcp", Action: "alert", Timestamp: now.Add(-70 * time.Minute)},
		// current window: 1 event on port 443
		{Port: 443, Protocol: "tcp", Action: "allow", Timestamp: now.Add(-5 * time.Minute)},
		// previous window: 1 event on port 443
		{Port: 443, Protocol: "tcp", Action: "allow", Timestamp: now.Add(-90 * time.Minute)},
		// old entry beyond lookback
		{Port: 22, Protocol: "tcp", Action: "alert", Timestamp: now.Add(-48 * time.Hour)},
	}
}

func TestVelocity_RateCalculated(t *testing.T) {
	entries := baseVelocityEntries()
	results := Velocity(entries, VelocityOptions{})
	if len(results) == 0 {
		t.Fatal("expected velocity results, got none")
	}
	var port80 *VelocityEntry
	for i := range results {
		if results[i].Port == 80 {
			port80 = &results[i]
			break
		}
	}
	if port80 == nil {
		t.Fatal("expected entry for port 80")
	}
	if port80.Rate <= 0 {
		t.Errorf("expected positive rate for port 80, got %v", port80.Rate)
	}
}

func TestVelocity_DeltaPositiveForIncreasing(t *testing.T) {
	entries := baseVelocityEntries()
	results := Velocity(entries, VelocityOptions{})
	for _, r := range results {
		if r.Port == 80 && r.Protocol == "tcp" {
			if r.Delta <= 0 {
				t.Errorf("expected positive delta for port 80 (more events in current window), got %v", r.Delta)
			}
			return
		}
	}
	t.Error("port 80 tcp not found in results")
}

func TestVelocity_FilterByAction(t *testing.T) {
	entries := baseVelocityEntries()
	results := Velocity(entries, VelocityOptions{Action: "allow"})
	for _, r := range results {
		if r.Action != "allow" {
			t.Errorf("expected only allow entries, got action=%q", r.Action)
		}
	}
}

func TestVelocity_MinEventsFilters(t *testing.T) {
	now := time.Now()
	entries := []Entry{
		{Port: 9999, Protocol: "tcp", Action: "alert", Timestamp: now.Add(-5 * time.Minute)},
	}
	results := Velocity(entries, VelocityOptions{MinEvents: 3})
	for _, r := range results {
		if r.Port == 9999 {
			t.Error("expected port 9999 to be filtered out by MinEvents")
		}
	}
}

func TestVelocity_OldEntriesExcluded(t *testing.T) {
	entries := baseVelocityEntries()
	results := Velocity(entries, VelocityOptions{})
	for _, r := range results {
		if r.Port == 22 {
			t.Error("expected port 22 (old entry) to be excluded from velocity")
		}
	}
}

func TestVelocity_OrderedByAbsDeltaDescending(t *testing.T) {
	entries := baseVelocityEntries()
	results := Velocity(entries, VelocityOptions{})
	for i := 1; i < len(results); i++ {
		prev := results[i-1]
		curr := results[i]
		prevAbs := prev.Delta
		if prevAbs < 0 {
			prevAbs = -prevAbs
		}
		currAbs := curr.Delta
		if currAbs < 0 {
			currAbs = -currAbs
		}
		if prevAbs < currAbs {
			t.Errorf("results not sorted by abs(delta) descending at index %d", i)
		}
	}
}
