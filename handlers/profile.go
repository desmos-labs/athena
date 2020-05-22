package handlers

import (
	"github.com/desmos-labs/desmos/x/profile"
	desmosdb "github.com/desmos-labs/djuno/db"
)

// HandleMsgCreateProfile handles a MsgCreateProfile and properly stores the new profile inside the database
func HandleMsgCreateProfile(msg profile.MsgCreateProfile, database desmosdb.DesmosDb) error {
	newProfile := profile.NewProfile(msg.Moniker, msg.Creator).
		WithName(msg.Name).
		WithSurname(msg.Surname).
		WithBio(msg.Bio).
		WithPictures(msg.Pictures)
	_, err := database.UpsertProfile(newProfile)
	return err
}

// HandleMsgEditProfile handles a MsgEditProfile updating the profile information that are already stored
// inside the database
func HandleMsgEditProfile(msg profile.MsgEditProfile, database desmosdb.DesmosDb) error {
	user, err := database.GetUserByAddress(msg.Creator)
	if err != nil {
		return err
	}

	var moniker string
	if user != nil && user.Moniker.Valid {
		moniker = user.Moniker.String
	}

	if msg.NewMoniker != nil {
		moniker = *msg.NewMoniker
	}

	newProfile := profile.NewProfile(moniker, msg.Creator).
		WithName(msg.Name).
		WithSurname(msg.Surname).
		WithBio(msg.Bio)
	_, err = database.UpsertProfile(newProfile)
	return err
}

// HandleMsgDeleteProfile handles a MsgDeleteProfile correctly deleting the account present inside the database
func HandleMsgDeleteProfile(msg profile.MsgDeleteProfile, database desmosdb.DesmosDb) error {
	return database.DeleteProfile(msg.Creator)
}
