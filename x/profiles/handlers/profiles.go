package handlers

import (
	"github.com/desmos-labs/desmos/x/profile"
	desmosdb "github.com/desmos-labs/djuno/database"
)

// HandleMsgSaveProfile handles a MsgCreateProfile and properly stores the new profile inside the database
func HandleMsgSaveProfile(msg profile.MsgSaveProfile, database desmosdb.DesmosDb) error {
	newProfile := profile.NewProfile(msg.Creator).
		WithMoniker(msg.Moniker).
		WithName(msg.Name).
		WithSurname(msg.Surname).
		WithBio(msg.Bio).
		WithPictures(msg.ProfileCov, msg.ProfileCov)
	_, err := database.SaveProfile(newProfile)
	return err
}

// HandleMsgDeleteProfile handles a MsgDeleteProfile correctly deleting the account present inside the database
func HandleMsgDeleteProfile(msg profile.MsgDeleteProfile, database desmosdb.DesmosDb) error {
	return database.DeleteProfile(msg.Creator)
}
