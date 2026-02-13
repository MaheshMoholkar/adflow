package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// AuthClaims defines the claims in JWT tokens
type AuthClaims struct {
	UserID int64  `json:"user_id"`
	Phone  string `json:"phone"`
	Plan   string `json:"plan"`
	jwt.RegisteredClaims
}

// RegisterRequest represents a request to register a new user
type RegisterRequest struct {
	Phone        string `json:"phone" validate:"required"`
	Password     string `json:"password" validate:"required,min=6"`
	Name         string `json:"name"`
	BusinessName string `json:"business_name"`
	City         string `json:"city"`
	Address      string `json:"address"`
}

// LoginRequest represents a request to log in
type LoginRequest struct {
	Phone    string `json:"phone" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// TokenResponse represents the response returned after authentication
type TokenResponse struct {
	AccessToken string    `json:"access_token"`
	TokenType   string    `json:"token_type"`
	ExpiresAt   time.Time `json:"expires_at"`
	User        *UserInfo `json:"user"`
}

// UserInfo represents user data included in auth responses
type UserInfo struct {
	ID            int64      `json:"id"`
	Phone         string     `json:"phone"`
	Name          string     `json:"name,omitempty"`
	BusinessName  string     `json:"business_name,omitempty"`
	City          string     `json:"city,omitempty"`
	Address       string     `json:"address,omitempty"`
	LocationURL   string     `json:"location_url,omitempty"`
	Plan          string     `json:"plan"`
	PlanStartedAt *time.Time `json:"plan_started_at,omitempty"`
	PlanExpiresAt *time.Time `json:"plan_expires_at,omitempty"`
	Status        string     `json:"status"`
}
