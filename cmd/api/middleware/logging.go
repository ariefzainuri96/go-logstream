package middleware

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ctxKey string

const (
	CtxRequestID ctxKey = "request-id"
)

func Logging(next http.Handler, logger *zap.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Skip Swagger & static assets
		swaggerPrefixes := []string{"/v1/swagger"}
		for _, p := range swaggerPrefixes {
			if strings.HasPrefix(r.URL.Path, p) {
				next.ServeHTTP(w, r)
				return
			}
		}

		// Get request-id from header or generate new
		reqID := r.Header.Get("X-Request-ID")
		if reqID == "" {
			reqID = uuid.New().String()
		}

		// Put into context
		ctx := context.WithValue(r.Context(), CtxRequestID, reqID)
		r = r.WithContext(ctx) // IMPORTANT

		// Limit size
		const maxSize = 2048

		// Ensure application/json header is set before any writes
		w.Header().Set("Content-Type", "application/json")

		var query strings.Builder

		first := true
		for key, values := range r.URL.Query() {
			for _, v := range values {
				if !first {
					query.WriteString("&")
				} else {
					query.WriteString("?")
				}
				first = false
				query.WriteString(fmt.Sprintf("%s=%s", key, v))
			}
		}

		// Use the wrapped writer to capture status and body
		wrapped := &wrappedWriter{w, http.StatusOK, bytes.Buffer{}}

		if r.Method != http.MethodGet {
			// --- REQUEST BODY LOGGING ---

			// Read the original body stream entirely
			requestBody, err := io.ReadAll(r.Body)
			if err != nil {
				logger.Error("Warning: Failed to read request body:", zap.Error(err))
			}

			if len(requestBody) > maxSize {
				requestBody = requestBody[:maxSize]
			}

			// Log the body
			logger.Info("REQUEST", zap.String("Method", r.Method), zap.String("Path", r.URL.Path), zap.String("Body", string(requestBody)))

			// Replace the Request Body
			// CRITICAL: Give the buffered body back to the request for downstream handlers
			r.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		} else {
			logger.Info("REQUEST", zap.String("Method", r.Method), zap.String("Path", fmt.Sprintf("%v%v", r.URL.Path, query.String())))
		}

		// --- Execute Handler Chain ---
		next.ServeHTTP(wrapped, r)

		// Log the captured status, path, latency, and captured response body
		logger.Info("RESPONSE",
			zap.Int("Code", wrapped.statusCode),
			zap.String("Method", r.Method),
			zap.String("Path", fmt.Sprintf("%v%v", r.URL.Path, query.String())),
			zap.String("Latency", time.Since(start).String()),
			zap.String("Body", wrapped.body.String()))
	})
}
