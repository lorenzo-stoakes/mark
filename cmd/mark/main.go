package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/lorenzo-stoakes/mark"
)

var fix = flag.Bool("fix", false, "Attempt to fix duplicate/unused references.")

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s <options> [markdown files...]\n\n",
		os.Args[0])
	flag.PrintDefaults()
}

func tryFix(document *mark.Document) {
	if err := document.Fix(); err != nil {
		fmt.Fprintf(os.Stderr, "--fix ERROR: %s: %s\n",
			document.Path, err)
	}
}

func main() {
	// Sync -h/--help with too few params usage output.
	flag.Usage = usage

	flag.Parse()
	args := flag.Args()

	if len(args) < 2 {
		usage()
		os.Exit(1)
	}

	bad := false
	for _, path := range args {
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

			if *fix {
				tryFix(doc)
			}
		}
	}

	if bad {
		os.Exit(1)
	}
}
