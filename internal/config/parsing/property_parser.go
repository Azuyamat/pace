package parsing

import (
	"fmt"
	"reflect"

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

var taskPropertyRegistry = map[string]PropertyDefinition{
	"command":           prop(PropString, "Command", "Command values must be strings, e.g., command \"echo hello\""),
	"inputs":            prop(PropStringArray, "Inputs", "Input values must be strings, e.g., [\"src/main.go\", \"src/util.go\"]"),
	"outputs":           prop(PropStringArray, "Outputs", "Output values must be strings, e.g., [\"bin/app\", \"bin/util\"]"),
	"depends-on":        prop(PropStringArray, "DependsOn", "Task names must be strings, e.g., [build, test]"),
	"dependencies":      prop(PropStringArray, "DependsOn", "Task names must be strings, e.g., [build, test]"),
	"env":               prop(PropStringMap, "Env", ""),
	"cache":             prop(PropBoolean, "Cache", ""),
	"working_dir":       prop(PropString, "WorkingDir", "Working directory value must be a string, e.g., \"/app\""),
	"requires":          prop(PropStringArray, "Requires", "Hook names must be strings, e.g., [setup, clean]"),
	"before":            prop(PropStringArray, "Requires", "Hook names must be strings, e.g., [setup, clean]"),
	"triggers":          prop(PropStringArray, "Triggers", "Hook names must be strings, e.g., [cleanup, notify]"),
	"after":             prop(PropStringArray, "Triggers", "Hook names must be strings, e.g., [cleanup, notify]"),
	"description":       prop(PropString, "Description", "Description values must be strings"),
	"watch":             prop(PropBoolean, "Watch", ""),
	"parallel":          prop(PropBoolean, "Parallel", ""),
	"silent":            prop(PropBoolean, "Silent", ""),
	"continue_on_error": prop(PropBoolean, "ContinueOnError", ""),
	"timeout":           prop(PropString, "Timeout", "Timeout values must be strings like \"5m\", \"30s\", \"1h\""),
	"retry":             prop(PropNumber, "Retry", ""),
	"retry_delay":       prop(PropString, "RetryDelay", "Retry_delay values must be strings like \"5s\", \"1m\""),
	"on_success":        prop(PropStringArray, "OnSuccess", "Hook names must be strings"),
	"on_failure":        prop(PropStringArray, "OnFailure", "Hook names must be strings"),
	"when":              prop(PropString, "When", "Condition value must be a string"),
	"args": {
		Type:         PropCustom,
		CustomParser: (*PropertyParser).parseArgs,
	},
}

var hookPropertyRegistry = map[string]PropertyDefinition{
	"command":     hookProp(PropString, "Command", "Command values must be strings, e.g., command \"echo setup\""),
	"env":         hookProp(PropStringMap, "Env", ""),
	"working_dir": hookProp(PropString, "WorkingDir", "Working directory value must be a string"),
	"description": hookProp(PropString, "Description", "Description values must be strings"),
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
