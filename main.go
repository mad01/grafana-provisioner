package main

import (
	"fmt"
)

var grafanasData = `
teams:
  - foo
  - bar
  - baz
`

type TeamsConfig struct {
	teams []string `yaml:teams`
}

type DatabaseConfig struct {
	dnsBase     string `yaml:"dnsBase"`
	databaseURL string `yaml:"databaseURL"`
}

func main() {
	err := runCmd()
	if err != nil {
		fmt.Println(err.Error())
	}
}
