package handlers

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/desmos-labs/desmos/x/posts"
	"github.com/desmos-labs/juno/db"
	"github.com/desmos-labs/juno/db/postgresql"
	"github.com/desmos-labs/juno/types"
	"github.com/rs/zerolog/log"
)

func MsgHandler(tx types.Tx, index int, msg sdk.Msg, db db.Database) error {
	if len(tx.Logs) == 0 {
		log.Info().Msg(fmt.Sprintf("Skipping message at index %d of tx hash %s as it was not successull",
			index, tx.TxHash))
		return nil
	}

	postgresqlDb, ok := db.(postgresql.Database)
	if !ok {
		return fmt.Errorf("database is not a PostgreSQL instance")
	}

	switch desmosMsg := msg.(type) {
	case posts.MsgCreatePost:
		return HandleMsgCreatePost(tx, index, desmosMsg, postgresqlDb)
	case posts.MsgEditPost:
		return HandleMsgEditPost(desmosMsg, postgresqlDb)
	case posts.MsgAddPostReaction:
		return HandleMsgAddPostReaction(desmosMsg, postgresqlDb)
	case posts.MsgRemovePostReaction:
		return HandleMsgRemovePostReaction(desmosMsg, postgresqlDb)
	case posts.MsgAnswerPoll:
		return HandleMsgAnswerPoll(desmosMsg, postgresqlDb)
	}

	return nil
}
