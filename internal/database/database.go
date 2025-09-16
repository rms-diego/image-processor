package database

import (
	"database/sql"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	_ "github.com/lib/pq"
	"github.com/rms-diego/image-processor/pkg/config"
)

var Db *goqu.Database

func Init() error {
	dialect := goqu.Dialect("postgres")
	pgDb, err := sql.Open("postgres", config.ServerEnv.DATABASE_URL)
	if err != nil {
		return err
	}

	if err := pgDb.Ping(); err != nil {
		return err
	}

	Db = dialect.DB(pgDb)
	return nil
}
