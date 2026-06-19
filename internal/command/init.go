package command

import (
	"encoding/json"
	"fmt"
	"os"

	gear "github.com/azuyamat/gear/command"
	"github.com/azuyamat/pace/internal/config"
	"github.com/azuyamat/pace/internal/logger"
	"github.com/azuyamat/pace/internal/models"
	"github.com/azuyamat/pace/internal/template/detector"
	"github.com/azuyamat/pace/internal/template/generator"
)

var initCommand = gear.NewExecutableCommand("init", "Initialize a new Pace project in the current directory").
	Flags(
		gear.NewStringFlag("type", "t", "Specify the project type (go, python, rust, etc.)", ""),
		gear.NewBoolFlag("yes", "y", "Run non-interactively and accept defaults", false),
		gear.NewBoolFlag("force", "f", "Overwrite config.pace if it already exists", false),
		gear.NewBoolFlag("stdout", "s", "Print generated config to stdout instead of writing config.pace", false),
		gear.NewBoolFlag("json", "j", "With --stdout, print generated config as JSON", false),
	).
	Handler(initHandler)

func init() {
	RootCommand.AddChild(initCommand)
}

func initHandler(ctx *gear.Context, args gear.ValidatedArgs) error {
	projectTypeFlag, err := args.GetFlagString("type")
	if err != nil {
		return err
	}
	yes := args.FlagBool("yes")
	force := args.FlagBool("force")
	printStdout := args.FlagBool("stdout")
	printJSON := args.FlagBool("json")
	if printJSON {
		printStdout = true
	}
	projectType := models.ProjectTypeUnknown
	if projectTypeFlag != "" {
		projectType = models.ParseProjectType(projectTypeFlag)
	} else {
		projectType = detector.DetectCurrentProjectType()
		if projectType == models.ProjectTypeUnknown && !yes {
			logger.Warning("Could not detect project type automatically.")
			typeList := detector.ListSupportedProjectTypes()
			logger.Info("Supported project types: %v", typeList)
			answer, err := logger.Prompt("Please specify the project type (or 'unknown' for default config): ")
			if err != nil {
				return err
			}
			projectType = models.ParseProjectType(answer)
		}
	}
	if projectType == models.ProjectTypeUnknown && !yes {
		logger.Info("Detected project type is unknown and would generate a default config.")
		answer, err := logger.Prompt("Continue with default config? (y/n): ")
		if err != nil {
			return err
		}
		if answer != "y" {
			logger.Warning("Initialization cancelled.")
			return nil
		}
	}
	if printStdout {
		logger.Default.SetEnabled(false)
	} else if _, err := os.Stat("config.pace"); err == nil && !force {
		return fmt.Errorf("config.pace already exists; use --force to overwrite")
	} else if err != nil && !os.IsNotExist(err) {
		return err
	}
	logger.Info("Detected project type: %s", projectType)
	generator := generator.GetGeneratorByProjectType(projectType)
	if generator == nil {
		return fmt.Errorf("unsupported project type: %s", projectType)
	}
	cfg, err := generator.Generate()
	if err != nil {
		return err
	}
	if printStdout {
		if printJSON {
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			return encoder.Encode(cfg)
		}
		fmt.Print(cfg.String())
		return nil
	}
	err = cfg.WriteToFile("config.pace")
	if err != nil {
		return err
	}
	logger.Info("Generated config.pace")

	cwd, err := os.Getwd()
	if err != nil {
		logger.Warning("Failed to get current directory: %v", err)
	} else {
		err = config.UpdateGitignore(cwd)
		if err != nil {
			logger.Warning("Failed to update .gitignore: %v", err)
		} else {
			logger.Info("Updated .gitignore to exclude .pace-cache/")
		}
	}

	return nil
}
