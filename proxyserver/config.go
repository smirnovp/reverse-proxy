package proxyserver

import (
	"path"

	"github.com/spf13/viper"
)

// Config ...
type Config struct {
	URL  string `toml:"url"`
	Port string `toml:"port"`
}

// NewConfig ...
func NewConfig() *Config {
	return &Config{
		URL:  "habr.com",
		Port: ":8080",
	}
}

// GetFromFile ...
func (c *Config) GetFromFile(f string) error {

	viper.SetConfigFile(path.Clean(f))
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	viper.SetDefault("proxyserver.url", "habr.com")
	viper.SetDefault("proxyserver.port", ":8080")

	c.URL = viper.GetString("proxyserver.url")
	c.Port = viper.GetString("proxyserver.port")
	return nil
}
