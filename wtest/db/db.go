package db

import "fmt"

type DB struct {
	DSN string
}

func NewDB(dsn string) *DB {
	fmt.Println("Connecting to DB:", dsn)
	return &DB{DSN: dsn}
}
