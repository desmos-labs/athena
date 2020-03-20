package main

import (
	"github.com/desmos-labs/desmos/app"
	setup "github.com/desmos-labs/djuno/config"
	desmosdb "github.com/desmos-labs/djuno/db"
	"github.com/desmos-labs/djuno/handlers"
	"github.com/desmos-labs/juno/config"
	"github.com/desmos-labs/juno/executor"
	"github.com/desmos-labs/juno/parse/worker"
	"github.com/desmos-labs/juno/version"
)

func main() {
	// Register custom handlers
	worker.RegisterGenesisHandler(handlers.GenesisHandler)
	worker.RegisterTxHandler(handlers.TxHandler)
	worker.RegisterMsgHandler(handlers.MsgHandler)

	// Build the executor
	rootCmd := executor.BuildRootCmd("djuno", setup.DesmosConfig)
	rootCmd.AddCommand(
		version.GetVersionCmd(),
		GetDesmosParseCmd(app.MakeCodec(), desmosdb.Builder),
	)

	command := config.PrepareMainCmd(rootCmd)

	// Run the commands and panic on any error
	err := command.Execute()
	if err != nil {
		panic(err)
	}
}
