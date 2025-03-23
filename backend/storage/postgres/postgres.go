package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/lib/pq"
	"github.com/sfu-teamproject/smartbuy/backend/storage"
)

type PostgresDB struct {
	*sql.DB
}

func NewPostgresDB() (*PostgresDB, error) {
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	pass := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, pass, dbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return &PostgresDB{db}, nil
}

func (db *PostgresDB) wrapError(err error) error {
	if err == nil {
		return nil
	}
	var pqErr *pq.Error
	var errWrapper = storage.ErrInternal
	if errors.As(err, &pqErr) {
		switch pqErr.Code {
		case "23505":
			errWrapper = storage.ErrAlreadyExists
		case "0200":
			errWrapper = storage.ErrNotFound
		}
	} else if errors.Is(err, sql.ErrNoRows) {
		errWrapper = storage.ErrNotFound
	}
	err = fmt.Errorf("%w: %w", errWrapper, err)
	return err
}
