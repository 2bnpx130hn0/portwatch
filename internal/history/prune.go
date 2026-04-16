package history

import "time"

// PruneOptions controls what gets removed during a prune operation.
type PruneOptions struct {
	MaxAge     time.Duration
	MaxEntries int
}

// Prune removes old entries from the history according to the given options
// and persists the result to disk. It combines age-based and count-based
// trimming in a single pass so callers don't need to call Cleanup twice.
func (h *History) Prune(opts PruneOptions) (removed int, err error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	before := len(h.entries)

	if opts.MaxAge > 0 {
		cutoff := time.Now().Add(-opts.MaxAge)
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

	removed = before - len(h.entries)

	if err = h.save(); err != nil {
		return 0, err
	}
	return removed, nil
}
