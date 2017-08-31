package main

import (
	"io/ioutil"
	"github.com/go-yaml/yaml"
	"github.com/davecgh/go-spew/spew"
)

func readConfiguration() {
	configbytes, err := ioutil.ReadFile("porkmeterconfig.yaml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(configbytes, &config)
	if err != nil {
		panic(err)
	}
	spew.Dump(config)
}

func writeConfiguration() {
	spew.Dump(config)
	configbytes, err := yaml.Marshal(&config)
	if err != nil {
		panic(err)
	}
    err = ioutil.WriteFile("porkmeterconfig.yaml", configbytes, 0644)
	if err != nil {
		panic(err)
	}
}