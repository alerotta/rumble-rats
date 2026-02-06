package db

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func Open (dsn string) (*sql.DB, error){
	conn, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	conn.SetMaxOpenConns(20)
	conn.SetMaxIdleConns(10)
	conn.SetConnMaxLifetime(30 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 2* time.Second)
	defer cancel()

	if err := conn.PingContext(ctx); err != nil {
		conn.Close()
		return nil, err
	}
	return conn, nil
}

