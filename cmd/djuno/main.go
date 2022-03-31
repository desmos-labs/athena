package main

import (
	desmosapp "github.com/desmos-labs/desmos/v2/app"
	junocmd "github.com/forbole/juno/v3/cmd"
	initcmd "github.com/forbole/juno/v3/cmd/init"
	parsecmd "github.com/forbole/juno/v3/cmd/parse"
	parsecmdtypes "github.com/forbole/juno/v3/cmd/parse/types"

	fixcmd "github.com/desmos-labs/djuno/v2/cmd/parse"
	desmosdb "github.com/desmos-labs/djuno/v2/database"
	"github.com/desmos-labs/djuno/v2/x"
)

func main() {
	// Setup the config
	parseCfg := parsecmdtypes.NewConfig().
		WithRegistrar(x.NewModulesRegistrar()).
		WithEncodingConfigBuilder(desmosapp.MakeTestEncodingConfig).
		WithDBBuilder(desmosdb.Builder)

	cfg := junocmd.NewConfig("djuno").
		WithParseConfig(parseCfg)

	// Run the command
	rootCmd := junocmd.RootCmd(cfg.GetName())

	rootCmd.AddCommand(
		junocmd.VersionCmd(),
		initcmd.NewInitCmd(cfg.GetInitConfig()),
		parsecmd.NewParseCmd(cfg.GetParseConfig()),
		fixcmd.NewFixCmd(cfg.GetParseConfig()),
	)

	executor := junocmd.PrepareRootCmd(cfg.GetName(), rootCmd)
	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}
