package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(databaseURL string) (*pgxpool.Pool, error) {
	dbpool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к базу данных: %w ", err)
	}

	if err = dbpool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("Не удалось выполнить ping: %w ", err)
	}

	return dbpool, nil

}
