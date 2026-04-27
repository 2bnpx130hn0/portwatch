package history

import (
	"testing"
	"time"
)

func baseImpactEntries() []Entry {
	now := time.Now()
	return []Entry{
		{Port: 22, Protocol: "TCP", Action: "alert", Timestamp: now.Add(-10 * time.Minute)},
		{Port: 22, Protocol: "TCP", Action: "alert", Timestamp: now.Add(-5 * time.Minute)},
		{Port: 22, Protocol: "TCP", Action: "warn", Timestamp: now.Add(-2 * time.Minute)},
		{Port: 80, Protocol: "tcp", Action: "allow", Timestamp: now.Add(-1 * time.Hour)},
		{Port: 80, Protocol: "tcp", Action: "allow", Timestamp: now.Add(-30 * time.Minute)},
		{Port: 443, Protocol: "tcp", Action: "alert", Timestamp: now.Add(-72 * time.Hour)},
	}
}

func TestImpact_RankedByScore(t *testing.T) {
	entries := baseImpactEntries()
	result := Impact(entries, ImpactOptions{})
	if len(result) == 0 {
		t.Fatal("expected results")
	}
	// port 22 should rank highest (2 alerts + recency boost)
	if result[0].Port != 22 {
		t.Errorf("expected port 22 first, got %d", result[0].Port)
	}
}

func TestImpact_TopN(t *testing.T) {
	entries := baseImpactEntries()
	result := Impact(entries, ImpactOptions{TopN: 2})
	if len(result) != 2 {
		t.Errorf("expected 2 results, got %d", len(result))
	}
}

func TestImpact_FilterByAction(t *testing.T) {
	entries := baseImpactEntries()
	result := Impact(entries, ImpactOptions{Action: "allow"})
	for _, r := range result {
		if r.Alerts > 0 || r.Warnings > 0 {
			t.Errorf("expected only allowed entries, got alerts=%d warnings=%d", r.Alerts, r.Warnings)
		}
	}
	if len(result) != 1 || result[0].Port != 80 {
		t.Errorf("expected port 80, got %+v", result)
	}
}

func TestImpact_SinceFilter(t *testing.T) {
	entries := baseImpactEntries()
	// exclude the old port 443 alert
	result := Impact(entries, ImpactOptions{Since: time.Now().Add(-48 * time.Hour)})
	for _, r := range result {
		if r.Port == 443 {
			t.Error("port 443 should be excluded by since filter")
		}
	}
}

func TestImpact_CountsCorrect(t *testing.T) {
	entries := baseImpactEntries()
	result := Impact(entries, ImpactOptions{})
	var p22 *ImpactEntry
	for i := range result {
		if result[i].Port == 22 {
			p22 = &result[i]
			break
		}
	}
	if p22 == nil {
		t.Fatal("port 22 not found")
	}
	if p22.Alerts != 2 {
		t.Errorf("expected 2 alerts, got %d", p22.Alerts)
	}
	if p22.Warnings != 1 {
		t.Errorf("expected 1 warning, got %d", p22.Warnings)
	}
	if p22.Total != 3 {
		t.Errorf("expected total 3, got %d", p22.Total)
	}
}

func TestImpact_EmptyEntries(t *testing.T) {
	result := Impact([]Entry{}, ImpactOptions{})
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d", len(result))
	}
}
