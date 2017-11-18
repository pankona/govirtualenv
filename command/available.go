package command

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/google/subcommands"
	"github.com/necomeshi/govirtualenv/config"
	"github.com/necomeshi/govirtualenv/manager"
)

type availableCommand struct {
}

// AvailableCommand returns 'available' command interface
func AvailableCommand() subcommands.Command {
	return &availableCommand{}
}

func (cmd *availableCommand) Name() string {
	return "available"
}

func (cmd *availableCommand) Synopsis() string {
	return "Show available golang version"
}

func (cmd *availableCommand) Help() string {
	return "available"
}

func (cmd *availableCommand) Usage() string {
	return "available"
}

func (cmd *availableCommand) SetFlags(f *flag.FlagSet) {
}

func (cmd *availableCommand) Execute(_ context.Context, f *flag.FlagSet, args ...interface{}) subcommands.ExitStatus {
	rc, ok := args[0].(*config.RCFile)
	if !ok {
		// assertion
		panic("First argument missing")
	}

	// Open GovenvManager
	govenvManager := manager.New(rc.GitPath, rc.ManagementDirectoryPath)

	// Get GoRootManager
	goRootManager := govenvManager.GetGoRootsManager()

	// Get available golang version(this is just a tag name)
	versions, err := goRootManager.Available()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error on goRootManager.Available(); %s", err)
		return subcommands.ExitFailure
	}

	for _, v := range versions {
		fmt.Printf("%s\n", v)
	}

	return subcommands.ExitSuccess
}
