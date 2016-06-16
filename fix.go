package mark

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	pathLib "path"
)

func (d *Document) fix() (err error) {
	// TODO: Remove duplication between this and ParseFile().
	// TODO: Reduce function size.

	// Since we're now modifying the file, we should be cautious.
	if !d.IsMarkdownFile() {
		return fmt.Errorf("Cowardly refusing to fix possibly non-markdown '%s'.",
			d.Path)
	}

	var file, tmpFile *os.File

	// Reopen path so we can filter out 'bad' lines.
	if file, err = os.Open(d.Path); err != nil {
		return
	}
	defer file.Close()

	// Create a temporary output file, we'll rename it over the target path
	// when we're done.
	dir := pathLib.Dir(d.Path)
	if tmpFile, err = ioutil.TempFile(dir, "mark"); err != nil {
		return
	}
	defer os.Remove(tmpFile.Name())

	// Get a hash of defunct lines we need to filter out.
	hash := d.defunctLineHash()

	line := 1
	reader := bufio.NewReader(file)
	var bytes []byte
	for bytes, err = reader.ReadBytes('\n'); err != io.EOF; bytes, err = reader.ReadBytes('\n') {
		if err != nil {
			return
		}

		// Simply output non-defunct lines.
		if !hash[line] {
			tmpFile.Write(bytes)
		}

		line++
	}
	// Clear EOF.
	err = nil

	return os.Rename(tmpFile.Name(), d.Path)
}
