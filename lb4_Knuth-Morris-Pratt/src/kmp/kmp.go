package kmp

var DEBUG = false

func EnableDebug() {
	DEBUG = true
}

func findPrefixFunction(pattern string) []int {
	pi := make([]int, len(pattern))
	j := 0

	for i := 1; i < len(pattern); i++ {
		for j > 0 && pattern[i] != pattern[j] {
			j = pi[j-1]
		}
		if pattern[i] == pattern[j] {
			j++
			pi[i] = j
		}
	}

	return pi
}

func FindPatternOccurrences(text, pattern string, firstOnly bool) []int {
	if len(pattern) > len(text) || len(pattern) == 0 {
		return []int{}
	}

	occurrences := []int{}
	pi := findPrefixFunction(pattern)

	for i, j := 0, 0; i < len(text); i++ {
		for j > 0 && text[i] != pattern[j] {
			j = pi[j-1]
		}
		if text[i] == pattern[j] {
			j++
		}
		if j == len(pattern) {
			occurrences = append(occurrences, i-j+1)
			if firstOnly {
				return occurrences
			}
			j = pi[j-1]
		}
	}

	return occurrences
}

func IsCyclicShift(text, pattern string) int {
	if len(pattern) != len(text) {
		return -1
	}

	text += text
	occurrences := FindPatternOccurrences(text, pattern, true)

	if len(occurrences) > 0 {
		return occurrences[0]
	}

	return -1
}
