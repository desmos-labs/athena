package config

import (
	"github.com/BurntSushi/toml"
	juno "github.com/desmos-labs/juno/types"
)

type configToml struct {
	NotificationsConfig *NotificationsConfig `toml:"notifications"`
}

// ParseCfg parses the given file contents into a configuration object
func ParseCfg(fileContents []byte) (juno.Config, error) {
	junoCfg, err := juno.DefaultConfigParser(fileContents)
	if err != nil {
		return nil, err
	}

	var djunoCfg configToml
	err = toml.Unmarshal(fileContents, &djunoCfg)
	if err != nil {
		return nil, err
	}

	return NewConfig(junoCfg, djunoCfg.NotificationsConfig), nil
}
