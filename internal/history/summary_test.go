package history

import (
	"testing"
	"time"
)

func buildEntries() []Entry {
	now := time.Now()
	return []Entry{
		{Timestamp: now, Protocol: "tcp", Port: 80, Action: "allow"},
		{Timestamp: now, Protocol: "tcp", Port: 80, Action: "allow"},
		{Timestamp: now, Protocol: "tcp", Port: 443, Action: "allow"},
		{Timestamp: now, Protocol: "udp", Port: 53, Action: "alert"},
		{Timestamp: now, Protocol: "udp", Port: 53, Action: "alert"},
		{Timestamp: now, Protocol: "tcp", Port: 8080, Action: "warn"},
	}
}

func TestSummarize_Total(t *testing.T) {
	s := Summarize(buildEntries())
	if s.Total != 6 {
		t.Errorf("expected total 6, got %d", s.Total)
	}
}

func TestSummarize_ByAction(t *testing.T) {
	s := Summarize(buildEntries())
	if s.ByAction["allow"] != 3 {
		t.Errorf("expected 3 allow, got %d", s.ByAction["allow"])
	}
	if s.ByAction["alert"] != 2 {
		t.Errorf("expected 2 alert, got %d", s.ByAction["alert"])
	}
	if s.ByAction["warn"] != 1 {
		t.Errorf("expected 1 warn, got %d", s.ByAction["warn"])
	}
}

func TestSummarize_ByProtocol(t *testing.T) {
	s := Summarize(buildEntries())
	if s.ByProtocol["tcp"] != 4 {
		t.Errorf("expected 4 tcp, got %d", s.ByProtocol["tcp"])
	}
	if s.ByProtocol["udp"] != 2 {
		t.Errorf("expected 2 udp, got %d", s.ByProtocol["udp"])
	}
}

func TestSummarize_TopPorts(t *testing.T) {
	s := Summarize(buildEntries())
	if len(s.TopPorts) == 0 {
		t.Fatal("expected top ports to be populated")
	}
	if s.TopPorts[0].Port != 80 && s.TopPorts[0].Port != 53 {
		t.Errorf("expected top port to be 80 or 53, got %d", s.TopPorts[0].Port)
	}
	if s.TopPorts[0].Count < 2 {
		t.Errorf("expected top port count >= 2, got %d", s.TopPorts[0].Count)
	}
}

func TestSummarize_Empty(t *testing.T) {
	s := Summarize([]Entry{})
	if s.Total != 0 {
		t.Errorf("expected total 0, got %d", s.Total)
	}
	if len(s.TopPorts) != 0 {
		t.Errorf("expected no top ports, got %d", len(s.TopPorts))
	}
}
