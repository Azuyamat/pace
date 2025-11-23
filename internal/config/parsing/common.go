package parsing

import (
	"fmt"
)

type ParseHelper struct {
	parser *Parser
}

func NewParseHelper(parser *Parser) *ParseHelper {
	return &ParseHelper{
		parser: parser,
	}
}

func (ph *ParseHelper) ParseStringArray(contextName, hintExample string) ([]string, error) {
	if err := ph.parser.expect(TOKEN_LBRACKET); err != nil {
		return nil, err
	}

	result := make([]string, 0)

	for !ph.parser.currentToken.Is(TOKEN_RBRACKET) && !ph.parser.isAtEnd() {
		ph.parser.skipInsignificantTokens()

		if ph.parser.currentToken.Is(TOKEN_RBRACKET) {
			break
		}

		if !ph.parser.currentToken.IsOneOf(TOKEN_STRING, TOKEN_IDENTIFIER) {
			return nil, ph.parser.createError(
				fmt.Sprintf("Expected string or identifier value but got %s", ph.parser.currentToken.Type.String()),
			).WithContext(contextName).WithHint(hintExample)
		}

		result = append(result, ph.parser.currentToken.Literal)
		ph.parser.advance()

		if ph.parser.currentToken.Is(TOKEN_COMMA) {
			ph.parser.advance()
		}
	}

	if err := ph.parser.expect(TOKEN_RBRACKET); err != nil {
		return nil, err
	}

	return result, nil
}

func (ph *ParseHelper) ParseStringMap(keyContext, valueContext string) (map[string]string, error) {
	if err := ph.parser.expect(TOKEN_LBRACE); err != nil {
		return nil, err
	}

	result := make(map[string]string)

	for !ph.parser.currentToken.Is(TOKEN_RBRACE) && !ph.parser.isAtEnd() {
		ph.parser.skipInsignificantTokens()

		if ph.parser.currentToken.Is(TOKEN_RBRACE) {
			break
		}

		if !ph.parser.currentToken.IsOneOf(TOKEN_STRING, TOKEN_IDENTIFIER) {
			return nil, ph.parser.createError(
				fmt.Sprintf("Expected %s (string or identifier) but got %s", keyContext, ph.parser.currentToken.Type.String()),
			)
		}
		key := ph.parser.currentToken.Literal
		ph.parser.advance()

		if err := ph.parser.expect(TOKEN_EQUALS); err != nil {
			return nil, ph.parser.createError(
				fmt.Sprintf("Expected '=' after %s", keyContext),
			).WithHint("Use format: KEY=value")
		}

		if !ph.parser.currentToken.IsOneOf(TOKEN_STRING, TOKEN_IDENTIFIER, TOKEN_BOOLEAN) {
			return nil, ph.parser.createError(
				fmt.Sprintf("Expected %s (string, identifier, or boolean) but got %s", valueContext, ph.parser.currentToken.Type.String()),
			)
		}
		value := ph.parser.currentToken.Literal
		result[key] = value
		ph.parser.advance()
	}

	if err := ph.parser.expect(TOKEN_RBRACE); err != nil {
		return nil, err
	}

	return result, nil
}

func (ph *ParseHelper) ParseBoolean(propertyName string) (bool, error) {
	if ph.parser.currentToken.IsTrue() {
		ph.parser.advance()
		return true, nil
	} else if ph.parser.currentToken.IsFalse() {
		ph.parser.advance()
		return false, nil
	}
	return false, ph.parser.createError(
		fmt.Sprintf("Expected boolean value for '%s' property but got %s", propertyName, ph.parser.currentToken.Type.String()),
	).WithContext(fmt.Sprintf("Parsing '%s' property", propertyName)).WithHint(fmt.Sprintf("%s property must be either true or false", propertyName))
}

func (ph *ParseHelper) ParseString(propertyName, hint string) (string, error) {
	if ph.parser.currentToken.Type != TOKEN_STRING && ph.parser.currentToken.Type != TOKEN_MULTILINE_STRING {
		return "", ph.parser.createError(
			fmt.Sprintf("Expected %s value (string) but got %s", propertyName, ph.parser.currentToken.Type.String()),
		).WithContext(fmt.Sprintf("Parsing '%s' property", propertyName)).WithHint(hint)
	}
	value := ph.parser.currentToken.Literal
	ph.parser.advance()
	return value, nil
}

func (ph *ParseHelper) ParseNumber(propertyName string) (int, error) {
	if ph.parser.currentToken.Type != TOKEN_NUMBER {
		return 0, ph.parser.createError(
			fmt.Sprintf("Expected %s value (number) but got %s", propertyName, ph.parser.currentToken.Type.String()),
		).WithContext(fmt.Sprintf("Parsing '%s' property", propertyName)).WithHint(fmt.Sprintf("%s must be a number", propertyName))
	}
	var value int
	fmt.Sscanf(ph.parser.currentToken.Literal, "%d", &value)
	ph.parser.advance()
	return value, nil
}

func isLetter(char byte) bool {
	return (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || char == '_' || char == '-'
}

func isDigit(char byte) bool {
	return char >= '0' && char <= '9'
}

func isWhitespace(char byte) bool {
	return char == ' ' || char == '\t'
}
