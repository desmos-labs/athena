package config

import juno "github.com/desmos-labs/juno/types"

var _ juno.Config = &Config{}

// Config contains the data used to configure DJuno
type Config struct {
	juno.Config
	Notifications *NotificationsConfig
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
	Enable                  bool   `toml:"enable"`
	FirebaseCredentialsFile string `toml:"firebase_credentials_file"`
	FirebaseProjectID       string `toml:"firebase_project_id"`
}

// NewNotificationsConfig returns a new NotificationsConfig instance
func NewNotificationsConfig(enable bool, firebaseFilePath, firebaseProjectID string) *NotificationsConfig {
	return &NotificationsConfig{
		Enable:                  enable,
		FirebaseCredentialsFile: firebaseFilePath,
		FirebaseProjectID:       firebaseProjectID,
	}
}
