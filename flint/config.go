package flint

import (
	"io"

	"gopkg.in/yaml.v3"
)

type Content struct {
	Name  string   `yaml:"name"`
	Paths []string `yaml:"paths"`
	Rules []string `yaml:"rules"`
}

type Config struct {
	Rules   map[string]Rule `yaml:"rules"`
	Content []Content       `yaml:"content"`
}

func ReadConfig(r io.Reader) (*Config, error) {
	var c Config

	err := yaml.NewDecoder(r).Decode(&c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
