package store

import (
	"context"

	"github.com/ariefzainuri96/go-logstream/cmd/api/dto/entity"
	db "github.com/ariefzainuri96/go-logstream/internal/db"
	"github.com/ariefzainuri96/go-logstream/cmd/api/dto/request"
)

type Storage struct {
	IAuth interface {
		Login(context.Context, request.LoginRequest) (entity.User, string, error)
		Register(context.Context, request.RegisterRequest) (uint, error)
		ForgotPassword(context.Context, request.LoginRequest) (string, error)
	}
}

func NewStorage(gorm *db.GormDB) Storage {
	return Storage{
		IAuth: &AuthStore{gorm},
	}
}
