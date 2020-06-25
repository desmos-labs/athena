package posts

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/desmos-labs/desmos/x/posts"
	"github.com/desmos-labs/djuno/database"
	"github.com/desmos-labs/djuno/x/posts/handlers"
	"github.com/desmos-labs/juno/parse/worker"
	juno "github.com/desmos-labs/juno/types"
	"github.com/rs/zerolog/log"
)

// MsgHandler allows to handle different message types from the posts module
func MsgHandler(tx juno.Tx, index int, msg sdk.Msg, w worker.Worker) error {
	if len(tx.Logs) == 0 {
		log.Info().
			Str("module", "posts").
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
	case posts.MsgCreatePost:
		return handlers.HandleMsgCreatePost(tx, index, desmosMsg, database)

	case posts.MsgEditPost:
		return handlers.HandleMsgEditPost(desmosMsg, database)

	// Reactions
	case posts.MsgRegisterReaction:
		return handlers.HandleMsgRegisterReaction(desmosMsg, database)

	case posts.MsgAddPostReaction:
		return handlers.HandleMsgAddPostReaction(tx, index, database)

	case posts.MsgRemovePostReaction:
		return handlers.HandleMsgRemovePostReaction(tx, index, database)

	// Polls
	case posts.MsgAnswerPoll:
		return handlers.HandleMsgAnswerPoll(desmosMsg, database)
	}

	return nil
}
