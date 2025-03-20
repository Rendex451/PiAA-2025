package aho_corasick

import "fmt"

var debug bool = false

func SetDebugFlag() {
	debug = true
}

func generateTrie(patterns []string) *Trie {
	trie := NewTrie()
	for _, pattern := range patterns {
		trie.addWord(pattern)
	}
	trie.generateLinks()
	if debug {
		fmt.Println("Trie generation completed for patterns:", patterns)
	}

	return trie
}

func findMatchesOnTrie(text string, trie *Trie, matches map[string][]int) {
	current := trie.root
	if debug {
		fmt.Printf("\n[Search] Processing text '%s':\n", text)
		fmt.Println("  Step-by-step traversal:")
	}
	for i, char := range text {
		if debug {
			fmt.Printf("  Pos %d ('%c'): ", i, char)
		}
		if next, exists := current.children[char]; exists {
			current = next
			if debug {
				fmt.Printf("Moved to '%c' (path: %s)\n",
					current.value, current.getPath())
			}
		} else {
			if debug {
				fmt.Printf("No direct transition, following suffix links:\n")
			}
			for current != trie.root && current.children[char] == nil {
				if debug {
					fmt.Printf("    From '%c' -> '%c'\n",
						current.value, current.suffixLink.value)
				}
				current = current.suffixLink
			}
			if next, exists := current.children[char]; exists {
				current = next
				if debug {
					fmt.Printf("    Found transition to '%c' (path: %s)\n",
						current.value, current.getPath())
				}
			} else if debug {
				fmt.Println("    Stayed at root")
			}
		}
		matchNode := current
		if matchNode.isEnd {
			pattern := matchNode.getPath()
			pos := i - len(pattern) + 1
			matches[pattern] = append(matches[pattern], pos)
			if debug {
				fmt.Printf("  -> Found match '%s' at pos %d\n", pattern, pos)
			}
		}
		for matchNode.terminalLink != nil && matchNode.terminalLink != trie.root {
			matchNode = matchNode.terminalLink
			pattern := matchNode.getPath()
			pos := i - len(pattern) + 1
			matches[pattern] = append(matches[pattern], pos)
			if debug {
				fmt.Printf("  -> Found via terminal '%s' at pos %d\n",
					pattern, pos)
			}
		}
	}
	if debug {
		fmt.Println("\n[Search] Final matches:", matches)
	}
}

func FindAllEntries(text string, patterns []string) map[string][]int {
	result := make(map[string][]int)
	trie := generateTrie(patterns)
	findMatchesOnTrie(text, trie, result)

	return result
}

func isValidWildcardMatch(text, pattern string, start int, wildcard, forbidden rune) bool {
	if forbidden == 0 {
		return true
	}

	for j, char := range pattern {
		if char == wildcard {
			textPos := start + j
			if textPos < len(text) && rune(text[textPos]) == forbidden {
				if debug {
					fmt.Printf("Invalid match at pos %d: forbidden char '%c' found\n", start, forbidden)
				}
				return false
			}
		}
	}

	return true
}

func FindEntriesWithWildcard(text string, pattern string, wildcard, forbidden rune) []int {
	parts := []string{}
	startPositions := []int{}
	current := ""
	start := 0

	if debug {
		fmt.Printf("\n[Wildcard Search] Pattern '%s' (wildcard: '%c', forbidden: '%c')\n",
			pattern, wildcard, forbidden)
	}
	for i, char := range pattern {
		if char == wildcard {
			if current != "" {
				parts = append(parts, current)
				startPositions = append(startPositions, start)
				current = ""
			}
		} else {
			if current == "" {
				start = i
			}
			current += string(char)
		}
	}
	if current != "" {
		parts = append(parts, current)
		startPositions = append(startPositions, start)
	}
	if debug {
		fmt.Println("  Parts:", parts)
		fmt.Println("  Start positions:", startPositions)
	}

	trie := generateTrie(parts)
	count := make([]int, len(text)+1)
	matches := make(map[string][]int)

	findMatchesOnTrie(text, trie, matches)

	if debug {
		fmt.Println("\n[Wildcard Search] Combining matches:")
	}
	for i, part := range parts {
		startPos := startPositions[i]
		for _, pos := range matches[part] {
			patternStart := pos - startPos
			if patternStart >= 0 && patternStart <= len(text)-len(pattern) {
				if isValidWildcardMatch(text, pattern, patternStart, wildcard, forbidden) {
					count[patternStart]++
					if debug {
						fmt.Printf("  Part '%s' matched at %d (count: %d)\n",
							part, patternStart, count[patternStart])
					}
				} else if debug {
					fmt.Printf("  Part '%s' at %d rejected (forbidden char)\n",
						part, patternStart)
				}
			}
		}
	}

	result := []int{}
	requiredMatches := len(parts)
	for i := 0; i <= len(text)-len(pattern); i++ {
		if count[i] == requiredMatches {
			result = append(result, i+1)
			if debug {
				fmt.Printf("  Complete match at position %d\n", i+1)
			}
		}
	}
	if debug {
		fmt.Println("[Wildcard Search] Final result:", result)
	}

	return result
}
