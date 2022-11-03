package twitter

import (
	scorersutils "github.com/desmos-labs/djuno/v2/x/profiles-score/scorers/utils"
)

type Config struct {
	Token string `yaml:"token"`
}

func ParseConfig(bz []byte) (*Config, error) {
	var cfg Config
	found, err := scorersutils.UnmarshalConfig(bz, "twitter", &cfg)
	if err != nil {
		return nil, err
	}

	if !found {
		return nil, nil
	}

	return &cfg, err
}
