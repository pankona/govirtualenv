package command

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/google/subcommands"
	"github.com/necomeshi/govirtualenv/config"
	"github.com/necomeshi/govirtualenv/manager"
)

type installCommand struct {
}

// InstallCommand create new initialize command struct
func InstallCommand() subcommands.Command {
	return &installCommand{}
}

func (cmd installCommand) Name() string {
	return "install"
}

func (cmd installCommand) Synopsis() string {
	return "Install new golang version"
}

func (cmd installCommand) Usage() string {
	return "install GO_VERSION"
}

func (cmd *installCommand) SetFlags(f *flag.FlagSet) {

}

func (cmd *installCommand) Execute(_ context.Context, f *flag.FlagSet, args ...interface{}) subcommands.ExitStatus {

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

	// if the version is already installed, exit and nothing to do
	_, found, err := goRootManager.LookupInstalledVersion(version)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err)
		return subcommands.ExitFailure
	}

	// If already installed. not a error, but just exit
	if found {
		fmt.Fprintf(os.Stdout, "Version %s is already installed", version)
		return subcommands.ExitSuccess
	}

	versions, err := goRootManager.Available()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error on goRootManager.Available();  %s", err)
		return subcommands.ExitFailure
	}
	found = false
	for _, v := range versions {
		if strings.Compare(v, version) == 0 {
			found = true
			break
		}
	}

	if !found {
		fmt.Fprintf(os.Stderr, "Requested version %s is not available\n", version)
		return subcommands.ExitFailure
	}

	fmt.Fprintf(os.Stdout, "Installing %s...", version)

	// install
	if err := goRootManager.Install(version); err != nil {
		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}
