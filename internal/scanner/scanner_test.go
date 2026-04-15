package scanner

import (
	"net"
	"strconv"
	"testing"
	"time"
)

// startTestServer starts a TCP listener on a random port and returns the port and a stop func.
func startTestServer(t *testing.T) (int, func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start test server: %v", err)
	}
	port, _ := strconv.Atoi(ln.Addr().(*net.TCPAddr).Port.String())
	_ = port
	actualPort := ln.Addr().(*net.TCPAddr).Port
	return actualPort, func() { ln.Close() }
}

func TestScan_OpenPort(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("could not start listener: %v", err)
	}
	defer ln.Close()
	port := ln.Addr().(*net.TCPAddr).Port

	s := New("127.0.0.1", time.Second)
	snap, err := s.Scan([]int{port})
	if err != nil {
		t.Fatalf("Scan returned error: %v", err)
	}
	if len(snap.Ports) != 1 {
		t.Fatalf("expected 1 port state, got %d", len(snap.Ports))
	}
	if !snap.Ports[0].Open {
		t.Errorf("expected port %d to be open", port)
	}
}

func TestScan_ClosedPort(t *testing.T) {
	// Port 1 is almost certainly closed and requires no privilege to connect to.
	s := New("127.0.0.1", 200*time.Millisecond)
	snap, err := s.Scan([]int{1})
	if err != nil {
		t.Fatalf("Scan returned error: %v", err)
	}
	if snap.Ports[0].Open {
		t.Errorf("expected port 1 to be closed")
	}
}

func TestScan_NoPorts(t *testing.T) {
	s := New("127.0.0.1", time.Second)
	_, err := s.Scan([]int{})
	if err == nil {
		t.Error("expected error when no ports are provided, got nil")
	}
}

func TestOpenPorts(t *testing.T) {
	snap := &Snapshot{
		Ports: []PortState{
			{Port: 80, Open: true},
			{Port: 81, Open: false},
			{Port: 443, Open: true},
		},
	}
	open := OpenPorts(snap)
	if len(open) != 2 {
		t.Errorf("expected 2 open ports, got %d", len(open))
	}
}
