package config

import (
	"github.com/desmos-labs/juno/cmd/init"
	juno "github.com/desmos-labs/juno/types"
	"github.com/spf13/cobra"
)

const (
	FlagNotificationsEnabled                = "notifications-enabled"
	FlagNotificationsFirebaseCredentialFile = "notifications-firebase-credential-file"
	FlagNotificationsFirebaseProjectID      = "notifications-firebase-project-id"
)

// SetupFlags adds the proper flags to the init command allowing to specify the notifications-related data
func SetupFlags(cmd *cobra.Command) {
	cmd.Flags().Bool(FlagNotificationsEnabled, false, "Enable the push notifications")
	cmd.Flags().String(FlagNotificationsFirebaseCredentialFile, "", "Path to the Firebase credentials file")
	cmd.Flags().String(FlagNotificationsFirebaseProjectID, "", "ID of the Firebase project to be used to send the notifications")
}

// CreateConfigFromFlags returns a new Config instance based on the values provided to the "init"
// command using the various flags
func CreateConfigFromFlags(cmd *cobra.Command) juno.Config {
	junoCfg := init.DefaultConfigCreator(cmd)

	notificationsEnabled, _ := cmd.Flags().GetBool(FlagNotificationsEnabled)
	notificationsFirebaseFile, _ := cmd.Flags().GetString(FlagNotificationsFirebaseCredentialFile)
	notificationsProjectID, _ := cmd.Flags().GetString(FlagNotificationsFirebaseProjectID)

	return NewConfig(
		junoCfg,
		NewNotificationsConfig(notificationsEnabled, notificationsFirebaseFile, notificationsProjectID),
	)
}
