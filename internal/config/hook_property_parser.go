package config

import (
	"fmt"

	"github.com/azuyamat/pace/internal/models"
)

func (pp *PropertyParser) ParseHookCommand(hook *models.Hook) error {
	pp.parser.advance()

	if pp.parser.currentToken.Type != TOKEN_STRING {
		return pp.parser.createError(
			fmt.Sprintf("Expected command value (string) but got %s", pp.parser.currentToken.Type.String()),
		).WithContext("Parsing 'command' property").WithHint("Command values must be strings, e.g., command \"echo setup\"")
	}

	hook.Command = pp.parser.currentToken.Literal
	pp.parser.advance()
	return nil
}

func (pp *PropertyParser) ParseHookEnvironment(hook *models.Hook) error {
	pp.parser.advance()

	if hook.Env == nil {
		hook.Env = make(map[string]string)
	}

	if err := pp.parser.expect(TOKEN_LBRACE); err != nil {
		return err
	}

	for !pp.parser.currentToken.Is(TOKEN_RBRACE) && !pp.parser.isAtEnd() {
		pp.parser.skipInsignificantTokens()

		if pp.parser.currentToken.Is(TOKEN_RBRACE) {
			break
		}

		if pp.parser.currentToken.Type != TOKEN_STRING {
			return pp.parser.createError(
				fmt.Sprintf("Expected environment variable name (string) but got %s", pp.parser.currentToken.Type.String()),
			).WithContext("Parsing 'env' property").WithHint("Environment variable names must be strings")
		}
		key := pp.parser.currentToken.Literal
		pp.parser.advance()

		if pp.parser.currentToken.Type != TOKEN_STRING {
			return pp.parser.createError(
				fmt.Sprintf("Expected environment variable value (string) but got %s", pp.parser.currentToken.Type.String()),
			).WithContext("Parsing 'env' property").WithHint("Environment variable values must be strings")
		}
		value := pp.parser.currentToken.Literal
		hook.Env[key] = value
		pp.parser.advance()
	}
	if err := pp.parser.expect(TOKEN_RBRACE); err != nil {
		return err
	}

	return nil
}

func (pp *PropertyParser) ParseHookWorkingDir(hook *models.Hook) error {
	pp.parser.advance()

	if pp.parser.currentToken.Type != TOKEN_STRING {
		return pp.parser.createError(
			fmt.Sprintf("Expected working directory value (string) but got %s", pp.parser.currentToken.Type.String()),
		).WithContext("Parsing 'working_dir' property").WithHint("Working directory value must be a string")
	}
	hook.WorkingDir = pp.parser.currentToken.Literal
	pp.parser.advance()
	return nil
}

func (pp *PropertyParser) ParseHookDescription(hook *models.Hook) error {
	pp.parser.advance()

	if pp.parser.currentToken.Type != TOKEN_STRING && pp.parser.currentToken.Type != TOKEN_MULTILINE_STRING {
		return pp.parser.createError(
			fmt.Sprintf("Expected description value (string) but got %s", pp.parser.currentToken.Type.String()),
		).WithContext("Parsing 'description' property").WithHint("Description values must be strings")
	}
	hook.Description = pp.parser.currentToken.Literal
	pp.parser.advance()
	return nil
}
