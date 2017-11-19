package command

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/google/subcommands"
	"github.com/necomeshi/govirtualenv/config"
	"github.com/necomeshi/govirtualenv/manager"
	"github.com/necomeshi/govirtualenv/project"
)

type changeCommand struct {
}

func ChangeCommand() subcommands.Command {
	return &changeCommand{}
}

func (cmd *changeCommand) Name() string {
	return "change"
}

func (cmd *changeCommand) Synopsis() string {
	return "change changes current using golang version."
}

func (cmd *changeCommand) Usage() string {
	return "change GOLANG_VERSION"
}

func (cmd *changeCommand) SetFlags(f *flag.FlagSet) {
}

func (cmd *changeCommand) Execute(_ context.Context, f *flag.FlagSet, args ...interface{}) subcommands.ExitStatus {
	rc, ok := args[0].(*config.RCFile)
	if !ok {
		// assertion
		panic("First argument missing")
	}

	if f.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "Less argument.")
		return subcommands.ExitFailure
	}

	if !project.IsActivated() {
		fmt.Fprintf(os.Stderr, "Virtual environment is not activated.")
		return subcommands.ExitFailure
	}

	version := f.Arg(0)

	// Open GovenvManager
	govenvManager := manager.New(rc.GitPath, rc.ManagementDirectoryPath)

	// Get GoRootManager
	goRootManager := govenvManager.GetGoRootsManager()

	// Check it is installed.
	goRootPath, found, err := goRootManager.LookupInstalledVersion(version)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error on goRootManager.LookupInstalledVersion; %s", err)
		return subcommands.ExitFailure
	}

	if !found {
		fmt.Fprintf(os.Stderr, "Requested version %s is not installed.", version)
		return subcommands.ExitFailure
	}

	pj, err := project.CreateFromCurrentEnvironment()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error on project.CreateFromCurrentEnvironment(); %s", err)
		return subcommands.ExitFailure
	}

	pj.SetGoRoot(goRootPath)

	if err := pj.Configure(); err != nil {
		fmt.Fprintf(os.Stderr, "Error on project.Configure(); %s", err)
		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}
