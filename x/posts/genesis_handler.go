package posts

import (
	"encoding/json"
	"sort"

	"github.com/cosmos/cosmos-sdk/codec"
	poststypes "github.com/desmos-labs/desmos/x/staging/posts/types"

	desmosdb "github.com/desmos-labs/djuno/database"
)

// HandleGenesis allows to properly handle the genesis state for the posts module
func HandleGenesis(appState map[string]json.RawMessage, codec *codec.LegacyAmino, db *desmosdb.DesmosDb) error {
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
	err := savePosts(genPosts, db)
	if err != nil {
		return err
	}

	// Save the registered reactions
	err = saveRegisteredReactions(genState.RegisteredReactions, db)
	if err != nil {
		return err
	}

	// Save the reactions
	err = savePostReactions(genState.PostsReactions, db)
	if err != nil {
		return err
	}

	// Save poll answers
	err = savePollAnswers(genState.UsersPollAnswers, db)
	if err != nil {
		return err
	}

	return nil
}

func savePosts(posts []poststypes.Post, db *desmosdb.DesmosDb) error {
	for _, post := range posts {
		err := db.SavePost(post)
		if err != nil {
			return err
		}
	}
	return nil
}

func saveRegisteredReactions(reactions []poststypes.RegisteredReaction, db *desmosdb.DesmosDb) error {
	for _, reaction := range reactions {
		err := db.RegisterReactionIfNotPresent(reaction)
		if err != nil {
			return err
		}
	}
	return nil
}

func savePostReactions(reactions []poststypes.PostReactionsEntry, db *desmosdb.DesmosDb) error {
	for _, entry := range reactions {
		for _, reaction := range entry.Reactions {
			err := db.SavePostReaction(entry.PostId, reaction)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func savePollAnswers(userAnswers []poststypes.UserAnswersEntry, db *desmosdb.DesmosDb) error {
	for _, entry := range userAnswers {
		for _, answer := range entry.UserAnswers {
			err := db.SaveUserPollAnswer(entry.PostId, answer)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
