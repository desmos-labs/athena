package links

import (
	"gopkg.in/yaml.v3"
)

type Config struct {
	FirebaseCredentialsFilePath string      `yaml:"firebase_credentials_file_path"`
	Domain                      string      `yaml:"domain"`
	Desktop                     DesktopInfo `yaml:"desktop"`
	Android                     AndroidInfo `yaml:"android"`
	Ios                         IosConfig   `yaml:"ios"`
}

type DesktopInfo struct {
	FallbackLink string `yaml:"fallback_link"`
}

type AndroidInfo struct {
	PackageName           string `yaml:"package_name"`
	MinPackageVersionCode string `yaml:"min_package_version_code"`
}

type IosConfig struct {
	BundleID       string `yaml:"bundle_id"`
	MinimumVersion string `yaml:"minimum_version"`
	AppStoreID     string `yaml:"app_store_id"`
}

func ParseConfig(bz []byte) (*Config, error) {
	type T struct {
		Config *Config `yaml:"dynamic_links"`
	}
	var cfg T
	err := yaml.Unmarshal(bz, &cfg)
	return cfg.Config, err
}
