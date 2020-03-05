package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/desmos-labs/desmos/x/posts"
	"github.com/desmos-labs/juno/db"
	"github.com/desmos-labs/juno/db/postgresql"
	tmtypes "github.com/tendermint/tendermint/types"
)

func GenesisHandler(codec *codec.Codec, genesisDoc *tmtypes.GenesisDoc, appState map[string]json.RawMessage, database db.Database) error {
	psqlDb, ok := database.(postgresql.Database)
	if !ok {
		return fmt.Errorf("database is not a PostgreSQL instance")
	}

	// Handle posts
	var genDocs posts.GenesisState
	codec.MustUnmarshalJSON(appState[posts.ModuleName], &genDocs)

	for _, post := range genDocs.Posts {

	}

	return nil
}
