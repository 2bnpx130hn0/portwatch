package scanner

import (
	"fmt"
	"net"
	"time"
)

// PortState represents the state of a single port.
type PortState struct {
	Port     int
	Protocol string
	Open     bool
}

// Snapshot is a collection of port states captured at a point in time.
type Snapshot struct {
	Timestamp time.Time
	Ports     []PortState
}

// Scanner scans a range of ports on a given host.
type Scanner struct {
	Host    string
	Timeout time.Duration
}

// New creates a new Scanner for the given host with the specified timeout.
func New(host string, timeout time.Duration) *Scanner {
	return &Scanner{
		Host:    host,
		Timeout: timeout,
	}
}

// Scan checks the given list of ports and returns a Snapshot of their states.
func (s *Scanner) Scan(ports []int) (*Snapshot, error) {
	if len(ports) == 0 {
		return nil, fmt.Errorf("no ports specified for scanning")
	}

	states := make([]PortState, 0, len(ports))

	for _, port := range ports {
		address := fmt.Sprintf("%s:%d", s.Host, port)
		conn, err := net.DialTimeout("tcp", address, s.Timeout)
		open := err == nil
		if conn != nil {
			conn.Close()
		}
		states = append(states, PortState{
			Port:     port,
			Protocol: "tcp",
			Open:     open,
		})
	}

	return &Snapshot{
		Timestamp: time.Now(),
		Ports:     states,
	}, nil
}

// OpenPorts returns only the open ports from a Snapshot.
func OpenPorts(snap *Snapshot) []PortState {
	var open []PortState
	for _, p := range snap.Ports {
		if p.Open {
			open = append(open, p)
		}
	}
	return open
}
