package mark

import (
	"fmt"
	"sort"
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

type References []*Reference

func (rs References) Len() int {
	return len(rs)
}

func (rs References) Less(i, j int) bool {
	return rs[i].Name < rs[j].Name
}

func (rs References) Swap(i, j int) {
	rs[i], rs[j] = rs[j], rs[i]
}

type Document struct {
	Path             string
	Lines            int
	References       []*Reference
	ReferencesByName map[string]*Reference
	Referenced       map[string][]Location
}

func newDocument(path string) *Document {
	return &Document{Path: path, Referenced: make(map[string][]Location),
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

func (d *Document) Duplicates() []*Reference {
	var ret []*Reference

	for _, ref := range d.References {
		if len(ref.Duplicates) > 0 {
			ret = append(ret, ref)
		}
	}

	return ret
}

func (d *Document) MissingDefines() []string {
	var ret []string

	for ref, _ := range d.Referenced {
		if _, has := d.ReferencesByName[ref]; !has {
			ret = append(ret, ref)
		}
	}

	return ret
}

func (d *Document) String() string {
	dupes := References(d.Duplicates())
	sort.Sort(dupes)
	missing := d.MissingDefines()
	sort.Strings(missing)

	lines := make([]string, 0, 3 + len(dupes) + len(missing))

	add := func(str string, args ...interface{}) {
		lines = append(lines, fmt.Sprintf(str, args...))
	}

	add(d.Path)

	add("  %d duplicate(s)", len(dupes))
	for _, dupe := range dupes {
		add("    %s", dupe.Name)
	}

	add("  %d missing reference(s)", len(missing))
	for _, def := range missing {
		add("    %s", def)
	}

	return strings.Join(lines, "\n")
}
