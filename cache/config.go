package cache

import (
	"path"

	"github.com/spf13/viper"
)

// Config ...
type Config struct {
	Size int
	Dir  string
}

// NewConfig ...
func NewConfig() *Config {
	return &Config{
		Size: 10,
		Dir:  "cache/files",
	}
}

// GetFromFile ...
func (c *Config) GetFromFile(f string) error {

	viper.SetConfigFile(path.Clean(f))
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	viper.SetDefault("cache.size", 10)
	viper.SetDefault("cache.dir", 10)

	c.Size = viper.GetInt("cache.size")
	c.Dir = viper.GetString("cache.dir")
	return nil
}
