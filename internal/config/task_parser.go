package config

import (
	"fmt"

	"github.com/azuyamat/pace/internal/models"
)

func (p *Parser) parseTask() (models.Task, error) {
	var task models.Task

	p.advance()

	if p.currentToken.Type != TOKEN_IDENTIFIER {
		return task, p.createError(
			fmt.Sprintf("Expected task name (identifier) but got %s", p.currentToken.Type.String()),
		).WithContext("Parsing task definition").WithHint("Task names must be identifiers, e.g., task build { ... }")
	}
	task.Name = p.currentToken.Literal
	p.advance()

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
		return p.propertyParser.ParseCommand(task)
	case p.currentToken.IsKeyword("inputs"):
		return p.propertyParser.ParseInputs(task)
	case p.currentToken.IsKeyword("outputs"):
		return p.propertyParser.ParseOutputs(task)
	case p.currentToken.IsKeyword("dependencies"):
		return p.propertyParser.ParseDependencies(task)
	case p.currentToken.IsKeyword("env"):
		return p.propertyParser.ParseEnvironment(task)
	case p.currentToken.IsKeyword("cache"):
		return p.propertyParser.ParseCache(task)
	case p.currentToken.IsKeyword("working_dir"):
		return p.propertyParser.ParseWorkingDir(task)
	case p.currentToken.IsKeyword("before"):
		return p.propertyParser.ParseBeforeHooks(task)
	case p.currentToken.IsKeyword("after"):
		return p.propertyParser.ParseAfterHooks(task)
	case p.currentToken.IsKeyword("description"):
		return p.propertyParser.ParseDescription(task)
	case p.currentToken.IsKeyword("watch"):
		return p.propertyParser.ParseWatch(task)
	case p.currentToken.IsKeyword("parallel"):
		return p.propertyParser.ParseParallel(task)
	case p.currentToken.IsKeyword("silent"):
		return p.propertyParser.ParseSilent(task)
	case p.currentToken.IsKeyword("continue_on_error"):
		return p.propertyParser.ParseContinueOnError(task)
	case p.currentToken.IsKeyword("timeout"):
		return p.propertyParser.ParseTimeout(task)
	case p.currentToken.IsKeyword("retry"):
		return p.propertyParser.ParseRetry(task)
	case p.currentToken.IsKeyword("retry_delay"):
		return p.propertyParser.ParseRetryDelay(task)
	case p.currentToken.IsKeyword("on_success"):
		return p.propertyParser.ParseOnSuccess(task)
	case p.currentToken.IsKeyword("on_failure"):
		return p.propertyParser.ParseOnFailure(task)
	case p.currentToken.IsKeyword("args"):
		return p.propertyParser.ParseArgs(task)
	default:
		p.advance()
		return nil
	}
}
