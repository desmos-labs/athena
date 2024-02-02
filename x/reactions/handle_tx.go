package reactions

import (
	abci "github.com/cometbft/cometbft/abci/types"
	reactionstypes "github.com/desmos-labs/desmos/v6/x/reactions/types"
	juno "github.com/forbole/juno/v5/types"

	"github.com/desmos-labs/athena/utils/events"
	"github.com/desmos-labs/athena/utils/transactions"

	"github.com/desmos-labs/athena/x/posts"
)

func (m *Module) HandleTx(tx *juno.Tx) error {
	return transactions.ParseTxEvents(tx, map[string]func(tx *juno.Tx, event abci.Event) error{
		reactionstypes.EventTypeAddReaction:              m.parseAddReactionEvent,
		reactionstypes.EventTypeRemoveReaction:           m.parseRemoveReactionEvent,
		reactionstypes.EventTypeAddRegisteredReaction:    m.parseAddRegisteredReactionEvent,
		reactionstypes.ActionEditRegisteredReaction:      m.parseEditRegisteredReactionEvent,
		reactionstypes.EventTypeRemoveRegisteredReaction: m.parseRemoveRegisteredReactionEvent,
		reactionstypes.EventTypeSetReactionsParams:       m.parseSetReactionsParamsEvent,
	})
}

// -------------------------------------------------------------------------------------------------------------------

// parseAddReactionEvent parses a reaction add event
func (m *Module) parseAddReactionEvent(tx *juno.Tx, event abci.Event) error {
	subspaceID, err := events.GetSubspaceIDFromEvent(event)
	if err != nil {
		return err
	}

	postID, err := posts.GetPostIDFromEvent(event)
	if err != nil {
		return err
	}

	reactionID, err := GetReactionIDFromEvent(event)
	if err != nil {
		return err
	}

	reaction, err := m.GetReaction(tx.Height, subspaceID, postID, reactionID)
	if err != nil {
		return err
	}

	return m.db.SaveReaction(reaction)
}

// parseRemoveReactionEvent parses a reaction remove event
func (m *Module) parseRemoveReactionEvent(tx *juno.Tx, event abci.Event) error {
	subspaceID, err := events.GetSubspaceIDFromEvent(event)
	if err != nil {
		return err
	}

	postID, err := posts.GetPostIDFromEvent(event)
	if err != nil {
		return err
	}

	reactionID, err := GetReactionIDFromEvent(event)
	if err != nil {
		return err
	}

	return m.db.DeleteReaction(tx.Height, subspaceID, postID, reactionID)
}

// parseAddRegisteredReactionEvent parses a registered reaction add event
func (m *Module) parseAddRegisteredReactionEvent(tx *juno.Tx, event abci.Event) error {
	subspaceID, err := events.GetSubspaceIDFromEvent(event)
	if err != nil {
		return err
	}

	registeredReactionID, err := GetRegisteredReactionIDFromEvent(event)
	if err != nil {
		return err
	}

	return m.updateRegisteredReaction(tx.Height, subspaceID, registeredReactionID)
}

// parseEditRegisteredReactionEvent parses a registered reaction edit event
func (m *Module) parseEditRegisteredReactionEvent(tx *juno.Tx, event abci.Event) error {
	subspaceID, err := events.GetSubspaceIDFromEvent(event)
	if err != nil {
		return err
	}

	registeredReactionID, err := GetRegisteredReactionIDFromEvent(event)
	if err != nil {
		return err
	}

	return m.updateRegisteredReaction(tx.Height, subspaceID, registeredReactionID)
}

// parseRemoveRegisteredReactionEvent parses a registered reaction remove event
func (m *Module) parseRemoveRegisteredReactionEvent(tx *juno.Tx, event abci.Event) error {
	subspaceID, err := events.GetSubspaceIDFromEvent(event)
	if err != nil {
		return err
	}

	registeredReactionID, err := GetRegisteredReactionIDFromEvent(event)
	if err != nil {
		return err
	}

	return m.db.DeleteRegisteredReaction(tx.Height, subspaceID, registeredReactionID)
}

// parseSetReactionsParamsEvent parses a set reactions params event
func (m *Module) parseSetReactionsParamsEvent(tx *juno.Tx, event abci.Event) error {
	subspaceID, err := events.GetSubspaceIDFromEvent(event)
	if err != nil {
		return err
	}

	return m.updateReactionParams(tx.Height, subspaceID)
}
