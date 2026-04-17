package history

import "time"

// Entry represents a single recorded port event in history.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Port      int       `json:"port"`
	Protocol  string    `json:"protocol"`
	Action    string    `json:"action"`
	Message   string    `json:"message,omitempty"`
	Tags      []string  `json:"tags,omitempty"`
}
