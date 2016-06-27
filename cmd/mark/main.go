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

// Check specified path, return true if error or issues detected.
func check(path string, doFix bool) bool {
	var (
		doc *mark.Document
		err error
	)

	if doc, err = mark.ParseFile(path); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return true
	}

	if doFix {
		tryFix(doc)

		// Re-parse.
		// TODO: Do this on the document side without
		//       having to re-read the file.
		return check(path, false)
	}

	if str := doc.String(); str != "" {
		fmt.Printf("%s\n", doc)

		return true
	}

	return false
}

func main() {
	// Sync -h/--help with too few params usage output.
	flag.Usage = usage

	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		usage()
		os.Exit(1)
	}

	someBad := false
	for _, path := range args {
		// Just skip files that don't look like markdown, more
		// convenient for running `mark *` in a mixed directory.
		if !mark.IsMarkdownFile(path) {
			continue
		}

		if bad := check(path, *fix); bad {
			someBad = true
		}
	}

	if someBad {
		os.Exit(1)
	}
}
