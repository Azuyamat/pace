package command

import (
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
	projectType := models.ProjectTypeUnknown
	if projectTypeFlag != "" {
		projectType = models.ParseProjectType(projectTypeFlag)
	} else {
		projectType = detector.DetectCurrentProjectType()
		if projectType == models.ProjectTypeUnknown {
			logger.Warning("Could not detect project type automatically.")
			typeList := detector.ListSupportedProjectTypes()
			logger.Info("Supported project types: %v", typeList)
			answer, err := logger.Prompt("Please specify the project type (or 'unknown' for default config): ")
			if err != nil {
				return err
			}
			projectType = models.ParseProjectType(answer)
			if answer != "y" {
				logger.Warning("Initialization cancelled.")
				return nil
			}
		}
	}
	if projectType == models.ProjectTypeUnknown {
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
	logger.Info("Detected project type: %s", projectType)
	generator := generator.GetGeneratorByProjectType(projectType)
	if generator == nil {
		return fmt.Errorf("unsupported project type: %s", projectType)
	}
	cfg, err := generator.Generate()
	if err != nil {
		return err
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
