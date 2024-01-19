package utils

import (
	abci "github.com/cometbft/cometbft/abci/types"
	subspacestypes "github.com/desmos-labs/desmos/v6/x/subspaces/types"
	juno "github.com/forbole/juno/v5/types"
	"github.com/rs/zerolog/log"

	"github.com/desmos-labs/athena/x/filters"
)

// GetSubspaceIDFromEvent returns the subspace ID from the given event
func GetSubspaceIDFromEvent(event abci.Event) (uint64, error) {
	attribute, err := juno.FindAttributeByKey(event, subspacestypes.AttributeKeySubspaceID)
	if err != nil {
		return 0, err
	}
	return subspacestypes.ParseSubspaceID(attribute.Value)
}

// ParseTxEvents parses the given events using the given parsers
func ParseTxEvents(tx *juno.Tx, eventsParsers map[string]func(tx *juno.Tx, event abci.Event) error) error {
	for _, event := range tx.Events {
		if !filters.ShouldEventBeParsed(event) {
			continue
		}

		parseEvent, canBeParsed := eventsParsers[event.Type]
		if !canBeParsed {
			continue
		}

		err := parseEvent(tx, event)
		if err != nil {
			return err
		}

		log.Debug().Str("event", event.Type).Msg("handled event")
	}

	return nil
}
