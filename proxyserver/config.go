package proxyserver

import (
	"path"

	"github.com/spf13/viper"
)

// Config ...
type Config struct {
	URL      string `toml:"url"`
	Port     string `toml:"port"`
	LogLevel string `toml:"loglevel"`
}

// NewConfig ...
func NewConfig() *Config {
	return &Config{
		URL:      "habr.com",
		Port:     ":8080",
		LogLevel: "debug",
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
	viper.SetDefault("logger.loglevel", "debug")

	c.URL = viper.GetString("proxyserver.url")
	c.Port = viper.GetString("proxyserver.port")
	c.LogLevel = viper.GetString("logger.loglevel")
	return nil
}
