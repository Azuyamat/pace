package parsing

import (
	"fmt"

	"github.com/azuyamat/pace/internal/config/types"
	"github.com/azuyamat/pace/internal/models"
)

type ExpectedTokenType int

const (
	ExpectIdentifier ExpectedTokenType = iota
	ExpectString
	ExpectIdentifierOrString
)

type SimpleStatementDef struct {
	Arg1Type    ExpectedTokenType
	Arg1Hint    string
	Arg2Type    *ExpectedTokenType
	Arg2Hint    string
	NeedsEquals bool
	Handler     func(p *Parser, config *types.Config, arg1, arg2 string) error
}

var simpleStatements = map[string]SimpleStatementDef{
	"var": {
		Arg1Type:    ExpectIdentifier,
		Arg1Hint:    "Variable names must be identifiers, e.g., var output = \"bin/app\"",
		Arg2Type:    ptr(ExpectString),
		Arg2Hint:    "Variable values must be strings",
		NeedsEquals: true,
		Handler: func(p *Parser, config *types.Config, name, value string) error {
			config.Constants[name] = value
			return nil
		},
	},
	"default": {
		Arg1Type: ExpectIdentifier,
		Arg1Hint: "Default task name must be an identifier, e.g., default build",
		Handler: func(p *Parser, config *types.Config, taskName, _ string) error {
			config.DefaultTask = taskName
			return nil
		},
	},
	"alias": {
		Arg1Type: ExpectIdentifier,
		Arg1Hint: "Alias names must be identifiers, e.g., alias b build",
		Arg2Type: ptr(ExpectIdentifier),
		Arg2Hint: "Task names must be identifiers, e.g., alias b build",
		Handler: func(p *Parser, config *types.Config, aliasName, taskName string) error {
			config.Aliases[aliasName] = taskName
			return nil
		},
	},
	"import": {
		Arg1Type: ExpectString,
		Arg1Hint: "Import paths must be strings, e.g., import \"tasks/build.pace\"",
		Handler: func(p *Parser, config *types.Config, importPath, _ string) error {
			config.Imports = append(config.Imports, importPath)
			return nil
		},
	},
}

func ptr[T any](v T) *T { return &v }

type StatementHandler func(p *Parser, config *types.Config) error

var statementRegistry = map[string]StatementHandler{
	"task":    (*Parser).parseTaskStatement,
	"hook":    (*Parser).parseHookStatement,
	"var":     (*Parser).parseSimpleStatement,
	"default": (*Parser).parseSimpleStatement,
	"alias":   (*Parser).parseSimpleStatement,
	"import":  (*Parser).parseSimpleStatement,
}

func (p *Parser) parseTopLevelStatement(config *types.Config) error {
	keyword := p.currentToken.Literal
	handler, exists := statementRegistry[keyword]

	if !exists {
		return p.unexpectedTokenError("top-level statement")
	}

	return handler(p, config)
}

func (p *Parser) parseSimpleStatement(config *types.Config) error {
	keyword := p.currentToken.Literal
	def, exists := simpleStatements[keyword]
	if !exists {
		return fmt.Errorf("no definition for statement: %s", keyword)
	}

	p.advance()

	arg1, err := p.expectToken(def.Arg1Type, def.Arg1Hint)
	if err != nil {
		return err
	}

	if def.NeedsEquals {
		if err := p.expect(TOKEN_EQUALS); err != nil {
			return p.createError(
				fmt.Sprintf("Expected '=' but got %s", p.currentToken.Type.String()),
			).WithContext(fmt.Sprintf("Parsing %s statement", keyword))
		}
	}

	var arg2 string
	if def.Arg2Type != nil {
		arg2, err = p.expectToken(*def.Arg2Type, def.Arg2Hint)
		if err != nil {
			return err
		}
	}

	return def.Handler(p, config, arg1, arg2)
}

func (p *Parser) parseTaskStatement(config *types.Config) error {
	task, err := p.parseTask()
	if err != nil {
		return err
	}
	config.Tasks[task.Name] = task
	// Register alias if task has one
	if task.Alias != "" {
		config.Aliases[task.Alias] = task.Name
	}
	return nil
}

func (p *Parser) parseHookStatement(config *types.Config) error {
	hook, err := p.parseHook()
	if err != nil {
		return err
	}
	config.Hooks[hook.Name] = hook
	return nil
}

func (p *Parser) expectToken(tokenType ExpectedTokenType, hint string) (string, error) {
	switch tokenType {
	case ExpectIdentifier:
		return p.expectIdentifier("identifier", hint)
	case ExpectString:
		return p.expectString("string", hint)
	case ExpectIdentifierOrString:
		return p.expectIdentifierOrString("identifier or string", hint)
	default:
		return "", fmt.Errorf("unknown expected token type: %d", tokenType)
	}
}

func (p *Parser) expectIdentifier(contextName, hint string) (string, error) {
	if p.currentToken.Type != TOKEN_IDENTIFIER {
		return "", p.createError(
			fmt.Sprintf("Expected %s (identifier) but got %s", contextName, p.currentToken.Type.String()),
		).WithContext(fmt.Sprintf("Parsing %s", contextName)).WithHint(hint)
	}
	value := p.currentToken.Literal
	p.advance()
	return value, nil
}

func (p *Parser) expectString(contextName, hint string) (string, error) {
	if p.currentToken.Type != TOKEN_STRING {
		return "", p.createError(
			fmt.Sprintf("Expected %s (string) but got %s", contextName, p.currentToken.Type.String()),
		).WithContext(fmt.Sprintf("Parsing %s", contextName)).WithHint(hint)
	}
	value := p.currentToken.Literal
	p.advance()
	return value, nil
}

func (p *Parser) expectIdentifierOrString(contextName, hint string) (string, error) {
	if p.currentToken.Type != TOKEN_IDENTIFIER && p.currentToken.Type != TOKEN_STRING {
		return "", p.createError(
			fmt.Sprintf("Expected %s but got %s", contextName, p.currentToken.Type.String()),
		).WithContext(fmt.Sprintf("Parsing %s", contextName)).WithHint(hint)
	}
	value := p.currentToken.Literal
	p.advance()
	return value, nil
}

func (p *Parser) parseTask() (models.Task, error) {
	var task models.Task

	p.advance()

	if p.currentToken.Type != TOKEN_IDENTIFIER {
		return task, p.createError(
			fmt.Sprintf("Expected task name (identifier) but got %s", p.currentToken.Type.String()),
		).WithContext("Parsing task definition").WithHint("Task names must be identifiers, e.g., task build [b] { ... }")
	}
	task.Name = p.currentToken.Literal
	p.advance()

	// Check for optional alias syntax: task name [alias]
	if p.currentToken.Type == TOKEN_LBRACKET {
		p.advance()
		if p.currentToken.Type != TOKEN_IDENTIFIER {
			return task, p.createError(
				fmt.Sprintf("Expected alias name (identifier) but got %s", p.currentToken.Type.String()),
			).WithContext("Parsing task alias").WithHint("Alias must be an identifier, e.g., task build [b] { ... }")
		}
		task.Alias = p.currentToken.Literal
		p.advance()
		if err := p.expect(TOKEN_RBRACKET); err != nil {
			return task, err
		}
	}

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

		if err := p.propertyParser.ParseTaskProperty(task); err != nil {
			return err
		}
	}
	return nil
}

func (p *Parser) parseHook() (models.Hook, error) {
	var hook models.Hook

	p.advance()

	if p.currentToken.Type != TOKEN_IDENTIFIER {
		return hook, p.createError(
			fmt.Sprintf("Expected hook name (identifier) but got %s", p.currentToken.Type.String()),
		).WithContext("Parsing hook definition").WithHint("Hook names must be identifiers, e.g., hook setup { ... }")
	}
	hook.Name = p.currentToken.Literal
	p.advance()

	if err := p.expect(TOKEN_LBRACE); err != nil {
		return hook, err
	}

	if err := p.parseHookBody(&hook); err != nil {
		return hook, err
	}

	if err := p.expect(TOKEN_RBRACE); err != nil {
		return hook, err
	}

	return hook, nil
}

func (p *Parser) parseHookBody(hook *models.Hook) error {
	for !p.currentToken.Is(TOKEN_RBRACE) && !p.isAtEnd() {
		p.skipInsignificantTokens()

		if p.currentToken.Is(TOKEN_RBRACE) {
			break
		}

		if err := p.propertyParser.ParseHookProperty(hook); err != nil {
			return err
		}
	}
	return nil
}
