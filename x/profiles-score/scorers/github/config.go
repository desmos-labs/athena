package github

import (
	scorersutils "github.com/desmos-labs/athena/x/profiles-score/scorers/utils"
)

type Config struct {
	AppID              int64  `yaml:"app_id"`
	InstallationID     int64  `yaml:"installation_id"`
	PrivateKeyFilePath string `yaml:"private_key_file_path"`
}

func UnmarshalConfig(bz []byte) (*Config, error) {
	var cfg Config
	found, err := scorersutils.UnmarshalConfig(bz, "github", &cfg)
	if err != nil {
		return nil, err
	}

	if !found {
		return nil, nil
	}

	return &cfg, nil
}
