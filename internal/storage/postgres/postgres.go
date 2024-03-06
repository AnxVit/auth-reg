package postgres

import (
	"Auth-Reg/internal/config"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type Storage struct {
	db *pgx.Conn
}

func New(storage config.Storage) (*Storage, error) {
	const op = "storage.postgres.New"
	psqlInfo := fmt.Sprintf("user=%s password=%s host=%s "+
		"port=%d dbname=%s sslmode=disable",
		storage.User, storage.Password, storage.Host, storage.Port, storage.DBName)
	db, err := pgx.Connect(context.Background(), psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &Storage{
		db: db,
	}, nil
}
