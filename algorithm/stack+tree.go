package main

import (
	"fmt"
	"log"
)

type Stack struct {
	node  []interface{}
	index int
}

func InitStack(n int) (Stack, bool) {
	if n < 1 {
		return Stack{}, false
	}
	return Stack{
		node:  make([]interface{}, n),
		index: -1,
	}, true
}
func (s *Stack) Push(n interface{}) bool {
	if s.IsFull() {
		return false
	}
	s.index += 1
	s.node[s.index] = n
	return true
}
func (s *Stack) Pop() (interface{}, bool) {
	if s.IsEmpty() {
		return -1, false
	}
	s.index -= 1
	return s.node[s.index+1], true
}

func (s *Stack) IsEmpty() bool {
	if s.index < 0 {
		return true
	}
	return false
}
func (s *Stack) IsFull() bool {
	if s.index+1 > len(s.node)-1 {
		return true
	}
	return false
}

type T struct {
	root *Tree
}
type Tree struct {
	key   int
	left  *Tree
	right *Tree
}

func (t T) Reserve() {
	stack, ok := InitStack(10)
	if !ok {
		log.Fatal("Can't initialize stack")
	}
	n := t.root
	for n != nil || !stack.IsEmpty() {
		for n != nil {
			fmt.Println(n.key)
			stack.Push(n)
			n = n.left
		}
		if !stack.IsEmpty() {
			i, ok := stack.Pop()
			if !ok {
				return
			}
			n = i.(*Tree).right
		}
	}
}
