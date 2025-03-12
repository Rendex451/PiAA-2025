package aho_corasick

import "fmt"

func generateTrie(patterns []string) *Trie {
	trie := NewTrie()
	for _, pattern := range patterns {
		trie.addWord(pattern)
	}
	trie.genSuffixLinks()
	if debug {
		fmt.Println("Trie generation completed for patterns:", patterns)
	}
	return trie
}

func findMatchesOnTrie(text string, trie *Trie, matches map[string][]int) {
	current := trie.root
	if debug {
		fmt.Printf("Starting search in text '%s'\n", text)
	}
	for i, char := range text {
		if next, exists := current.children[char]; exists {
			current = next
		} else {
			for current != trie.root && current.children[char] == nil {
				if debug {
					fmt.Printf("At pos %d ('%c'), following suffix link from '%c'\n", i, char, current.value)
				}
				current = current.suffixLink
			}
			if next, exists := current.children[char]; exists {
				current = next
			}
		}

		if debug {
			fmt.Printf("At pos %d ('%c'), current node: %s\n", i, char, current.String())
		}

		matchNode := current
		if matchNode.isEnd {
			pattern := matchNode.getPath()
			matches[pattern] = append(matches[pattern], i-len(pattern)+1)
			if debug {
				fmt.Printf("Found match '%s' at position %d\n", pattern, i-len(pattern)+1)
			}
		}
		for matchNode.terminalLink != nil {
			matchNode = matchNode.terminalLink
			pattern := matchNode.getPath()
			matches[pattern] = append(matches[pattern], i-len(pattern)+1)
			if debug {
				fmt.Printf("Found match via terminal link '%s' at position %d\n", pattern, i-len(pattern)+1)
			}
		}
	}
	if debug {
		fmt.Println("Search completed, matches:", matches)
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
		fmt.Printf("Processing pattern '%s' with wildcard '%c' and forbidden '%c'\n", pattern, wildcard, forbidden)
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
		fmt.Println("Pattern split into parts:", parts, "with start positions:", startPositions)
	}

	trie := generateTrie(parts)

	count := make([]int, len(text)+1)
	matches := make(map[string][]int)

	findMatchesOnTrie(text, trie, matches)

	for i, part := range parts {
		startPos := startPositions[i]
		for _, pos := range matches[part] {
			patternStart := pos - startPos
			if patternStart >= 0 && patternStart <= len(text)-len(pattern) {
				if isValidWildcardMatch(text, pattern, patternStart, wildcard, forbidden) {
					count[patternStart]++
					if debug {
						fmt.Printf("Valid partial match for '%s' at pos %d, count now %d\n", part, patternStart, count[patternStart])
					}
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
				fmt.Printf("Full match found at position %d\n", i+1)
			}
		}
	}

	if debug {
		fmt.Println("Wildcard search completed, result:", result)
	}
	return result
}
