package history

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

var baseFlowEntries = func() []Entry {
	now := time.Now().Truncate(time.Second)
	return []Entry{
		{Port: 80, Protocol: "tcp", Action: "allow", Timestamp: now},
		{Port: 443, Protocol: "tcp", Action: "allow", Timestamp: now.Add(30 * time.Second)},
		{Port: 80, Protocol: "tcp", Action: "allow", Timestamp: now.Add(60 * time.Second)},
		{Port: 443, Protocol: "tcp", Action: "allow", Timestamp: now.Add(90 * time.Second)},
		{Port: 8080, Protocol: "tcp", Action: "alert", Timestamp: now.Add(200 * time.Second)},
		// gap > 5m — should NOT form an edge with previous
		{Port: 9090, Protocol: "tcp", Action: "alert", Timestamp: now.Add(10 * time.Minute)},
	}
}()

func TestBuildFlow_DetectsEdge(t *testing.T) {
	edges := BuildFlow(baseFlowEntries, FlowOptions{})
	if len(edges) == 0 {
		t.Fatal("expected at least one edge")
	}
	// 80->443 should appear twice
	var found *FlowEdge
	for i := range edges {
		if edges[i].FromPort == 80 && edges[i].ToPort == 443 {
			found = &edges[i]
		}
	}
	if found == nil {
		t.Fatal("expected edge 80->443")
	}
	if found.Count != 2 {
		t.Errorf("expected count 2, got %d", found.Count)
	}
}

func TestBuildFlow_RespectsWindow(t *testing.T) {
	edges := BuildFlow(baseFlowEntries, FlowOptions{})
	// 8080->9090 gap is ~8m, default window 5m — should NOT appear
	for _, e := range edges {
		if e.FromPort == 8080 && e.ToPort == 9090 {
			t.Error("edge 8080->9090 should be excluded by window")
		}
	}
}

func TestBuildFlow_FilterByAction(t *testing.T) {
	edges := BuildFlow(baseFlowEntries, FlowOptions{Action: "alert"})
	for _, e := range edges {
		if e.FromPort == 80 || e.ToPort == 80 {
			t.Error("port 80 (allow) should be excluded")
		}
	}
}

func TestBuildFlow_MinCount(t *testing.T) {
	edges := BuildFlow(baseFlowEntries, FlowOptions{MinCount: 2})
	for _, e := range edges {
		if e.Count < 2 {
			t.Errorf("expected count >= 2, got %d", e.Count)
		}
	}
}

func TestBuildFlow_OrderedByCountDesc(t *testing.T) {
	edges := BuildFlow(baseFlowEntries, FlowOptions{})
	for i := 1; i < len(edges); i++ {
		if edges[i].Count > edges[i-1].Count {
			t.Error("edges not sorted by count descending")
		}
	}
}

func TestRenderFlow_Text(t *testing.T) {
	edges := []FlowEdge{{FromPort: 80, ToPort: 443, Protocol: "tcp", Count: 3}}
	var buf bytes.Buffer
	if err := RenderFlow(edges, "text", &buf); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "80") || !strings.Contains(buf.String(), "443") {
		t.Error("expected ports in output")
	}
}

func TestRenderFlow_JSON(t *testing.T) {
	edges := []FlowEdge{{FromPort: 22, ToPort: 80, Protocol: "tcp", Count: 1}}
	var buf bytes.Buffer
	if err := RenderFlow(edges, "json", &buf); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "FromPort") {
		t.Error("expected JSON field in output")
	}
}

func TestRenderFlow_Empty(t *testing.T) {
	var buf bytes.Buffer
	if err := RenderFlow(nil, "text", &buf); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "no flow") {
		t.Error("expected empty message")
	}
}
