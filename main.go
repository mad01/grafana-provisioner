package main

import "fmt"

var configData = `
dnsBase: example.com
databaseURL: username:password@example.com/dashboards
`

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

type DB struct{}

func (d *DB) connect()  {}
func (d *DB) createDB() {}
func (d *DB) removeDB() {}
func (d *DB) listDB()   {}

func main() {
	err := runCmd()
	if err != nil {
		fmt.Println(err.Error())
	}
}
