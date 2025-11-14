package command

import "azuyamat.dev/pace/internal/logger"

func Help() {
	logger.Println("Pace - A modern task runner")
	logger.Println()
	logger.Println("Usage:")
	logger.Println("  pace <command> [arguments]")
	logger.Println()
	logger.Println("Commands:")
	logger.Println("  run <task>         Run a specific task")
	logger.Println("    --dry-run        Show what would be executed")
	logger.Println("    --force          Force rebuild, ignoring cache")
	logger.Println("  list               List all available tasks")
	logger.Println("    --tree           Show task dependency tree")
	logger.Println("  watch <task>       Watch task inputs and rerun on changes")
	logger.Println("  <task>             Run task directly (shorthand for 'run')")
	logger.Println("  help               Show this help message")
	logger.Println()
	logger.Println("Examples:")
	logger.Println("  pace run build")
	logger.Println("  pace build")
	logger.Println("  pace run build --dry-run")
	logger.Println("  pace list --tree")
	logger.Println("  pace watch dev")
}
