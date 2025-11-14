package config

import (
	"fmt"

	"azuyamat.dev/pace/internal/models"
)

type Parser struct {
	lexer        *Lexer
	currentToken Token
	peekToken    Token
	errors       []error
	input        string
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
	p.advance() // Initialize currentToken
	p.advance() // Initialize peekToken
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
		Tasks: make(map[string]models.Task),
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

func (p *Parser) parseTopLevelStatement(config *Config) error {
	if p.currentToken.IsKeyword("task") {
		task, err := p.parseTask()
		if err != nil {
			return err
		}
		config.Tasks[task.Name] = task
		return nil
	}

	return p.unexpectedTokenError("top-level statement")
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

func (p *Parser) parseTask() (models.Task, error) {
	var task models.Task

	// The 'task' keyword has already been verified by parseTopLevelStatement
	// so we just need to consume it
	p.advance() // consume 'task'

	// Parse task name
	if p.currentToken.Type != TOKEN_STRING {
		return task, p.createError(
			fmt.Sprintf("Expected task name (string) but got %s", p.currentToken.Type.String()),
		).WithContext("Parsing task definition").WithHint("Task names must be enclosed in double quotes, e.g., task \"build\" { ... }")
	}
	task.Name = p.currentToken.Literal
	p.advance()

	// Parse task body
	if err := p.expect(TOKEN_LBRACE); err != nil {
		return task, err
	}

	if err := p.parseTaskBody(&task); err != nil {
		return task, err
	}

	if err := p.expect(TOKEN_RBRACE); err != nil {
		return task, err
	}

	return task, nil
}

func (p *Parser) parseTaskBody(task *models.Task) error {
	for !p.currentToken.Is(TOKEN_RBRACE) && !p.isAtEnd() {
		p.skipInsignificantTokens()

		if p.currentToken.Is(TOKEN_RBRACE) {
			break
		}

		if err := p.parseTaskProperty(task); err != nil {
			return err
		}
	}
	return nil
}

func (p *Parser) parseTaskProperty(task *models.Task) error {
	if !p.currentToken.Is(TOKEN_IDENTIFIER) {
		return p.createError(
			fmt.Sprintf("Expected task property name but got %s", p.currentToken.Type.String()),
		).WithContext("Parsing task body").WithHint("Valid properties include: command, inputs, outputs, dependencies, etc.")
	}

	switch {
	case p.currentToken.IsKeyword("command"):
		return p.parseCommandProperty(task)
	case p.currentToken.IsKeyword("inputs"):
		return p.parseInputsProperty(task)
	case p.currentToken.IsKeyword("outputs"):
		return p.parseOutputsProperty(task)
	case p.currentToken.IsKeyword("dependencies"):
		return p.parseDependenciesProperty(task)
	case p.currentToken.IsKeyword("env"):
		return p.parseEnvironmentProperty(task)
	case p.currentToken.IsKeyword("cache"):
		return p.parseCacheProperty(task)
	case p.currentToken.IsKeyword("working_dir"):
		return p.parseWorkingDirProperty(task)
	default:
		p.advance()
		return nil
	}
}

func (p *Parser) parseCommandProperty(task *models.Task) error {
	p.advance() // consume 'command'

	if p.currentToken.Type != TOKEN_STRING {
		return p.createError(
			fmt.Sprintf("Expected command value (string) but got %s", p.currentToken.Type.String()),
		).WithContext("Parsing 'command' property").WithHint("Command values must be strings, e.g., command \"echo hello\"")
	}

	task.Command = p.currentToken.Literal
	p.advance()
	return nil
}

func (p *Parser) parseInputsProperty(task *models.Task) error {
	p.advance() // consume 'inputs'

	if err := p.expect(TOKEN_LBRACKET); err != nil {
		return err
	}

	for !p.currentToken.Is(TOKEN_RBRACKET) && !p.isAtEnd() {
		p.skipInsignificantTokens()

		if p.currentToken.Is(TOKEN_RBRACKET) {
			break
		}

		if p.currentToken.Type != TOKEN_STRING {
			return p.createError(
				fmt.Sprintf("Expected input value (string) but got %s", p.currentToken.Type.String()),
			).WithContext("Parsing 'inputs' property").WithHint("Input values must be strings, e.g., [\"src/main.go\", \"src/util.go\"]")
		}

		task.Inputs = append(task.Inputs, p.currentToken.Literal)
		p.advance()

		if p.currentToken.Is(TOKEN_COMMA) {
			p.advance() // consume comma
		}
	}

	if err := p.expect(TOKEN_RBRACKET); err != nil {
		return err
	}

	return nil
}

func (p *Parser) parseOutputsProperty(task *models.Task) error {
	p.advance() // consume 'outputs'

	if err := p.expect(TOKEN_LBRACKET); err != nil {
		return err
	}

	for !p.currentToken.Is(TOKEN_RBRACKET) && !p.isAtEnd() {
		p.skipInsignificantTokens()

		if p.currentToken.Is(TOKEN_RBRACKET) {
			break
		}

		if p.currentToken.Type != TOKEN_STRING {
			return p.createError(
				fmt.Sprintf("Expected output value (string) but got %s", p.currentToken.Type.String()),
			).WithContext("Parsing 'outputs' property").WithHint("Output values must be strings, e.g., [\"bin/app\", \"bin/util\"]")
		}
		task.Outputs = append(task.Outputs, p.currentToken.Literal)
		p.advance()
		if p.currentToken.Is(TOKEN_COMMA) {
			p.advance() // consume comma
		}
	}

	if err := p.expect(TOKEN_RBRACKET); err != nil {
		return err
	}

	return nil
}

func (p *Parser) parseDependenciesProperty(task *models.Task) error {
	p.advance() // consume 'dependencies'
	if err := p.expect(TOKEN_LBRACKET); err != nil {
		return err
	}
	for !p.currentToken.Is(TOKEN_RBRACKET) && !p.isAtEnd() {
		p.skipInsignificantTokens()

		if p.currentToken.Is(TOKEN_RBRACKET) {
			break
		}

		if p.currentToken.Type != TOKEN_STRING {
			return p.createError(
				fmt.Sprintf("Expected dependency value (string) but got %s", p.currentToken.Type.String()),
			).WithContext("Parsing 'dependencies' property").WithHint("Dependency values must be strings, e.g., [\"build\", \"test\"]")
		}
		task.Dependencies = append(task.Dependencies, p.currentToken.Literal)
		p.advance()
		if p.currentToken.Is(TOKEN_COMMA) {
			p.advance() // consume comma
		}
	}

	if err := p.expect(TOKEN_RBRACKET); err != nil {
		return err
	}

	return nil
}

func (p *Parser) parseEnvironmentProperty(task *models.Task) error {
	p.advance()

	if err := p.expect(TOKEN_LBRACE); err != nil {
		return err
	}

	task.Env = make(map[string]string)

	for !p.currentToken.Is(TOKEN_RBRACE) && !p.isAtEnd() {
		p.skipInsignificantTokens()

		if p.currentToken.Is(TOKEN_RBRACE) {
			break
		}

		if p.currentToken.Type != TOKEN_STRING {
			return p.createError(
				fmt.Sprintf("Expected environment variable name (string) but got %s", p.currentToken.Type.String()),
			).WithContext("Parsing 'env' property").WithHint("Environment variable names must be strings, e.g., \"PATH\"")
		}
		key := p.currentToken.Literal
		p.advance()

		if p.currentToken.Type != TOKEN_STRING {
			return p.createError(
				fmt.Sprintf("Expected environment variable value (string) but got %s", p.currentToken.Type.String()),
			).WithContext("Parsing 'env' property").WithHint("Environment variable values must be strings")
		}
		value := p.currentToken.Literal
		task.Env[key] = value
		p.advance()
	}
	if err := p.expect(TOKEN_RBRACE); err != nil {
		return err
	}

	return nil
}

func (p *Parser) parseCacheProperty(task *models.Task) error {
	p.advance() // consume 'cache'

	if p.currentToken.IsTrue() {
		task.Cache = true
		p.advance()
	} else if p.currentToken.IsFalse() {
		task.Cache = false
		p.advance()
	} else {
		return p.createError(
			fmt.Sprintf("Expected boolean value for 'cache' property but got %s", p.currentToken.Type.String()),
		).WithContext("Parsing 'cache' property").WithHint("Cache property must be either true or false")
	}
	return nil
}

func (p *Parser) parseWorkingDirProperty(task *models.Task) error {
	p.advance() // consume 'working_dir'

	if p.currentToken.Type != TOKEN_STRING {
		return p.createError(
			fmt.Sprintf("Expected working directory value (string) but got %s", p.currentToken.Type.String()),
		).WithContext("Parsing 'working_dir' property").WithHint("Working directory value must be a string, e.g., \"/app\"")
	}
	task.WorkingDir = p.currentToken.Literal
	p.advance()
	return nil
}
