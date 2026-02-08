package auth

import (
	"context"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	store *Store
}

func NewService(store *Store) *Service {
	return &Service{store: store}
}

func (s *Service) Register(ctx context.Context , req RegisterRequest) (RegisterResponse, error){
	
	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	if req.Username == "" || req.Email == "" || req.Password == "" {
		return RegisterResponse{}, ErrValidation("username, email and password are required")
	}
	if len(req.Password) < 8 {
		return RegisterResponse{}, ErrValidation("password must be at least 8 characters")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return RegisterResponse{}, err
	}

	user, err := s.store.CreateUser(ctx, req.Username, req.Email, string(hash))
	if err != nil {
		return RegisterResponse{}, err
	}

	return RegisterResponse{Username: user.Username, Email: user.Email}, nil
	
}