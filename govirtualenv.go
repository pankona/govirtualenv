package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/subcommands"
	"github.com/necomeshi/govirtualenv/command"
	"github.com/necomeshi/govirtualenv/config"
)

func main() {
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")
	subcommands.Register(command.InstallCommand(), "")
	subcommands.Register(command.RemoveCommand(), "")
	subcommands.Register(command.AvailableCommand(), "")
	subcommands.Register(command.ListCommand(), "")
	subcommands.Register(command.StartprojectCommand(), "")
	subcommands.Register(command.ChangeCommand(), "")

	flag.Parse()

	// TODO: Check here no argument passed

	ctx := context.Background()

	// Read RC file
	var rcFilePath string

	if homeDirPath, found := os.LookupEnv("HOME"); found {
		rcFilePath = filepath.Join(homeDirPath, config.GovenvEnvironmentFile)
	} else {
		rcFilePath = filepath.Join(config.GovenvEnvironmentFile)
	}

	rcFile, err := os.Open(rcFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while opening %s; %s\n", rcFilePath, err)
		os.Exit(1)
	}

	rc, err := config.ReadRCFile(rcFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while reading %s; %s\n", rcFilePath, err)
		os.Exit(1)
	}

	os.Exit(int(subcommands.Execute(ctx, rc)))
}
