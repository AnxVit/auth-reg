package postgres

import (
	"Auth-Reg/internal/config"
	"Auth-Reg/internal/storage"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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

func (s *Storage) SaveUser(name, email, password string) error {
	const op = "storage.postgres.SaveUser"

	_, err := s.db.Exec(context.Background(), "INSERT INTO Users (name, email, password) values ($1, $2, $3)",
		name, email, password)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return fmt.Errorf("%s: %w", op, storage.ErrUserExists)
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
