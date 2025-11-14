package config

import (
	"fmt"
	"strings"
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

	sb.WriteString("\n")
	sb.WriteString("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	sb.WriteString(fmt.Sprintf("Parse Error: %s\n", e.Message))
	sb.WriteString("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")

	if e.Input != "" {
		sb.WriteString(e.formatSourceContext())
		sb.WriteString("\n")
	}

	sb.WriteString(fmt.Sprintf("  Location: line %d, column %d\n", e.Line, e.Column))

	if e.Context != "" {
		sb.WriteString(fmt.Sprintf("  Context: %s\n", e.Context))
	}

	if e.Hint != "" {
		sb.WriteString(fmt.Sprintf("\n  💡 Hint: %s\n", e.Hint))
	}

	sb.WriteString("\n")

	return sb.String()
}

// formatSourceContext creates a visual representation of where the error occurred
func (e *ParseError) formatSourceContext() string {
	lines := strings.Split(e.Input, "\n")

	if e.Line < 1 || e.Line > len(lines) {
		return ""
	}

	var sb strings.Builder

	// Show line numbers with padding
	lineNumWidth := len(fmt.Sprintf("%d", e.Line+1))

	// Show previous line if available
	if e.Line > 1 {
		sb.WriteString(fmt.Sprintf("  %*d | %s\n", lineNumWidth, e.Line-1, lines[e.Line-2]))
	}

	// Show the error line
	errorLine := lines[e.Line-1]
	sb.WriteString(fmt.Sprintf("  %*d | %s\n", lineNumWidth, e.Line, errorLine))

	// Show the error indicator
	sb.WriteString(fmt.Sprintf("  %*s | ", lineNumWidth, ""))

	// Add spaces to align with the column
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
	sb.WriteString("^")

	// Add wavy underline for multi-character tokens
	underlineLength := 1
	if column < len(errorLine) {
		// Try to underline the whole token
		for i := column + 1; i < len(errorLine) && i < column+20; i++ {
			ch := errorLine[i]
			if ch == ' ' || ch == '\t' || ch == '\n' || ch == '{' || ch == '}' || ch == '[' || ch == ']' || ch == ',' || ch == '"' {
				break
			}
			underlineLength++
		}
	}

	for i := 1; i < underlineLength; i++ {
		sb.WriteString("~")
	}

	sb.WriteString(" ❌\n")

	// Show next line if available
	if e.Line < len(lines) {
		sb.WriteString(fmt.Sprintf("  %*d | %s\n", lineNumWidth, e.Line+1, lines[e.Line]))
	}

	return sb.String()
}

// newParseError creates a new ParseError
func newParseError(message string, line, column int, input string) *ParseError {
	return &ParseError{
		Message: message,
		Line:    line,
		Column:  column,
		Input:   input,
	}
}

// WithContext adds context information to the error
func (e *ParseError) WithContext(context string) *ParseError {
	e.Context = context
	return e
}

// WithHint adds a helpful hint to the error
func (e *ParseError) WithHint(hint string) *ParseError {
	e.Hint = hint
	return e
}
