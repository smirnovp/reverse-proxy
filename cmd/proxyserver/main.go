package main

import (
	"flag"
	"log"
	"path"
	"reverse-proxy/cache"
	"reverse-proxy/proxyserver"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	var f string
	flag.StringVar(&f, "config-file", "config/config.toml", "config file")
	flag.Parse()

	logger := logrus.New()
	if err := configureLogger(logger, f); err != nil {
		log.Fatal("Can`t configure logger: ", err)
	}

	config := proxyserver.NewConfig()
	if err := config.GetFromFile(f); err != nil {
		log.Fatal("Can`t get config: ", err)
	}

	cacheConfig := cache.NewConfig()
	if err := cacheConfig.GetFromFile(f); err != nil {
		log.Fatal("Can`t get cache config: ", err)
	}

	cache := cache.New(logger, cacheConfig)

	pServer := proxyserver.New(logger, config, cache)
	err := pServer.Start()
	if err != nil {
		log.Fatal(err)
	}
}

func configureLogger(l *logrus.Logger, f string) error {

	viper.SetConfigFile(path.Clean(f))
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	viper.SetDefault("logger.level", "debug")
	ls := viper.GetString("logger.level")

	level, err := logrus.ParseLevel(ls)
	if err != nil {
		return err
	}

	l.SetLevel(level)

	return nil
}
