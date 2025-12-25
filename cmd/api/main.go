// @title           LogStream API
// @version         1.0
// @description     This is the documentation for the main e-commerce service.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@example.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @securitydefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/ariefzainuri96/go-logstream/cmd/api/controller"
	"github.com/ariefzainuri96/go-logstream/cmd/api/docs"
	"github.com/ariefzainuri96/go-logstream/internal/db"
	"github.com/ariefzainuri96/go-logstream/internal/logger"
	"github.com/ariefzainuri96/go-logstream/internal/service"
	"github.com/ariefzainuri96/go-logstream/internal/store"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

// loadConfig loads config from environment (simple)
func loadConfig() controller.Config {
	// In real project use env parsing lib (envconfig/viper)
	httpPort := 8080
	ttl := 15 * time.Second

	if v := os.Getenv("HTTP_PORT"); v != "" {
		fmt.Sscanf(v, "%d", &httpPort)
	}
	if v := os.Getenv("SHUTDOWN_TTL"); v != "" {
		var s int
		fmt.Sscanf(v, "%d", &s)
		ttl = time.Duration(s) * time.Second
	}

	return controller.Config{
		HTTPPort:    httpPort,
		ShutdownTTL: ttl,
	}
}

func main() {
	// setup zap logger
	logger := logger.NewLogger()
	defer logger.Sync()
	
	if os.Getenv("APP_ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			logger.Fatal("Error loading .env file", zap.Error(err))
			return
		}
	}

	cfg := loadConfig()

	db.RunMigrations(logger)

	docs.SwaggerInfo.Schemes = []string{"https", "http"}
	docs.SwaggerInfo.Host = fmt.Sprintf("%v", os.Getenv("SWAGGER_HOST"))
	docs.SwaggerInfo.BasePath = fmt.Sprintf("%v", os.Getenv("SWAGGER_PATH"))

	// get gormdb
	gorm, errGorm := db.NewGorm(os.Getenv("DB_ADDR"), logger)

	if errGorm != nil {
		logger.Fatal("Error connecting to gorm database", zap.Error(errGorm))		
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// WaitGroup to wait for servers to stop
	var wg sync.WaitGroup

	store := store.NewStorage(gorm)
	service := service.NewService(store, logger)

	application := &controller.Application{
		Config:    cfg,
		Service:   service,
		Validator: validator.New(),
	}

	// run server
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := application.RunServer(ctx, cfg, logger); err != nil {
			logger.Error("http server stopped with error", zap.Error(err))
			cancel()
		} else {
			logger.Info("http server stopped")
		}
	}()

	// ---------------------------------------------------------
	// Graceful shutdown on OS signals
	// ---------------------------------------------------------
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-ctx.Done():
		logger.Info("context canceled, starting shutdown")
	case sig := <-sigCh:
		logger.Info("signal received, starting shutdown", zap.String("signal", sig.String()))
	}

	// start shutdown procedure
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.ShutdownTTL)
	defer shutdownCancel()

	// Let goroutines handle ctx cancellation — we call cancel to notify them
	cancel()

	// Wait for background goroutines to finish or timeout
	doneCh := make(chan struct{})
	go func() {
		wg.Wait()
		close(doneCh)
	}()

	select {
	case <-doneCh:
		logger.Info("all servers stopped gracefully")
	case <-shutdownCtx.Done():
		logger.Warn("shutdown timed out, forcing exit")
	}

	logger.Info("shutdown complete")
}
