package service

import (
	"github.com/ariefzainuri96/go-logstream/internal/store"
	"go.uber.org/zap"
)

type Service struct {
	AuthService *AuthServiceImpl
}

func NewService(store store.Storage, logger *zap.Logger) Service {
	return Service{
		AuthService: NewAuthService(store, logger),
	}
}