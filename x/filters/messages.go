package filters

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	subspacestypes "github.com/desmos-labs/desmos/v4/x/subspaces/types"
	"github.com/forbole/juno/v3/types/config"
)

var (
	initialized bool
	cfg         *Config
)

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
