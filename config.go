package main

import (
	"log"
	"path/filepath"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	Token string
	Listen string
	Channel map[string]ChannelDef
}

type ChannelDef struct {
	Id string
}

func readConfigFile() Config {
	var configfile = os.Getenv("ALERTBOT_CONFIG")
	if configfile == "" {
		configfile, _  = filepath.Abs("./config.yml")
	}
	yamlFile, err := ioutil.ReadFile(configfile)

	if err != nil {
		log.Panic(err)
	}

	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Panic(err)
	}
	return config
}