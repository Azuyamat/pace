package parsing

import (
	"fmt"
	"strings"

	"github.com/azuyamat/pace/internal/logger"
)

type ParseError struct {
	Message string
	Line    int
	Column  int
	Input   string
	Context string
	Hint    string
}

func (e *ParseError) Error() string {
	var sb strings.Builder

	sb.WriteString("\n\n")
	sb.WriteString(logger.ColorRed.Bright().Wrap(e.Message))
	sb.WriteString("\n\n")

	if e.Input != "" {
		sb.WriteString(e.formatSourceContext())
		sb.WriteString("\n")
	}

	sb.WriteString(logger.ColorYellow.Wrap(fmt.Sprintf("line %d, column %d", e.Line, e.Column)))
	sb.WriteString("\n")

	if e.Context != "" {
		sb.WriteString("\n")
		sb.WriteString(logger.ColorGray.Wrap(e.Context))
		sb.WriteString("\n")
	}

	if e.Hint != "" {
		sb.WriteString("\n")
		sb.WriteString(logger.ColorGreen.Bright().Wrap(e.Hint))
		sb.WriteString("\n")
	}

	sb.WriteString("\n")

	return sb.String()
}

func (e *ParseError) formatSourceContext() string {
	lines := strings.Split(e.Input, "\n")

	if e.Line < 1 || e.Line > len(lines) {
		return ""
	}

	var sb strings.Builder

	lineNumWidth := len(fmt.Sprintf("%d", e.Line+1))

	if e.Line > 1 {
		lineNum := fmt.Sprintf("%*d", lineNumWidth, e.Line-1)
		sb.WriteString("  " + logger.ColorGray.Wrap(lineNum) + " " + logger.ColorGray.Wrap("│") + " ")
		sb.WriteString(logger.ColorGray.Wrap(lines[e.Line-2]))
		sb.WriteString("\n")
	}

	lineNum := fmt.Sprintf("%*d", lineNumWidth, e.Line)
	errorLine := lines[e.Line-1]
	sb.WriteString("  " + logger.ColorYellow.Bright().Wrap(lineNum) + " " + logger.ColorRed.Wrap("│") + " ")
	sb.WriteString(errorLine)
	sb.WriteString("\n")

	sb.WriteString("  " + logger.ColorGray.Wrap(strings.Repeat(" ", lineNumWidth)) + " " + logger.ColorRed.Wrap("│") + " ")

	column := e.Column - 1
	if column < 0 {
		column = 0
	}
	if column > len(errorLine) {
		column = len(errorLine)
	}

	for i := 0; i < column; i++ {
		if errorLine[i] == '\t' {
			sb.WriteString("\t")
		} else {
			sb.WriteString(" ")
		}
	}

	underlineLength := 1
	if column < len(errorLine) {
		for i := column + 1; i < len(errorLine) && i < column+20; i++ {
			ch := errorLine[i]
			if ch == ' ' || ch == '\t' || ch == '\n' || ch == '{' || ch == '}' || ch == '[' || ch == ']' || ch == ',' || ch == '"' {
				break
			}
			underlineLength++
		}
	}

	sb.WriteString(logger.ColorRed.Bright().Wrap("^"))
	for i := 1; i < underlineLength; i++ {
		sb.WriteString(logger.ColorRed.Bright().Wrap("~"))
	}
	sb.WriteString(" " + logger.ColorRed.Bright().Wrap("✗"))
	sb.WriteString("\n")

	if e.Line < len(lines) {
		lineNum := fmt.Sprintf("%*d", lineNumWidth, e.Line+1)
		sb.WriteString("  " + logger.ColorGray.Wrap(lineNum) + " " + logger.ColorGray.Wrap("│") + " ")
		sb.WriteString(logger.ColorGray.Wrap(lines[e.Line]))
		sb.WriteString("\n")
	}

	return sb.String()
}

func newParseError(message string, line, column int, input string) *ParseError {
	return &ParseError{
		Message: message,
		Line:    line,
		Column:  column,
		Input:   input,
	}
}

func (e *ParseError) WithContext(context string) *ParseError {
	e.Context = context
	return e
}

func (e *ParseError) WithHint(hint string) *ParseError {
	e.Hint = hint
	return e
}
