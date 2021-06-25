package main

import (
	"flag"
	"log"
	"reverse-proxy/cache"
	"reverse-proxy/proxyserver"
)

func main() {
	var f string
	flag.StringVar(&f, "config-file", "config/config.toml", "config file")
	flag.Parse()

	config := proxyserver.NewConfig()
	if err := config.GetFromFile(f); err != nil {
		log.Fatal("Can`t get config: ", err)
	}

	cacheConfig := cache.NewConfig()
	if err := cacheConfig.GetFromFile(f); err != nil {
		log.Fatal("Can`t get cache config: ", err)
	}

	cache := cache.New(cacheConfig)

	pServer := proxyserver.New(config, cache)
	err := pServer.Start()
	if err != nil {
		log.Fatal(err)
	}
}
