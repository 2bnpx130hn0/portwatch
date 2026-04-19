package history

import (
	"testing"
	"time"
)

func baseClusterEntries() []Entry {
	now := time.Now()
	return []Entry{
		{Port: 80, Protocol: "tcp", Action: "alert", Timestamp: now},
		{Port: 80, Protocol: "tcp", Action: "alert", Timestamp: now},
		{Port: 80, Protocol: "tcp", Action: "allow", Timestamp: now},
		{Port: 443, Protocol: "tcp", Action: "allow", Timestamp: now},
		{Port: 443, Protocol: "tcp", Action: "allow", Timestamp: now},
		{Port: 443, Protocol: "tcp", Action: "allow", Timestamp: now},
		{Port: 22, Protocol: "tcp", Action: "alert", Timestamp: now},
		{Port: 53, Protocol: "udp", Action: "allow", Timestamp: now},
	}
}

func TestCluster_GroupsByPortProtocol(t *testing.T) {
	results := Cluster(baseClusterEntries(), ClusterOptions{})
	if len(results) != 4 {
		t.Fatalf("expected 4 clusters, got %d", len(results))
	}
	// highest count first
	if results[0].Port != 443 || results[0].Count != 3 {
		t.Errorf("expected port 443 first with count 3, got port %d count %d", results[0].Port, results[0].Count)
	}
}

func TestCluster_FilterByAction(t *testing.T) {
	results := Cluster(baseClusterEntries(), ClusterOptions{Action: "alert"})
	for _, r := range results {
		for _, e := range r.Entries {
			if !equalFold(e.Action, "alert") {
				t.Errorf("expected only alert entries, got %s", e.Action)
			}
		}
	}
}

func TestCluster_MinCount(t *testing.T) {
	results := Cluster(baseClusterEntries(), ClusterOptions{MinCount: 3})
	for _, r := range results {
		if r.Count < 3 {
			t.Errorf("expected count >= 3, got %d", r.Count)
		}
	}
	if len(results) != 1 {
		t.Errorf("expected 1 cluster with MinCount=3, got %d", len(results))
	}
}

func TestCluster_Empty(t *testing.T) {
	results := Cluster([]Entry{}, ClusterOptions{})
	if len(results) != 0 {
		t.Errorf("expected empty results, got %d", len(results))
	}
}

func TestCluster_CaseInsensitiveAction(t *testing.T) {
	entries := []Entry{
		{Port: 80, Protocol: "tcp", Action: "ALERT"},
		{Port: 80, Protocol: "tcp", Action: "alert"},
	}
	results := Cluster(entries, ClusterOptions{Action: "alert"})
	if len(results) != 1 || results[0].Count != 2 {
		t.Errorf("expected 1 cluster with count 2, got %d clusters", len(results))
	}
}
