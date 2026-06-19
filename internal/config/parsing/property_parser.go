package parsing

import (
	"fmt"
	"reflect"

	"github.com/azuyamat/pace/internal/config/schema"
	"github.com/azuyamat/pace/internal/models"
)

type PropertyType int

const (
	PropString PropertyType = iota
	PropStringArray
	PropStringMap
	PropBoolean
	PropNumber
	PropCustom
)

type PropertyDefinition struct {
	Type         PropertyType
	TaskField    string
	HookField    string
	Hint         string
	CustomParser func(pp *PropertyParser, task *models.Task) error
}

func prop(propType PropertyType, taskField, hint string) PropertyDefinition {
	return PropertyDefinition{Type: propType, TaskField: taskField, Hint: hint}
}

func hookProp(propType PropertyType, hookField, hint string) PropertyDefinition {
	return PropertyDefinition{Type: propType, HookField: hookField, Hint: hint}
}

var taskPropertyRegistry = buildTaskPropertyRegistry(schema.Get().TaskProps)
var hookPropertyRegistry = buildHookPropertyRegistry(schema.Get().HookProps)

func buildTaskPropertyRegistry(properties []schema.Property) map[string]PropertyDefinition {
	registry := make(map[string]PropertyDefinition, len(properties))
	for _, property := range properties {
		if property.Type == "argsBlock" {
			registry[property.Name] = PropertyDefinition{Type: PropCustom, CustomParser: (*PropertyParser).parseArgs}
			continue
		}
		registry[property.Name] = prop(schemaPropType(property.Type), property.Field, hintForProperty(property))
	}
	return registry
}

func buildHookPropertyRegistry(properties []schema.Property) map[string]PropertyDefinition {
	registry := make(map[string]PropertyDefinition, len(properties))
	for _, property := range properties {
		registry[property.Name] = hookProp(schemaPropType(property.Type), property.Field, hintForProperty(property))
	}
	return registry
}

func schemaPropType(propType string) PropertyType {
	switch propType {
	case "string":
		return PropString
	case "stringArray":
		return PropStringArray
	case "stringMap":
		return PropStringMap
	case "boolean":
		return PropBoolean
	case "number":
		return PropNumber
	default:
		return PropCustom
	}
}

func hintForProperty(property schema.Property) string {
	switch property.Type {
	case "string":
		return fmt.Sprintf("%s values must be strings", property.Name)
	case "stringArray":
		return fmt.Sprintf("%s values must be strings or identifiers in an array", property.Name)
	default:
		return ""
	}
}

type PropertyParser struct {
	parser *Parser
}

func NewPropertyParser(parser *Parser) *PropertyParser {
	return &PropertyParser{
		parser: parser,
	}
}

func (pp *PropertyParser) ParseTaskProperty(task *models.Task) error {
	if !pp.parser.currentToken.Is(TOKEN_IDENTIFIER) {
		return pp.parser.createError(
			fmt.Sprintf("Expected task property name but got %s", pp.parser.currentToken.Type.String()),
		).WithContext("Parsing task body").WithHint("Valid properties include: command, inputs, outputs, dependencies, etc.")
	}

	propName := pp.parser.currentToken.Literal
	propDef, exists := taskPropertyRegistry[propName]

	if !exists {
		pp.parser.advance()
		return nil
	}

	pp.parser.advance()

	if propDef.Type == PropCustom {
		return propDef.CustomParser(pp, task)
	}

	value, err := pp.parseByType(propDef.Type, propName, propDef.Hint)
	if err != nil {
		return err
	}

	return pp.setFieldValue(reflect.ValueOf(task).Elem(), propDef.TaskField, value)
}

func (pp *PropertyParser) ParseHookProperty(hook *models.Hook) error {
	if !pp.parser.currentToken.Is(TOKEN_IDENTIFIER) {
		return pp.parser.createError(
			fmt.Sprintf("Expected hook property name but got %s", pp.parser.currentToken.Type.String()),
		).WithContext("Parsing hook body").WithHint("Valid properties include: command, env, working_dir")
	}

	propName := pp.parser.currentToken.Literal
	propDef, exists := hookPropertyRegistry[propName]

	if !exists {
		pp.parser.advance()
		return nil
	}

	pp.parser.advance()

	value, err := pp.parseByType(propDef.Type, propName, propDef.Hint)
	if err != nil {
		return err
	}

	return pp.setFieldValue(reflect.ValueOf(hook).Elem(), propDef.HookField, value)
}

func (pp *PropertyParser) setFieldValue(structValue reflect.Value, fieldName string, value any) error {
	field := structValue.FieldByName(fieldName)
	if !field.IsValid() {
		return fmt.Errorf("field %s not found", fieldName)
	}
	if !field.CanSet() {
		return fmt.Errorf("field %s cannot be set", fieldName)
	}

	fieldValue := reflect.ValueOf(value)
	if field.Type() != fieldValue.Type() {
		return fmt.Errorf("type mismatch for field %s: expected %s, got %s", fieldName, field.Type(), fieldValue.Type())
	}

	field.Set(fieldValue)
	return nil
}

func (pp *PropertyParser) parseByType(propType PropertyType, propName, hint string) (any, error) {
	switch propType {
	case PropString:
		return pp.parser.helper.ParseString(propName, hint)
	case PropStringArray:
		return pp.parser.helper.ParseStringArray(
			fmt.Sprintf("Parsing '%s' property", propName),
			hint,
		)
	case PropStringMap:
		return pp.parser.helper.ParseStringMap("environment variable name", "environment variable value")
	case PropBoolean:
		return pp.parser.helper.ParseBoolean(propName)
	case PropNumber:
		return pp.parser.helper.ParseNumber(propName)
	default:
		return nil, fmt.Errorf("unknown property type")
	}
}

func (pp *PropertyParser) parseArgs(task *models.Task) error {
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
			required, err := pp.parser.helper.ParseStringArray("Parsing 'required' args", "Argument names must be strings, e.g., [\"name\", \"version\"]")
			if err != nil {
				return err
			}
			task.Args.Required = required

		case "optional":
			pp.parser.advance()
			optional, err := pp.parser.helper.ParseStringArray("Parsing 'optional' args", "Argument names must be strings, e.g., [\"region\", \"verbose\"]")
			if err != nil {
				return err
			}
			task.Args.Optional = optional

		default:
			return pp.parser.createError(
				fmt.Sprintf("Unknown args property: %s", keyword),
			).WithContext("Parsing 'args' block").WithHint("Valid properties are 'required' and 'optional'")
		}
	}

	return pp.parser.expect(TOKEN_RBRACE)
}
