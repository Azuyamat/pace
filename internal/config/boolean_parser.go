package config

import (
	"fmt"

	"github.com/azuyamat/pace/internal/models"
)

func (pp *PropertyParser) ParseCache(task *models.Task) error {
	pp.parser.advance()

	if pp.parser.currentToken.IsTrue() {
		task.Cache = true
		pp.parser.advance()
	} else if pp.parser.currentToken.IsFalse() {
		task.Cache = false
		pp.parser.advance()
	} else {
		return pp.parser.createError(
			fmt.Sprintf("Expected boolean value for 'cache' property but got %s", pp.parser.currentToken.Type.String()),
		).WithContext("Parsing 'cache' property").WithHint("Cache property must be either true or false")
	}
	return nil
}

func (pp *PropertyParser) ParseWatch(task *models.Task) error {
	pp.parser.advance()

	if pp.parser.currentToken.IsTrue() {
		task.Watch = true
		pp.parser.advance()
	} else if pp.parser.currentToken.IsFalse() {
		task.Watch = false
		pp.parser.advance()
	} else {
		return pp.parser.createError(
			fmt.Sprintf("Expected boolean value for 'watch' property but got %s", pp.parser.currentToken.Type.String()),
		).WithContext("Parsing 'watch' property").WithHint("Watch property must be either true or false")
	}
	return nil
}

func (pp *PropertyParser) ParseParallel(task *models.Task) error {
	pp.parser.advance()

	if pp.parser.currentToken.IsTrue() {
		task.Parallel = true
		pp.parser.advance()
	} else if pp.parser.currentToken.IsFalse() {
		task.Parallel = false
		pp.parser.advance()
	} else {
		return pp.parser.createError(
			fmt.Sprintf("Expected boolean value for 'parallel' property but got %s", pp.parser.currentToken.Type.String()),
		).WithContext("Parsing 'parallel' property").WithHint("Parallel property must be either true or false")
	}
	return nil
}

func (pp *PropertyParser) ParseSilent(task *models.Task) error {
	pp.parser.advance()

	if pp.parser.currentToken.IsTrue() {
		task.Silent = true
		pp.parser.advance()
	} else if pp.parser.currentToken.IsFalse() {
		task.Silent = false
		pp.parser.advance()
	} else {
		return pp.parser.createError(
			fmt.Sprintf("Expected boolean value for 'silent' property but got %s", pp.parser.currentToken.Type.String()),
		).WithContext("Parsing 'silent' property").WithHint("Silent property must be either true or false")
	}
	return nil
}

func (pp *PropertyParser) ParseContinueOnError(task *models.Task) error {
	pp.parser.advance()

	if pp.parser.currentToken.IsTrue() {
		task.ContinueOnError = true
		pp.parser.advance()
	} else if pp.parser.currentToken.IsFalse() {
		task.ContinueOnError = false
		pp.parser.advance()
	} else {
		return pp.parser.createError(
			fmt.Sprintf("Expected boolean value for 'continue_on_error' property but got %s", pp.parser.currentToken.Type.String()),
		).WithContext("Parsing 'continue_on_error' property").WithHint("Continue_on_error property must be either true or false")
	}
	return nil
}
