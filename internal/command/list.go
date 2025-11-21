package command

import (
	"sort"

	gear "github.com/azuyamat/gear/command"
	"github.com/azuyamat/pace/internal/config"
	"github.com/azuyamat/pace/internal/logger"
)

var listCommand = gear.NewExecutableCommand("list", "List all available tasks and their details").
	Flags(
		gear.NewBoolFlag("tree", "t", "Display tasks in a tree view showing dependencies", false)).
	Handler(listHandler)

func init() {
	RootCommand.AddChild(listCommand)
}

func listHandler(ctx *gear.Context, args gear.ValidatedArgs) error {
	config, err := config.GetConfig()
	if err != nil {
		return err
	}
	treeView := args.FlagBool("tree")

	if treeView {
		printTaskTree(config)
	} else {
		printTaskList(config)
	}

	return nil
}

func printTaskList(cfg *config.Config) {
	logger.Println("Available tasks:")
	logger.Println()

	taskNames := make([]string, 0, len(cfg.Tasks))
	for name := range cfg.Tasks {
		taskNames = append(taskNames, name)
	}
	sort.Strings(taskNames)

	for _, name := range taskNames {
		task := cfg.Tasks[name]
		defaultMarker := ""
		if cfg.DefaultTask == name {
			defaultMarker = " (default)"
		}

		if task.Description != "" {
			logger.Printf("  %-20s %s%s\n", name, task.Description, defaultMarker)
		} else {
			logger.Printf("  %-20s %s%s\n", name, task.Command, defaultMarker)
		}
	}

	if len(cfg.Aliases) > 0 {
		logger.Println("\nAliases:")
		aliasNames := make([]string, 0, len(cfg.Aliases))
		for alias := range cfg.Aliases {
			aliasNames = append(aliasNames, alias)
		}
		sort.Strings(aliasNames)

		for _, alias := range aliasNames {
			logger.Printf("  %-20s -> %s\n", alias, cfg.Aliases[alias])
		}
	}

	if len(cfg.Hooks) > 0 {
		logger.Println("\nAvailable hooks:")
		hookNames := make([]string, 0, len(cfg.Hooks))
		for name := range cfg.Hooks {
			hookNames = append(hookNames, name)
		}
		sort.Strings(hookNames)

		for _, name := range hookNames {
			hook := cfg.Hooks[name]
			if hook.Description != "" {
				logger.Printf("  %-20s %s\n", name, hook.Description)
			} else {
				logger.Printf("  %-20s %s\n", name, hook.Command)
			}
		}
	}
}

func printTaskTree(cfg *config.Config) {
	logger.Println("Task dependency tree:")
	logger.Println()

	taskNames := make([]string, 0, len(cfg.Tasks))
	for name := range cfg.Tasks {
		taskNames = append(taskNames, name)
	}
	sort.Strings(taskNames)

	visited := make(map[string]bool)

	for _, name := range taskNames {
		if !visited[name] {
			printTaskNode(cfg, name, "", visited, make(map[string]bool))
		}
	}
}

func printTaskNode(cfg *config.Config, taskName string, prefix string, visited map[string]bool, ancestry map[string]bool) {
	task, exists := cfg.Tasks[taskName]
	if !exists {
		return
	}

	defaultMarker := ""
	if cfg.DefaultTask == taskName {
		defaultMarker = " (default)"
	}

	logger.Printf("%s%s%s\n", prefix, taskName, defaultMarker)
	visited[taskName] = true

	if len(task.Dependencies) > 0 {
		ancestry[taskName] = true
		for i, dep := range task.Dependencies {
			isLast := i == len(task.Dependencies)-1
			var newPrefix string
			if isLast {
				logger.Printf("%s  └── ", prefix)
				newPrefix = prefix + "      "
			} else {
				logger.Printf("%s  ├── ", prefix)
				newPrefix = prefix + "  │   "
			}

			if ancestry[dep] {
				logger.Printf("%s (circular)\n", dep)
			} else {
				printTaskNode(cfg, dep, newPrefix, visited, ancestry)
			}
		}
		delete(ancestry, taskName)
	}
}
