package types

import (
	initcmd "github.com/desmos-labs/juno/cmd/init"
	juno "github.com/desmos-labs/juno/types"
	"github.com/pelletier/go-toml"
	"github.com/spf13/cobra"
)

var _ juno.Config = &Config{}

// Config contains the data used to configure DJuno
type Config struct {
	juno.Config
	Notifications *NotificationsConfig `toml:"notifications"`
}

// NewConfig allows to build a new Config instance
func NewConfig(config juno.Config, notificationsConfig *NotificationsConfig) *Config {
	return &Config{
		Config:        config,
		Notifications: notificationsConfig,
	}
}

// NotificationsConfig contains the configuration for the notifications of DJuno
type NotificationsConfig struct {
	FirebaseCredentialsFile string `toml:"firebase_credentials_file"`
	FirebaseProjectID       string `toml:"firebase_project_id"`
}

// NewNotificationsConfig returns a new NotificationsConfig instance
func NewNotificationsConfig(firebaseFilePath, firebaseProjectID string) *NotificationsConfig {
	return &NotificationsConfig{
		FirebaseCredentialsFile: firebaseFilePath,
		FirebaseProjectID:       firebaseProjectID,
	}
}

// -------------------------------------------------------------------------------------------------------------------

type configToml struct {
	NotificationsConfig *NotificationsConfig `toml:"notifications"`
}

// ParseCfg parses the given file contents into a configuration object
func ParseCfg(fileContents []byte) (juno.Config, error) {
	junoCfg, err := juno.DefaultConfigParser(fileContents)
	if err != nil {
		return nil, err
	}

	var djunoCfg configToml
	err = toml.Unmarshal(fileContents, &djunoCfg)
	if err != nil {
		return nil, err
	}

	return NewConfig(junoCfg, djunoCfg.NotificationsConfig), nil
}

// -------------------------------------------------------------------------------------------------------------------

const (
	FlagNotificationsFirebaseCredentialFile = "notifications-firebase-credential-file"
	FlagNotificationsFirebaseProjectID      = "notifications-firebase-project-id"
)

// SetupFlags adds the proper flags to the init command allowing to specify the notifications-related data
func SetupFlags(cmd *cobra.Command) {
	cmd.Flags().String(FlagNotificationsFirebaseCredentialFile, "", "Path to the Firebase credentials file")
	cmd.Flags().String(FlagNotificationsFirebaseProjectID, "", "ID of the Firebase project to be used to send the notifications")
}

// CreateConfigFromFlags returns a new Config instance based on the values provided to the "init"
// command using the various flags
func CreateConfigFromFlags(cmd *cobra.Command) juno.Config {
	junoCfg := initcmd.DefaultConfigCreator(cmd)

	notificationsFirebaseFile, _ := cmd.Flags().GetString(FlagNotificationsFirebaseCredentialFile)
	notificationsProjectID, _ := cmd.Flags().GetString(FlagNotificationsFirebaseProjectID)

	return NewConfig(
		junoCfg,
		NewNotificationsConfig(notificationsFirebaseFile, notificationsProjectID),
	)
}
