package main

import (
	"encoding/json"
	"fmt"
)

type ValTree struct {
	L      *ValTree
	R      *ValTree
	V      int
	RDepth int
	LDepth int
}

func main() {
	tree := &ValTree{
		L: &ValTree{V: 5},
		V: 6,
		R: &ValTree{
			V: 8,
			L: &ValTree{V: 7},
			R: &ValTree{V: 10},
		},
	}
	b1, _ := json.Marshal(tree)
	fmt.Println(string(b1))
	FindR(tree, tree)
	FindL(tree, tree)
	b2, _ := json.Marshal(tree)
	fmt.Println(string(b2))

}

func FindR(t *ValTree, node *ValTree) {
	if t.R != nil || t.L != nil {
		node.RDepth += 1
		FindR(t.R, node)
	}
	return
}
func FindL(t *ValTree, node *ValTree) {
	if t.R != nil || t.L != nil {
		node.LDepth += 1
		FindL(t.R, node)
	}
	return
}
