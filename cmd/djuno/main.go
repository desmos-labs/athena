package main

import (
	"github.com/desmos-labs/desmos/app"
	setup "github.com/desmos-labs/djuno/config"
	"github.com/desmos-labs/djuno/handlers"
	"github.com/desmos-labs/juno/db/postgresql"
	"github.com/desmos-labs/juno/executor"
	"github.com/desmos-labs/juno/parse/worker"
)

func main() {
	// Register custom handlers
	worker.RegisterGenesisHandler(handlers.GenesisHandler)
	worker.RegisterMsgHandler(handlers.MsgHandler)

	// Build the executor
	command := executor.BuildExecutor("djuno", setup.DesmosConfig, app.MakeCodec, postgresql.Builder)

	// Run the commands and panic on any error
	err := command.Execute()
	if err != nil {
		panic(err)
	}
}
