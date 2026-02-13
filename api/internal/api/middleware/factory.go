package middleware

import (
	"callflow/internal/domain/auth"

	"github.com/gin-gonic/gin"
)

// MiddlewareFactory creates and organizes middleware with proper dependencies
type MiddlewareFactory struct {
	authMiddleware *AuthMiddleware
}

// NewMiddlewareFactory creates a new middleware factory
func NewMiddlewareFactory(authService auth.Service) *MiddlewareFactory {
	return &MiddlewareFactory{
		authMiddleware: NewAuthMiddleware(authService),
	}
}

// AuthChain returns the authentication middleware chain
func (f *MiddlewareFactory) AuthChain() gin.HandlerFunc {
	return f.authMiddleware.RequireAuth()
}
