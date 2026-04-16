package history

import (
	"testing"
	"time"
)

func buildStats() []Entry {
	now := time.Now()
	return []Entry{
		{Port: 80, Protocol: "tcp", Action: "allow", Timestamp: now.Add(-2 * time.Hour)},
		{Port: 80, Protocol: "tcp", Action: "allow", Timestamp: now.Add(-1 * time.Hour)},
		{Port: 80, Protocol: "tcp", Action: "alert", Timestamp: now},
		{Port: 443, Protocol: "tcp", Action: "allow", Timestamp: now.Add(-30 * time.Minute)},
		{Port: 53, Protocol: "udp", Action: "warn", Timestamp: now.Add(-10 * time.Minute)},
	}
}

func TestStats_CountsEntries(t *testing.T) {
	stats := Stats(buildStats())
	if len(stats) != 3 {
		t.Fatalf("expected 3 port stats, got %d", len(stats))
	}
}

func TestStats_SeenCount(t *testing.T) {
	stats := Stats(buildStats())
	for _, s := range stats {
		if s.Port == 80 && s.Protocol == "tcp" {
			if s.Seen != 3 {
				t.Errorf("expected seen=3 for port 80, got %d", s.Seen)
			}
			return
		}
	}
	t.Error("port 80/tcp not found in stats")
}

func TestStats_ActionBreakdown(t *testing.T) {
	stats := Stats(buildStats())
	for _, s := range stats {
		if s.Port == 80 {
			if s.Actions["allow"] != 2 {
				t.Errorf("expected 2 allow actions, got %d", s.Actions["allow"])
			}
			if s.Actions["alert"] != 1 {
				t.Errorf("expected 1 alert action, got %d", s.Actions["alert"])
			}
			return
		}
	}
	t.Error("port 80 not found")
}

func TestStats_LastSeen(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	entries := []Entry{
		{Port: 22, Protocol: "tcp", Action: "alert", Timestamp: now.Add(-5 * time.Minute)},
		{Port: 22, Protocol: "tcp", Action: "alert", Timestamp: now},
	}
	stats := Stats(entries)
	if len(stats) != 1 {
		t.Fatalf("expected 1 stat")
	}
	if !stats[0].LastSeen.Equal(now) {
		t.Errorf("expected LastSeen=%v, got %v", now, stats[0].LastSeen)
	}
}

func TestStats_EmptyInput(t *testing.T) {
	stats := Stats([]Entry{})
	if len(stats) != 0 {
		t.Errorf("expected empty stats, got %d", len(stats))
	}
}
