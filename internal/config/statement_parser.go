package config

import (
	"fmt"
)

func (p *Parser) parseTopLevelStatement(config *Config) error {
	if p.currentToken.IsKeyword("task") {
		task, err := p.parseTask()
		if err != nil {
			return err
		}
		config.Tasks[task.Name] = task
		return nil
	}

	if p.currentToken.IsKeyword("hook") {
		hook, err := p.parseHook()
		if err != nil {
			return err
		}
		config.Hooks[hook.Name] = hook
		return nil
	}

	if p.currentToken.IsKeyword("globals") {
		return p.parseGlobal(config)
	}

	if p.currentToken.IsKeyword("set") {
		return p.parseSet(config)
	}

	if p.currentToken.IsKeyword("default") {
		return p.parseDefault(config)
	}

	if p.currentToken.IsKeyword("alias") {
		return p.parseAlias(config)
	}

	if p.currentToken.IsKeyword("import") {
		return p.parseImport(config)
	}

	return p.unexpectedTokenError("top-level statement")
}

func (p *Parser) parseGlobal(config *Config) error {
	p.advance()

	if err := p.expect(TOKEN_LBRACE); err != nil {
		return err
	}

	for !p.currentToken.Is(TOKEN_RBRACE) && !p.isAtEnd() {
		p.skipInsignificantTokens()

		if p.currentToken.Is(TOKEN_RBRACE) {
			break
		}

		if p.currentToken.Type != TOKEN_STRING {
			return p.createError(
				fmt.Sprintf("Expected global variable name (string) but got %s", p.currentToken.Type.String()),
			).WithContext("Parsing globals block").WithHint("Global names must be strings")
		}
		name := p.currentToken.Literal
		p.advance()

		if p.currentToken.Type != TOKEN_STRING {
			return p.createError(
				fmt.Sprintf("Expected global variable value (string) but got %s", p.currentToken.Type.String()),
			).WithContext("Parsing globals block").WithHint("Global values must be strings")
		}
		value := p.currentToken.Literal
		p.advance()

		config.Globals[name] = value
	}

	if err := p.expect(TOKEN_RBRACE); err != nil {
		return err
	}

	return nil
}

func (p *Parser) parseSet(config *Config) error {
	p.advance()

	if p.currentToken.Type != TOKEN_IDENTIFIER && p.currentToken.Type != TOKEN_STRING {
		return p.createError(
			fmt.Sprintf("Expected constant name but got %s", p.currentToken.Type.String()),
		).WithContext("Parsing set statement").WithHint("Constant names can be identifiers or strings, e.g., set BUILD_DIR \"./build\"")
	}
	name := p.currentToken.Literal
	p.advance()

	if p.currentToken.Type != TOKEN_STRING {
		return p.createError(
			fmt.Sprintf("Expected constant value (string) but got %s", p.currentToken.Type.String()),
		).WithContext("Parsing set statement").WithHint("Constant values must be strings")
	}
	value := p.currentToken.Literal
	p.advance()

	config.Constants[name] = value
	return nil
}

func (p *Parser) parseDefault(config *Config) error {
	p.advance()

	if p.currentToken.Type != TOKEN_STRING {
		return p.createError(
			fmt.Sprintf("Expected default task name (string) but got %s", p.currentToken.Type.String()),
		).WithContext("Parsing default statement").WithHint("Default task name must be a string, e.g., default \"build\"")
	}
	config.DefaultTask = p.currentToken.Literal
	p.advance()
	return nil
}

func (p *Parser) parseAlias(config *Config) error {
	p.advance()

	if p.currentToken.Type != TOKEN_STRING {
		return p.createError(
			fmt.Sprintf("Expected alias name (string) but got %s", p.currentToken.Type.String()),
		).WithContext("Parsing alias statement").WithHint("Alias names must be strings, e.g., alias \"b\" \"build\"")
	}
	aliasName := p.currentToken.Literal
	p.advance()

	if p.currentToken.Type != TOKEN_STRING {
		return p.createError(
			fmt.Sprintf("Expected task name (string) but got %s", p.currentToken.Type.String()),
		).WithContext("Parsing alias statement").WithHint("Task names must be strings, e.g., alias \"b\" \"build\"")
	}
	taskName := p.currentToken.Literal
	p.advance()

	config.Aliases[aliasName] = taskName
	return nil
}

func (p *Parser) parseImport(config *Config) error {
	p.advance()

	if p.currentToken.Type != TOKEN_STRING {
		return p.createError(
			fmt.Sprintf("Expected import path (string) but got %s", p.currentToken.Type.String()),
		).WithContext("Parsing import statement").WithHint("Import paths must be strings, e.g., import \"tasks/build.pace\"")
	}
	importPath := p.currentToken.Literal
	p.advance()

	config.Imports = append(config.Imports, importPath)
	return nil
}
