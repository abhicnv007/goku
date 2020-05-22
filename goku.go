package main

// Goku is the core data strcture of the library
type Goku struct {
	items map[string]string
}

// New creates a new instance of Goku
func New(p string) Goku {
	return Goku{items: make(map[string]string)}
}

// Add a key value pair to the Goku instance
func (g *Goku) Add(key string, value string) {
	g.items[key] = value
}

// Get the value saved for a key
func (g *Goku) Get(key string) (val string, ok bool) {
	val, ok = g.items[key]
	return
}

// Count returns the number of items saved
func (g *Goku) Count() int {
	return len(g.items)
}
