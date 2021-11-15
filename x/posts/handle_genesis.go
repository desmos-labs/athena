package posts

import (
	"encoding/json"
	"sort"

	"github.com/desmos-labs/djuno/v2/types"

	tmtypes "github.com/tendermint/tendermint/types"

	poststypes "github.com/desmos-labs/desmos/x/staging/posts/types"
)

// HandleGenesis implements modules.GenesisModule
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	// Get the posts state
	var genState poststypes.GenesisState
	m.cdc.MustUnmarshalJSON(appState[poststypes.ModuleName], &genState)

	// Order the posts based on the ids
	genPosts := genState.Posts
	sort.SliceStable(genPosts, func(i, j int) bool {
		first, second := genPosts[i], genPosts[j]
		return first.Created.Before(second.Created)
	})

	// Save the posts
	err := m.savePosts(doc.InitialHeight, genPosts)
	if err != nil {
		return err
	}

	// Save the registered reactions
	err = m.saveRegisteredReactions(doc.InitialHeight, genState.RegisteredReactions)
	if err != nil {
		return err
	}

	// Save the reactions
	err = m.savePostReactions(doc.InitialHeight, genState.PostsReactions)
	if err != nil {
		return err
	}

	// Save poll answers
	err = m.savePollAnswers(doc.InitialHeight, genState.UsersPollAnswers)
	if err != nil {
		return err
	}

	return nil
}

func (m *Module) savePosts(height int64, posts []poststypes.Post) error {
	for index := range posts {
		err := m.db.SavePost(types.NewPost(posts[index], height))
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Module) saveRegisteredReactions(height int64, reactions []poststypes.RegisteredReaction) error {
	for _, reaction := range reactions {
		err := m.db.RegisterReactionIfNotPresent(types.NewRegisteredReaction(reaction, height))
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Module) savePostReactions(height int64, reactions []poststypes.PostReaction) error {
	for _, reaction := range reactions {
		err := m.db.SavePostReaction(types.NewPostReaction(reaction.PostID, reaction, height))
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Module) savePollAnswers(height int64, userAnswers []poststypes.UserAnswer) error {
	for _, answer := range userAnswers {
		err := m.db.SaveUserPollAnswer(types.NewUserPollAnswer(answer, height))
		if err != nil {
			return err
		}
	}

	return nil
}
