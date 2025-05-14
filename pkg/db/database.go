package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type Database struct {
	*sql.DB
}

func New(host string, port int, user, password, dbName string) (*Database, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", user, password, host, port, dbName)
	sqlDB, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	return &Database{sqlDB}, nil
}

func (d *Database) ExecQuery(query string) error {
	_, err := d.Exec(query)
	return err
}
