# mark

Mark is a tool which detects duplicate and missing references in markdown files.

## Installation

```
go get github.com/lorenzo-stoakes/mark/...
```

## Usage

```
$ mark [markdown files...]
```

The tool outputs duplicates/missing references if they exist, otherwise it
outputs nothing.

If any of the specified files have duplicate or missing entries, or if an error
occurs, the tool exits with status code 1.
