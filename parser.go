package mark

import (
	"bufio"
	"io"
	"os"
)

type parseState struct {
	*Document
	line                          string
	index, lineNum                int
	start, startName, startUri    int
	endDefName                    int
	canDefine                     bool
	opened, referencing, defining bool
	code, multiCode               bool
}

func (d *Document) newParseState(lineNum int, line string, multi bool) *parseState {
	// Reference definitions are only (possibly) allowed if they sit alone
	// on a line.
	return &parseState{lineNum: lineNum, line: line, Document: d,
		canDefine: line[0] == '[', multiCode: multi}
}

func (s *parseState) chr() uint8 {
	return s.line[s.index]
}

func (s *parseState) rest() string {
	return s.line[s.index:]
}

func (s *parseState) eol() bool {
	return s.index >= len(s.line)
}

func (s *parseState) isWhitespace() bool {
	switch s.chr() {
	case ' ', '\t', '\n':
		return true

	default:
		return false
	}
}

func (s *parseState) skipWhitespace() bool {
	for !s.eol() && s.isWhitespace() {
		s.index++
	}

	// Was there anything left after the whitespace?
	return !s.eol()
}

func (s *parseState) startDefine() {
	if !s.canDefine {
		return
	}

	s.defining = true
	s.endDefName = s.index - 1

	// Skip ':'.
	s.index++
	if s.skipWhitespace() {
		s.startUri = s.index
	} else {
		// Abort, nothing after the whitespace.
		s.defining = false
	}
}

func (s *parseState) finishDefine() {
	// Has to be on its own line starting with '['.
	name := s.line[1:s.endDefName]
	// Skip newline (we have to be at the end of the line to finish a define.
	uri := s.line[s.startUri : len(s.line)-1]

	s.Document.define(name, uri, Location{s.lineNum, 1})
}

func (s *parseState) startReference() {
	// Even [foo]:[bar][baz] is theoretically valid.
	if s.defining {
		return
	}

	s.referencing = true
	// Re-open, since we're moving into 'refname' in '[text][refname]'.
	s.opened = true
	// Skip '['.
	s.startName = s.index + 1
}

func (s *parseState) finishReference() {
	s.referencing = false

	name := s.line[s.startName:s.index]
	// Columns are base-1 in Location.
	s.Document.referTo(name, Location{s.lineNum, s.start + 1})
}

func (s *parseState) codeBlock() {
	// Note that s.code => !s.multiCode.
	if s.code {
		// We've closed out our code section.
		s.code = false
		return
	}

	rest := s.rest()
	if len(rest) >= 3 && rest[:3] == "```" {
		// Skip remaining 2 quote chars.
		s.index += 2

		s.multiCode = !s.multiCode
		return
	}

	if !s.multiCode {
		// Otherwise we just opened a single quote code block.
		s.code = true
	}
}

func (s *parseState) close() {
	if !s.opened || s.code || s.multiCode {
		return
	}
	s.opened = false

	if s.referencing {
		s.finishReference()
		return
	}

	// If we're at the end of the line, we can't peek.
	if s.index == len(s.line)-1 {
		return
	}

	// Peek.
	s.index++
	switch s.chr() {
	case '[':
		// If we do anything but start defining when it's possible,
		// defining is no longer allowed.
		s.canDefine = false
		s.startReference()
	case ':':
		s.startDefine()
	default:
		s.canDefine = false
	}
}

func (s *parseState) open() {
	if s.opened || s.code || s.multiCode {
		return
	}

	s.start = s.index
	s.opened = true
}

func (s *parseState) Next() {
	switch s.chr() {
	case '`':
		s.codeBlock()
	case ']':
		s.close()
	case '[':
		s.open()
	}

	s.index++
}

func (s *parseState) Parse() {
	// newline is included, and we want to peek 1 char ahead so only look at
	// len(line)-1 chars.
	for s.index < len(s.line)-1 {
		s.Next()
	}

	// Define ends at end of line (ignoring newline.)
	if s.defining {
		s.finishDefine()
	}
}

func (d *Document) parseLine(lineNum int, line string, multi bool) *parseState {
	parseState := d.newParseState(lineNum, line, multi)
	parseState.Parse()

	return parseState
}

func ParseFile(path string) (ret *Document, err error) {
	var file *os.File

	if file, err = os.Open(path); err != nil {
		return
	}

	defer file.Close()
	reader := bufio.NewReader(file)

	multi := false
	lineNum := 1
	ret = newDocument(path)
	var bytes []byte
	for bytes, err = reader.ReadBytes('\n'); err != io.EOF; bytes, err = reader.ReadBytes('\n') {
		if err != nil {
			return
		}

		// We intentionally include the newline.
		state := ret.parseLine(lineNum, string(bytes), multi)
		multi = state.multiCode
		lineNum++
	}
	// Clear EOF.
	err = nil

	ret.Lines = lineNum - 1

	return
}
