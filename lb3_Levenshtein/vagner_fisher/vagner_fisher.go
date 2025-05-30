package vagner_fisher

import (
	"fmt"

	"lb3_Levenshtein/logger"
)

var debug bool = false

type Logger interface {
	LogMsg(title, message string)
	LogRuneMatrix(title string, data [][]rune)
	LogCostMatrix(title string, data [][]int)
	SetDebugMode()
}

const (
	Match   = 'M'
	Replace = 'R'
	Insert  = 'I'
	Delete  = 'D'
)

type OperationCosts struct {
	Replace        int
	Insert         int
	Delete         int
	SpecialReplace int
	SpecialInsert  int
}

type SpecialRunes struct {
	Replace rune
	Insert  rune
}

func buildPath(n, m int, opCosts *OperationCosts, ops [][]rune, dp [][]int, s1, s2 string, specRunes *SpecialRunes, log *logger.Logger) string {
	var path []rune
	i, j := n, m

	if log != nil {
		log.LogMsg("BuildPath", fmt.Sprintf("Start backtracking from (%d, %d)", i, j),
			logger.ColorCyan)
	}

	for i > 0 || j > 0 {
		if i > 0 && j > 0 && ops[i][j] == Match {
			path = append(path, Match)
			if log != nil {
				log.LogMsg("BuildPath", fmt.Sprintf("Match at (%d, %d): %c == %c", i, j, s1[i-1], s2[j-1]),
					logger.ColorGreen)
			}
			i--
			j--
		} else if i > 0 && j > 0 && ops[i][j] == Replace {
			path = append(path, Replace)
			if log != nil {
				log.LogMsg("BuildPath", fmt.Sprintf("Replace at (%d, %d): %c -> %c", i, j, s1[i-1], s2[j-1]),
					logger.ColorYellow)
			}
			i--
			j--
		} else if j > 0 && ops[i][j] == Insert {
			path = append(path, Insert)
			if log != nil {
				log.LogMsg("BuildPath", fmt.Sprintf("Insert at (%d, %d): %c", i, j, s2[j-1]),
					logger.ColorBlue)
			}
			j--
		} else if i > 0 && ops[i][j] == Delete {
			path = append(path, Delete)
			if log != nil {
				log.LogMsg("BuildPath", fmt.Sprintf("Delete at (%d, %d): %c", i, j, s1[i-1]),
					logger.ColorRed)
			}
			i--
		} else {
			replaceCost := opCosts.Replace
			if i > 0 && j > 0 && rune(s1[i-1]) == specRunes.Replace {
				replaceCost = opCosts.SpecialReplace
			}

			insertCost := opCosts.Insert
			if j > 0 && rune(s2[j-1]) == specRunes.Insert {
				insertCost = opCosts.SpecialInsert
			}

			replaceTotal := dp[i-1][j-1] + replaceCost
			insertTotal := dp[i][j-1] + insertCost
			deleteTotal := dp[i-1][j] + opCosts.Delete

			minOp, minCost := Replace, replaceTotal
			if insertTotal < minCost {
				minOp, minCost = Insert, insertTotal
			}
			if deleteTotal < minCost {
				minOp, minCost = Delete, deleteTotal
			}

			path = append(path, minOp)
			if log != nil {
				log.LogMsg("BuildPath", fmt.Sprintf("Fallback at (%d, %d): chose %c (replace=%d, insert=%d, delete=%d)", i, j, minOp, replaceTotal, insertTotal, deleteTotal),
					logger.ColorWhite)
			}

			switch minOp {
			case Replace:
				i--
				j--
			case Insert:
				j--
			case Delete:
				i--
			}
		}
	}

	for k := 0; k < len(path)/2; k++ {
		path[k], path[len(path)-1-k] = path[len(path)-1-k], path[k]
	}

	if log != nil {
		log.LogMsg("BuildPath", fmt.Sprintf("Final path: %s", string(path)), logger.ColorGreen)
	}

	return string(path)
}

func minOperation(replaceTotal, insertTotal, deleteTotal int) (rune, int) {
	minCost := replaceTotal
	minOp := Replace

	if insertTotal < minCost {
		minCost = insertTotal
		minOp = Insert
	}
	if deleteTotal < minCost {
		minCost = deleteTotal
		minOp = Delete
	}

	return minOp, minCost
}

func FindLevenshteinDistance(s1, s2 string, opCosts *OperationCosts, specRunes *SpecialRunes, log *logger.Logger) (int, string) {
	n, m := len(s1), len(s2)

	dp := make([][]int, n+1)
	ops := make([][]rune, n+1)
	for i := range dp {
		dp[i] = make([]int, m+1)
		ops[i] = make([]rune, m+1)
	}

	log.LogMsg("Init", fmt.Sprintf("Calculating distance between '%s' (%d) and '%s' (%d)", s1, n, s2, m),
		logger.ColorCyan)
	log.LogMsg("Costs", fmt.Sprintf("Replace: %d, Insert: %d, Delete: %d, SpecialReplace: %d, SpecialInsert: %d",
		opCosts.Replace, opCosts.Insert, opCosts.Delete, opCosts.SpecialReplace, opCosts.SpecialInsert),
		logger.ColorCyan)
	log.LogMsg("SpecialRunes", fmt.Sprintf("Replace: %c, Insert: %c", specRunes.Replace, specRunes.Insert),
		logger.ColorCyan)

	dp[0][0] = 0
	ops[0][0] = Match

	for j := 1; j <= m; j++ {
		dp[0][j] = dp[0][j-1] + opCosts.Insert
		ops[0][j] = Insert
	}

	for i := 1; i <= n; i++ {
		dp[i][0] = dp[i-1][0] + opCosts.Delete
		ops[i][0] = Delete
	}

	log.LogCostMatrix("Initial DP", dp, logger.ColorRed)
	log.LogRuneMatrix("Initial Ops", ops, logger.ColorBlue)

	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			if s1[i-1] == s2[j-1] {
				dp[i][j] = dp[i-1][j-1]
				ops[i][j] = Match
				log.LogMsg("Match", fmt.Sprintf("Characters match at (%d,%d): %c", i, j, s1[i-1]),
					logger.ColorGreen)
			} else {
				replaceCost := opCosts.Replace
				if rune(s1[i-1]) == specRunes.Replace {
					replaceCost = opCosts.SpecialReplace
					log.LogMsg("SpecialReplace", fmt.Sprintf("Special replace at (%d,%d): %c", i, j, s1[i-1]),
						logger.ColorPurple)
				}

				insertCost := opCosts.Insert
				if rune(s2[j-1]) == specRunes.Insert {
					insertCost = opCosts.SpecialInsert
					log.LogMsg("SpecialInsert", fmt.Sprintf("Special insert at (%d,%d): %c", i, j, s2[j-1]),
						logger.ColorPurple)
				}

				replaceTotal := dp[i-1][j-1] + replaceCost
				insertTotal := dp[i][j-1] + insertCost
				deleteTotal := dp[i-1][j] + opCosts.Delete

				minOp, minCost := minOperation(replaceTotal, insertTotal, deleteTotal)

				ops[i][j] = minOp
				dp[i][j] = minCost

				log.LogMsg("Operation", fmt.Sprintf("Cell (%d,%d): chose %c with cost %d (replace=%d, insert=%d, delete=%d)",
					i, j, minOp, minCost, replaceTotal, insertTotal, deleteTotal),
					logger.ColorYellow)
			}
		}
	}

	log.LogCostMatrix("Final DP", dp, logger.ColorRed)
	log.LogRuneMatrix("Final Ops", ops, logger.ColorBlue)

	path := buildPath(n, m, opCosts, ops, dp, s1, s2, specRunes, log)
	log.LogMsg("Result", fmt.Sprintf("Final distance: %d, Path: %s", dp[n][m], path),
		logger.ColorGreen)

	return dp[n][m], path
}
