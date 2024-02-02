package filters

import (
	"gopkg.in/yaml.v3"

	"github.com/desmos-labs/athena/v2/x/contracts/tips"
)

type Config struct {
	// SupportedSubspaceIDs represents the list of supported subspaces
	SupportedSubspaceIDs []uint64 `yaml:"supported_subspace_ids"`

	// Contracts contains the list of contracts to be parsed
	Contracts *ContractsConfig `yaml:"contracts"`
}

type ContractsConfig struct {
	Tips *tips.Config `yaml:"tips"`
}

// isSubspaceSupported tells whether the given subspace is supported from this config
func (c *Config) isSubspaceSupported(subspaceID uint64) bool {
	for _, id := range cfg.SupportedSubspaceIDs {
		if id == subspaceID {
			return true
		}
	}
	return false
}

func ParseConfig(bz []byte) (*Config, error) {
	type T struct {
		Config *Config `yaml:"filters,omitempty"`
	}
	var cfg T
	err := yaml.Unmarshal(bz, &cfg)
	return cfg.Config, err
}
