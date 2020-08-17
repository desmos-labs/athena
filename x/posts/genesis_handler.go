package posts

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/cosmos/cosmos-sdk/codec"
	poststypes "github.com/desmos-labs/desmos/x/posts/types"
	desmosdb "github.com/desmos-labs/djuno/database"
	"github.com/desmos-labs/juno/parse/worker"
	tmtypes "github.com/tendermint/tendermint/types"
)

// GenesisHandler allows to properly handle the genesis state for the poststypes module
func GenesisHandler(
	codec *codec.Codec, _ *tmtypes.GenesisDoc, appState map[string]json.RawMessage, w worker.Worker,
) error {
	db, ok := w.Db.(desmosdb.DesmosDb)
	if !ok {
		return fmt.Errorf("database is not a DesmosDB instance")
	}

	// Get the poststypes state
	var genState poststypes.GenesisState
	codec.MustUnmarshalJSON(appState[poststypes.ModuleName], &genState)

	// Order the poststypes based on the ids
	genPosts := genState.Posts
	sort.SliceStable(genPosts, func(i, j int) bool {
		first, second := genPosts[i], genPosts[j]
		return first.Created.Before(second.Created)
	})

	// Save the poststypes
	if err := savePosts(genPosts, db); err != nil {
		return err
	}

	// Save the registered reactions
	if err := saveRegisteredReactions(genState.RegisteredReactions, db); err != nil {
		return err
	}

	// Save the reactions
	if err := savePostReactions(genState.PostReactions, db); err != nil {
		return err
	}

	// Save poll answers
	if err := savePollAnswers(genState.UsersPollAnswers, db); err != nil {
		return err
	}

	return nil
}

func savePosts(genPosts poststypes.Posts, db desmosdb.DesmosDb) error {
	for _, post := range genPosts {
		if err := db.SavePost(post); err != nil {
			return err
		}
	}
	return nil
}

func saveRegisteredReactions(reactions poststypes.Reactions, db desmosdb.DesmosDb) error {
	for _, reaction := range reactions {
		if _, err := db.SaveRegisteredReactionIfNotPresent(reaction); err != nil {
			return err
		}
	}
	return nil
}

func savePostReactions(reactions map[string]poststypes.PostReactions, db desmosdb.DesmosDb) error {
	for postIDKey, reactions := range reactions {
		postID, err := poststypes.ParsePostID(postIDKey)
		if err != nil {
			return err
		}

		for _, reaction := range reactions {
			if err := db.SaveReaction(postID, &reaction); err != nil {
				return err
			}
		}
	}
	return nil
}

func savePollAnswers(userAnswers map[string]poststypes.UserAnswers, db desmosdb.DesmosDb) error {
	for postIDKey, answers := range userAnswers {
		postID, err := poststypes.ParsePostID(postIDKey)
		if err != nil {
			return err
		}

		for _, answer := range answers {
			if err := db.SaveUserPollAnswer(postID, answer); err != nil {
				return err
			}
		}
	}

	return nil
}
