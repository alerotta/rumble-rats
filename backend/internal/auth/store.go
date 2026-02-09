package auth

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
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
type UserPasswordRow struct {
	PasswordHash string
}

func (s *Store) CreateUser(ctx context.Context, username, email, passwordHash string) (UserRow, error) {
	var u UserRow
	err := s.db.QueryRowContext(ctx, `
		INSERT INTO users (username, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING username, email
	`, username, email, passwordHash).Scan(&u.Username, &u.Email)

	if err == nil {
		return u, nil
	}

	if  isUniqueViolation(err){
		return UserRow{}, ErrUserAlreadyExists
	}
	return UserRow{}, err
}

func (s *Store) GetPasswordHashByUsername (ctx context.Context, username string) (UserPasswordRow,error){
	var u UserPasswordRow
	err := s.db.QueryRowContext(ctx, `
	SELECT password_hash
	FROM users
	WHERE username = $1`,
	username).Scan(&u.PasswordHash)

	if err == nil {
		return u , nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return UserPasswordRow{}, ErrUserNotFound
	}
	return UserPasswordRow{}, err

}

// error from db management

func isUniqueViolation (err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr){
		return pgErr.Code == "23505"
	}
	return false
}
