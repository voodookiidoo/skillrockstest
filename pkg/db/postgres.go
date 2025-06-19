package db

import (
	"context"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	_ "github.com/lib/pq"
	"os"
)

func MustConnect() (*pgx.Conn, *migrate.Migrate) {

	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	)
	config, err := pgx.ParseConfig(url)
	if err != nil {
		panic(err)
	}

	c, err := pgx.ConnectConfig(context.Background(), config)
	if err != nil {
		panic(err)
	}
	m, err := migrate.New("file:///opt/skill-rocks/migrations", url)
	if err != nil {
		panic(err)
	}
	version, dirty, err := m.Version()
	if dirty {
		m.Force(0)
		m.Up()
	} else if version == 0 {
		m.Up()
	}
	return c, m
}
