package render

import (
	"fmt"
	"le-grimoire/internal/state"
)

func BuildPlaceholderHTML(p state.Page) string {
	return fmt.Sprintf("<h1>%s</h1><p>cursor: %d</p>", p.Type, p.State["cursor"])
}
