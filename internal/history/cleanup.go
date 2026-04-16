package history

import (
	"time"
)

// CleanupOptions configures how old history entries are pruned.
type CleanupOptions struct {
	MaxAge     time.Duration
	MaxEntries int
}

// Cleanup removes entries from h that exceed the configured limits.
// MaxAge removes entries older than the duration (0 = no limit).
// MaxEntries keeps only the most recent N entries (0 = no limit).
// Returns the number of entries removed.
func (h *History) Cleanup(opts CleanupOptions) (int, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	before := len(h.entries)

	if opts.MaxAge > 0 {
		cutoff := time.Now().UTC().Add(-opts.MaxAge)
		filtered := h.entries[:0]
		for _, e := range h.entries {
			if e.Timestamp.After(cutoff) {
				filtered = append(filtered, e)
			}
		}
		h.entries = filtered
	}

	if opts.MaxEntries > 0 && len(h.entries) > opts.MaxEntries {
		h.entries = h.entries[len(h.entries)-opts.MaxEntries:]
	}

	removed := before - len(h.entries)
	if removed == 0 {
		return 0, nil
	}

	if err := h.persist(); err != nil {
		return 0, err
	}
	return removed, nil
}
