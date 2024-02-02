package relationships

import (
	abci "github.com/cometbft/cometbft/abci/types"
	relationshipstypes "github.com/desmos-labs/desmos/v6/x/relationships/types"
	juno "github.com/forbole/juno/v5/types"

	"github.com/desmos-labs/athena/utils/transactions"

	"github.com/desmos-labs/athena/types"
)

func (m *Module) HandleTx(tx *juno.Tx) error {
	return transactions.ParseTxEvents(tx, map[string]func(tx *juno.Tx, event abci.Event) error{
		relationshipstypes.EventTypeRelationshipCreated: m.parseCreateRelationshipEvent,
		relationshipstypes.EventTypeRelationshipDeleted: m.parseDeleteRelationshipEvent,
		relationshipstypes.EventTypeBlockUser:           m.parseBlockUserEvent,
		relationshipstypes.EventTypeUnblockUser:         m.parseUnblockUserEvent,
	})
}

// --------------------------------------------------------------------------------------------------------------------

// parseCreateRelationshipEvent allows to properly handle a relationship creation event
func (m *Module) parseCreateRelationshipEvent(tx *juno.Tx, event abci.Event) error {
	subspace, err := GetSubspaceFromEvent(event)
	if err != nil {
		return err
	}

	creator, err := GetCreatorFromEvent(event)
	if err != nil {
		return err
	}

	counterparty, err := GetCounterpartyFromEvent(event)
	if err != nil {
		return err
	}

	return m.db.SaveRelationship(types.NewRelationship(
		relationshipstypes.NewRelationship(creator, counterparty, subspace),
		tx.Height,
	))
}

// parseDeleteRelationshipEvent allows to properly handle a relationship deletion event
func (m *Module) parseDeleteRelationshipEvent(tx *juno.Tx, event abci.Event) error {
	subspace, err := GetSubspaceFromEvent(event)
	if err != nil {
		return err
	}

	creator, err := GetCreatorFromEvent(event)
	if err != nil {
		return err
	}

	counterparty, err := GetCounterpartyFromEvent(event)
	if err != nil {
		return err
	}

	return m.db.DeleteRelationship(types.NewRelationship(
		relationshipstypes.NewRelationship(creator, counterparty, subspace),
		tx.Height,
	))
}

// parseBlockUserEvent allows to properly handle a user block event
func (m *Module) parseBlockUserEvent(tx *juno.Tx, event abci.Event) error {
	subspace, err := GetSubspaceFromEvent(event)
	if err != nil {
		return err
	}

	blocker, err := GetBlockerFromEvent(event)
	if err != nil {
		return err
	}

	blocked, err := GetBlockedFromEvent(event)
	if err != nil {
		return err
	}

	return m.updateUserBlock(tx.Height, subspace, blocker, blocked)
}

// parseUnblockUserEvent allows to properly handle a user unblock event
func (m *Module) parseUnblockUserEvent(tx *juno.Tx, event abci.Event) error {
	subspace, err := GetSubspaceFromEvent(event)
	if err != nil {
		return err
	}

	blocker, err := GetBlockerFromEvent(event)
	if err != nil {
		return err
	}

	blocked, err := GetBlockedFromEvent(event)
	if err != nil {
		return err
	}

	return m.db.DeleteBlockage(types.NewBlockage(
		relationshipstypes.NewUserBlock(blocker, blocked, "", subspace),
		tx.Height,
	))
}
