package notifications

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/desmos-labs/juno/client"
	"github.com/desmos-labs/juno/config"
	"github.com/desmos-labs/juno/db"
	"github.com/desmos-labs/juno/modules"
	juno "github.com/desmos-labs/juno/types"
	"github.com/go-co-op/gocron"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

var _ modules.Module = Module{}

type Module struct{}

// Name implements Module
func (m Module) Name() string {
	return "notifications"
}

// RunAdditionalOperations implements Module
func (m Module) RunAdditionalOperations(
	_ *config.Config, _ *codec.LegacyAmino, _ *client.Proxy, _ db.Database,
) error {
	return nil
}

// RegisterPeriodicOperations implements Module
func (m Module) RegisterPeriodicOperations(
	_ *gocron.Scheduler, _ *codec.LegacyAmino, _ *client.Proxy, _ db.Database,
) error {
	return nil
}

// HandleGenesis implements Module
func (m Module) HandleGenesis(
	_ *tmtypes.GenesisDoc, _ map[string]json.RawMessage, _ *codec.LegacyAmino, _ *client.Proxy, _ db.Database,
) error {
	return nil
}

// HandleBlock implements Module
func (m Module) HandleBlock(
	_ *coretypes.ResultBlock, _ []*juno.Tx, _ *coretypes.ResultValidators,
	_ *codec.LegacyAmino, _ *client.Proxy, _ db.Database,
) error {
	return nil
}

// HandleTx implements Module
func (m Module) HandleTx(tx *juno.Tx, _ *codec.LegacyAmino, _ *client.Proxy, _ db.Database) error {
	return TxHandler(tx)
}

// HandleMsg implements Module
func (m Module) HandleMsg(_ int, _ sdk.Msg, _ *juno.Tx, _ *codec.LegacyAmino, _ *client.Proxy, _ db.Database) error {
	return nil
}
