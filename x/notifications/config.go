package notifications

import (
	"gopkg.in/yaml.v3"
)

type Config struct {
	FirebaseCredentialsFilePath string `yaml:"firebase_credentials_file_path"`
	FirebaseProjectID           string `yaml:"firebase_project_id"`
	AndroidChannelID            string `yaml:"android_channel_id"`
}

func ParseConfig(bz []byte) (*Config, error) {
	type T struct {
		Config *Config `yaml:"notifications"`
	}
	var cfg T
	err := yaml.Unmarshal(bz, &cfg)
	return cfg.Config, err
}
