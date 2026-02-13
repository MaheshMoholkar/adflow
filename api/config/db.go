package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DBConfig contains database connection parameters
type DBConfig struct {
	Host             string
	Port             int
	Username         string
	Password         string
	Database         string
	SSLMode          string
	ChannelBinding   string
	MaxConns         int
	MinConns         int
	MaxConnLifetime  time.Duration
	MaxConnIdleTime  time.Duration
	StatementTimeout time.Duration
}

// NewDBConfigFromEnv creates a database configuration from environment variables
func NewDBConfigFromEnv() *DBConfig {
	maxConns, _ := strconv.Atoi(getEnv("DB_MAX_CONNS", "20"))
	minConns, _ := strconv.Atoi(getEnv("DB_MIN_CONNS", "5"))
	maxConnLifetimeSecs, _ := strconv.Atoi(getEnv("DB_MAX_CONN_LIFETIME_SECS", "1800"))
	maxConnIdleTimeSecs, _ := strconv.Atoi(getEnv("DB_MAX_CONN_IDLE_TIME_SECS", "600"))
	statementTimeoutMs, _ := strconv.Atoi(getEnv("DB_STATEMENT_TIMEOUT_MS", "30000"))
	port, _ := strconv.Atoi(getEnv("DB_PORT", "5432"))

	return &DBConfig{
		Host:             getEnv("DB_HOST", "localhost"),
		Port:             port,
		Username:         getEnv("DB_USER", "postgres"),
		Password:         getEnv("DB_PASSWORD", "postgres"),
		Database:         getEnv("DB_NAME", "callflow_db"),
		SSLMode:          getEnv("DB_SSL_MODE", "disable"),
		ChannelBinding:   getEnv("DB_CHANNEL_BINDING", ""),
		MaxConns:         maxConns,
		MinConns:         minConns,
		MaxConnLifetime:  time.Duration(maxConnLifetimeSecs) * time.Second,
		MaxConnIdleTime:  time.Duration(maxConnIdleTimeSecs) * time.Second,
		StatementTimeout: time.Duration(statementTimeoutMs) * time.Millisecond,
	}
}

// ConnectionString returns a PostgreSQL connection string
func (cfg *DBConfig) ConnectionString() string {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		cfg.SSLMode,
	)
	if cfg.ChannelBinding != "" {
		connStr += "&channel_binding=" + cfg.ChannelBinding
	}
	return connStr
}

// InitDBPool initializes and returns a postgres connection pool
func InitDBPool() (*pgxpool.Pool, error) {
	cfg := NewDBConfigFromEnv()

	poolConfig, err := pgxpool.ParseConfig(cfg.ConnectionString())
	if err != nil {
		return nil, fmt.Errorf("failed to parse database connection string: %w", err)
	}

	poolConfig.MaxConns = int32(cfg.MaxConns)
	poolConfig.MinConns = int32(cfg.MinConns)
	poolConfig.MaxConnLifetime = cfg.MaxConnLifetime
	poolConfig.MaxConnIdleTime = cfg.MaxConnIdleTime

	if cfg.StatementTimeout > 0 {
		poolConfig.ConnConfig.RuntimeParams["statement_timeout"] = strconv.FormatInt(cfg.StatementTimeout.Milliseconds(), 10)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create database connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to database")
	return pool, nil
}


func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
