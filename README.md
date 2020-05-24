# Goku

![Build Status](https://github.com/abhicnv007/goku/workflows/Test/badge.svg)
[![Go Report](https://goreportcard.com/badge/github.com/abhicnv007/goku)](https://goreportcard.com/badge/github.com/abhicnv007/goku)

Goku is a library written in Go to create a simple in-memory datastore that is persisted on disk.

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


## TODO

1. Test in multithreaded applications
2. Allow for dump and restore
3. Test file corruption
4. Lock the db file, preventing multiple instances from writing to the same file
