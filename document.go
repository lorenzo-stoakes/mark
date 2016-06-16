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

// TODO: Add statistics fields.

type Document struct {
	Path             string
	Lines            int
	References       []*Reference
	ReferencesByName map[string]*Reference
	ReferencedByName map[string][]Location
}

func newDocument(path string) *Document {
	return &Document{Path: path,
		ReferencesByName: make(map[string]*Reference),
		ReferencedByName: make(map[string][]Location)}
}

func (d *Document) referTo(name string, loc Location) {
	d.ReferencedByName[name] = append(d.ReferencedByName[name], loc)
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

func (d *Document) Duplicates() References {
	var ret []*Reference

	for _, ref := range d.References {
		if len(ref.Duplicates) > 0 {
			ret = append(ret, ref)
		}
	}

	return ret
}

func (d *Document) Missing() []string {
	var ret []string

	for ref, _ := range d.ReferencedByName {
		if _, has := d.ReferencesByName[ref]; !has {
			ret = append(ret, ref)
		}
	}

	return ret
}

func (d *Document) Unused() References {
	var ret []*Reference

	for _, ref := range d.References {
		if _, has := d.ReferencedByName[ref.Name]; !has {
			ret = append(ret, ref)
		}
	}

	return ret
}

func (d *Document) String() string {
	dupes := d.Duplicates()
	sort.Sort(dupes)
	missing := d.Missing()
	sort.Strings(missing)
	unused := d.Unused()
	sort.Sort(unused)

	if len(dupes)+len(missing)+len(unused) == 0 {
		return ""
	}

	lines := make([]string, 0, 4+len(dupes)+len(missing)+len(unused))

	add := func(str string, args ...interface{}) {
		lines = append(lines, fmt.Sprintf(str, args...))
	}

	add("%s:", d.Path)

	if len(dupes) > 0 {
		add("   %d duplicate reference(s):", len(dupes))
		for _, ref := range dupes {
			add("      %s", ref.Name)
		}
	}

	if len(missing) > 0 {
		add("   %d missing reference(s):", len(missing))
		for _, refName := range missing {
			add("      %s", refName)
		}
	}

	if len(unused) > 0 {
		add("   %d unused reference(s):", len(unused))
		for _, ref := range unused {
			add("      %s", ref.Name)
		}
	}

	return strings.Join(lines, "\n")
}

func (d *Document) IsMarkdownFile() bool {
	return IsMarkdownFile(d.Path)
}
