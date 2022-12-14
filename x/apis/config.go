package apis

import (
	"gopkg.in/yaml.v3"
)

type Config struct {
	Address string `yaml:"address,omitempty"`
	Port    uint   `yaml:"port"`
}

func ParseConfig(bz []byte) (*Config, error) {
	type T struct {
		Config *Config `yaml:"apis"`
	}
	var cfg T
	err := yaml.Unmarshal(bz, &cfg)
	return cfg.Config, err
}
