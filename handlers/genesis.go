package handlers

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/desmos-labs/desmos/x/posts"
	desmosdb "github.com/desmos-labs/djuno/db"
	"github.com/desmos-labs/juno/db"
	tmtypes "github.com/tendermint/tendermint/types"
)

// GenesisHandler allows to handle properly the parsing of the genesis file, by storing the data present inside it
func GenesisHandler(codec *codec.Codec, _ *tmtypes.GenesisDoc, appState map[string]json.RawMessage, database db.Database) error {
	desmosDb, ok := database.(desmosdb.DesmosDb)
	if !ok {
		return fmt.Errorf("database is not a DesmosDB instance")
	}

	// Get the posts state
	var postsGenState posts.GenesisState
	codec.MustUnmarshalJSON(appState[posts.ModuleName], &postsGenState)
	if err := handlePostsGenesis(desmosDb, postsGenState); err != nil {
		return err
	}

	// TODO: Add other modules here

	return nil
}

// handlePostsGenesis allows to properly handle the genesis state of the posts module
func handlePostsGenesis(db desmosdb.DesmosDb, genState posts.GenesisState) error {
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
	// TODO: Re-implement once the Desmos issue has been implemented
	// https://github.com/desmos-labs/desmos/issues/157
	//for postIDKey, reactions := range genState.PostReactions {
	//	postID, err := posts.ParsePostID(postIDKey)
	//	if err != nil {
	//		return err
	//	}
	//
	//	for _, reaction := range reactions {
	//		err := db.SaveReaction(reaction)
	//		if err != nil {
	//			return err
	//		}
	//	}
	//}

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
