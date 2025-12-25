package controller

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ariefzainuri96/go-logstream/cmd/api/middleware"
	"github.com/ariefzainuri96/go-logstream/internal/service"
	"github.com/go-playground/validator/v10"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

type Config struct {
	HTTPPort    int
	ShutdownTTL time.Duration
}

type Application struct {
	Config    Config
	Service   service.Service
	Validator *validator.Validate
}

func (app *Application) RunServer(ctx context.Context, cfg Config, logger *zap.Logger) error {
	mux := http.NewServeMux()

	stack := middleware.CreateStack(
		middleware.Logging,
		middleware.Recoverer,
	)

	mux.Handle("/v1/auth/", http.StripPrefix("/v1/auth", app.AuthController()))
	
	mux.Handle("/v1/swagger/", httpSwagger.Handler(
		httpSwagger.URL("doc.json"),
	))

	// mux.Handle("/v1/product/", middleware.Authentication(http.StripPrefix("/v1/product", app.ProductRouter())))	

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%v", cfg.HTTPPort),
		Handler:      stack(mux),
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  1 * time.Minute,
	}	

	// Start server in goroutine so we can watch ctx
	serverErrCh := make(chan error, 1)
	go func() {
		logger.Info("starting http server", zap.Int("port", cfg.HTTPPort))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErrCh <- err
			return
		}
		serverErrCh <- nil
	}()

	select {
	case <-ctx.Done():
		// Shutdown with timeout
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		logger.Info("http server shutting down")
		if err := srv.Shutdown(shutdownCtx); err != nil {
			logger.Error("http server shutdown error", zap.Error(err))
			return err
		}
		logger.Info("http server stopped")
		return nil
	case err := <-serverErrCh:
		return err
	}
}
