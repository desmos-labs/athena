package tips

import (
	"gopkg.in/yaml.v3"
)

type ContractsConfig struct {
	Tips *Config `yaml:"tips"`
}

type Config struct {
	CodeID uint64 `yaml:"code_id"`
}

func ParseConfig(bz []byte) (*Config, error) {
	type T struct {
		ContractsCfg *ContractsConfig `yaml:"contracts"`
	}
	var cfg T
	err := yaml.Unmarshal(bz, &cfg)
	if err != nil {
		return nil, err
	}

	if cfg.ContractsCfg != nil {
		return cfg.ContractsCfg.Tips, nil
	}

	return nil, nil
}
