package main

import (
	"os"

	desmosapp "github.com/desmos-labs/desmos/app"
	junocmd "github.com/forbole/juno/v2/cmd"
	parsecmd "github.com/forbole/juno/v2/cmd/parse"

	desmosdb "github.com/desmos-labs/djuno/v2/database"
	"github.com/desmos-labs/djuno/v2/x"
)

func main() {
	// Setup the config
	parseCfg := parsecmd.NewConfig().
		WithRegistrar(x.NewModulesRegistrar()).
		WithEncodingConfigBuilder(desmosapp.MakeTestEncodingConfig).
		WithDBBuilder(desmosdb.Builder)

	cfg := junocmd.NewConfig("djuno").
		WithParseConfig(parseCfg)

	// Run the commands and panic on any error
	executor := junocmd.BuildDefaultExecutor(cfg)
	err := executor.Execute()
	if err != nil {
		os.Exit(1)
	}
}
