package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"lb3_Levenshtein/logger"
	"lb3_Levenshtein/vagner_fisher"
)

func readInputStrings(reader *bufio.Reader, writer *bufio.Writer) (string, string) {
	fmt.Fprint(writer, "Enter the first string: ")
	writer.Flush()
	s1, _ := reader.ReadString('\n')
	s1 = strings.TrimSpace(s1)

	fmt.Fprint(writer, "Enter the second string: ")
	writer.Flush()
	s2, _ := reader.ReadString('\n')
	s2 = strings.TrimSpace(s2)

	return s1, s2
}

func readInputConfig(reader *bufio.Reader, writer *bufio.Writer) ([]string, []string, []string) {
	fmt.Fprint(writer, "Enter the costs (replace, insert, delete): ")
	writer.Flush()
	costsInput, _ := reader.ReadString('\n')
	costs := strings.Split(strings.TrimSpace(costsInput), " ")
	if len(costs) != 3 {
		fmt.Fprintln(os.Stderr, "Invalid input. Please enter 3 costs separated by spaces.")
		os.Exit(1)
	}

	fmt.Fprint(writer, "Enter special runes (replace, insert): ")
	writer.Flush()
	specialRunesInput, _ := reader.ReadString('\n')
	specialRunesStrs := strings.Split(strings.TrimSpace(specialRunesInput), " ")
	if len(specialRunesStrs) != 2 {
		fmt.Fprintln(os.Stderr, "Invalid input. Please enter 2 special runes separated by spaces.")
		os.Exit(1)
	}

	fmt.Fprint(writer, "Enter special runes costs (replace, insert): ")
	writer.Flush()
	specialRunesCostsInput, _ := reader.ReadString('\n')
	specialRunesCosts := strings.Split(strings.TrimSpace(specialRunesCostsInput), " ")
	if len(specialRunesCosts) != 2 {
		fmt.Fprintln(os.Stderr, "Invalid input. Please enter 2 special runes costs separated by spaces.")
		os.Exit(1)
	}

	return costs, specialRunesStrs, specialRunesCosts
}

func main() {
	debugMode := flag.Bool("debug", false, "Enable debug mode.")
	flag.Parse()

	if *debugMode {
		fmt.Println("Debug mode enabled.")
	}

	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	costs, specialRunesStrs, specialRunesCosts := readInputConfig(reader, writer)
	s1, s2 := readInputStrings(reader, writer)

	parseCost := func(costStr string) int {
		cost, err := strconv.Atoi(costStr)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error parsing cost:", err)
			os.Exit(1)
		}
		return cost
	}

	parseRunes := func(runeStr string) rune {
		if len(runeStr) != 1 {
			fmt.Fprintln(os.Stderr, "Invalid rune input. Please enter a single character.")
			os.Exit(1)
		}
		return rune(runeStr[0])
	}

	specialRunes := vagner_fisher.SpecialRunes{
		Replace: parseRunes(specialRunesStrs[0]),
		Insert:  parseRunes(specialRunesStrs[1]),
	}

	opCosts := vagner_fisher.OperationCosts{
		Replace:        parseCost(costs[0]),
		Insert:         parseCost(costs[1]),
		Delete:         parseCost(costs[2]),
		SpecialReplace: parseCost(specialRunesCosts[0]),
		SpecialInsert:  parseCost(specialRunesCosts[1]),
	}

	log := logger.NewLogger(writer)
	if *debugMode {
		log.SetDebugMode()
	}

	distance, operations := vagner_fisher.FindLevenshteinDistance(s1, s2, &opCosts, &specialRunes, log)

	fmt.Fprintln(writer, "\nResults:")
	fmt.Fprintln(writer, "Levenshtein distance: "+strconv.Itoa(distance))
	fmt.Fprintln(writer, "Operations sequence: "+operations)
}
