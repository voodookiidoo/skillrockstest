package main

import (
	"context"
	"github.com/jackc/pgx/v5"
	"testing"
)

func TestConnection(t *testing.T) {
	url := "host=localhost port=5432 user=user password=pass sslmode=disable"
	//url := "postgresql://postgres:postgres@localhost/postgres?sslmode=disable"
	c, err := pgx.Connect(context.Background(), url)
	
	if err != nil {
		t.Error(err)
	} else {
		t.Log(c)
	}
	
}
