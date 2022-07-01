package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

func (f *Config) loadFromConfigFile() error {
	confFileBuff, err := os.ReadFile(f.ConfigFile)
	if err != nil {
		return fmt.Errorf("could not read config file %w", err)
	}

	if err := yaml.Unmarshal(confFileBuff, f); err != nil {
		return fmt.Errorf("could not unmarshal yaml file: %w", err)
	}

	return nil
}

func (f *Config) createConfigFileWithDefaults() error {
	buff, err := yaml.Marshal(f)
	if err != nil {
		return fmt.Errorf("could not marshal Config struct: %w", err)
	}

	if err := os.WriteFile("config.yaml", buff, 655); err != nil {
		return fmt.Errorf("could not write config file: %w", err)
	}

	return nil
}
