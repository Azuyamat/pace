package parsing

import (
	"fmt"
	"strings"
)

type TokenType int

const (
	TOKEN_EOF TokenType = iota

	TOKEN_IDENTIFIER
	TOKEN_STRING
	TOKEN_MULTILINE_STRING
	TOKEN_NUMBER
	TOKEN_BOOLEAN

	TOKEN_LBRACE
	TOKEN_RBRACE
	TOKEN_LBRACKET
	TOKEN_RBRACKET
	TOKEN_COMMA
	TOKEN_DOLLAR
	TOKEN_LPAREN
	TOKEN_RPAREN
	TOKEN_EQUALS

	TOKEN_COMMENT
	TOKEN_NEWLINE
	TOKEN_ILLEGAL
)

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

func newToken(tokenType TokenType, ch byte, line, column int) Token {
	return Token{
		Type:    tokenType,
		Literal: string(ch),
		Line:    line,
		Column:  column,
	}
}

func NewTokenWithLiteral(tokenType TokenType, literal string, line, column int) Token {
	return Token{
		Type:    tokenType,
		Literal: literal,
		Line:    line,
		Column:  column,
	}
}

func NewEOFToken(line, column int) Token {
	return Token{
		Type:    TOKEN_EOF,
		Literal: "",
		Line:    line,
		Column:  column,
	}
}

func (tokenType TokenType) String() string {
	switch tokenType {
	case TOKEN_IDENTIFIER:
		return "IDENTIFIER"
	case TOKEN_STRING:
		return "STRING"
	case TOKEN_MULTILINE_STRING:
		return "MULTILINE_STRING"
	case TOKEN_NUMBER:
		return "NUMBER"
	case TOKEN_BOOLEAN:
		return "BOOLEAN"
	case TOKEN_LBRACE:
		return "LBRACE"
	case TOKEN_RBRACE:
		return "RBRACE"
	case TOKEN_LBRACKET:
		return "LBRACKET"
	case TOKEN_RBRACKET:
		return "RBRACKET"
	case TOKEN_COMMA:
		return "COMMA"
	case TOKEN_DOLLAR:
		return "DOLLAR"
	case TOKEN_LPAREN:
		return "LPAREN"
	case TOKEN_RPAREN:
		return "RPAREN"
	case TOKEN_EQUALS:
		return "EQUALS"
	case TOKEN_COMMENT:
		return "COMMENT"
	case TOKEN_NEWLINE:
		return "NEWLINE"
	case TOKEN_ILLEGAL:
		return "ILLEGAL"
	case TOKEN_EOF:
		return "EOF"
	default:
		return "UNKNOWN"
	}
}

func (t Token) String() string {
	return fmt.Sprintf("Token{Type: %s, Literal: %q, Line: %d, Column: %d}", t.Type.String(), t.Literal, t.Line, t.Column)
}

func (t Token) Is(tokenType TokenType) bool {
	return t.Type == tokenType
}

func (t Token) IsOneOf(tokenTypes ...TokenType) bool {
	for _, tt := range tokenTypes {
		if t.Type == tt {
			return true
		}
	}
	return false
}

func (t Token) Expect(tokenType TokenType) error {
	if t.Type != tokenType {
		return fmt.Errorf("expected token %s, got %s at line %d, column %d", tokenType.String(), t.Type.String(), t.Line, t.Column)
	}
	return nil
}

func (t Token) LiteralIs(expected string, ignoreCase bool) bool {
	if ignoreCase {
		return strings.EqualFold(t.Literal, expected)
	}
	return strings.Compare(t.Literal, expected) == 0
}

func (t Token) IsKeyword(keyword string) bool {
	return t.Is(TOKEN_IDENTIFIER) && t.LiteralIs(keyword, true)
}

func (t Token) ExpectKeyword(keyword string) error {
	if !t.IsKeyword(keyword) {
		return fmt.Errorf("expected keyword '%s', got %s with value %q at line %d, column %d",
			keyword, t.Type.String(), t.Literal, t.Line, t.Column)
	}
	return nil
}

func (t Token) IsTrue() bool {
	return t.Is(TOKEN_BOOLEAN) && strings.EqualFold(t.Literal, "true")
}

func (t Token) IsFalse() bool {
	return t.Is(TOKEN_BOOLEAN) && strings.EqualFold(t.Literal, "false")
}
