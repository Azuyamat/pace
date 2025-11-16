package command

import "fmt"

func init() {
	CommandRegistry.Register(testCommand())
}

func testCommand() *Command {
	return NewCommand("test", "Test command").
		Arg(NewStringArg("name", "Name argument", true)).
		SetHandler(NewHandler(
			func(ctx *CommandContext, args *ValidatedArgs) {
				name := args.String("name")
				fmt.Printf("Test command executed with name: %s\n", name)
			}))
}
