package main

import (
	"github.com/desmos-labs/desmos/app"
	setup "github.com/desmos-labs/djuno/config"
	desmosdb "github.com/desmos-labs/djuno/db"
	"github.com/desmos-labs/djuno/handlers"
	"github.com/desmos-labs/juno/executor"
	"github.com/desmos-labs/juno/parse/worker"
)

func main() {
	// Register custom handlers
	worker.RegisterGenesisHandler(handlers.GenesisHandler)
	worker.RegisterMsgHandler(handlers.MsgHandler)

	// Build the executor
	command := executor.BuildExecutor("djuno", setup.DesmosConfig, app.MakeCodec, desmosdb.Builder)

	// Run the commands and panic on any error
	err := command.Execute()
	if err != nil {
		panic(err)
	}
}
