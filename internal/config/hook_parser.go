package config

import (
	"fmt"

	"azuyamat.dev/pace/internal/models"
)

func (p *Parser) parseHook() (models.Hook, error) {
	var hook models.Hook

	p.advance()

	if p.currentToken.Type != TOKEN_STRING {
		return hook, p.createError(
			fmt.Sprintf("Expected hook name (string) but got %s", p.currentToken.Type.String()),
		).WithContext("Parsing hook definition").WithHint("Hook names must be enclosed in double quotes, e.g., hook \"setup\" { ... }")
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

		if err := p.parseHookProperty(hook); err != nil {
			return err
		}
	}
	return nil
}

func (p *Parser) parseHookProperty(hook *models.Hook) error {
	if !p.currentToken.Is(TOKEN_IDENTIFIER) {
		return p.createError(
			fmt.Sprintf("Expected hook property name but got %s", p.currentToken.Type.String()),
		).WithContext("Parsing hook body").WithHint("Valid properties include: command, env, working_dir")
	}

	switch {
	case p.currentToken.IsKeyword("command"):
		return p.propertyParser.ParseHookCommand(hook)
	case p.currentToken.IsKeyword("env"):
		return p.propertyParser.ParseHookEnvironment(hook)
	case p.currentToken.IsKeyword("working_dir"):
		return p.propertyParser.ParseHookWorkingDir(hook)
	case p.currentToken.IsKeyword("description"):
		return p.propertyParser.ParseHookDescription(hook)
	default:
		p.advance()
		return nil
	}
}
