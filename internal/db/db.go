package db

import (
	"context"
	"time"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GormDB struct {
	gormDb *gorm.DB
}

func (d *GormDB) ExecWithTimeoutErr(ctx context.Context, fn func(tx *gorm.DB) error) error {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	return fn(d.gormDb.WithContext(ctx))
}

func (d *GormDB) ExecWithTimeoutVal(ctx context.Context, fn func(tx *gorm.DB) *gorm.DB) *gorm.DB {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	return fn(d.gormDb.WithContext(ctx))
}

func NewGorm(addr string, logger *zap.Logger) (*GormDB, error) {
	// Use your existing DSN (Data Source Name) / connection string
	// Example DSN: "host=localhost user=user password=pass dbname=ecommerce-db port=5432 sslmode=disable"

	db, err := gorm.Open(postgres.Open(addr), &gorm.Config{})

	if err != nil {		
		return &GormDB{}, err
	}

	logger.Info("Database connection successfully established with GORM.")

	return &GormDB{
		gormDb: db,
	}, nil
}
