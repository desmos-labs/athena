package posts

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	poststypes "github.com/desmos-labs/desmos/x/posts/types"
	"github.com/desmos-labs/djuno/database"
	"github.com/desmos-labs/djuno/x/posts/handlers"
	"github.com/desmos-labs/juno/parse/worker"
	juno "github.com/desmos-labs/juno/types"
	"github.com/rs/zerolog/log"
)

// MsgHandler allows to handle different message types from the poststypes module
func MsgHandler(tx juno.Tx, index int, msg sdk.Msg, w worker.Worker) error {
	if len(tx.Logs) == 0 {
		log.Info().
			Str("module", "poststypes").
			Str("tx_hash", tx.TxHash).Int("msg_index", index).
			Msg("skipping message as it was not successful")
		return nil
	}

	database, ok := w.Db.(database.DesmosDb)
	if !ok {
		return fmt.Errorf("invalid BigDipper database provided")
	}

	switch desmosMsg := msg.(type) {
	// Posts
	case poststypes.MsgCreatePost:
		return handlers.HandleMsgCreatePost(tx, index, desmosMsg, database)

	case poststypes.MsgEditPost:
		return handlers.HandleMsgEditPost(tx, index, desmosMsg, database)

	// Reactions
	case poststypes.MsgRegisterReaction:
		return handlers.HandleMsgRegisterReaction(desmosMsg, database)

	case poststypes.MsgAddPostReaction:
		return handlers.HandleMsgAddPostReaction(tx, index, database)

	case poststypes.MsgRemovePostReaction:
		return handlers.HandleMsgRemovePostReaction(tx, index, database)

	// Polls
	case poststypes.MsgAnswerPoll:
		return handlers.HandleMsgAnswerPoll(desmosMsg, database)
	}

	return nil
}
