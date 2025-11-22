package command

import (
	"fmt"

	gear "github.com/azuyamat/gear/command"
	"github.com/azuyamat/pace/internal/version"
)

var versionCommand = gear.NewExecutableCommand("version", "Show version information").
	Handler(versionHandler)

func init() {
	RootCommand.AddChild(versionCommand)
}

func versionHandler(ctx *gear.Context, args gear.ValidatedArgs) error {
	fmt.Printf("pace version %s\n", version.Version)
	fmt.Printf("commit: %s\n", version.Commit)
	fmt.Printf("built at: %s\n", version.Date)
	return nil
}
