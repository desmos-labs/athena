package bank

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	juno "github.com/desmos-labs/juno/types"

	desmosdb "github.com/desmos-labs/djuno/database"
)

// HandleMsg handles properly all the Cosmos x/bank modules messages
func HandleMsg(tx *juno.Tx, msg sdk.Msg, db *desmosdb.DesmosDb) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch cosmosMsg := msg.(type) {

	case *banktypes.MsgSend:
		return handleMsgSend(tx.Height, cosmosMsg, db)

	case *banktypes.MsgMultiSend:
		return handleMsgMultiSend(tx.Height, cosmosMsg, db)
	}

	return nil
}

func handleMsgSend(height int64, msg *banktypes.MsgSend, database *desmosdb.DesmosDb) error {
	err := database.SaveUserIfNotExisting(msg.FromAddress, height)
	if err != nil {
		return err
	}

	return database.SaveUserIfNotExisting(msg.ToAddress, height)
}

func handleMsgMultiSend(height int64, msg *banktypes.MsgMultiSend, database *desmosdb.DesmosDb) error {
	for _, input := range msg.Inputs {
		err := database.SaveUserIfNotExisting(input.Address, height)
		if err != nil {
			return err
		}
	}

	for _, output := range msg.Outputs {
		err := database.SaveUserIfNotExisting(output.Address, height)
		if err != nil {
			return err
		}
	}

	return nil
}
