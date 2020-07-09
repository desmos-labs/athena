package handlers

import (
	"time"

	profilestypes "github.com/desmos-labs/desmos/x/profiles"
	desmosdb "github.com/desmos-labs/djuno/database"
	juno "github.com/desmos-labs/juno/types"
)

// HandleMsgSaveProfile handles a MsgCreateProfile and properly stores the new profile inside the database
func HandleMsgSaveProfile(tx juno.Tx, index int, msg profilestypes.MsgSaveProfile, database desmosdb.DesmosDb) error {
	// Get the creation date
	event, err := tx.FindEventByType(index, profilestypes.EventTypeProfileSaved)
	if err != nil {
		return err
	}
	creationDateStr, err := tx.FindAttributeByKey(event, profilestypes.AttributeProfileCreationTime)
	if err != nil {
		return err
	}

	creationDate, err := time.Parse(time.RFC3339, creationDateStr)
	if err != nil {
		return err
	}

	newProfile := profilestypes.NewProfile(msg.Dtag, msg.Creator, creationDate).
		WithMoniker(msg.Moniker).
		WithBio(msg.Bio).
		WithPictures(msg.ProfilePic, msg.CoverPic)
	return database.SaveProfile(newProfile)
}

// HandleMsgDeleteProfile handles a MsgDeleteProfile correctly deleting the account present inside the database
func HandleMsgDeleteProfile(msg profilestypes.MsgDeleteProfile, database desmosdb.DesmosDb) error {
	return database.DeleteProfile(msg.Creator)
}
