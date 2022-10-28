package utils

import (
	"gopkg.in/yaml.v3"
)

func UnmarshalConfig(cfgBz []byte, nodeName string, value interface{}) (bool, error) {
	type T struct {
		Scorers map[interface{}]interface{} `yaml:"scorers"`
	}

	var cfg T
	err := yaml.Unmarshal(cfgBz, &cfg)
	if err != nil {
		return false, err
	}

	if len(cfg.Scorers) == 0 {
		return false, nil
	}

	nodeValue, ok := cfg.Scorers[nodeName]
	if !ok {
		return false, nil
	}

	nodeValueBz, err := yaml.Marshal(nodeValue)
	if err != nil {
		return false, err
	}
	return true, yaml.Unmarshal(nodeValueBz, value)
}
