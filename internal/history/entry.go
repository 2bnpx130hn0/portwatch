package history

import "time"

// Entry represents a single recorded port event.
type Entry struct {
	Timestamp time.Time         `json:"timestamp"`
	Port      int               `json:"port"`
	Protocol  string            `json:"protocol"`
	Action    string            `json:"action"`
	Tags      []string          `json:"tags,omitempty"`
	Meta      map[string]string `json:"meta,omitempty"`
}
