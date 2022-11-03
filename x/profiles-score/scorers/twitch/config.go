package twitch

import (
	scorersutils "github.com/desmos-labs/djuno/v2/x/profiles-score/scorers/utils"
)

type Config struct {
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
}

func UnmarshalConfig(bz []byte) (*Config, error) {
	var cfg Config
	found, err := scorersutils.UnmarshalConfig(bz, "twitch", &cfg)
	if err != nil {
		return nil, err
	}

	if !found {
		return nil, nil
	}

	return &cfg, nil
}
