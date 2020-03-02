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
		if err := handleMsgCreatePost(tx, index, desmosMsg, postgresqlDb); err != nil {
			return err
		}
	case posts.MsgEditPost:
		if err := handleMsgEditPost(desmosMsg, postgresqlDb); err != nil {
			return err
		}
	case posts.MsgAddPostReaction:
		if err := handleMsgAddPostReaction(desmosMsg, postgresqlDb); err != nil {
			return err
		}
	case posts.MsgRemovePostReaction:
		if err := handleMsgRemovePostReaction(desmosMsg, postgresqlDb); err != nil {
			return err
		}
	case posts.MsgAnswerPoll:
		if err := handleMsgAnswerPoll(desmosMsg, postgresqlDb); err != nil {
			return err
		}
	}

	return nil
}
