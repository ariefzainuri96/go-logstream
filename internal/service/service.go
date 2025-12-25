package service

import (
	"github.com/ariefzainuri96/go-logstream/internal/interfaces"
	"github.com/ariefzainuri96/go-logstream/internal/store"
	"go.uber.org/zap"
)

type Service struct {
	IAuth interfaces.IAuth
}

func NewService(store store.Storage, logger *zap.Logger) Service {
	return Service{
		IAuth: NewAuthService(store, logger),
	}
}