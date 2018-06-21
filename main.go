package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var configData = `
databaseURL: root:qwerty@localhost/
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

type DB struct {
	URL  string
	conn *sql.DB
}

func (d *DB) connect() {
	db, err := sql.Open("mysql", d.URL)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	d.conn = db
}

func (d *DB) createDB(name string) error {
	query := "CREATE DATABASE IF NOT EXISTS name = $1;"
	rows, err := d.conn.Query(
		query,
		name,
	)
	defer rows.Close()
	return err
}

func main() {
	err := runCmd()
	if err != nil {
		fmt.Println(err.Error())
	}
}
