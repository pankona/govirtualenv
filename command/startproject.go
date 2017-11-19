package command

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/necomeshi/govirtualenv/project"

	"github.com/google/subcommands"
	"github.com/necomeshi/govirtualenv/config"
	"github.com/necomeshi/govirtualenv/manager"
)

type startprojectCommand struct {
}

// StartprojectCommand returns new startproject command interface
func StartprojectCommand() subcommands.Command {
	return &startprojectCommand{}
}

func (cmd *startprojectCommand) Name() string {
	return "startproject"
}

func (cmd *startprojectCommand) Synopsis() string {
	return "Start new project with given golang verison"
}

func (cmd *startprojectCommand) Usage() string {
	return "startproject GOLANG_VERSION PROJECTNAME"
}

func (cmd *startprojectCommand) SetFlags(f *flag.FlagSet) {
}

func (cmd *startprojectCommand) Execute(_ context.Context, f *flag.FlagSet, args ...interface{}) subcommands.ExitStatus {
	rc, ok := args[0].(*config.RCFile)
	if !ok {
		// assertion
		panic("First argument missing")
	}

	if f.NArg() < 2 {
		fmt.Fprintf(os.Stderr, "Less argument.")
		return subcommands.ExitFailure
	}

	version := f.Arg(0)
	projectName := f.Arg(1)

	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error on os.Getwd(); %s", err)
		return subcommands.ExitFailure
	}

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

	if _, err := project.Create(projectName, goRootPath, filepath.Join(wd, projectName)); err != nil {
		fmt.Fprintf(os.Stderr, "Error on project.Create(); %s", err)
		return subcommands.ExitFailure
	}

	fmt.Printf("Project %s created.\n", projectName)

	return subcommands.ExitSuccess
}
