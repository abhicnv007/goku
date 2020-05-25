package goku

import (
	"math/rand"
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	dbPath := ".goku"
	g := New(dbPath)

	g.Add("Q", "there")
	g.Add("W", "here")
	g.Add("E", "everywhere")

	if got := g.Count(); got != 3 {
		t.Errorf("Goku New; expected 3, got %d", got)
	}

	// file is supposed to exist
	_, err := os.Stat(dbPath)
	if os.IsNotExist(err) {
		t.Errorf("Goku New, datafile not present")
	}

	g.Close()

	// file is still supposed to exist
	f, err := os.Stat(dbPath)
	t.Logf("The file is %d bytes long", f.Size())

	if os.IsNotExist(err) {
		t.Errorf("Goku New, datafile not present after close")
	}

	// open with the same datafile again
	g = New(dbPath)
	if got := g.Count(); got != 3 {
		t.Errorf("Goku New; expected 3 after reload, got %d", got)
	}

	// g.Clear()
	g.Close()
}

func TestAdd(t *testing.T) {
	dbPath := ".goku"
	g := New(dbPath)
	defer g.Close()
	defer g.Clear()

	g.Add("Q", "there")
	g.Add("W", "here")
	g.Add("E", "everywhere")
	if got := g.Count(); got != 3 {
		t.Errorf("Goku Add = %d; want 3", got)
	}

	g.Add("W", "again")
	if got := g.Count(); got != 3 {
		t.Errorf("Goku Add = %d; want 3", got)
	}
}

func TestGet(t *testing.T) {
	dbPath := ".goku"
	g := New(dbPath)
	defer g.Close()
	defer g.Clear()

	keys := []string{"A", "B", "C", "D"}
	values := []string{"1", "2", "3", "4"}

	for i, k := range keys {
		g.Add(k, values[i])
	}

	for i, k := range keys {
		v, ok := g.Get(k)
		if !ok || v != values[i] {
			t.Errorf("Goku Get = %s, want %s, got %s", k, values[i], v)
		}
	}
}

func TestClear(t *testing.T) {
	dbPath := ".goku"
	g := New(dbPath)
	defer g.Close()

	keys := []string{"A", "B", "C", "D"}
	values := []string{"1", "2", "3", "4"}

	for i, k := range keys {
		g.Add(k, values[i])
	}

	g.Clear()

	if c := g.Count(); c != 0 {
		t.Errorf("Goku clear, expected 0 elements after clear, got %d", c)
	}

	keys = []string{"A", "B", "C", "D"}
	values = []string{"1", "2", "3", "4"}

	for i, k := range keys {
		g.Add(k, values[i])
	}

	if c := g.Count(); c != 4 {
		t.Errorf("Goku clear, expected 4 elements after insertion, got %d", c)
	}
	g.Clear()

}

func TestGokuClose(t *testing.T) {
	dbPath := ".goku"
	g := New(dbPath)
	keys := []string{"A", "B", "C", "D"}
	values := []string{"1", "2", "3", "4"}

	for i, k := range keys {
		g.Add(k, values[i])
	}

	g.Clear()
	g.Close()

	// datafile not supposed to exist
	_, err := os.Stat(dbPath)
	if !os.IsNotExist(err) {
		t.Errorf("Goku clear, datafile not cleared")
	}

}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func BenchmarkAdd(b *testing.B) {
	dbPath := ".goku"
	g := New(dbPath)
	defer g.Close()
	defer g.Clear()

	// create some random strings to be used as keys and values
	randStrs := make([][]string, 0, b.N)
	for i := 0; i < b.N; i++ {
		randStrs = append(randStrs, []string{RandStringBytes(10), RandStringBytes(10_000)})
	}
	// reset timer and benchmark add
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.Add(randStrs[i][0], randStrs[i][1])
	}

}

func BenchmarkGet(b *testing.B) {
	dbPath := ".goku"
	g := New(dbPath)
	defer g.Close()
	defer g.Clear()
	// create some random strings to be used as keys and values
	randStrs := make([][]string, 0, b.N)
	for i := 0; i < b.N; i++ {
		k, v := RandStringBytes(10), RandStringBytes(100)
		randStrs = append(randStrs, []string{k, v})
		g.Add(k, v)
	}
	// reset timer and benchmark add
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.Get(randStrs[i][0])
	}
}
