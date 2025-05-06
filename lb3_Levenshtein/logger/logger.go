package logger

import (
	"bufio"
	"fmt"
)

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
)

type Logger struct {
	Writer *bufio.Writer
	Debug  bool
}

func NewLogger(writer *bufio.Writer) *Logger {
	return &Logger{
		Writer: writer,
		Debug:  false,
	}
}

func (l *Logger) SetDebugMode() {
	l.Debug = true
}

func (l *Logger) LogMsg(title, message, color string) {
	if l.Debug {
		fmt.Fprintf(l.Writer, "%s[%s]:%s %s\n", color, title, ColorReset, message)
		l.Writer.Flush()
	}
}

func (l *Logger) LogRuneMatrix(title string, data [][]rune, color string) {
	if l.Debug {
		fmt.Fprintf(l.Writer, "%s[%s]:%s \n", color, title, ColorReset)
		for _, row := range data {
			for _, val := range row {
				fmt.Fprint(l.Writer, string(val), " ")
			}
			fmt.Fprint(l.Writer, "\n")
		}
		l.Writer.Flush()
	}
}

func (l *Logger) LogCostMatrix(title string, data [][]int, color string) {
	if l.Debug {
		fmt.Fprintf(l.Writer, "%s[%s]:%s \n", color, title, ColorReset)
		for _, row := range data {
			for _, val := range row {
				fmt.Fprint(l.Writer, val, " ")
			}
			fmt.Fprint(l.Writer, "\n")
		}
		l.Writer.Flush()
	}
}
