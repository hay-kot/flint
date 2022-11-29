package flint

import (
	"encoding/json"
	"io"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

type Content struct {
	Name  string   `yaml:"name" toml:"name" json:"name"`
	Paths []string `yaml:"paths" toml:"paths" json:"paths"`
	Rules []string `yaml:"rules" toml:"rules" json:"rules"`
}

type Config struct {
	Rules   map[string]Rule `yaml:"rules" toml:"rules" json:"rules"`
	Content []Content       `yaml:"content" toml:"content" json:"content"`
}

type ConfigFormat string

const (
	JSON ConfigFormat = "json"
	TOML ConfigFormat = "toml"
	YAML ConfigFormat = "yaml"
)

func ReadConfig(r io.Reader, format ConfigFormat) (*Config, error) {
	var c Config

	switch format {
	case JSON:
		if err := json.NewDecoder(r).Decode(&c); err != nil {
			return nil, err
		}
	case TOML:
		if _, err := toml.DecodeReader(r, &c); err != nil {
			return nil, err
		}
	case YAML:
		if err := yaml.NewDecoder(r).Decode(&c); err != nil {
			return nil, err
		}
	}

	return &c, nil
}
