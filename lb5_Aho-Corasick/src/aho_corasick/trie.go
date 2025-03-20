package aho_corasick

import "fmt"

type Trie struct {
	root *Node
}

func NewTrie() *Trie {
	root := NewNode(0, nil)
	root.parent = root
	root.suffixLink = root
	if debug {
		fmt.Println("\nCreated root node:", root.String())
	}
	return &Trie{root: root}
}

func (t *Trie) addWord(word string) {
	current := t.root
	if debug {
		fmt.Printf("\n[Trie] Adding pattern '%s':\n", word)
	}
	for i, char := range word {
		if _, exists := current.children[char]; !exists {
			newNode := NewNode(char, current)
			current.addChild(char, newNode)
			if debug {
				fmt.Printf("  Step %d: Added '%c' to '%c'\n", i+1, char,
					current.value)
			}
		}
		current = current.children[char]
		if debug {
			fmt.Printf("  Step %d: Moved to '%c' (path: %s)\n", i+1,
				char, current.getPath())
		}
	}
	current.setEnd()
	if debug {
		fmt.Printf("  Completed: Marked '%s' as pattern end\n", current.getPath())
	}
}

func (t *Trie) generateLinks() {
	queue := []*Node{t.root}
	t.root.suffixLink = t.root

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		for _, child := range node.children {
			queue = append(queue, child)
		}
		if node == t.root {
			continue
		}

		// Set suffix link
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

		// Set terminal link
		if node.suffixLink.isEnd {
			node.terminalLink = node.suffixLink
			if debug {
				fmt.Printf("    Set terminal link to '%s' (direct)\n", node.terminalLink.getPath())
			}
		} else {
			node.terminalLink = node.suffixLink.terminalLink
			if debug {
				if node.terminalLink != nil {
					fmt.Printf("    Set terminal link to '%s' (inherited)\n", node.terminalLink.getPath())
				} else {
					fmt.Println("    Terminal link is nil")
				}
			}
		}

		if debug {
			fmt.Printf("  Processing '%c' (path: %s):\n", node.value, node.getPath())
			fmt.Printf("    Set suffix link to '%s'\n", node.suffixLink.getPath())
		}
	}
}
