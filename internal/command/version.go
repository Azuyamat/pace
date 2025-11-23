package command

import (
	gear "github.com/azuyamat/gear/command"
	"github.com/azuyamat/pace/internal/logger"
	"github.com/azuyamat/pace/internal/version"
)

var versionCommand = gear.NewExecutableCommand("version", "Show version information").
	Handler(versionHandler)

func init() {
	RootCommand.AddChild(versionCommand)
}

func versionHandler(ctx *gear.Context, args gear.ValidatedArgs) error {
	logger.Printf("pace version %s\n", version.Version)
	logger.Printf("commit: %s\n", version.Commit)
	logger.Printf("built at: %s\n", version.Date)
	return nil
}
