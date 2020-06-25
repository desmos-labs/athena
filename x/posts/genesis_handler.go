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
	for _, post := range genPosts {
		if err := db.SavePost(post); err != nil {
			return err
		}
	}

	// Save the registered reactions
	for _, reaction := range genState.RegisteredReactions {
		if _, err := db.RegisterReactionIfNotPresent(reaction); err != nil {
			return err
		}
	}

	// Save the reactions
	for postIDKey, reactions := range genState.PostReactions {
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

	// Save poll answers
	for postIDKey, answers := range genState.UsersPollAnswers {
		postID, err := posts.ParsePostID(postIDKey)
		if err != nil {
			return err
		}

		for _, answer := range answers {
			if err := db.SavePollAnswer(postID, answer); err != nil {
				return err
			}
		}
	}

	return nil
}
