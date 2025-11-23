package command

import (
	"fmt"

	gear "github.com/azuyamat/gear/command"
	"github.com/azuyamat/pace/internal/detector"
	"github.com/azuyamat/pace/internal/logger"
	"github.com/azuyamat/pace/internal/models"
	"github.com/azuyamat/pace/internal/template"
)

var initCommand = gear.NewExecutableCommand("init", "Initialize a new Pace project in the current directory").
	Flags().
	Handler(initHandler)

func init() {
	RootCommand.AddChild(initCommand)
}

func initHandler(ctx *gear.Context, args gear.ValidatedArgs) error {
	projectType := detector.DetectCurrentProjectType()
	if projectType == models.ProjectTypeUnknown {
		return fmt.Errorf("unsupported project type: %s", projectType)
	}
	logger.Info("Detected project type: %s", projectType)
	generator := template.GetGeneratorByProjectType(projectType)
	if generator == nil {
		return fmt.Errorf("unsupported project type: %s", projectType)
	}
	config, err := generator.Generate()
	if err != nil {
		return err
	}
	err = config.WriteToFile("config.pace")
	if err != nil {
		return err
	}
	logger.Info("Generated config.pace")

	return nil
}
