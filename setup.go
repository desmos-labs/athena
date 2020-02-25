package Djuno

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/desmos-labs/desmos/app"
)

func desmosConfig(cfg *sdk.Config) {
	cfg.SetBech32PrefixForAccount(
		app.Bech32MainPrefix,
		app.Bech32MainPrefix+sdk.PrefixPublic,
	)
	cfg.SetBech32PrefixForValidator(
		app.Bech32MainPrefix+sdk.PrefixValidator+sdk.PrefixOperator,
		app.Bech32MainPrefix+sdk.PrefixValidator+sdk.PrefixOperator+sdk.PrefixPublic,
	)
	cfg.SetBech32PrefixForConsensusNode(
		app.Bech32MainPrefix+sdk.PrefixValidator+sdk.PrefixConsensus,
		app.Bech32MainPrefix+sdk.PrefixValidator+sdk.PrefixConsensus+sdk.PrefixPublic,
	)
}
