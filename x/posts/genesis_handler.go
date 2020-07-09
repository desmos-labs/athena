package posts

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/desmos-labs/desmos/x/posts"
	desmosdb "github.com/desmos-labs/djuno/database"
	"github.com/desmos-labs/juno/parse/worker"
	tmtypes "github.com/tendermint/tendermint/types"
)

// GenesisHandler allows to properly handle the genesis state for the posts module
func GenesisHandler(
	codec *codec.Codec, _ *tmtypes.GenesisDoc, appState map[string]json.RawMessage, w worker.Worker,
) error {
	db, ok := w.Db.(desmosdb.DesmosDb)
	if !ok {
		return fmt.Errorf("database is not a DesmosDB instance")
	}

	// Get the posts state
	var genState posts.GenesisState
	codec.MustUnmarshalJSON(appState[posts.ModuleName], &genState)

	// Order the posts based on the ids
	genPosts := genState.Posts
	sort.SliceStable(genPosts, func(i, j int) bool {
		first, second := genPosts[i], genPosts[j]
		return first.Created.Before(second.Created)
	})

	// Save the posts
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

func savePosts(genPosts posts.Posts, db desmosdb.DesmosDb) error {
	for _, post := range genPosts {
		if err := db.SavePost(post); err != nil {
			return err
		}
	}
	return nil
}

func saveRegisteredReactions(reactions posts.Reactions, db desmosdb.DesmosDb) error {
	for _, reaction := range reactions {
		if _, err := db.SaveRegisteredReactionIfNotPresent(reaction); err != nil {
			return err
		}
	}
	return nil
}

func savePostReactions(reactions map[string]posts.PostReactions, db desmosdb.DesmosDb) error {
	for postIDKey, reactions := range reactions {
		postID, err := posts.ParsePostID(postIDKey)
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

func savePollAnswers(userAnswers map[string]posts.UserAnswers, db desmosdb.DesmosDb) error {
	for postIDKey, answers := range userAnswers {
		postID, err := posts.ParsePostID(postIDKey)
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
