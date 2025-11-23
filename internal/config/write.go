package config

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/azuyamat/pace/internal/models"
)

func (c *Config) WriteToFile(path string) error {
	content := c.String()
	return os.WriteFile(path, []byte(content), 0644)
}

func (c *Config) String() string {
	var builder strings.Builder

	if c.DefaultTask != "" {
		builder.WriteString(fmt.Sprintf("default %s\n", c.DefaultTask))
	}

	if len(c.Imports) > 0 {
		if builder.Len() > 0 {
			builder.WriteString("\n")
		}
		for _, imp := range c.Imports {
			builder.WriteString(fmt.Sprintf("import \"%s\"\n", imp))
		}
	}

	if len(c.Constants) > 0 {
		if builder.Len() > 0 {
			builder.WriteString("\n")
		}
		keys := sortedKeys(c.Constants)
		for _, key := range keys {
			builder.WriteString(fmt.Sprintf("var %s = \"%s\"\n", key, c.Constants[key]))
		}
	}

	if len(c.Globals) > 0 {
		if builder.Len() > 0 {
			builder.WriteString("\n")
		}
		builder.WriteString("globals {\n")
		keys := sortedKeys(c.Globals)
		for _, key := range keys {
			builder.WriteString(fmt.Sprintf("    \"%s\" \"%s\"\n", key, c.Globals[key]))
		}
		builder.WriteString("}\n")
	}

	if len(c.Aliases) > 0 {
		if builder.Len() > 0 {
			builder.WriteString("\n")
		}
		keys := sortedKeys(c.Aliases)
		for _, key := range keys {
			builder.WriteString(fmt.Sprintf("alias %s %s\n", key, c.Aliases[key]))
		}
	}

	if len(c.Tasks) > 0 {
		keys := sortedKeys(c.Tasks)
		for _, name := range keys {
			if builder.Len() > 0 {
				builder.WriteString("\n")
			}
			builder.WriteString(taskString(c.Tasks[name]))
		}
	}

	if len(c.Hooks) > 0 {
		keys := sortedKeys(c.Hooks)
		for _, name := range keys {
			if builder.Len() > 0 {
				builder.WriteString("\n")
			}
			builder.WriteString(hookString(c.Hooks[name]))
		}
	}

	return builder.String()
}

func taskString(task models.Task) string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("task %s {\n", task.Name))

	if task.Command != "" {
		builder.WriteString(fmt.Sprintf("    command \"%s\"\n", task.Command))
	}

	if task.Description != "" {
		builder.WriteString(fmt.Sprintf("    description \"%s\"\n", task.Description))
	}

	if task.WorkingDir != "" {
		builder.WriteString(fmt.Sprintf("    working_dir \"%s\"\n", task.WorkingDir))
	}

	if len(task.Inputs) > 0 {
		builder.WriteString(fmt.Sprintf("    inputs %s\n", formatStringSlice(task.Inputs)))
	}

	if len(task.Outputs) > 0 {
		builder.WriteString(fmt.Sprintf("    outputs %s\n", formatStringSlice(task.Outputs)))
	}

	if len(task.Dependencies) > 0 {
		builder.WriteString(fmt.Sprintf("    dependencies %s\n", formatStringSlice(task.Dependencies)))
	}

	if len(task.BeforeHooks) > 0 {
		builder.WriteString(fmt.Sprintf("    before %s\n", formatStringSlice(task.BeforeHooks)))
	}

	if len(task.AfterHooks) > 0 {
		builder.WriteString(fmt.Sprintf("    after %s\n", formatStringSlice(task.AfterHooks)))
	}

	if len(task.OnSuccess) > 0 {
		builder.WriteString(fmt.Sprintf("    on_success %s\n", formatStringSlice(task.OnSuccess)))
	}

	if len(task.OnFailure) > 0 {
		builder.WriteString(fmt.Sprintf("    on_failure %s\n", formatStringSlice(task.OnFailure)))
	}

	if len(task.Env) > 0 {
		builder.WriteString(fmt.Sprintf("    env %s\n", formatStringMap(task.Env)))
	}

	if task.Cache {
		builder.WriteString("    cache true\n")
	}

	if task.Watch {
		builder.WriteString("    watch true\n")
	}

	if task.Parallel {
		builder.WriteString("    parallel true\n")
	}

	if task.Silent {
		builder.WriteString("    silent true\n")
	}

	if task.ContinueOnError {
		builder.WriteString("    continue_on_error true\n")
	}

	if task.Timeout != "" {
		builder.WriteString(fmt.Sprintf("    timeout \"%s\"\n", task.Timeout))
	}

	if task.Retry > 0 {
		builder.WriteString(fmt.Sprintf("    retry %d\n", task.Retry))
	}

	if task.RetryDelay != "" {
		builder.WriteString(fmt.Sprintf("    retry_delay \"%s\"\n", task.RetryDelay))
	}

	if task.Args != nil {
		if len(task.Args.Required) > 0 {
			builder.WriteString(fmt.Sprintf("    args.required %s\n", formatStringSlice(task.Args.Required)))
		}
		if len(task.Args.Optional) > 0 {
			builder.WriteString(fmt.Sprintf("    args.optional %s\n", formatStringSlice(task.Args.Optional)))
		}
	}

	builder.WriteString("}\n")
	return builder.String()
}

func hookString(hook models.Hook) string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("hook %s {\n", hook.Name))

	if hook.Command != "" {
		builder.WriteString(fmt.Sprintf("    command \"%s\"\n", hook.Command))
	}

	if hook.Description != "" {
		builder.WriteString(fmt.Sprintf("    description \"%s\"\n", hook.Description))
	}

	if hook.WorkingDir != "" {
		builder.WriteString(fmt.Sprintf("    working_dir \"%s\"\n", hook.WorkingDir))
	}

	if len(hook.Env) > 0 {
		builder.WriteString(fmt.Sprintf("    env %s\n", formatStringMap(hook.Env)))
	}

	builder.WriteString("}\n")
	return builder.String()
}

func formatStringSlice(items []string) string {
	if len(items) == 0 {
		return "[]"
	}
	quoted := make([]string, len(items))
	for i, item := range items {
		quoted[i] = fmt.Sprintf("\"%s\"", escapeString(item))
	}
	return "[" + strings.Join(quoted, ", ") + "]"
}

func formatStringMap(m map[string]string) string {
	if len(m) == 0 {
		return "{}"
	}
	var pairs []string
	keys := sortedKeys(m)
	for _, key := range keys {
		pairs = append(pairs, fmt.Sprintf("%s: \"%s\"", key, escapeString(m[key])))
	}
	return "{" + strings.Join(pairs, ", ") + "}"
}

func escapeString(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\t", "\\t")
	return s
}

func sortedKeys[T any](m map[string]T) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func UpdateGitignore(projectPath string) error {
	gitignorePath := fmt.Sprintf("%s/.gitignore", projectPath)

	content, err := os.ReadFile(gitignorePath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	gitignoreContent := string(content)
	cacheEntry := ".pace-cache/"

	if strings.Contains(gitignoreContent, cacheEntry) {
		return nil
	}

	if len(gitignoreContent) > 0 && !strings.HasSuffix(gitignoreContent, "\n") {
		gitignoreContent += "\n"
	}

	gitignoreContent += cacheEntry + "\n"

	return os.WriteFile(gitignorePath, []byte(gitignoreContent), 0644)
}
