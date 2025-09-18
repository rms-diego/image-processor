package database

import (
	"database/sql"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	_ "github.com/lib/pq"
)

var DB *goqu.Database

func Init(dsn string) error {
	dialect := goqu.Dialect("postgres")
	pgDb, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}

	if err := pgDb.Ping(); err != nil {
		return err
	}

	DB = dialect.DB(pgDb)
	return nil
}
