package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds the application configuration.
type Config struct {
	URL               string `yaml:"url"`
	Username          string `yaml:"username"`
	Password          string `yaml:"password"`
	PasswordFile      string `yaml:"password_file"`
	Insecure          bool   `yaml:"insecure"`
	EnableLabelValues bool   `yaml:"enable_label_values"`
	HistoryFile       string `yaml:"history_file"`
	PersistHistory    bool   `yaml:"persist_history"`
	Debug             bool   `yaml:"debug"`
	Tips              bool   `yaml:"tips"`
	Graph             bool   `yaml:"graph"`
	Start             string `yaml:"start"`
	End               string `yaml:"end"`
	Step              string `yaml:"step"`
}

// NewConfig returns a Config with default values.
func NewConfig() *Config {
	return &Config{
		URL:               "http://localhost:9090",
		EnableLabelValues: true,
		Tips:              false,
	}
}

// LoadFromFile reads the configuration from a YAML file.
func LoadFromFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := NewConfig() // Start with defaults
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, err
	}

	return config, nil
}
