package repository

import (
	"context"
	"errors"
	"time"

	"callflow/internal/domain/token"
	db "callflow/internal/sql/db"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TokenRepository implements token.Repository
type TokenRepository struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

// NewTokenRepository creates a new token repository
func NewTokenRepository(pool *pgxpool.Pool) *TokenRepository {
	return &TokenRepository{
		pool:    pool,
		queries: db.New(pool),
	}
}

func (r *TokenRepository) CreateRefreshToken(ctx context.Context, data token.TokenCreate) (*token.Token, error) {
	row, err := r.queries.CreateToken(ctx, db.CreateTokenParams{
		UserID:    data.UserID,
		Token:     data.Token,
		ExpiresAt: pgtype.Timestamptz{Time: data.ExpiresAt, Valid: true},
		TokenType: data.TokenType,
		ClientIp:  data.ClientIP,
		UserAgent: data.UserAgent,
	})
	if err != nil {
		return nil, err
	}

	return dbTokenToModel(row), nil
}

func (r *TokenRepository) GetRefreshTokenByToken(ctx context.Context, tokenStr string) (*token.Token, error) {
	row, err := r.queries.GetTokenByToken(ctx, tokenStr)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, token.ErrTokenNotFound
		}
		return nil, err
	}

	t := dbTokenToModel(row)

	if t.IsRevoked {
		return t, token.ErrTokenRevoked
	}
	if time.Now().After(t.ExpiresAt) {
		return t, token.ErrTokenExpired
	}

	return t, nil
}

func (r *TokenRepository) RevokeToken(ctx context.Context, tokenStr string) error {
	return r.queries.RevokeToken(ctx, tokenStr)
}

func (r *TokenRepository) RevokeAllUserTokens(ctx context.Context, userID int64) error {
	return r.queries.RevokeAllUserTokens(ctx, userID)
}

func (r *TokenRepository) UpdateTokenLastUsed(ctx context.Context, id int64) error {
	return r.queries.UpdateTokenLastUsed(ctx, id)
}

func (r *TokenRepository) DeleteExpiredTokens(ctx context.Context, before time.Time) error {
	return r.queries.DeleteExpiredTokens(ctx, pgtype.Timestamptz{Time: before, Valid: true})
}

func (r *TokenRepository) StoreTokenID(ctx context.Context, tokenID string, userID int64, expiresAt time.Time, tokenType string) error {
	_, err := r.queries.CreateToken(ctx, db.CreateTokenParams{
		UserID:    userID,
		Token:     tokenID,
		ExpiresAt: pgtype.Timestamptz{Time: expiresAt, Valid: true},
		TokenType: tokenType,
	})
	return err
}

func (r *TokenRepository) IsTokenRevoked(ctx context.Context, tokenID string) (bool, error) {
	row, err := r.queries.GetTokenByToken(ctx, tokenID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, token.ErrTokenNotFound
		}
		return false, err
	}
	return row.IsRevoked, nil
}

func (r *TokenRepository) RevokeTokenByID(ctx context.Context, tokenID string) error {
	return r.queries.RevokeToken(ctx, tokenID)
}

func (r *TokenRepository) RevokeAllTokensByType(ctx context.Context, userID int64, tokenType string) error {
	return r.queries.RevokeAllUserTokensByType(ctx, db.RevokeAllUserTokensByTypeParams{
		UserID:    userID,
		TokenType: tokenType,
	})
}

func dbTokenToModel(row db.Token) *token.Token {
	return &token.Token{
		ID:         row.ID,
		UserID:     row.UserID,
		Token:      row.Token,
		ExpiresAt:  row.ExpiresAt.Time,
		IsRevoked:  row.IsRevoked,
		TokenType:  row.TokenType,
		ClientIP:   row.ClientIp,
		UserAgent:  row.UserAgent,
		CreatedAt:  row.CreatedAt.Time,
		LastUsedAt: row.LastUsedAt,
	}
}
