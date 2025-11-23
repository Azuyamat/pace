package runner

import (
	"fmt"
	"runtime"
	"strings"
)

type ConditionEvaluator struct {
	platform string
}

func NewConditionEvaluator() *ConditionEvaluator {
	return &ConditionEvaluator{
		platform: runtime.GOOS,
	}
}

func (ce *ConditionEvaluator) Evaluate(condition string) (bool, error) {
	if condition == "" {
		return true, nil
	}

	condition = strings.TrimSpace(condition)

	if strings.Contains(condition, "==") {
		parts := strings.SplitN(condition, "==", 2)
		if len(parts) != 2 {
			return false, fmt.Errorf("invalid condition format: %s", condition)
		}
		left := strings.TrimSpace(parts[0])
		right := strings.TrimSpace(strings.Trim(parts[1], "\""))

		leftValue := ce.resolveValue(left)
		return leftValue == right, nil
	}

	if strings.Contains(condition, "!=") {
		parts := strings.SplitN(condition, "!=", 2)
		if len(parts) != 2 {
			return false, fmt.Errorf("invalid condition format: %s", condition)
		}
		left := strings.TrimSpace(parts[0])
		right := strings.TrimSpace(strings.Trim(parts[1], "\""))

		leftValue := ce.resolveValue(left)
		return leftValue != right, nil
	}

	if strings.Contains(condition, " in ") {
		parts := strings.SplitN(condition, " in ", 2)
		if len(parts) != 2 {
			return false, fmt.Errorf("invalid condition format: %s", condition)
		}
		left := strings.TrimSpace(parts[0])
		right := strings.TrimSpace(parts[1])

		leftValue := ce.resolveValue(left)

		right = strings.Trim(right, "[]")
		values := strings.Split(right, ",")
		for _, v := range values {
			v = strings.TrimSpace(strings.Trim(v, "\""))
			if leftValue == v {
				return true, nil
			}
		}
		return false, nil
	}

	return false, fmt.Errorf("unsupported condition format: %s", condition)
}

func (ce *ConditionEvaluator) resolveValue(varName string) string {
	switch varName {
	case "platform":
		return ce.platform
	default:
		return varName
	}
}
