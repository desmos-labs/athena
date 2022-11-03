package youtube

import (
	scorersutils "github.com/desmos-labs/djuno/v2/x/profiles-score/scorers/utils"
)

type Config struct {
	APIKey string `yaml:"api_key"`
}

func ParseConfig(bz []byte) (*Config, error) {
	var cfg Config
	found, err := scorersutils.UnmarshalConfig(bz, "youtube", &cfg)
	if err != nil {
		return nil, err
	}

	if !found {
		return nil, nil
	}

	return &cfg, err
}
