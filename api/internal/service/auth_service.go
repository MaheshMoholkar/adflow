package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"callflow/internal/domain/auth"
	"callflow/internal/domain/user"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const accessTokenExpiry = 365 * 24 * time.Hour // 1 year

// AuthService provides authentication functionality
type AuthService struct {
	userRepo  user.Repository
	jwtSecret []byte
}

// NewAuthService creates a new auth service instance
func NewAuthService(userRepo user.Repository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: []byte(jwtSecret),
	}
}

// Register creates a new user with phone and password
func (s *AuthService) Register(ctx context.Context, req auth.RegisterRequest) (*auth.TokenResponse, error) {
	_, err := s.userRepo.GetByPhone(ctx, req.Phone)
	if err == nil {
		return nil, auth.ErrPhoneTaken
	}
	if !errors.Is(err, user.ErrUserNotFound) {
		return nil, fmt.Errorf("failed to lookup user: %w", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	u, err := s.userRepo.Create(ctx, user.UserCreate{
		Phone:        req.Phone,
		PasswordHash: string(hash),
		Name:         req.Name,
		BusinessName: req.BusinessName,
		City:         req.City,
		Address:      req.Address,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return s.tokenResponseForUser(u)
}

// Login authenticates a user with phone and password
func (s *AuthService) Login(ctx context.Context, req auth.LoginRequest) (*auth.TokenResponse, error) {
	u, err := s.userRepo.GetByPhone(ctx, req.Phone)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			return nil, auth.ErrInvalidCredentials
		}
		return nil, fmt.Errorf("failed to lookup user: %w", err)
	}

	if u.Status != user.StatusActive {
		return nil, auth.ErrUserInactive
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)); err != nil {
		return nil, auth.ErrInvalidCredentials
	}

	return s.tokenResponseForUser(u)
}

// VerifyToken verifies a JWT token and returns the claims
func (s *AuthService) VerifyToken(_ context.Context, tokenString string) (*auth.AuthClaims, error) {
	t, err := jwt.ParseWithClaims(tokenString, &auth.AuthClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, auth.ErrExpiredToken
		}
		return nil, auth.ErrInvalidToken
	}

	if !t.Valid {
		return nil, auth.ErrInvalidToken
	}

	claims, ok := t.Claims.(*auth.AuthClaims)
	if !ok {
		return nil, auth.ErrInvalidToken
	}

	return claims, nil
}

func (s *AuthService) tokenResponseForUser(u *user.User) (*auth.TokenResponse, error) {
	expiresAt := time.Now().Add(accessTokenExpiry)

	claims := &auth.AuthClaims{
		UserID: u.ID,
		Phone:  u.Phone,
		Plan:   u.Plan,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "callflow-api",
			Subject:   fmt.Sprintf("%d", u.ID),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := t.SignedString(s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to sign token: %w", err)
	}

	return &auth.TokenResponse{
		AccessToken: tokenString,
		TokenType:   "Bearer",
		ExpiresAt:   expiresAt,
		User: &auth.UserInfo{
			ID:            u.ID,
			Phone:         u.Phone,
			Name:          u.Name,
			BusinessName:  u.BusinessName,
			City:          u.City,
			Address:       u.Address,
			LocationURL:   u.LocationURL,
			Plan:          u.Plan,
			PlanStartedAt: u.PlanStartedAt,
			PlanExpiresAt: u.PlanExpiresAt,
			Status:        u.Status,
		},
	}, nil
}
