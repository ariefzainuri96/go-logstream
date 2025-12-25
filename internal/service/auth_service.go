package service

import (
	"context"

	"github.com/ariefzainuri96/go-logstream/cmd/api/dto/entity"
	"github.com/ariefzainuri96/go-logstream/cmd/api/dto/request"
	"github.com/ariefzainuri96/go-logstream/internal/store"
	"go.uber.org/zap"
)

// AuthServiceImpl implements the RPC methods defined in auth.proto
type AuthServiceImpl struct {
	logger *zap.Logger
	store  store.Storage
}

func NewAuthService(store store.Storage, logger *zap.Logger) *AuthServiceImpl {
	return &AuthServiceImpl{
		logger: logger,
		store: store,
	}
}

func (s *AuthServiceImpl) Register(ctx context.Context, req request.RegisterRequest) (error) {
	_, err := s.store.IAuth.Register(ctx, req)

	if err != nil {
		return err
	}

	return nil
}

func (s *AuthServiceImpl) Login(ctx context.Context, req request.LoginRequest) (entity.User, string, error) {		
	user, token, err := s.store.IAuth.Login(ctx, req)

	if err != nil {
		return entity.User{}, "", err
	}

	return user, token, nil
}

func (s *AuthServiceImpl) ForgotPassword(ctx context.Context, req request.LoginRequest) (string, error) {	
	msg, err := s.store.IAuth.ForgotPassword(ctx, req)

	if err != nil {
		return "", err
	}

	return msg, nil
}