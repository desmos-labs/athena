package main

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	desmosapp "github.com/desmos-labs/desmos/app"
	desmosdb "github.com/desmos-labs/djuno/database"
	"github.com/desmos-labs/djuno/x/bank"
	"github.com/desmos-labs/djuno/x/notifications"
	"github.com/desmos-labs/djuno/x/posts"
	"github.com/desmos-labs/djuno/x/profiles"
	"github.com/desmos-labs/djuno/x/relationships"
	"github.com/desmos-labs/juno/config"
	"github.com/desmos-labs/juno/executor"
	"github.com/desmos-labs/juno/parse/worker"
	"github.com/desmos-labs/juno/version"
)

func main() {
	// Register custom modules
	SetupModules()

	// Build the executor
	rootCmd := executor.BuildRootCmd("djuno", SetupConfig)
	rootCmd.AddCommand(
		version.GetVersionCmd(),
		GetDesmosParseCmd(desmosapp.MakeCodec(), desmosdb.Builder),
	)

	command := config.PrepareMainCmd(rootCmd)

	// Run the commands and panic on any error
	err := command.Execute()
	if err != nil {
		panic(err)
	}
}

func SetupConfig(cfg *sdk.Config) {
	cfg.SetBech32PrefixForAccount(
		desmosapp.Bech32MainPrefix,
		desmosapp.Bech32MainPrefix+sdk.PrefixPublic,
	)
	cfg.SetBech32PrefixForValidator(
		desmosapp.Bech32MainPrefix+sdk.PrefixValidator+sdk.PrefixOperator,
		desmosapp.Bech32MainPrefix+sdk.PrefixValidator+sdk.PrefixOperator+sdk.PrefixPublic,
	)
	cfg.SetBech32PrefixForConsensusNode(
		desmosapp.Bech32MainPrefix+sdk.PrefixValidator+sdk.PrefixConsensus,
		desmosapp.Bech32MainPrefix+sdk.PrefixValidator+sdk.PrefixConsensus+sdk.PrefixPublic,
	)
}

func SetupModules() {
	// Register genesis handlers
	worker.RegisterGenesisHandler(posts.GenesisHandler)
	worker.RegisterGenesisHandler(profiles.GenesisHandler)

	// Register tx handlers
	worker.RegisterTxHandler(notifications.TxHandler)

	// Register message handlers
	worker.RegisterMsgHandler(bank.MsgHandler)
	worker.RegisterMsgHandler(posts.MsgHandler)
	worker.RegisterMsgHandler(profiles.MsgHandler)
	worker.RegisterMsgHandler(relationships.MsgHandler)
}
