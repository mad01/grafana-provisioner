package main

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Teams []string `yaml:teams`
}

func GetConfig(filepath string) *Config {
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Printf("yamlFile.Get err   #%v ", err)
	}
	c := &Config{}
	err = yaml.Unmarshal(file, c)
	if err != nil {
		fmt.Printf("Unmarshal: %v", err)
	}

	return c
}
