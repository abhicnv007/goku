# Goku

![Build Status](https://github.com/abhicnv007/goku/workflows/Test/badge.svg)
[![Go Report](https://goreportcard.com/badge/github.com/abhicnv007/goku)](https://goreportcard.com/badge/github.com/abhicnv007/goku)

Goku is a library written in Go to create a simple in-memory datastore that is persisted on disk.

## Usage

```go
package main

import "github.com/abhicnv007/goku"

func main() {
    g := goku.New(".goku_data")
    g.Add("Key", "Some Value")
    val := g.Get("Key") // val = "Some Value"
    g.Close()
}
```

## TODO

1. Test in multithreaded applications
2. Allow for dump and restore
3. Test file corruption
