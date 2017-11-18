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

type listCommand struct {
}

func ListCommand() subcommands.Command {
	return &listCommand{}
}

func (cmd *listCommand) Name() string {
	return "list"
}

func (cmd *listCommand) Synopsis() string {
	return "Show all installed golang version"
}

func (cmd *listCommand) Usage() string {
	return "list"
}

func (cmd *listCommand) SetFlags(f *flag.FlagSet) {
}

func (cmd *listCommand) Execute(_ context.Context, _ *flag.FlagSet, args ...interface{}) subcommands.ExitStatus {
	rc, ok := args[0].(*config.RCFile)
	if !ok {
		// assertion
		panic("First argument missing")
	}

	// Open GovenvManager
	govenvManager := manager.New(rc.GitPath, rc.ManagementDirectoryPath)

	// Get GoRootManager
	goRootManager := govenvManager.GetGoRootsManager()

	versions, err := goRootManager.Installed()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error on goRootManager.Installed(); %s\n", err)
		return subcommands.ExitFailure
	}

	for _, v := range versions {
		fmt.Println(v)
	}

	return subcommands.ExitSuccess
}
