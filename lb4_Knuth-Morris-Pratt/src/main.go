package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"lb4_Knuth-Morris-Pratt/kmp"
)

func runPatternSearch(reader *bufio.Reader) {
	fmt.Println("Enter pattern:")
	pattern, _ := reader.ReadString('\n')

	fmt.Println("Enter text:")
	text, _ := reader.ReadString('\n')

	pattern = strings.TrimSpace(pattern)
	text = strings.TrimSpace(text)

	occurrences := kmp.FindPatternOccurrences(text, pattern, false)
	fmt.Println("Result:")
	printOccurrences(occurrences)
}

func runCyclicShiftCheck(reader *bufio.Reader) {
	fmt.Println("Enter text A:")
	textA, _ := reader.ReadString('\n')

	fmt.Println("Enter text B:")
	textB, _ := reader.ReadString('\n')

	textA = strings.TrimSpace(textA)
	textB = strings.TrimSpace(textB)

	res := kmp.IsCyclicShift(textA, textB)
	fmt.Println("Result:")
	fmt.Println(res)
}

func printOccurrences(occurrences []int) {
	if len(occurrences) == 0 {
		fmt.Println(-1)
	} else {
		for i, o := range occurrences {
			if i == len(occurrences)-1 {
				fmt.Printf("%d\n", o)
			} else {
				fmt.Printf("%d,", o)
			}
		}
	}
}

func main() {
	debugMode := flag.Bool("debug", false, "enable debug mode")
	cyclicShiftCheck := flag.Bool("cyclic", false, "check cyclic shift")
	flag.Parse()

	reader := bufio.NewReader(os.Stdin)

	if *debugMode {
		fmt.Println("Debug mode enabled")
		kmp.EnableDebug()
	}

	if *cyclicShiftCheck {
		runCyclicShiftCheck(reader)
	} else {
		runPatternSearch(reader)
	}
}
