package filters

import (
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	subspacestypes "github.com/desmos-labs/desmos/v6/x/subspaces/types"
	"github.com/forbole/juno/v5/types/config"

	"github.com/desmos-labs/athena/utils"
)

var (
	initialized bool
	cfg         *Config
)

// ShouldEventBeParsed tells whether the given event should be parsed
func ShouldEventBeParsed(event abci.Event) bool {
	parseCfg()

	subspaceID, err := utils.GetSubspaceIDFromEvent(event)
	if err != nil {
		return false
	}

	return cfg.isSubspaceSupported(subspaceID)
}

// ShouldMsgBeParsed tells whether the given subspace is currently supported and its messages should be parsed
func ShouldMsgBeParsed(msg sdk.Msg) bool {
	parseCfg()
	if subspaceMsg, ok := msg.(subspacestypes.SubspaceMsg); ok && cfg != nil {
		return cfg.isSubspaceSupported(subspaceMsg.GetSubspaceID())
	}
	return true
}

// parseCfg parses the filter configuration
func parseCfg() {
	if initialized {
		return
	}

	junoCfgBz, err := config.Cfg.GetBytes()
	if err != nil {
		panic(err)
	}
	parsedCfg, err := ParseConfig(junoCfgBz)
	if err != nil {
		panic(err)
	}
	cfg = parsedCfg
	initialized = true
}
