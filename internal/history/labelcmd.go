package history

import "fmt"

// LabelCmd provides high-level label operations on a History store.
type LabelCmd struct {
	h *History
}

// NewLabelCmd creates a LabelCmd backed by the given History.
func NewLabelCmd(h *History) *LabelCmd {
	return &LabelCmd{h: h}
}

// Set applies key=value label to all entries matching port+protocol and persists.
func (c *LabelCmd) Set(port int, protocol, key, value string) error {
	c.h.mu.Lock()
	defer c.h.mu.Unlock()
	c.h.entries = Label(c.h.entries, port, protocol, key, value)
	return c.h.save()
}

// Remove removes a label key from entries matching port+protocol and persists.
func (c *LabelCmd) Remove(port int, protocol, key string) error {
	c.h.mu.Lock()
	defer c.h.mu.Unlock()
	c.h.entries = RemoveLabel(c.h.entries, port, protocol, key)
	return c.h.save()
}

// List prints all distinct label keys and values present across all entries.
func (c *LabelCmd) List() []string {
	c.h.mu.Lock()
	defer c.h.mu.Unlock()
	seen := map[string]struct{}{}
	var out []string
	for _, e := range c.h.entries {
		for k, v := range e.Labels {
			token := fmt.Sprintf("%s=%s", k, v)
			if _, ok := seen[token]; !ok {
				seen[token] = struct{}{}
				out = append(out, token)
			}
		}
	}
	return out
}
