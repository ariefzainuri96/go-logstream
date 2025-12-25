package store

import (
	db "github.com/ariefzainuri96/go-logstream/internal/db"
	"github.com/ariefzainuri96/go-logstream/internal/interfaces"
)

type Storage struct {
	IAuth interfaces.IAuth
}

func NewStorage(gorm *db.GormDB) Storage {
	return Storage{
		IAuth: &AuthStore{gorm},
	}
}
