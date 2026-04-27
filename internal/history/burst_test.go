package history

import (
	"testing"
	"time"
)

func baseBurstEntries() []Entry {
	now := time.Now()
	entries := []Entry{}
	// 6 events on port 443/tcp within 30 minutes — should burst
	for i := 0; i < 6; i++ {
		entries = append(entries, Entry{
			Port:      443,
			Protocol:  "tcp",
			Action:    "alert",
			Timestamp: now.Add(time.Duration(i*4) * time.Minute),
		})
	}
	// 2 events on port 80/tcp spread over 2 hours — should NOT burst
	entries = append(entries, Entry{Port: 80, Protocol: "tcp", Action: "alert", Timestamp: now})
	entries = append(entries, Entry{Port: 80, Protocol: "tcp", Action: "alert", Timestamp: now.Add(2 * time.Hour)})
	return entries
}

func TestDetectBursts_FindsBurst(t *testing.T) {
	entries := baseBurstEntries()
	results := DetectBursts(entries, BurstOptions{
		Window:    time.Hour,
		Threshold: 5,
	})
	if len(results) != 1 {
		t.Fatalf("expected 1 burst result, got %d", len(results))
	}
	if results[0].Port != 443 {
		t.Errorf("expected port 443, got %d", results[0].Port)
	}
	if results[0].Count < 5 {
		t.Errorf("expected count >= 5, got %d", results[0].Count)
	}
}

func TestDetectBursts_BelowThreshold(t *testing.T) {
	entries := baseBurstEntries()
	results := DetectBursts(entries, BurstOptions{
		Window:    time.Hour,
		Threshold: 10,
	})
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestDetectBursts_FilterByAction(t *testing.T) {
	entries := baseBurstEntries()
	results := DetectBursts(entries, BurstOptions{
		Window:    time.Hour,
		Threshold: 5,
		Action:    "allow",
	})
	if len(results) != 0 {
		t.Errorf("expected 0 results for action=allow, got %d", len(results))
	}
}

func TestDetectBursts_FilterByProtocol(t *testing.T) {
	entries := baseBurstEntries()
	results := DetectBursts(entries, BurstOptions{
		Window:    time.Hour,
		Threshold: 5,
		Protocol:  "udp",
	})
	if len(results) != 0 {
		t.Errorf("expected 0 results for protocol=udp, got %d", len(results))
	}
}

func TestDetectBursts_SinceFilter(t *testing.T) {
	entries := baseBurstEntries()
	// Set since to future — all entries excluded.
	results := DetectBursts(entries, BurstOptions{
		Window:    time.Hour,
		Threshold: 5,
		Since:     time.Now().Add(24 * time.Hour),
	})
	if len(results) != 0 {
		t.Errorf("expected 0 results with future since, got %d", len(results))
	}
}

func TestDetectBursts_DefaultOptions(t *testing.T) {
	now := time.Now()
	var entries []Entry
	for i := 0; i < 6; i++ {
		entries = append(entries, Entry{
			Port:      8080,
			Protocol:  "tcp",
			Action:    "alert",
			Timestamp: now.Add(time.Duration(i*5) * time.Minute),
		})
	}
	// Zero-value options should use defaults (window=1h, threshold=5).
	results := DetectBursts(entries, BurstOptions{})
	if len(results) != 1 {
		t.Fatalf("expected 1 result with default options, got %d", len(results))
	}
}
