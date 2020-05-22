package main

import "testing"

func TestAdd(t *testing.T) {
	g := New()
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
	g := New()
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
