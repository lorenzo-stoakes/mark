package mark

import (
	pathLib "path"
	"strings"
)

// Does the path look like a markdown file?
func IsMarkdownFile(path string) bool {
	switch strings.ToLower(pathLib.Ext(path)) {
	case "md", "markdown":
		return true
	}

	return false
}
