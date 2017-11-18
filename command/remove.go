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

type removeCommand struct {
}

// RemoveCommand create new initialize command struct
func RemoveCommand() subcommands.Command {
	return &removeCommand{}
}

func (cmd removeCommand) Name() string {
	return "remove"
}

func (cmd removeCommand) Synopsis() string {
	return "Remove installed golang version"
}

func (cmd removeCommand) Usage() string {
	return "remove GO_VERSION"
}

func (cmd *removeCommand) SetFlags(f *flag.FlagSet) {

}

func (cmd *removeCommand) Execute(_ context.Context, f *flag.FlagSet, args ...interface{}) subcommands.ExitStatus {

	rc, ok := args[0].(*config.RCFile)
	if !ok {
		// assertion
		panic("First argument missing")
	}

	version := f.Arg(0)

	// Open GovenvManager
	govenvManager := manager.New(rc.GitPath, rc.ManagementDirectoryPath)

	// Get GoRootManager
	goRootManager := govenvManager.GetGoRootsManager()

	// if the version is already removeed, exit and nothing to do
	_, found, err := goRootManager.LookupInstalledVersion(version)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error on goRootManager.LookupInstalledVersion(); %s\n", err)
		return subcommands.ExitFailure
	}

	if !found {
		fmt.Fprintf(os.Stderr, "Requested version %s is not installed\n", version)
		return subcommands.ExitFailure
	}

	fmt.Printf("Removing %s ...", version)

	// remove
	if err := goRootManager.Remove(version); err != nil {
		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}
