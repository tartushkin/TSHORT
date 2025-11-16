package db

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type repo struct {
	conn *sql.DB
}

func NewConnection(ps string) (*sql.DB, error) {
	repo := repo{}

	db, err := sql.Open("pgx", ps)
	if err != nil {
		return nil, err
	}
	//defer db.Close()

	repo.conn = db
	//проверкаа подключения
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err = db.PingContext(ctx); err != nil {
		panic(err)
	}

	return db, nil
}
