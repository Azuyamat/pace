package config

import (
	"fmt"

	"github.com/azuyamat/pace/internal/models"
)

type Parser struct {
	lexer          *Lexer
	currentToken   Token
	peekToken      Token
	errors         []error
	input          string
	propertyParser *PropertyParser
}

func Parse(input string) (*Config, error) {
	lexer := NewLexer(input)
	parser := NewParser(lexer)
	return parser.Parse()
}

func NewParser(lexer *Lexer) *Parser {
	p := &Parser{
		lexer:  lexer,
		errors: make([]error, 0),
		input:  lexer.GetInput(),
	}
	p.propertyParser = NewPropertyParser(p)
	p.advance()
	p.advance()
	return p
}

func (p *Parser) advance() {
	p.currentToken = p.peekToken
	p.peekToken = p.lexer.NextToken()

	for p.peekToken.IsOneOf(TOKEN_COMMENT, TOKEN_NEWLINE) {
		p.peekToken = p.lexer.NextToken()
	}
}

func (p *Parser) skipInsignificantTokens() {
	for p.currentToken.IsOneOf(TOKEN_COMMENT, TOKEN_NEWLINE) {
		p.advance()
	}
}

func (p *Parser) expect(tokenType TokenType) error {
	if p.currentToken.Type != tokenType {
		return p.createError(
			fmt.Sprintf("Expected %s but got %s", tokenType.String(), p.currentToken.Type.String()),
		).WithHint(p.getTokenTypeHint(tokenType))
	}
	p.advance()
	return nil
}

func (p *Parser) expectKeyword(keyword string) error {
	if !p.currentToken.IsKeyword(keyword) {
		return p.createError(
			fmt.Sprintf("Expected keyword '%s' but got %s with value %q", keyword, p.currentToken.Type.String(), p.currentToken.Literal),
		).WithContext(fmt.Sprintf("Looking for '%s' keyword", keyword))
	}
	p.advance()
	return nil
}

func (p *Parser) consume() {
	p.advance()
}

func (p *Parser) isAtEnd() bool {
	return p.currentToken.Is(TOKEN_EOF)
}

func (p *Parser) Parse() (*Config, error) {
	config := &Config{
		Tasks:     make(map[string]models.Task),
		Hooks:     make(map[string]models.Hook),
		Globals:   make(map[string]string),
		Constants: make(map[string]string),
		Aliases:   make(map[string]string),
		Imports:   make([]string, 0),
	}

	for !p.isAtEnd() {
		p.skipInsignificantTokens()

		if p.isAtEnd() {
			break
		}

		if err := p.parseTopLevelStatement(config); err != nil {
			return nil, err
		}
	}

	return config, nil
}

func (p *Parser) unexpectedTokenError(context string) error {
	return p.createError(
		fmt.Sprintf("Unexpected %s", p.currentToken.Type.String()),
	).WithContext(fmt.Sprintf("Expected %s", context)).WithHint("Check the syntax of your configuration file")
}

func (p *Parser) createError(message string) *ParseError {
	return newParseError(message, p.currentToken.Line, p.currentToken.Column, p.input)
}

func (p *Parser) getTokenTypeHint(tokenType TokenType) string {
	switch tokenType {
	case TOKEN_LBRACE:
		return "Opening brace '{' is required to start a block"
	case TOKEN_RBRACE:
		return "Closing brace '}' is required to end a block"
	case TOKEN_LBRACKET:
		return "Opening bracket '[' is required to start an array"
	case TOKEN_RBRACKET:
		return "Closing bracket ']' is required to end an array"
	case TOKEN_STRING:
		return "String values must be enclosed in double quotes"
	case TOKEN_IDENTIFIER:
		return "Expected a property name or keyword"
	case TOKEN_COMMA:
		return "Items in a list should be separated by commas"
	default:
		return ""
	}
}
