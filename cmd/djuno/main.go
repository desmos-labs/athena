package main

import (
	"os"

	"github.com/desmos-labs/djuno/config"

	desmosapp "github.com/desmos-labs/desmos/app"
	junocmd "github.com/desmos-labs/juno/cmd"
	"github.com/desmos-labs/juno/cmd/parse"

	desmosdb "github.com/desmos-labs/djuno/database"
	"github.com/desmos-labs/djuno/x"
)

func main() {
	// Setup the config
	cfg := parse.NewConfig("djuno").
		WithRegistrar(x.NewModulesRegistrar()).
		WithEncodingConfigBuilder(desmosapp.MakeTestEncodingConfig).
		WithDBBuilder(desmosdb.Builder).
		WithConfigParser(config.ParseCfg)

	// Run the commands and panic on any error
	executor := junocmd.BuildDefaultExecutor(cfg)
	err := executor.Execute()
	if err != nil {
		os.Exit(1)
	}
}
