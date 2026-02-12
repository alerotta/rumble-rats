package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RefreshRequest struct {
	RefreshToken string 
}

type AuthResponse struct {
	Username string `json:"username"`
	Token string `json:"token"`
}

type RefreshToken struct {
	Token string 
	exp time.Time 
}

type VerifyTokenRequest struct {
	Token string `json:"token"`
}

type CustomClaims struct {
	Type string `json:"typ"`
	jwt.RegisteredClaims
}