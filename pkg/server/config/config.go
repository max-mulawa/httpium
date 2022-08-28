package config

import (
	"errors"
	"os"

	"github.com/BurntSushi/toml"
)

type ServerOptions struct {
	Port uint
}

type ContentOptions struct {
	StaticDir string
	Default   []string
}

type HttpiumConfig struct {
	Server  ServerOptions
	Content ContentOptions
}

var defaultPort uint = 8080

func NewHttpiumConfig() *HttpiumConfig {
	return &HttpiumConfig{
		Server: ServerOptions{
			Port: defaultPort,
		},
		Content: ContentOptions{
			StaticDir: "./static",
		},
	}
}

func (c *HttpiumConfig) Load() error {
	configPath := "./config.toml"
	_, err := os.Stat(configPath)

	if errors.Is(err, os.ErrNotExist) {
		return nil
	}

	if _, err := toml.DecodeFile(configPath, c); err != nil {
		return err
	}

	return nil
}
