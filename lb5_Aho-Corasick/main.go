package main

import (
	"flag"
	"fmt"
	"sort"

	"lb5_Aho-Corasick/aho_corasick"
)

func basicStart() {
	var text string
	var n int

	fmt.Print("Text: ")
	fmt.Scan(&text)
	fmt.Print("Number of patterns: ")
	fmt.Scan(&n)
	fmt.Print("Patterns: \n")

	patterns := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Scan(&patterns[i])
	}

	result := aho_corasick.FindAllEntries(text, patterns)
	var ans [][2]int

	for i, pattern := range patterns {
		if len(result[pattern]) > 0 {
			for _, pos := range result[pattern] {
				ans = append(ans, [2]int{pos + 1, i + 1})
			}
		}
	}

	sort.Slice(ans, func(i, j int) bool {
		if ans[i][0] == ans[j][0] {
			return ans[i][1] < ans[j][1]
		}
		return ans[i][0] < ans[j][0]
	})

	fmt.Print("\nResult: \n")
	for _, pair := range ans {
		fmt.Printf("%d %d\n", pair[0], pair[1])
	}
}

func wildcardStart(withForbidden bool) {
	var text, pattern, wildcardStr string

	fmt.Print("Text: ")
	fmt.Scan(&text)
	fmt.Print("Pattern: ")
	fmt.Scan(&pattern)
	fmt.Print("Wildcard: ")
	fmt.Scan(&wildcardStr)
	wildcard := rune(wildcardStr[0])

	var forbidden rune
	if withForbidden {
		var forbiddenStr string
		fmt.Print("Forbidden: ")
		fmt.Scan(&forbiddenStr)

		forbidden = rune(forbiddenStr[0])
	} else {
		forbidden = 0
	}

	result := aho_corasick.FindEntriesWithWildcard(text, pattern, wildcard, forbidden)

	fmt.Print("\nResult: \n")
	sort.Ints(result)
	for _, pos := range result {
		fmt.Println(pos)
	}
}

func main() {
	debugMode := flag.Bool("debug", false, "Enable debug mode.")
	withWildcard := flag.Bool("wildcard", false, "Find entries of pattern with wildcard.")
	withForbidden := flag.Bool("forbidden", false, "Find entries of pattern with wildcard that is not forbidden.")
	flag.Parse()

	if *debugMode {
		fmt.Println("Debug mode enabled")
		aho_corasick.SetDebugFlag()
	}

	if *withForbidden {
		wildcardStart(true)
	} else if *withWildcard {
		wildcardStart(false)
	} else {
		basicStart()
	}
}
