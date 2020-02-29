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

		postgrDb, ok := db.(postgresql.Database)
		if !ok {
			return fmt.Errorf("database is not a MongoDB instance")
		}

		if err := handleMsgCreatePost(postID, createPostMsg, postgrDb); err != nil {
			return err
		}
	}

	if editPostMsg, ok := msg.(posts.MsgEditPost); ok {

	}

	if addReactionMsg, ok := msg.(posts.MsgAddPostReaction); ok {

	}

	if removeReactionMsg, ok := msg.(posts.MsgRemovePostReaction); ok {

	}

	if answerPollMsg, ok := msg.(posts.MsgAnswerPoll); ok {

	}

	return nil
}
