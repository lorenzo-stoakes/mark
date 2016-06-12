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

	for _, path := range os.Args[1:] {
		if doc, err := mark.ParseFile(path); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		} else {
			fmt.Printf("%s\n", doc)
		}
	}
}
