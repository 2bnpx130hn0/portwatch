package history

import "time"

// NoteEntry attaches a free-text note to all entries matching port+protocol.
// The note and its timestamp are stored in the entry's Labels map under the
// keys "note" and "note_at" respectively.
func NoteEntry(entries []Entry, port int, protocol, note string) []Entry {
	updated := make([]Entry, len(entries))
	for i, e := range entries {
		if e.Port == port && equalFold(e.Protocol, protocol) {
			if e.Labels == nil {
				e.Labels = map[string]string{}
			}
			e.Labels["note"] = note
			e.Labels["note_at"] = time.Now().UTC().Format(time.RFC3339)
		}
		updated[i] = e
	}
	return updated
}

// RemoveNote clears the note label from matching entries.
func RemoveNote(entries []Entry, port int, protocol string) []Entry {
	updated := make([]Entry, len(entries))
	for i, e := range entries {
		if e.Port == port && equalFold(e.Protocol, protocol) {
			if e.Labels != nil {
				delete(e.Labels, "note")
				delete(e.Labels, "note_at")
			}
		}
		updated[i] = e
	}
	return updated
}

// FilterNoted returns only entries that have a note label set.
func FilterNoted(entries []Entry) []Entry {
	var out []Entry
	for _, e := range entries {
		if e.Labels != nil {
			if v, ok := e.Labels["note"]; ok && v != "" {
				out = append(out, e)
			}
		}
	}
	return out
}

// GetNote returns the note text and timestamp for the first entry matching
// port+protocol, along with a boolean indicating whether a note was found.
func GetNote(entries []Entry, port int, protocol string) (note, noteAt string, ok bool) {
	for _, e := range entries {
		if e.Port == port && equalFold(e.Protocol, protocol) {
			if e.Labels != nil {
				if v := e.Labels["note"]; v != "" {
					return v, e.Labels["note_at"], true
				}
			}
		}
	}
	return "", "", false
}
