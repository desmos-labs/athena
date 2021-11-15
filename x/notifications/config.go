package notifications

import "gopkg.in/yaml.v3"

// Config contains the configuration for the notifications of DJuno
type Config struct {
	FirebaseCredentialsFile string `yaml:"firebase_credentials_file"`
	FirebaseProjectID       string `yaml:"firebase_project_id"`
}

func ParseConfig(data []byte) (*Config, error) {
	type cfgType struct {
		Notifications *Config `yaml:"notifications"`
	}

	var cfg cfgType
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return cfg.Notifications, nil
}
