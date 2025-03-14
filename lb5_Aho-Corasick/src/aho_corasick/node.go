package aho_corasick

import (
	"fmt"
	"slices"
)

type Node struct {
	value        rune
	isEnd        bool
	children     map[rune]*Node
	suffixLink   *Node
	terminalLink *Node
	parent       *Node
}

func NewNode(value rune, parent *Node) *Node {
	return &Node{
		value:        value,
		isEnd:        false,
		children:     make(map[rune]*Node),
		suffixLink:   nil,
		terminalLink: nil,
		parent:       parent,
	}
}

func (n *Node) addChild(value rune, node *Node) {
	n.children[value] = node
	if debug {
		fmt.Printf("Added child '%c' to node '%c'\n", value, n.value)
	}
}

func (n *Node) setEnd() {
	n.isEnd = true
	if debug {
		fmt.Printf("Marked node '%c' as end of pattern\n", n.value)
	}
}

func (n *Node) getPath() string {
	current := n
	path := []rune{}

	for current.value != 0 {
		path = append(path, current.value)
		current = current.parent
	}

	slices.Reverse(path)

	return string(path)
}

func (n *Node) String() string {
	suffixVal := '0'
	terminalVal := '0'
	if n.suffixLink != nil {
		suffixVal = n.suffixLink.value
	}
	if n.terminalLink != nil {
		terminalVal = n.terminalLink.value
	}

	return fmt.Sprintf("Node('%c', isEnd=%t, suffixLink='%c', terminalLink='%c', children=%v)",
		n.value, n.isEnd, suffixVal, terminalVal, len(n.children))
}
