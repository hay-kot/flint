package flint

import (
	"encoding/json"
	"io"

	"github.com/pelletier/go-toml/v2"
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
	var err error

	switch format {
	case JSON:
		err = json.NewDecoder(r).Decode(&c)
	case TOML:
		err = toml.NewDecoder(r).Decode(&c)
	case YAML:
		err = yaml.NewDecoder(r).Decode(&c)
	}

	return &c, err
}
