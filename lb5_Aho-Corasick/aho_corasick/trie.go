package aho_corasick

import "fmt"

var debug bool = false

func SetDebugFlag() {
	debug = true
}

type Trie struct {
	root *Node
}

func NewTrie() *Trie {
	root := NewNode(0, nil)
	root.parent = root
	root.suffixLink = root
	if debug {
		fmt.Println("Created root node:", root.String())
	}
	return &Trie{root: root}
}

func (t *Trie) addWord(word string) {
	current := t.root
	if debug {
		fmt.Printf("Adding word '%s' to trie\n", word)
	}
	for _, char := range word {
		if _, exists := current.children[char]; !exists {
			newNode := NewNode(char, current)
			current.addChild(char, newNode)
		}
		current = current.children[char]
	}
	current.setEnd()
	if debug {
		fmt.Printf("Finished adding '%s', current node: %s\n", word, current.String())
	}
}

func (t *Trie) genSuffixLinks() {
	queue := []*Node{t.root}
	if debug {
		fmt.Println("Generating suffix links...")
	}
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]

		for _, child := range node.children {
			queue = append(queue, child)
		}

		parent := node.parent
		for parent != parent.suffixLink {
			if child, exists := parent.suffixLink.children[node.value]; exists {
				node.suffixLink = child
				break
			}
			parent = parent.suffixLink
		}
		if node.suffixLink == nil {
			node.suffixLink = t.root
		}

		if node.suffixLink.isEnd {
			node.terminalLink = node.suffixLink
		} else {
			node.terminalLink = node.suffixLink.terminalLink
		}

		if debug {
			fmt.Printf("Processed node: %s\n", node.String())
		}
	}
	if debug {
		fmt.Println("Suffix links generation completed")
	}
}
