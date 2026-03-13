package output

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"
)

// PrintJSON prints data as indented JSON.
func PrintJSON(data any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(data)
}

// Table is a simple table printer backed by tabwriter.
type Table struct {
	w *tabwriter.Writer
}

// NewTable creates a new table and prints the header row.
func NewTable(headers ...string) *Table {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	for i, h := range headers {
		if i > 0 {
			fmt.Fprint(w, "\t")
		}
		fmt.Fprint(w, h)
	}
	fmt.Fprintln(w)
	return &Table{w: w}
}

// AddRow adds a row to the table.
func (t *Table) AddRow(cols ...any) {
	for i, c := range cols {
		if i > 0 {
			fmt.Fprint(t.w, "\t")
		}
		fmt.Fprint(t.w, c)
	}
	fmt.Fprintln(t.w)
}

// Flush writes the table to stdout.
func (t *Table) Flush() {
	t.w.Flush()
}
