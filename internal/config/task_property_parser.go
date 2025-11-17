package config

import (
	"fmt"

	"github.com/azuyamat/pace/internal/models"
)

type PropertyParser struct {
	parser *Parser
}

func NewPropertyParser(parser *Parser) *PropertyParser {
	return &PropertyParser{
		parser: parser,
	}
}

func (pp *PropertyParser) ParseCommand(task *models.Task) error {
	pp.parser.advance()

	if pp.parser.currentToken.Type != TOKEN_STRING && pp.parser.currentToken.Type != TOKEN_MULTILINE_STRING {
		return pp.parser.createError(
			fmt.Sprintf("Expected command value (string) but got %s", pp.parser.currentToken.Type.String()),
		).WithContext("Parsing 'command' property").WithHint("Command values must be strings, e.g., command \"echo hello\"")
	}

	task.Command = pp.parser.currentToken.Literal
	pp.parser.advance()
	return nil
}

func (pp *PropertyParser) ParseInputs(task *models.Task) error {
	pp.parser.advance()

	if err := pp.parser.expect(TOKEN_LBRACKET); err != nil {
		return err
	}

	for !pp.parser.currentToken.Is(TOKEN_RBRACKET) && !pp.parser.isAtEnd() {
		pp.parser.skipInsignificantTokens()

		if pp.parser.currentToken.Is(TOKEN_RBRACKET) {
			break
		}

		if pp.parser.currentToken.Type != TOKEN_STRING {
			return pp.parser.createError(
				fmt.Sprintf("Expected input value (string) but got %s", pp.parser.currentToken.Type.String()),
			).WithContext("Parsing 'inputs' property").WithHint("Input values must be strings, e.g., [\"src/main.go\", \"src/util.go\"]")
		}

		task.Inputs = append(task.Inputs, pp.parser.currentToken.Literal)
		pp.parser.advance()

		if pp.parser.currentToken.Is(TOKEN_COMMA) {
			pp.parser.advance()
		}
	}

	if err := pp.parser.expect(TOKEN_RBRACKET); err != nil {
		return err
	}

	return nil
}

func (pp *PropertyParser) ParseOutputs(task *models.Task) error {
	pp.parser.advance()

	if err := pp.parser.expect(TOKEN_LBRACKET); err != nil {
		return err
	}

	for !pp.parser.currentToken.Is(TOKEN_RBRACKET) && !pp.parser.isAtEnd() {
		pp.parser.skipInsignificantTokens()

		if pp.parser.currentToken.Is(TOKEN_RBRACKET) {
			break
		}

		if pp.parser.currentToken.Type != TOKEN_STRING {
			return pp.parser.createError(
				fmt.Sprintf("Expected output value (string) but got %s", pp.parser.currentToken.Type.String()),
			).WithContext("Parsing 'outputs' property").WithHint("Output values must be strings, e.g., [\"bin/app\", \"bin/util\"]")
		}
		task.Outputs = append(task.Outputs, pp.parser.currentToken.Literal)
		pp.parser.advance()
		if pp.parser.currentToken.Is(TOKEN_COMMA) {
			pp.parser.advance()
		}
	}

	if err := pp.parser.expect(TOKEN_RBRACKET); err != nil {
		return err
	}

	return nil
}

func (pp *PropertyParser) ParseDependencies(task *models.Task) error {
	pp.parser.advance()
	if err := pp.parser.expect(TOKEN_LBRACKET); err != nil {
		return err
	}
	for !pp.parser.currentToken.Is(TOKEN_RBRACKET) && !pp.parser.isAtEnd() {
		pp.parser.skipInsignificantTokens()

		if pp.parser.currentToken.Is(TOKEN_RBRACKET) {
			break
		}

		if pp.parser.currentToken.Type != TOKEN_STRING {
			return pp.parser.createError(
				fmt.Sprintf("Expected dependency value (string) but got %s", pp.parser.currentToken.Type.String()),
			).WithContext("Parsing 'dependencies' property").WithHint("Dependency values must be strings, e.g., [\"build\", \"test\"]")
		}
		task.Dependencies = append(task.Dependencies, pp.parser.currentToken.Literal)
		pp.parser.advance()
		if pp.parser.currentToken.Is(TOKEN_COMMA) {
			pp.parser.advance()
		}
	}

	if err := pp.parser.expect(TOKEN_RBRACKET); err != nil {
		return err
	}

	return nil
}

func (pp *PropertyParser) ParseEnvironment(task *models.Task) error {
	pp.parser.advance()

	if err := pp.parser.expect(TOKEN_LBRACE); err != nil {
		return err
	}

	task.Env = make(map[string]string)

	for !pp.parser.currentToken.Is(TOKEN_RBRACE) && !pp.parser.isAtEnd() {
		pp.parser.skipInsignificantTokens()

		if pp.parser.currentToken.Is(TOKEN_RBRACE) {
			break
		}

		if pp.parser.currentToken.Type != TOKEN_STRING {
			return pp.parser.createError(
				fmt.Sprintf("Expected environment variable name (string) but got %s", pp.parser.currentToken.Type.String()),
			).WithContext("Parsing 'env' property").WithHint("Environment variable names must be strings, e.g., \"PATH\"")
		}
		key := pp.parser.currentToken.Literal
		pp.parser.advance()

		if pp.parser.currentToken.Type != TOKEN_STRING {
			return pp.parser.createError(
				fmt.Sprintf("Expected environment variable value (string) but got %s", pp.parser.currentToken.Type.String()),
			).WithContext("Parsing 'env' property").WithHint("Environment variable values must be strings")
		}
		value := pp.parser.currentToken.Literal
		task.Env[key] = value
		pp.parser.advance()
	}
	if err := pp.parser.expect(TOKEN_RBRACE); err != nil {
		return err
	}

	return nil
}

func (pp *PropertyParser) ParseWorkingDir(task *models.Task) error {
	pp.parser.advance()

	if pp.parser.currentToken.Type != TOKEN_STRING {
		return pp.parser.createError(
			fmt.Sprintf("Expected working directory value (string) but got %s", pp.parser.currentToken.Type.String()),
		).WithContext("Parsing 'working_dir' property").WithHint("Working directory value must be a string, e.g., \"/app\"")
	}
	task.WorkingDir = pp.parser.currentToken.Literal
	pp.parser.advance()
	return nil
}

func (pp *PropertyParser) ParseBeforeHooks(task *models.Task) error {
	pp.parser.advance()

	if err := pp.parser.expect(TOKEN_LBRACKET); err != nil {
		return err
	}

	for !pp.parser.currentToken.Is(TOKEN_RBRACKET) && !pp.parser.isAtEnd() {
		pp.parser.skipInsignificantTokens()

		if pp.parser.currentToken.Is(TOKEN_RBRACKET) {
			break
		}

		if pp.parser.currentToken.Type != TOKEN_STRING {
			return pp.parser.createError(
				fmt.Sprintf("Expected hook name (string) but got %s", pp.parser.currentToken.Type.String()),
			).WithContext("Parsing 'before' property").WithHint("Hook names must be strings, e.g., [\"setup\", \"clean\"]")
		}

		task.BeforeHooks = append(task.BeforeHooks, pp.parser.currentToken.Literal)
		pp.parser.advance()

		if pp.parser.currentToken.Is(TOKEN_COMMA) {
			pp.parser.advance()
		}
	}

	if err := pp.parser.expect(TOKEN_RBRACKET); err != nil {
		return err
	}

	return nil
}

func (pp *PropertyParser) ParseAfterHooks(task *models.Task) error {
	pp.parser.advance()

	if err := pp.parser.expect(TOKEN_LBRACKET); err != nil {
		return err
	}

	for !pp.parser.currentToken.Is(TOKEN_RBRACKET) && !pp.parser.isAtEnd() {
		pp.parser.skipInsignificantTokens()

		if pp.parser.currentToken.Is(TOKEN_RBRACKET) {
			break
		}

		if pp.parser.currentToken.Type != TOKEN_STRING {
			return pp.parser.createError(
				fmt.Sprintf("Expected hook name (string) but got %s", pp.parser.currentToken.Type.String()),
			).WithContext("Parsing 'after' property").WithHint("Hook names must be strings, e.g., [\"cleanup\", \"notify\"]")
		}

		task.AfterHooks = append(task.AfterHooks, pp.parser.currentToken.Literal)
		pp.parser.advance()

		if pp.parser.currentToken.Is(TOKEN_COMMA) {
			pp.parser.advance()
		}
	}

	if err := pp.parser.expect(TOKEN_RBRACKET); err != nil {
		return err
	}

	return nil
}

func (pp *PropertyParser) ParseDescription(task *models.Task) error {
	pp.parser.advance()

	if pp.parser.currentToken.Type != TOKEN_STRING && pp.parser.currentToken.Type != TOKEN_MULTILINE_STRING {
		return pp.parser.createError(
			fmt.Sprintf("Expected description value (string) but got %s", pp.parser.currentToken.Type.String()),
		).WithContext("Parsing 'description' property").WithHint("Description values must be strings")
	}
	task.Description = pp.parser.currentToken.Literal
	pp.parser.advance()
	return nil
}

func (pp *PropertyParser) ParseTimeout(task *models.Task) error {
	pp.parser.advance()

	if pp.parser.currentToken.Type != TOKEN_STRING {
		return pp.parser.createError(
			fmt.Sprintf("Expected timeout value (string) but got %s", pp.parser.currentToken.Type.String()),
		).WithContext("Parsing 'timeout' property").WithHint("Timeout values must be strings like \"5m\", \"30s\", \"1h\"")
	}
	task.Timeout = pp.parser.currentToken.Literal
	pp.parser.advance()
	return nil
}

func (pp *PropertyParser) ParseRetry(task *models.Task) error {
	pp.parser.advance()

	if pp.parser.currentToken.Type != TOKEN_NUMBER {
		return pp.parser.createError(
			fmt.Sprintf("Expected retry count (number) but got %s", pp.parser.currentToken.Type.String()),
		).WithContext("Parsing 'retry' property").WithHint("Retry must be a number, e.g., retry 3")
	}

	var retry int
	fmt.Sscanf(pp.parser.currentToken.Literal, "%d", &retry)
	task.Retry = retry
	pp.parser.advance()
	return nil
}

func (pp *PropertyParser) ParseRetryDelay(task *models.Task) error {
	pp.parser.advance()

	if pp.parser.currentToken.Type != TOKEN_STRING {
		return pp.parser.createError(
			fmt.Sprintf("Expected retry_delay value (string) but got %s", pp.parser.currentToken.Type.String()),
		).WithContext("Parsing 'retry_delay' property").WithHint("Retry_delay values must be strings like \"5s\", \"1m\"")
	}
	task.RetryDelay = pp.parser.currentToken.Literal
	pp.parser.advance()
	return nil
}

func (pp *PropertyParser) ParseOnSuccess(task *models.Task) error {
	pp.parser.advance()

	if err := pp.parser.expect(TOKEN_LBRACKET); err != nil {
		return err
	}

	for !pp.parser.currentToken.Is(TOKEN_RBRACKET) && !pp.parser.isAtEnd() {
		pp.parser.skipInsignificantTokens()

		if pp.parser.currentToken.Is(TOKEN_RBRACKET) {
			break
		}

		if pp.parser.currentToken.Type != TOKEN_STRING {
			return pp.parser.createError(
				fmt.Sprintf("Expected hook name (string) but got %s", pp.parser.currentToken.Type.String()),
			).WithContext("Parsing 'on_success' property").WithHint("Hook names must be strings")
		}

		task.OnSuccess = append(task.OnSuccess, pp.parser.currentToken.Literal)
		pp.parser.advance()

		if pp.parser.currentToken.Is(TOKEN_COMMA) {
			pp.parser.advance()
		}
	}

	if err := pp.parser.expect(TOKEN_RBRACKET); err != nil {
		return err
	}

	return nil
}

func (pp *PropertyParser) ParseOnFailure(task *models.Task) error {
	pp.parser.advance()

	if err := pp.parser.expect(TOKEN_LBRACKET); err != nil {
		return err
	}

	for !pp.parser.currentToken.Is(TOKEN_RBRACKET) && !pp.parser.isAtEnd() {
		pp.parser.skipInsignificantTokens()

		if pp.parser.currentToken.Is(TOKEN_RBRACKET) {
			break
		}

		if pp.parser.currentToken.Type != TOKEN_STRING {
			return pp.parser.createError(
				fmt.Sprintf("Expected hook name (string) but got %s", pp.parser.currentToken.Type.String()),
			).WithContext("Parsing 'on_failure' property").WithHint("Hook names must be strings")
		}

		task.OnFailure = append(task.OnFailure, pp.parser.currentToken.Literal)
		pp.parser.advance()

		if pp.parser.currentToken.Is(TOKEN_COMMA) {
			pp.parser.advance()
		}
	}

	if err := pp.parser.expect(TOKEN_RBRACKET); err != nil {
		return err
	}

	return nil
}

func (pp *PropertyParser) ParseArgs(task *models.Task) error {
	pp.parser.advance()

	if err := pp.parser.expect(TOKEN_LBRACE); err != nil {
		return err
	}

	task.Args = &models.TaskArgs{
		Required: []string{},
		Optional: []string{},
	}

	for !pp.parser.currentToken.Is(TOKEN_RBRACE) && !pp.parser.isAtEnd() {
		pp.parser.skipInsignificantTokens()

		if pp.parser.currentToken.Is(TOKEN_RBRACE) {
			break
		}

		if !pp.parser.currentToken.Is(TOKEN_IDENTIFIER) {
			return pp.parser.createError(
				fmt.Sprintf("Expected 'required' or 'optional' but got %s", pp.parser.currentToken.Type.String()),
			).WithContext("Parsing 'args' block").WithHint("Args block should contain 'required' and/or 'optional' lists")
		}

		keyword := pp.parser.currentToken.Literal

		switch keyword {
		case "required":
			pp.parser.advance()
			if err := pp.parser.expect(TOKEN_LBRACKET); err != nil {
				return err
			}

			for !pp.parser.currentToken.Is(TOKEN_RBRACKET) && !pp.parser.isAtEnd() {
				pp.parser.skipInsignificantTokens()

				if pp.parser.currentToken.Is(TOKEN_RBRACKET) {
					break
				}

				if pp.parser.currentToken.Type != TOKEN_STRING {
					return pp.parser.createError(
						fmt.Sprintf("Expected argument name (string) but got %s", pp.parser.currentToken.Type.String()),
					).WithContext("Parsing 'required' args").WithHint("Argument names must be strings, e.g., [\"name\", \"version\"]")
				}

				task.Args.Required = append(task.Args.Required, pp.parser.currentToken.Literal)
				pp.parser.advance()

				if pp.parser.currentToken.Is(TOKEN_COMMA) {
					pp.parser.advance()
				}
			}

			if err := pp.parser.expect(TOKEN_RBRACKET); err != nil {
				return err
			}

		case "optional":
			pp.parser.advance()
			if err := pp.parser.expect(TOKEN_LBRACKET); err != nil {
				return err
			}

			for !pp.parser.currentToken.Is(TOKEN_RBRACKET) && !pp.parser.isAtEnd() {
				pp.parser.skipInsignificantTokens()

				if pp.parser.currentToken.Is(TOKEN_RBRACKET) {
					break
				}

				if pp.parser.currentToken.Type != TOKEN_STRING {
					return pp.parser.createError(
						fmt.Sprintf("Expected argument name (string) but got %s", pp.parser.currentToken.Type.String()),
					).WithContext("Parsing 'optional' args").WithHint("Argument names must be strings, e.g., [\"region\", \"verbose\"]")
				}

				task.Args.Optional = append(task.Args.Optional, pp.parser.currentToken.Literal)
				pp.parser.advance()

				if pp.parser.currentToken.Is(TOKEN_COMMA) {
					pp.parser.advance()
				}
			}

			if err := pp.parser.expect(TOKEN_RBRACKET); err != nil {
				return err
			}

		default:
			return pp.parser.createError(
				fmt.Sprintf("Unknown args property: %s", keyword),
			).WithContext("Parsing 'args' block").WithHint("Valid properties are 'required' and 'optional'")
		}
	}

	if err := pp.parser.expect(TOKEN_RBRACE); err != nil {
		return err
	}

	return nil
}
