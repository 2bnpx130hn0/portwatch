package history

import "fmt"

// PinCmd provides pin/unpin/list operations on a History file.
type PinCmd struct {
	history *History
}

// NewPinCmd creates a PinCmd backed by the given History.
func NewPinCmd(h *History) *PinCmd {
	return &PinCmd{history: h}
}

// PinPort pins all entries matching port+protocol and persists.
func (c *PinCmd) PinPort(port int, protocol string) error {
	c.history.mu.Lock()
	defer c.history.mu.Unlock()
	c.history.entries = Pin(c.history.entries, port, protocol)
	return c.history.save()
}

// UnpinPort unpins all entries matching port+protocol and persists.
func (c *PinCmd) UnpinPort(port int, protocol string) error {
	c.history.mu.Lock()
	defer c.history.mu.Unlock()
	c.history.entries = Unpin(c.history.entries, port, protocol)
	return c.history.save()
}

// ListPinned prints all pinned entries to stdout.
func (c *PinCmd) ListPinned() {
	c.history.mu.Lock()
	pinned := FilterPinned(c.history.entries)
	c.history.mu.Unlock()
	if len(pinned) == 0 {
		fmt.Println("no pinned entries")
		return
	}
	for _, e := range pinned {
		fmt.Printf("port=%d protocol=%s action=%s pinned_at=%s\n",
			e.Port, e.Protocol, e.Action, e.Labels["pinned_at"])
	}
}
