package relationships

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	relationshipstypes "github.com/desmos-labs/desmos/x/relationships/types"
	desmosdb "github.com/desmos-labs/djuno/database"
	"github.com/desmos-labs/djuno/x/relationships/handlers"
	"github.com/desmos-labs/juno/parse/worker"
	juno "github.com/desmos-labs/juno/types"
	"github.com/rs/zerolog/log"
)

func MsgHandler(tx juno.Tx, index int, msg sdk.Msg, w worker.Worker) error {
	if len(tx.Logs) == 0 {
		log.Info().
			Str("module", "relationships").
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
	case relationshipstypes.MsgCreateRelationship:
		return handlers.HandleMsgCreateRelationship(desmosMsg, database)
	case relationshipstypes.MsgDeleteRelationship:
		return handlers.HandleMsgDeleteRelationship(desmosMsg, database)
	case relationshipstypes.MsgBlockUser:
		return handlers.HandleMsgBlockUser(desmosMsg, database)
	case relationshipstypes.MsgUnblockUser:
		return handlers.HandleMsgUnblockUser(desmosMsg, database)
	}

	return nil
}
