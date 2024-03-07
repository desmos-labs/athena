package main

import (
	desmosapp "github.com/desmos-labs/desmos/v7/app"
	junocmd "github.com/forbole/juno/v5/cmd"
	initcmd "github.com/forbole/juno/v5/cmd/init"
	migratecmd "github.com/forbole/juno/v5/cmd/migrate"
	parsecmdtypes "github.com/forbole/juno/v5/cmd/parse/types"
	startcmd "github.com/forbole/juno/v5/cmd/start"
	"github.com/forbole/juno/v5/types/params"

	parsecmd "github.com/desmos-labs/athena/v2/cmd/parse"
	desmosdb "github.com/desmos-labs/athena/v2/database"
	"github.com/desmos-labs/athena/v2/x"
)

func main() {
	// Setup the config
	parseCfg := parsecmdtypes.NewConfig().
		WithRegistrar(x.NewModulesRegistrar()).
		WithEncodingConfigBuilder(func() params.EncodingConfig {
			config := desmosapp.MakeEncodingConfig()
			return params.EncodingConfig(config)
		}).
		WithDBBuilder(desmosdb.Builder)

	cfg := junocmd.NewConfig("athena").
		WithParseConfig(parseCfg)

	// Run the command
	rootCmd := junocmd.RootCmd(cfg.GetName())

	rootCmd.AddCommand(
		junocmd.VersionCmd(),
		initcmd.NewInitCmd(cfg.GetInitConfig()),
		startcmd.NewStartCmd(cfg.GetParseConfig()),
		parsecmd.NewParseCmd(cfg.GetParseConfig()),
		migratecmd.NewMigrateCmd(cfg.GetName(), cfg.GetParseConfig()),
	)

	executor := junocmd.PrepareRootCmd(cfg.GetName(), rootCmd)
	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}
