# mark

Mark is a markdown analysis tool. It lists references and detects duplicate and
missing references.

## Installation

Mark is an idiomatic go tool, so you need only run:

```
go get github.com/lorenzo-stoakes/mark/...
```

Making sure that your `$GOPATH/bin` is on your `PATH`.

## Usage

```
$ mark [markdown files...]
```

The tool outputs duplicates/missing references if they exist, otherwise it
outputs nothing.

If any of the specified files have duplicate or missing entries, or if an error
occurs, the tool exits with status code 1.
