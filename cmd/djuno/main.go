package main

import (
	"os"

	"github.com/desmos-labs/djuno/types"

	"github.com/desmos-labs/juno/cmd/init"

	desmosapp "github.com/desmos-labs/desmos/app"
	junocmd "github.com/desmos-labs/juno/cmd"
	"github.com/desmos-labs/juno/cmd/parse"

	desmosdb "github.com/desmos-labs/djuno/database"
	"github.com/desmos-labs/djuno/x"
)

func main() {
	// Setup the config
	initCfg := init.NewConfig().
		WithConfigFlagSetup(types.SetupFlags).
		WithConfigCreator(types.CreateConfigFromFlags)

	parseCfg := parse.NewConfig().
		WithRegistrar(x.NewModulesRegistrar()).
		WithEncodingConfigBuilder(desmosapp.MakeTestEncodingConfig).
		WithDBBuilder(desmosdb.Builder).
		WithConfigParser(types.ParseCfg)

	cfg := junocmd.NewConfig("djuno").
		WithInitConfig(initCfg).
		WithParseConfig(parseCfg)

	// Run the commands and panic on any error
	executor := junocmd.BuildDefaultExecutor(cfg)
	err := executor.Execute()
	if err != nil {
		os.Exit(1)
	}
}
