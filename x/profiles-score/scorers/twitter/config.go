package twitter

import "gopkg.in/yaml.v3"

type ScorersConfig struct {
	Scorer *Config `yaml:"twitter"`
}

type Config struct {
	Token string `yaml:"token"`
}

func ParseConfig(bz []byte) (*Config, error) {
	type T struct {
		Scorers *ScorersConfig `yaml:"scorers"`
	}

	var cfg T
	err := yaml.Unmarshal(bz, &cfg)
	if err != nil {
		return nil, err
	}

	if cfg.Scorers == nil {
		return nil, nil
	}

	return cfg.Scorers.Scorer, err
}
