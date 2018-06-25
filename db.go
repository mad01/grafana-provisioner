package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

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
	fmt.Printf("connected to database: %v\n", d.URL) // todo: doint print full connnect string with pw/user
	return nil
}

func (d *DB) createDB(name string) error {
	_, err := d.conn.Exec("CREATE DATABASE IF NOT EXISTS " + name)
	return err
}
