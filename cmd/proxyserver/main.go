package main

import (
	"flag"
	"log"
	"reverse-proxy/proxyserver"
)

func main() {
	var f string
	flag.StringVar(&f, "config-file", "config/config.toml", "config file")
	flag.Parse()

	config := proxyserver.NewConfig()
	if err := config.GetFromFile(f); err != nil {
		log.Fatal("Can`t get config from file: ", err)
	}

	pServer := proxyserver.New(config)
	err := pServer.Start()
	if err != nil {
		log.Fatal(err)
	}
}
