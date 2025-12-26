package store

import (
	db "github.com/ariefzainuri96/go-logstream/internal/db"
	"github.com/ariefzainuri96/go-logstream/internal/interfaces"
	"go.uber.org/zap"
)

type Storage struct {
	IAuth    interfaces.IAuth
	IProject interfaces.IProject
	IPost    interfaces.IPost
}

func NewStorage(gorm *db.GormDB, logger *zap.Logger) Storage {
	return Storage{
		IAuth:    &AuthStore{gorm},
		IProject: &ProjectStore{gorm, logger},
		IPost:    &PostStore{gorm, logger},
	}
}
