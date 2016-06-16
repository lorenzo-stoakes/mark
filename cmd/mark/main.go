package main

import (
	"fmt"
	"os"

	"github.com/lorenzo-stoakes/mark"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [markdown files...]\n", os.Args[0])
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}

	bad := false
	for _, path := range os.Args[1:] {
		// Just skip files that don't look like markdown, more
		// convenient for running `mark *` in a mixed directory.
		if !mark.IsMarkdownFile(path) {
			continue
		}

		if doc, err := mark.ParseFile(path); err != nil {
			bad = true
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		} else if str := doc.String(); str != "" {
			bad = true
			fmt.Printf("%s\n", doc)
		}
	}

	if bad {
		os.Exit(1)
	}
}
