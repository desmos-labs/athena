package main

import (
	desmosapp "github.com/desmos-labs/desmos/app"
	desmosdb "github.com/desmos-labs/djuno/database"
	"github.com/desmos-labs/djuno/x/bank"
	"github.com/desmos-labs/djuno/x/notifications"
	"github.com/desmos-labs/djuno/x/posts"
	"github.com/desmos-labs/djuno/x/profiles"
	"github.com/desmos-labs/djuno/x/relationships"
	junocmd "github.com/desmos-labs/juno/cmd"
	"github.com/desmos-labs/juno/config"
	"github.com/desmos-labs/juno/modules/registrar"
)

func main() {
	registrar.RegisterModules(
		bank.Module{},
		notifications.Module{},
		posts.Module{},
		profiles.Module{},
		relationships.Module{},
	)

	// Build the root command
	rootCmd := junocmd.RootCmd("djuno")
	rootCmd.AddCommand(
		junocmd.VersionCmd(),
		ParseCmd(desmosapp.MakeCodecs, config.DefaultSetup, desmosdb.Builder),
	)

	// Run the commands and panic on any error
	command := junocmd.PrepareMainCmd(rootCmd)
	err := command.Execute()
	if err != nil {
		panic(err)
	}
}
