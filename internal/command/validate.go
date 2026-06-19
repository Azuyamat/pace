package command

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	gear "github.com/azuyamat/gear/command"
	"github.com/azuyamat/pace/internal/config"
	"github.com/azuyamat/pace/internal/config/parsing"
	"github.com/azuyamat/pace/internal/logger"
)

type validateDiagnostic struct {
	File     string        `json:"file"`
	Range    validateRange `json:"range"`
	Severity string        `json:"severity"`
	Message  string        `json:"message"`
	Hint     string        `json:"hint,omitempty"`
}

type validateRange struct {
	Start validatePosition `json:"start"`
	End   validatePosition `json:"end"`
}

type validatePosition struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

var validateCommand = gear.NewExecutableCommand("validate", "Validate a Pace config without running tasks").
	Args(gear.NewStringArg("file", "Path to the Pace config file").AsOptional()).
	Flags(gear.NewBoolFlag("json", "j", "Output machine-readable JSON diagnostics", false)).
	Handler(validateHandler)

func init() {
	RootCommand.AddChild(validateCommand)
}

func validateHandler(ctx *gear.Context, args gear.ValidatedArgs) error {
	path := args.String("file")
	if path == "" {
		path = config.ConfigFile
	}

	_, err := config.ParseFile(path)
	if err == nil {
		if args.FlagBool("json") {
			return writeValidationJSON([]validateDiagnostic{})
		}
		fmt.Fprintf(os.Stdout, "%s is valid\n", path)
		return nil
	}

	diagnostics := diagnosticsFromError(path, err)
	if args.FlagBool("json") {
		if jsonErr := writeValidationJSON(diagnostics); jsonErr != nil {
			return jsonErr
		}
		logger.Default.SetEnabled(false)
		return fmt.Errorf("validation failed")
	}

	return err
}

func diagnosticsFromError(path string, err error) []validateDiagnostic {
	var parseErr *parsing.ParseError
	if errors.As(err, &parseErr) {
		line := parseErr.Line
		column := parseErr.Column
		if line < 1 {
			line = 1
		}
		if column < 1 {
			column = 1
		}
		return []validateDiagnostic{{
			File:     path,
			Range:    rangeAt(line, column),
			Severity: "error",
			Message:  parseErr.Message,
			Hint:     parseErr.Hint,
		}}
	}

	return []validateDiagnostic{{
		File:     path,
		Range:    rangeAt(1, 1),
		Severity: "error",
		Message:  err.Error(),
	}}
}

func rangeAt(line, column int) validateRange {
	return validateRange{
		Start: validatePosition{Line: line, Column: column},
		End:   validatePosition{Line: line, Column: column + 1},
	}
}

func writeValidationJSON(diagnostics []validateDiagnostic) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(diagnostics)
}
