package profiles

import (
	"time"

	profilestypes "github.com/desmos-labs/desmos/x/profiles/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	desmosdb "github.com/desmos-labs/djuno/database"
	juno "github.com/desmos-labs/juno/types"
)

// HandleMsg allows to handle different messages types for the profiles module
func HandleMsg(tx *juno.Tx, index int, msg sdk.Msg, db *desmosdb.DesmosDb) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch desmosMsg := msg.(type) {
	case *profilestypes.MsgSaveProfile:
		return handleMsgSaveProfile(tx, index, desmosMsg, db)

	case *profilestypes.MsgDeleteProfile:
		return handleMsgDeleteProfile(desmosMsg, db)

		// TODO: Handle DTag transfer requests

	}

	return nil
}

// handleMsgSaveProfile handles a MsgCreateProfile and properly stores the new profile inside the database
func handleMsgSaveProfile(tx *juno.Tx, index int, msg *profilestypes.MsgSaveProfile, database *desmosdb.DesmosDb) error {
	event, err := tx.FindEventByType(index, profilestypes.EventTypeProfileSaved)
	if err != nil {
		return err
	}

	// Get creation date
	creationDateStr, err := tx.FindAttributeByKey(event, profilestypes.AttributeProfileCreationTime)
	if err != nil {
		return err
	}
	creationDate, err := time.Parse(time.RFC3339, creationDateStr)
	if err != nil {
		return err
	}

	newProfile := profilestypes.NewProfile(
		msg.Dtag,
		msg.Moniker,
		msg.Bio,
		profilestypes.NewPictures(msg.ProfilePicture, msg.CoverPicture),
		creationDate,
		msg.Creator,
	)
	return database.SaveProfile(newProfile)
}

// handleMsgDeleteProfile handles a MsgDeleteProfile correctly deleting the account present inside the database
func handleMsgDeleteProfile(msg *profilestypes.MsgDeleteProfile, database *desmosdb.DesmosDb) error {
	return database.DeleteProfile(msg.Creator)
}
