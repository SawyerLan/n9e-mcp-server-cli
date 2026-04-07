package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// Render writes data to stdout in the requested format.
// Currently only JSON is supported; additional formats (table, etc.) will be added later.
func Render(format string, data any) error {
	return RenderTo(os.Stdout, format, data)
}

// RenderTo writes data to the given writer in the requested format.
func RenderTo(w io.Writer, format string, data any) error {
	switch format {
	case "json", "":
		return renderJSON(w, data)
	default:
		return fmt.Errorf("unsupported output format: %s", format)
	}
}

func renderJSON(w io.Writer, data any) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(data)
}
