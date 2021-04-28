package main

import (
	desmosapp "github.com/desmos-labs/desmos/app"
	desmosdb "github.com/desmos-labs/djuno/database"
	"github.com/desmos-labs/djuno/x"
	junocmd "github.com/desmos-labs/juno/cmd"
	junotypes "github.com/desmos-labs/juno/types"
)

func main() {
	// Build the root command
	rootCmd := junocmd.RootCmd("djuno")
	rootCmd.AddCommand(
		junocmd.VersionCmd(),
		junocmd.InitCmd(),
		ParseCmd(
			x.NewModulesRegistrar(),
			desmosapp.MakeTestEncodingConfig,
			junotypes.DefaultConfigSetup,
			desmosdb.Builder,
		),
	)

	// Run the commands and panic on any error
	command := junocmd.PrepareRootCmd(rootCmd)
	err := command.Execute()
	if err != nil {
		panic(err)
	}
}
