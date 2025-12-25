package interfaces

import (
	"context"

	"github.com/ariefzainuri96/go-logstream/cmd/api/dto/entity"
	"github.com/ariefzainuri96/go-logstream/cmd/api/dto/request"
)

type IAuth interface {
	Login(context.Context, request.LoginRequest) (entity.User, string, error)
	Register(context.Context, request.RegisterRequest) (uint, error)
	ForgotPassword(context.Context, request.LoginRequest) (string, error)
}