package posts

import (
	"encoding/json"
	"sort"

	poststypes "github.com/desmos-labs/desmos/x/posts/types"

	"github.com/cosmos/cosmos-sdk/codec"
	desmosdb "github.com/desmos-labs/djuno/database"
)

// HandleGenesis allows to properly handle the genesis state for the posts module
func HandleGenesis(codec *codec.LegacyAmino, appState map[string]json.RawMessage, db *desmosdb.DesmosDb) error {
	// Get the posts state
	var genState poststypes.GenesisState
	codec.MustUnmarshalJSON(appState[poststypes.ModuleName], &genState)

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
		err := db.RegisterReactionIfNotPresent(reaction)
		if err != nil {
			return err
		}
	}

	// Save the reactions
	for _, entry := range genState.PostsReactions {
		for _, reaction := range entry.Reactions {
			err := db.SaveReaction(entry.PostId, &reaction)
			if err != nil {
				return err
			}
		}
	}

	// Save poll answers
	for _, entry := range genState.UsersPollAnswers {
		for _, answer := range entry.UserAnswers {
			err := db.SaveUserPollAnswer(entry.PostId, answer)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
