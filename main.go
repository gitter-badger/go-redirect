package main

import (
	"flag"
	"fmt"

	"github.com/0rax/go-redirect/backend"
	"github.com/BurntSushi/toml"
)

type softwareConfig struct {
	WebServer string
	Backend   backend.Config
}

func main() {

	var config softwareConfig

	var cfgFile = flag.String("config", "config.toml", "path to config file")
	flag.Parse()

	if _, err := toml.DecodeFile(*cfgFile, &config); err != nil {
		fmt.Printf("[go-redirect] %s\n", err)
		return
	}
	backend.Configure(&config.Backend)
	runWebServer(config.WebServer)
}
