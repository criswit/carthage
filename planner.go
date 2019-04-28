package main

import (
	"bytes"
	"fmt"
)

type Action func() error

type Unit struct {
	rdeps  nodeSet
	deps   nodeSet
	Name   string
	Action Action
}

func (u *Unit) AddDependency(dep *Unit) {
	u.deps = u.deps.add(dep)
	dep.rdeps = dep.rdeps.add(u)
}

type nodeSet []*Unit

func (ns nodeSet) add(unit *Unit) nodeSet {
	for _, n := range ns {
		if n == unit {
			return ns
		}
	}
	return append(ns, unit)
}

type Cursor struct {
	Node   *Unit
	Parent *Unit

	Depth       int
	LastAtDepth bool
}

type WalkFunc func(*Cursor) bool

func Walk(unit *Unit, visitor WalkFunc) {
	Visit(unit, Visitor{
		Pre:  visitor,
		Post: func(*Cursor) bool { return true },
	})
}

func Visit(unit *Unit, visitor Visitor) {
	visit(&Cursor{
		Node:        unit,
		Depth:       0,
		LastAtDepth: true,
	}, visitor)
}

type Visitor struct {
	Pre  WalkFunc
	Post WalkFunc
}

func visit(cursor *Cursor, visitor Visitor) {
	if !visitor.Pre(cursor) {
		return
	}

	nextDepth := cursor.Depth + 1
	for i := 0; i < len(cursor.Node.deps); i++ {
		next := cursor.Node.deps[i]
		last := i+1 == len(cursor.Node.deps)
		visit(&Cursor{
			Node:        next,
			Parent:      cursor.Node,
			Depth:       nextDepth,
			LastAtDepth: last,
		}, visitor)
	}
	visitor.Post(cursor)
}

func NewExecutionPlan(root *Unit) (*ExecutionPlan, error) {
	var err error
	var units []*Unit

	inProgress := map[*Unit]struct{}{}
	done := map[*Unit]struct{}{}

	Visit(root, Visitor{
		Pre: func(cursor *Cursor) bool {
			if _, ok := done[cursor.Node]; err != nil || ok {
				return false
			}
			if _, ok := inProgress[cursor.Node]; ok {
				err = fmt.Errorf("detected cycle")
				return false
			}
			inProgress[cursor.Node] = struct{}{}
			return true
		},
		Post: func(cursor *Cursor) bool {
			fmt.Printf("%v %v\n", cursor.Node.Name, done)
			if _, ok := done[cursor.Node]; ok {
				return false
			}
			delete(inProgress, cursor.Node)
			done[cursor.Node] = struct{}{}
			units = append(units, cursor.Node)
			return true
		},
	})
	if err != nil {
		return nil, err
	}
	return &ExecutionPlan{
		root:  root,
		units: units,
	}, nil
}

type ExecutionPlan struct {
	root  *Unit
	units []*Unit
}

func (p *ExecutionPlan) String() {
	const (
		A = "| "
		B = "├─"
		C = "└─"
		D = "  "
	)

	var buf bytes.Buffer

	prevLastAtDepth := false
	prevDepth := 0
	prefix := ""

	Walk(p.root, func(cur *Cursor) bool {
		if prevDepth < cur.Depth {
			if prevLastAtDepth {
				prefix = prefix + D
			} else {
				prefix = prefix + A
			}
		} else {
			backTrackCount := (prevDepth - cur.Depth) * 2
			prefix = prefix[:len(prefix)-backTrackCount]
		}

		buf.WriteString(prefix)
		if cur.LastAtDepth {
			buf.WriteString(C)
		} else {
			buf.WriteString(B)
		}
		buf.WriteString(cur.Node.Name)
		buf.WriteString(fmt.Sprintf("[%p]", cur.Node))
		buf.WriteRune('\n')

		prevLastAtDepth = cur.LastAtDepth
		prevDepth = cur.Depth
		return true
	})

	fmt.Print(buf.String())
}
