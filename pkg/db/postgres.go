package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"os"
)

var prod = os.Getenv("POSTGRES_DB") != ""

func MustConnect() *pgx.Conn {
	var err error
	var url string
	if prod {
		url = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			os.Getenv("POSTGRES_HOST"),
			os.Getenv("POSTGRES_PORT"),
			os.Getenv("POSTGRES_USER"),
			os.Getenv("POSTGRES_PASSWORD"),
			os.Getenv("POSTGRES_DB"),
		)
	} else {
		url = "host=localhost port=5432 user=user password=pass dbname=postgres sslmode=disable"
	}
	config, err := pgx.ParseConfig(url)
	if err != nil {
		panic(err)
	}
	c, err := pgx.ConnectConfig(context.Background(), config)
	if err != nil {
		panic(err)
	}
	return c
}
