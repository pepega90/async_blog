package utils

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

func ConnectDB() *pgx.Conn {
	dsn := "postgres://learn:learn@localhost:5432/db_learn"
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable connect to databse: %v\n", err)
		os.Exit(1)
	}

	return conn
}
