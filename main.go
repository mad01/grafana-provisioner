package main

import (
	"database/sql"
	"fmt"

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

func (d *DB) connect() error {
	db, err := sql.Open("mysql", d.URL)
	if err != nil {
		return err
	}
	d.conn = db
	fmt.Printf("connected to database: %v\n", d.URL)
	return nil
}

func (d *DB) createDB(name string) error {
	query := "CREATE DATABASE IF NOT EXISTS name = $1;"
	_, err := d.conn.Query(
		query,
		name,
	)
	return err
}

func main() {
	err := runCmd()
	if err != nil {
		fmt.Println(err.Error())
	}
}
