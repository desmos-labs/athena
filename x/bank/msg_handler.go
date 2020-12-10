package bank

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	desmosdb "github.com/desmos-labs/djuno/database"
	juno "github.com/desmos-labs/juno/types"
)

// HandleMsg handles properly all the Cosmos x/bank modules messages
func HandleMsg(tx *juno.Tx, msg sdk.Msg, db *desmosdb.DesmosDb) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch cosmosMsg := msg.(type) {

	case *banktypes.MsgSend:
		return handleMsgSend(cosmosMsg, db)

	case *banktypes.MsgMultiSend:
		return handleMsgMultiSend(cosmosMsg, db)
	}

	return nil
}

func handleMsgSend(msg *banktypes.MsgSend, database *desmosdb.DesmosDb) error {
	err := database.SaveUserIfNotExisting(msg.FromAddress)
	if err != nil {
		return err
	}

	return database.SaveUserIfNotExisting(msg.ToAddress)
}

func handleMsgMultiSend(msg *banktypes.MsgMultiSend, database *desmosdb.DesmosDb) error {
	for _, input := range msg.Inputs {
		err := database.SaveUserIfNotExisting(input.Address)
		if err != nil {
			return err
		}
	}

	for _, output := range msg.Outputs {
		err := database.SaveUserIfNotExisting(output.Address)
		if err != nil {
			return err
		}
	}

	return nil
}
