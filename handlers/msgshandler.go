package handlers

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/desmos-labs/desmos/x/posts"
	"github.com/desmos-labs/desmos/x/profile"
	desmosdb "github.com/desmos-labs/djuno/db"
	"github.com/desmos-labs/juno/db"
	"github.com/desmos-labs/juno/types"
	"github.com/rs/zerolog/log"
)

func MsgHandler(tx types.Tx, index int, msg sdk.Msg, db db.Database) error {
	log.Info().Str("tx_hash", tx.TxHash).Int("msg_index", index).Str("msg_type", msg.Type()).Msg("found message")

	if len(tx.Logs) == 0 {
		log.Info().Msg(fmt.Sprintf("Skipping message at index %d of tx hash %s as it was not successull",
			index, tx.TxHash))
		return nil
	}

	database, ok := db.(desmosdb.DesmosDb)
	if !ok {
		return fmt.Errorf("database is not a DesmosDb instance")
	}

	switch desmosMsg := msg.(type) {
	// Posts
	case posts.MsgCreatePost:
		return HandleMsgCreatePost(tx, index, desmosMsg, database)
	case posts.MsgEditPost:
		return HandleMsgEditPost(desmosMsg, database)
	case posts.MsgAddPostReaction:
		return HandleMsgAddPostReaction(tx, index, database)
	case posts.MsgRemovePostReaction:
		return HandleMsgRemovePostReaction(tx, index, database)

	// Polls
	case posts.MsgAnswerPoll:
		return HandleMsgAnswerPoll(desmosMsg, database)

	// Reactions
	case posts.MsgRegisterReaction:
		return HandleMsgRegisterReaction(desmosMsg, database)

	// Users
	case profile.MsgSaveProfile:
		return HandleMsgSaveProfile(desmosMsg, database)
	case profile.MsgDeleteProfile:
		return HandleMsgDeleteProfile(desmosMsg, database)
	}

	return nil
}
