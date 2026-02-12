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
	refreshTTL time.Duration
}

func NewService(store *Store, jwtSecret string) *Service {
	return &Service{
		store: store,
		jwtSecret: []byte(jwtSecret),
		issuer:    "rumble-rats",
		audience:  "rumble-rats-web",
		accessTTL: 15 * time.Minute,
		refreshTTL: 4 * time.Hour,
	}
}

func (s *Service) Register(ctx context.Context , req RegisterRequest) (AuthResponse,RefreshToken, error){
	
	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	if req.Username == "" || req.Email == "" || req.Password == "" {
		return AuthResponse{},RefreshToken{}, ErrValidation("username, email and password are required")
	}
	if len(req.Password) < 8 {
		return AuthResponse{},RefreshToken{}, ErrValidation("password must be at least 8 characters")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return AuthResponse{},RefreshToken{}, fmt.Errorf("bcrypt generate hash: %w", err)
	}

	user, err := s.store.CreateUser(ctx, req.Username, req.Email, string(hash))
		if err != nil {
		if errors.Is(err, ErrUserAlreadyExists) {
			return AuthResponse{},RefreshToken{}, err
		}
		return AuthResponse{},RefreshToken{}, fmt.Errorf("create user: %w", err)
	}

	token, err := s.signAccessToken(user.Username) 
	if err != nil {
		return AuthResponse{},RefreshToken{}, fmt.Errorf("sign access token: %w", err)
	}

	refreshToken , exp , err := s.signRefreshToken(req.Username)
	if err != nil {
		return AuthResponse{},RefreshToken{}, fmt.Errorf("sign refresh token: %w", err)
	}


	return AuthResponse{Username: user.Username, Token: token},RefreshToken{refreshToken,exp}, nil
	
}

func (s *Service) Login(ctx context.Context , req LoginRequest) (AuthResponse,RefreshToken, error){

	req.Username = strings.TrimSpace(req.Username)
	if req.Username == "" || req.Password == "" {
		return AuthResponse{},RefreshToken{}, ErrValidation("username, password are required")
	}
	if len(req.Password) < 8 {
		return AuthResponse{},RefreshToken{}, ErrValidation("password must be at least 8 characters")
	}

	stored, err := s.store.GetPasswordHashByUsername(ctx, req.Username)
		if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return AuthResponse{},RefreshToken{}, ErrUnauthorized
		}
		return AuthResponse{},RefreshToken{}, fmt.Errorf("get password hash: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(stored.PasswordHash),
		[]byte(req.Password),
		); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return AuthResponse{},RefreshToken{}, ErrUnauthorized
		}
		return AuthResponse{},RefreshToken{}, fmt.Errorf("bcrypt compare: %w", err)
	}

	accessToken, err := s.signAccessToken(req.Username) 
	if err != nil {
		return AuthResponse{},RefreshToken{}, fmt.Errorf("sign access token: %w", err)
	}

	refreshToken , exp , err := s.signRefreshToken(req.Username)
	if err != nil {
		return AuthResponse{},RefreshToken{}, fmt.Errorf("sign refresh token: %w", err)
	}


	return AuthResponse{Username: req.Username , Token: accessToken },RefreshToken{refreshToken,exp}, nil

}

func (s *Service) Refresh(ctx context.Context,req RefreshRequest )(AuthResponse,RefreshToken, error){

	claims , err :=s.ValidateToken(req.RefreshToken)
	if err != nil {
		return AuthResponse{},RefreshToken{}, fmt.Errorf("invalid refresh token: %w", err)
	}
	if claims.Type != "refresh" {
		return AuthResponse{},RefreshToken{}, fmt.Errorf("invalid refresh token: %w", err)
	}
	newAccessToken, err := s.signAccessToken(claims.Subject) 
	if err != nil {
		return AuthResponse{},RefreshToken{}, fmt.Errorf("sign access token: %w", err)
	}

	return AuthResponse{Username: claims.Subject , Token: newAccessToken },RefreshToken{}, nil

}
// helper functions TODO: move away  

func (s *Service) ValidateToken (tokenString string) (*CustomClaims, error){

	claims := &CustomClaims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func (t *jwt.Token)(any,error){
			if t.Method != jwt.SigningMethodHS256{
				return nil , fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return s.jwtSecret, nil
		},
		jwt.WithIssuer(s.issuer),
		jwt.WithAudience(s.audience),
		jwt.WithLeeway(30*time.Second),
	)
	if err != nil {
		return nil, err 
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	if claims.Subject == "" {
		return nil, errors.New("missing sub")
	}

	if claims.Type == "" {
		return nil, errors.New("missing token type")
	}
	return claims, nil
}

func (s *Service) signAccessToken (username string)(string ,error){
	now := time.Now()

	claims := jwt.MapClaims{                
		"sub": username,               
		"iss": s.issuer,               
		"aud": s.audience,             
		"iat": now.Unix(),             
		"exp": now.Add(s.accessTTL).Unix(),
		"typ": "access" ,
	}

	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tok.SignedString(s.jwtSecret)
}

func (s *Service) signRefreshToken (username string)(string,time.Time,error){
	now := time.Now()
	exp := now.Add(s.refreshTTL)
	claims := jwt.MapClaims{
		"sub": username,               
		"iss": s.issuer,               
		"aud": s.audience,  
		"iat": now.Unix(),
		"exp": exp.Unix(),
		"typ": "refresh",
		}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := tok.SignedString(s.jwtSecret)
	if err != nil {
		return "", time.Time{}, err
	}
	return signedToken, exp, nil
}