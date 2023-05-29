package contracts

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/forbole/juno/v5/modules"
	"github.com/rs/zerolog/log"

	juno "github.com/forbole/juno/v5/types"
)

// SmartContractModule represents a generic smart contract module
type SmartContractModule interface {
	modules.Module
	modules.MessageModule
	modules.AuthzMessageModule

	// RefreshData refreshes the smart contract data for the given height and subspace id
	RefreshData(height int64, subspaceID uint64) error
}

var (
	_ SmartContractModule = &Module{}
)

// Module represents the module that allows to handle all smart contracts modules easily
type Module struct {
	modules []SmartContractModule
}

// NewModule returns a new Module instance
func NewModule(modules []SmartContractModule) *Module {
	return &Module{
		modules: modules,
	}
}

// Name implements modules.Module
func (m Module) Name() string {
	return "contracts"
}

// HandleMsg implements modules.MessageModule
func (m Module) HandleMsg(index int, msg sdk.Msg, tx *juno.Tx) error {
	for _, module := range m.modules {
		if module == nil {
			continue
		}

		err := module.HandleMsg(index, msg, tx)
		if err != nil {
			return err
		}
	}

	return nil
}

// HandleMsgExec implements modules.AuthzMessageModule
func (m Module) HandleMsgExec(index int, msgExec *authz.MsgExec, authzMsgIndex int, executedMsg sdk.Msg, tx *juno.Tx) error {
	for _, module := range m.modules {
		if module == nil {
			continue
		}

		err := module.HandleMsgExec(index, msgExec, authzMsgIndex, executedMsg, tx)
		if err != nil {
			return err
		}
	}

	return nil
}

// RefreshData implements SmartContractModule
func (m Module) RefreshData(height int64, subspaceID uint64) error {
	for _, module := range m.modules {
		if module == nil {
			continue
		}

		log.Info().Int64("height", height).Uint64("subspace id", subspaceID).Msgf("refreshing %s", module.Name())
		err := module.RefreshData(height, subspaceID)
		if err != nil {
			return err
		}
	}

	return nil
}
