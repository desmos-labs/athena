package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/desmos-labs/desmos/x/posts"
	desmosdb "github.com/desmos-labs/djuno/db"
	"github.com/desmos-labs/juno/db"
	tmtypes "github.com/tendermint/tendermint/types"
)

func GenesisHandler(codec *codec.Codec, genesisDoc *tmtypes.GenesisDoc, appState map[string]json.RawMessage, database db.Database) error {
	desmosDb, ok := database.(desmosdb.DesmosDb)
	if !ok {
		return fmt.Errorf("database is not a DesmosDB instance")
	}

	// Handle posts
	var genDocs posts.GenesisState
	codec.MustUnmarshalJSON(appState[posts.ModuleName], &genDocs)

	for _, post := range genDocs.Posts {
		if err := desmosDb.SavePost(post); err != nil {
			return err
		}
	}

	return nil
}
