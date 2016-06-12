package mark

import (
	"fmt"
	"strings"
)

type Location struct {
	Line, Col int
}

func (l *Location) String() string {
	return fmt.Sprintf("(%d, %d)", l.Line, l.Col)
}

type Reference struct {
	Location
	Name, Uri string

	Duplicates []*Reference
}

type Document struct {
	Lines            int
	References       []*Reference
	ReferencesByName map[string]*Reference
	Referenced       map[string][]Location
}

func newDocument() *Document {
	return &Document{Referenced: make(map[string][]Location),
		ReferencesByName: make(map[string]*Reference)}
}

func (d *Document) referTo(name string, loc Location) {
	d.Referenced[name] = append(d.Referenced[name], loc)
}

func (d *Document) define(name, uri string, loc Location) {
	ref := &Reference{Location: loc, Name: name, Uri: uri}

	d.References = append(d.References, ref)

	// We choose the first reference to be canonical.
	if first, exists := d.ReferencesByName[name]; exists {
		first.Duplicates = append(first.Duplicates, ref)
	} else {
		d.ReferencesByName[name] = ref
	}
}

func (d *Document) String() string {
	// TODO: Output duplicates/missing refs.

	// TODO: Output refs in alphabetical order.

	lines := make([]string, 0, len(d.Referenced))
	for name, locs := range d.Referenced {
		strlocs := make([]string, len(locs))

		for i, loc := range locs {
			strlocs[i] = loc.String()
		}

		line := fmt.Sprintf("%s: %s", name, strings.Join(strlocs, ", "))
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}
