package handlers

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/desmos-labs/desmos/x/posts"
	desmosdb "github.com/desmos-labs/djuno/db"
	"github.com/desmos-labs/juno/db"
	"github.com/desmos-labs/juno/types"
	"github.com/rs/zerolog/log"
)

func MsgHandler(tx types.Tx, index int, msg sdk.Msg, db db.Database) error {
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
	case posts.MsgCreatePost:
		return HandleMsgCreatePost(tx, index, desmosMsg, database)
	case posts.MsgEditPost:
		return HandleMsgEditPost(desmosMsg, database)
	case posts.MsgAddPostReaction:
		return HandleMsgAddPostReaction(desmosMsg, database)
	case posts.MsgRemovePostReaction:
		return HandleMsgRemovePostReaction(desmosMsg, database)
	case posts.MsgAnswerPoll:
		return HandleMsgAnswerPoll(desmosMsg, database)
	}

	return nil
}
