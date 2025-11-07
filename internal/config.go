package app

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
)

type AppConfig struct {
	BaseUrl string `toml:"base_url"`
}

type DatabaseConfig struct {
	Url string
}

type FrontendConfig struct {
	BaseUrl string `toml:"base_url"`
}

type ServerConfig struct {
	Bind string
	Port uint16
}

func (c *ServerConfig) BindAddress() string {
	return fmt.Sprintf("%s:%d", c.Bind, c.Port)
}

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Frontend FrontendConfig
	Server   ServerConfig
}

func LoadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open configuration file %s: %s", path, err.Error())
	}
	defer f.Close()

	var config Config
	decoder := toml.NewDecoder(f)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse configuration file %s: %s", path, err.Error())
	}

	return &config, nil
}
