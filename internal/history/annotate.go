package history

import "time"

// Annotation holds a user-defined note attached to a history entry.
type Annotation struct {
	Note      string    `json:"note"`
	CreatedAt time.Time `json:"created_at"`
}

// Annotate adds or replaces the annotation on entries matching port and protocol.
// Returns the number of entries updated.
func Annotate(entries []Entry, port int, protocol, note string) ([]Entry, int) {
	updated := 0
	for i, e := range entries {
		if e.Port == port && equalFold(e.Protocol, protocol) {
			if entries[i].Meta == nil {
				entries[i].Meta = map[string]string{}
			}
			entries[i].Meta["annotation"] = note
			entries[i].Meta["annotation_at"] = time.Now().UTC().Format(time.RFC3339)
			updated++
		}
	}
	return entries, updated
}

// ClearAnnotation removes the annotation from entries matching port and protocol.
func ClearAnnotation(entries []Entry, port int, protocol string) ([]Entry, int) {
	updated := 0
	for i, e := range entries {
		if e.Port == port && equalFold(e.Protocol, protocol) {
			if entries[i].Meta != nil {
				delete(entries[i].Meta, "annotation")
				delete(entries[i].Meta, "annotation_at")
				updated++
			}
		}
	}
	return entries, updated
}

// FilterAnnotated returns only entries that have an annotation.
func FilterAnnotated(entries []Entry) []Entry {
	var out []Entry
	for _, e := range entries {
		if e.Meta != nil {
			if _, ok := e.Meta["annotation"]; ok {
				out = append(out, e)
			}
		}
	}
	return out
}

func equalFold(a, b string) bool {
	return len(a) == len(b) && strings.EqualFold(a, b)
}
