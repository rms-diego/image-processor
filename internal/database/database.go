package database

import (
	"database/sql"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	_ "github.com/lib/pq"
	"github.com/rms-diego/image-processor/internal/config"
)

var Db *goqu.Database

func Init() error {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%v sslmode=disable",
		config.Env.DB_HOST,
		config.Env.DB_PORT,
		config.Env.DB_USER,
		config.Env.DB_PASSWORD,
		config.Env.DB_NAME,
	)

	dialect := goqu.Dialect("postgres")
	pgDb, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}

	if err := pgDb.Ping(); err != nil {
		return err
	}

	Db = dialect.DB(pgDb)
	return nil
}
