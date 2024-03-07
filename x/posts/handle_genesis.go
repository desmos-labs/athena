package posts

import (
	"encoding/json"

	poststypes "github.com/desmos-labs/desmos/v7/x/posts/types"

	"github.com/desmos-labs/athena/v2/types"

	tmtypes "github.com/cometbft/cometbft/types"
)

// HandleGenesis implements modules.GenesisModule
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	var genState poststypes.GenesisState
	m.cdc.MustUnmarshalJSON(appState[poststypes.ModuleName], &genState)

	// Save posts
	for _, post := range genState.Posts {
		err := m.db.SavePost(types.NewPost(post, doc.InitialHeight))
		if err != nil {
			return err
		}
	}

	// Save attachments
	for _, attachment := range genState.Attachments {
		err := m.db.SavePostAttachment(types.NewPostAttachment(attachment, doc.InitialHeight))
		if err != nil {
			return err
		}
	}

	// Save poll answers
	for _, answer := range genState.UserAnswers {
		err := m.db.SavePollAnswer(types.NewPollAnswer(answer, doc.InitialHeight))
		if err != nil {
			return err
		}
	}

	return m.db.SavePostsParams(types.NewPostsParams(genState.Params, doc.InitialHeight))
}
