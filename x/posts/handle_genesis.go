package posts

import (
	"encoding/hex"
	"encoding/json"

	poststypes "github.com/desmos-labs/desmos/v4/x/posts/types"

	"github.com/desmos-labs/djuno/v2/types"

	tmtypes "github.com/tendermint/tendermint/types"
)

// HandleGenesis implements modules.GenesisModule
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	var genState poststypes.GenesisState
	m.cdc.MustUnmarshalJSON(appState[poststypes.ModuleName], &genState)

	// Save posts
	for _, post := range genState.Posts {
		txHashes := []string{hex.EncodeToString(doc.AppHash)}
		err := m.db.SavePost(types.NewPost(post, txHashes, doc.InitialHeight))
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
