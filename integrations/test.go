package integrations

import (
	"flag"

	"github.com/AndreiGoStorm/previewer/internal/config"
)

var configFile string

var nginxHost = "100.90.15.4"

func init() {
	flag.StringVar(&configFile, "config", "../configs/config-testing.yml", "Path to configuration file")
}

func SetupSuite() (conf *config.Config) {
	flag.Parse()
	conf = config.New(configFile)
	if conf == nil {
		panic("config file is invalid")
	}
	return
}
