package bank

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	desmosdb "github.com/desmos-labs/djuno/database"
	"github.com/desmos-labs/juno/parse/worker"
	juno "github.com/desmos-labs/juno/types"
	"github.com/rs/zerolog/log"
)

// MsgHandler handles properly all the Cosmos x/bank modules messages
func MsgHandler(tx juno.Tx, index int, msg sdk.Msg, w worker.Worker) error {
	if len(tx.Logs) == 0 {
		log.Info().
			Str("module", "bank").
			Str("tx_hash", tx.TxHash).Int("msg_index", index).
			Msg("skipping message as it was not successful")
		return nil
	}

	database, ok := w.Db.(desmosdb.DesmosDb)
	if !ok {
		return fmt.Errorf("database is not a DesmosDb instance")
	}

	switch cosmosMsg := msg.(type) {

	// Users
	case bank.MsgSend:
		return handleMsgSend(cosmosMsg, database)
	case bank.MsgMultiSend:
		return handleMsgMultiSend(cosmosMsg, database)
	}

	return nil
}

func handleMsgMultiSend(cosmosMsg bank.MsgMultiSend, database desmosdb.DesmosDb) error {
	for _, input := range cosmosMsg.Inputs {
		if _, err := database.SaveUserIfNotExisting(input.Address); err != nil {
			return err
		}
	}

	for _, output := range cosmosMsg.Outputs {
		if _, err := database.SaveUserIfNotExisting(output.Address); err != nil {
			return err
		}
	}

	return nil
}

func handleMsgSend(cosmosMsg bank.MsgSend, database desmosdb.DesmosDb) error {
	_, err := database.SaveUserIfNotExisting(cosmosMsg.FromAddress)
	if err != nil {
		return err
	}

	_, err = database.SaveUserIfNotExisting(cosmosMsg.ToAddress)
	return err
}
