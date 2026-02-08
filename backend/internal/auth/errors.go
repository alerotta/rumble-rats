package auth

import "errors"

var (
	ErrConflict = errors.New("conflict")
	ErrNotFound   = errors.New("not found") 
	ErrUnauthorized = errors.New("unauthorized")
)

type ValidationError struct{ Msg string }
func (e ValidationError) Error() string { return e.Msg }

func ErrValidation(msg string) error { return ValidationError{Msg: msg} }