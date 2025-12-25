package db

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(logger *zap.Logger) {
	dbURL := os.Getenv("DB_ADDR")
	if dbURL == "" {
		logger.Error("⚠️  No DB_ADDR provided, skip migrations")
		return
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		logger.Error("❌ Cannot connect to DB", zap.Error(err))
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		logger.Error("❌ Cannot create migration driver", zap.Error(err))
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)

	if err != nil {
		logger.Error("❌ Cannot run migrations", zap.Error(err))
	}

	// Try to run all migrations
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		logger.Error("❌ Migration error", zap.Error(err))
	}

	logger.Info("✅ Migrations applied successfully")
}
