package main

import (
	"math/rand"
	"testing"
)

func TestAdd(t *testing.T) {
	g := New("")
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
	g := New("")
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

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func BenchmarkAdd(b *testing.B) {
	g := New("")
	// create some random strings to be used as keys and values
	randStrs := make([][]string, 0, b.N)
	for i := 0; i < b.N; i++ {
		randStrs = append(randStrs, []string{RandStringBytes(10), RandStringBytes(100)})
	}
	// reset timer and benchmark add
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.Add(randStrs[i][0], randStrs[i][1])
	}
}

func BenchmarkGet(b *testing.B) {
	g := New("")
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
