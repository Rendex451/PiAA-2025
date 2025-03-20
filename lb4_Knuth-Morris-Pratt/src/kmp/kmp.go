package kmp

import "fmt"

var DEBUG = false

const (
	COLOR_RESET  = "\033[0m"
	COLOR_GREEN  = "\033[32m"
	COLOR_BLUE   = "\033[34m"
	COLOR_YELLOW = "\033[33m"
)

func EnableDebug() {
	DEBUG = true
}

func findPrefixFunction(pattern string) []int {
	pi := make([]int, len(pattern))
	j := 0

	if DEBUG {
		fmt.Printf("%s[PrefixFunc]%s Computing prefix function for pattern: %s\n", COLOR_YELLOW, COLOR_RESET, pattern)
		fmt.Printf("%s[PrefixFunc]%s Initialized pi: %v\n", COLOR_YELLOW, COLOR_RESET, pi)
	}

	for i := 1; i < len(pattern); i++ {
		if DEBUG {
			fmt.Printf("%s[PrefixFunc]%s i=%d, character: %c, j=%d\n", COLOR_YELLOW, COLOR_RESET, i, pattern[i], j)
		}

		for j > 0 && pattern[i] != pattern[j] {
			if DEBUG {
				fmt.Printf("%s[PrefixFunc]%s Mismatch: %c != %c, rolling back j from %d to %d\n",
					COLOR_YELLOW, COLOR_RESET, pattern[i], pattern[j], j, pi[j-1])
			}
			j = pi[j-1]
		}

		if pattern[i] == pattern[j] {
			j++
			pi[i] = j
			if DEBUG {
				fmt.Printf("%s[PrefixFunc]%s Match: %c = %c, j increased to %d, pi[%d] = %d\n",
					COLOR_YELLOW, COLOR_RESET, pattern[i], pattern[j-1], j, i, pi[i])
			}
		}
	}

	if DEBUG {
		fmt.Printf("%s[PrefixFunc]%s Final prefix function: %v\n", COLOR_YELLOW, COLOR_RESET, pi)
	}

	return pi
}

func FindPatternOccurrences(text, pattern string, firstOnly bool) []int {
	if len(pattern) > len(text) || len(pattern) == 0 {
		if DEBUG {
			fmt.Printf("%s[Search]%s Early return: pattern longer than text (%d > %d) or empty\n",
				COLOR_GREEN, COLOR_RESET, len(pattern), len(text))
		}
		return []int{}
	}

	occurrences := []int{}
	pi := findPrefixFunction(pattern)

	if DEBUG {
		fmt.Printf("%s[Search]%s Searching for occurrences of pattern '%s' in text '%s'\n",
			COLOR_GREEN, COLOR_RESET, pattern, text)
		fmt.Printf("%s[Search]%s Using prefix function: %v\n", COLOR_GREEN, COLOR_RESET, pi)
	}

	for i, j := 0, 0; i < len(text); i++ {
		if DEBUG {
			fmt.Printf("%s[Search]%s i=%d, text character: %c, j=%d\n", COLOR_GREEN, COLOR_RESET, i, text[i], j)
		}

		for j > 0 && text[i] != pattern[j] {
			if DEBUG {
				fmt.Printf("%s[Search]%s Mismatch: %c != %c, rolling back j from %d to %d\n",
					COLOR_GREEN, COLOR_RESET, text[i], pattern[j], j, pi[j-1])
			}
			j = pi[j-1]
		}

		if text[i] == pattern[j] {
			j++
			if DEBUG {
				fmt.Printf("%s[Search]%s Match: %c = %c, j increased to %d\n",
					COLOR_GREEN, COLOR_RESET, text[i], pattern[j-1], j)
			}
		}

		if j == len(pattern) {
			occ := i - j + 1
			occurrences = append(occurrences, occ)
			if DEBUG {
				fmt.Printf("%s[Search]%s Found occurrence at position %d\n", COLOR_GREEN, COLOR_RESET, occ)
			}
			if firstOnly {
				if DEBUG {
					fmt.Printf("%s[Search]%s Returning first occurrence: %v\n", COLOR_GREEN, COLOR_RESET, occurrences)
				}
				return occurrences
			}
			j = pi[j-1]
			if DEBUG {
				fmt.Printf("%s[Search]%s Continuing search, rolling back j to %d\n", COLOR_GREEN, COLOR_RESET, j)
			}
		}
	}

	if DEBUG {
		fmt.Printf("%s[Search]%s All occurrences: %v\n", COLOR_GREEN, COLOR_RESET, occurrences)
	}
	return occurrences
}

func IsCyclicShift(text, pattern string) int {
	if len(pattern) != len(text) {
		if DEBUG {
			fmt.Printf("%s[CyclicCheck]%s Lengths don't match: text(%d) != pattern(%d)\n",
				COLOR_BLUE, COLOR_RESET, len(text), len(pattern))
		}
		return -1
	}

	text += text
	if DEBUG {
		fmt.Printf("%s[CyclicCheck]%s Checking cyclic shift\n", COLOR_BLUE, COLOR_RESET)
		fmt.Printf("%s[CyclicCheck]%s Doubled text: %s\n", COLOR_BLUE, COLOR_RESET, text)
		fmt.Printf("%s[CyclicCheck]%s pattern: %s\n", COLOR_BLUE, COLOR_RESET, pattern)
	}

	occurrences := FindPatternOccurrences(text, pattern, true)

	if len(occurrences) > 0 {
		if DEBUG {
			fmt.Printf("%s[CyclicCheck]%s Found cyclic shift at position %d\n", COLOR_BLUE, COLOR_RESET, occurrences[0])
		}
		return occurrences[0]
	}

	if DEBUG {
		fmt.Printf("%s[CyclicCheck]%s Cyclic shift not found\n", COLOR_BLUE, COLOR_RESET)
	}
	return -1
}
