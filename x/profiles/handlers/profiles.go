package handlers

import (
	"github.com/desmos-labs/desmos/x/profiles"
	desmosdb "github.com/desmos-labs/djuno/database"
	juno "github.com/desmos-labs/juno/types"
)

// HandleMsgSaveProfile handles a MsgCreateProfile and properly stores the new profile inside the database
func HandleMsgSaveProfile(tx juno.Tx, index int, msg profiles.MsgSaveProfile, database desmosdb.DesmosDb) error {
	// Get the creation date
	event, err := tx.FindEventByType(index, profiles.EventTypeProfileSaved)
	if err != nil {
		return err
	}
	creationDate = tx.FindEventByType(index, profiles.AttributeProfileCreator)

	newProfile := profiles.NewProfile(msg.Dtag, msg.Creator).
		WithMoniker(msg.Moniker).
		WithBio(msg.Bio).
		WithPictures(msg.ProfilePic, msg.CoverPic)
	_, err := database.SaveProfile(newProfile)
	return err
}

// HandleMsgDeleteProfile handles a MsgDeleteProfile correctly deleting the account present inside the database
func HandleMsgDeleteProfile(msg profiles.MsgDeleteProfile, database desmosdb.DesmosDb) error {
	return database.DeleteProfile(msg.Creator)
}
