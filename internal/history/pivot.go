package history

import (
	"io"
	"sort"
	"strings"
)

// PivotCell holds the count of events for a specific row/column intersection.
type PivotCell struct {
	Count int
}

// PivotTable represents a 2D cross-tabulation of history entries.
type PivotTable struct {
	RowKey  string
	ColKey  string
	Rows    []string
	Cols    []string
	Cells   map[string]map[string]int
}

// PivotOptions controls how the pivot table is built.
type PivotOptions struct {
	RowField  string // "port", "protocol", "action"
	ColField  string // "port", "protocol", "action"
	Action    string
	Protocol  string
	Since     int64 // unix seconds; 0 means no filter
}

// Pivot builds a cross-tabulation of entries by two fields.
func Pivot(entries []Entry, opts PivotOptions) PivotTable {
	rowField := strings.ToLower(opts.RowField)
	colField := strings.ToLower(opts.ColField)

	rowSet := map[string]struct{}{}
	colSet := map[string]struct{}{}
	cells := map[string]map[string]int{}

	for _, e := range entries {
		if opts.Since > 0 && e.Timestamp.Unix() < opts.Since {
			continue
		}
		if opts.Action != "" && !strings.EqualFold(e.Action, opts.Action) {
			continue
		}
		if opts.Protocol != "" && !strings.EqualFold(e.Protocol, opts.Protocol) {
			continue
		}

		row := fieldValue(e, rowField)
		col := fieldValue(e, colField)

		rowSet[row] = struct{}{}
		colSet[col] = struct{}{}

		if cells[row] == nil {
			cells[row] = map[string]int{}
		}
		cells[row][col]++
	}

	rows := sortedKeys(rowSet)
	cols := sortedKeys(colSet)

	return PivotTable{
		RowKey: rowField,
		ColKey: colField,
		Rows:   rows,
		Cols:   cols,
		Cells:  cells,
	}
}

func fieldValue(e Entry, field string) string {
	switch field {
	case "port":
		return itoa(e.Port)
	case "protocol":
		return strings.ToLower(e.Protocol)
	case "action":
		return strings.ToLower(e.Action)
	default:
		return "unknown"
	}
}

func sortedKeys(m map[string]struct{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// RenderPivot writes the pivot table to w in the given format.
func RenderPivot(pt PivotTable, format string, w io.Writer) {
	if w == nil {
		return
	}
	switch strings.ToLower(format) {
	case "json":
		renderPivotJSON(pt, w)
	default:
		renderPivotText(pt, w)
	}
}
