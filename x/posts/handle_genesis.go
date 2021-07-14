package posts

import (
	"encoding/json"
	"sort"

	"github.com/desmos-labs/djuno/types"

	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/codec"
	poststypes "github.com/desmos-labs/desmos/x/staging/posts/types"

	desmosdb "github.com/desmos-labs/djuno/database"
)

// HandleGenesis allows to properly handle the genesis state for the posts module
func HandleGenesis(
	doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage, codec codec.Marshaler, db *desmosdb.Db,
) error {
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
	err := savePosts(doc.InitialHeight, genPosts, db)
	if err != nil {
		return err
	}

	// Save the registered reactions
	err = saveRegisteredReactions(doc.InitialHeight, genState.RegisteredReactions, db)
	if err != nil {
		return err
	}

	// Save the reactions
	err = savePostReactions(doc.InitialHeight, genState.PostsReactions, db)
	if err != nil {
		return err
	}

	// Save poll answers
	err = savePollAnswers(doc.InitialHeight, genState.UsersPollAnswers, db)
	if err != nil {
		return err
	}

	return nil
}

func savePosts(height int64, posts []poststypes.Post, db *desmosdb.Db) error {
	for index := range posts {
		err := db.SavePost(types.NewPost(posts[index], height))
		if err != nil {
			return err
		}
	}
	return nil
}

func saveRegisteredReactions(height int64, reactions []poststypes.RegisteredReaction, db *desmosdb.Db) error {
	for _, reaction := range reactions {
		err := db.RegisterReactionIfNotPresent(types.NewRegisteredReaction(reaction, height))
		if err != nil {
			return err
		}
	}
	return nil
}

func savePostReactions(height int64, reactions []poststypes.PostReaction, db *desmosdb.Db) error {
	for _, reaction := range reactions {
		err := db.SavePostReaction(types.NewPostReaction(reaction.PostID, reaction, height))
		if err != nil {
			return err
		}
	}
	return nil
}

func savePollAnswers(height int64, userAnswers []poststypes.UserAnswer, db *desmosdb.Db) error {
	for _, answer := range userAnswers {
		err := db.SaveUserPollAnswer(types.NewUserPollAnswer(answer, height))
		if err != nil {
			return err
		}
	}

	return nil
}
