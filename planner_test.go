package main

import (
	"testing"
)

func TestPlanner(t *testing.T) {
	root := &Unit{
		Name: "root",
	}
	a := &Unit{
		Name: "a",
	}
	b := &Unit{
		Name: "b",
	}
	root.AddDependency(a)
	root.AddDependency(b)
	a.AddDependency(b)

	Print(root)
}
