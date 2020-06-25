package profiles

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/desmos-labs/desmos/x/profile"
	desmosdb "github.com/desmos-labs/djuno/database"
	"github.com/desmos-labs/djuno/x/profiles/handlers"
	"github.com/desmos-labs/juno/parse/worker"
	juno "github.com/desmos-labs/juno/types"
	"github.com/rs/zerolog/log"
)

// MsgHandler allows to handle different messages types for the profiles module
func MsgHandler(tx juno.Tx, index int, msg sdk.Msg, w worker.Worker) error {
	if len(tx.Logs) == 0 {
		log.Info().
			Str("module", "profiles").
			Str("tx_hash", tx.TxHash).Int("msg_index", index).
			Msg("skipping message as it was not successful")
		return nil
	}

	database, ok := w.Db.(desmosdb.DesmosDb)
	if !ok {
		return fmt.Errorf("database is not a DesmosDb instance")
	}

	switch desmosMsg := msg.(type) {

	// Users
	case profile.MsgSaveProfile:
		return handlers.HandleMsgSaveProfile(desmosMsg, database)
	case profile.MsgDeleteProfile:
		return handlers.HandleMsgDeleteProfile(desmosMsg, database)
	}

	return nil
}
