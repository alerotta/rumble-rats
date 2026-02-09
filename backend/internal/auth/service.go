package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	store *Store
	jwtSecret []byte
	issuer string
	audience  string
	accessTTL time.Duration
}

func NewService(store *Store, jwtSecret string) *Service {
	return &Service{
		store: store,
		jwtSecret: []byte(jwtSecret),
		issuer:    "rumble-rats",
		audience:  "rumble-rats-web",
		accessTTL: 15 * time.Minute,
	}
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
		return RegisterResponse{}, fmt.Errorf("bcrypt generate hash: %w", err)
	}

	user, err := s.store.CreateUser(ctx, req.Username, req.Email, string(hash))
		if err != nil {
		if errors.Is(err, ErrUserAlreadyExists) {
			return RegisterResponse{}, err
		}
		return RegisterResponse{}, fmt.Errorf("create user: %w", err)
	}

	return RegisterResponse{Username: user.Username, Email: user.Email}, nil
	
}

func (s *Service) Login(ctx context.Context , req LoginRequest) (AuthResponse, error){

	req.Username = strings.TrimSpace(req.Username)
	if req.Username == "" || req.Password == "" {
		return AuthResponse{}, ErrValidation("username, password are required")
	}
	if len(req.Password) < 8 {
		return AuthResponse{}, ErrValidation("password must be at least 8 characters")
	}

	stored, err := s.store.GetPasswordHashByUsername(ctx, req.Username)
		if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return AuthResponse{}, ErrUnauthorized
		}
		return AuthResponse{}, fmt.Errorf("get password hash: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(stored.PasswordHash),
		[]byte(req.Password),
		); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return AuthResponse{}, ErrUnauthorized
		}
		return AuthResponse{}, fmt.Errorf("bcrypt compare: %w", err)
	}

	token, err := s.signAccessToken(req.Username) 
	if err != nil {
		return AuthResponse{}, fmt.Errorf("sign access token: %w", err)
	}

	return AuthResponse{Username: req.Username , Token: token}, nil

}

func (s *Service) signAccessToken (username string)(string ,error){
	now := time.Now()

		claims := jwt.MapClaims{                
		"sub": username,               
		"iss":      s.issuer,               
		"aud":      s.audience,             
		"iat":      now.Unix(),             
		"exp":      now.Add(s.accessTTL).Unix(), 
	}

	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tok.SignedString(s.jwtSecret)
}