package handlers

import (
	"fmt"
	"github.com/desmos-labs/juno/db"
	"github.com/desmos-labs/juno/db/postgresql"
	"github.com/desmos-labs/juno/types"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/desmos-labs/desmos/x/posts"
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

	// MsgCreatePost
	if createPostMsg, ok := msg.(posts.MsgCreatePost); ok {
		log.Info().Str("tx_hash", tx.TxHash).Int("msg_index", index).Msg("Found MsgCreatePost")

		var postID uint64

		// Get the post id
		// TODO: test with multiple MsgCreatePost
		for _, ev := range tx.Logs[index].Events {
			for _, attr := range ev.Attributes {
				if attr.Key == "post_id" {
					postID, _ = strconv.ParseUint(attr.Value, 10, 64)
				}
			}
		}

		if err := handleMsgCreatePost(postID, createPostMsg, postgresqlDb); err != nil {
			return err
		}
	}

	if editPostMsg, ok := msg.(posts.MsgEditPost); ok {
		if err := handleMsgEditPost(editPostMsg, postgresqlDb); err != nil {
			return err
		}
	}

	if addReactionMsg, ok := msg.(posts.MsgAddPostReaction); ok {
		if err := handleMsgAddPostReaction(addReactionMsg, postgresqlDb); err != nil {
			return err
		}
	}

	if removeReactionMsg, ok := msg.(posts.MsgRemovePostReaction); ok {
		if err := handleMsgRemovePostReaction(removeReactionMsg, postgresqlDb); err != nil {
			return err
		}
	}

	if answerPollMsg, ok := msg.(posts.MsgAnswerPoll); ok {
		if err := handleMsgAnswerPoll(answerPollMsg, postgresqlDb); err != nil {
			return err
		}
	}

	return nil
}
