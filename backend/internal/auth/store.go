package auth

import (
	"context"
	"database/sql"
)

type Store struct {
	db *sql.DB
}

func NewStore (db *sql.DB) *Store{
	return &Store{ db: db }
}
type UserRow struct {
	Username string
	Email    string
}

func (s *Store) CreateUser(ctx context.Context, username, email, passwordHash string) (UserRow, error) {
	var u UserRow
	err := s.db.QueryRowContext(ctx, `
		INSERT INTO users (username, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING username, email
	`, username, email, passwordHash).Scan(&u.Username, &u.Email)
	return u, err
}