# Goku

![Build Status](https://github.com/abhicnv007/goku/workflows/Test/badge.svg)
[![Go Report](https://goreportcard.com/badge/github.com/abhicnv007/goku)](https://goreportcard.com/badge/github.com/abhicnv007/goku)
[![codecov](https://codecov.io/gh/abhicnv007/goku/branch/master/graph/badge.svg)](https://codecov.io/gh/abhicnv007/goku)

Goku is a library written in Go to create a simple thread safe in-memory datastore that is persisted to disk.

## Usage

```go
package main

import (
    "fmt"
    "github.com/abhicnv007/goku"
)

func main() {
    g := goku.New(".db")
    defer g.Close()

    g.Add("foo", "bar")
    if val, ok := g.Get("foo"); ok {
        // use val
        fmt.Println(val) // Output: bar
    }
}
```

## How does it work?

Goku creates an append only file at the location and logs out all operations to this file. In the event of a crash,
the events from log are replayed and the data is reconstructed in memory.

```go
g := goku.New(".db")
g.Add("foo", "bar")
// many more adds

// program crash

g = goku.New(".db") // load the data back from the file
g.Get("foo") // bar
```

The `.db` file can be freely trnasferred.

## Benchmarks

```bash
goos: darwin
goarch: amd64
pkg: github.com/abhicnv007/goku
BenchmarkAdd-12    	  275427	     19766 ns/op	   10483 B/op	       2 allocs/op
BenchmarkGet-12    	44926000	       401 ns/op	       0 B/op	       0 allocs/op
```

## TODO

1. Allow for dump and restore
2. Test file corruption
