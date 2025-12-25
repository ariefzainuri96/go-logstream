package middleware

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type ctxKey string

const (
	CtxRequestID ctxKey = "request-id"
)

func Logging(next http.Handler) http.Handler {
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

		log.Printf("REQUEST_ID: %s", reqID)

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
				log.Println("Warning: Failed to read request body:", err)
				// Don't interrupt flow, but proceed with empty body
			}

			if len(requestBody) > maxSize {
				requestBody = requestBody[:maxSize]
			}

			// Log the body
			log.Printf("REQUEST [%s] %s\n\nBody: %s", r.Method, r.URL.Path, string(requestBody))

			// Replace the Request Body
			// CRITICAL: Give the buffered body back to the request for downstream handlers
			r.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		} else {
			log.Printf("REQUEST [%s] %s", r.Method, fmt.Sprintf("%v%v", r.URL.Path, query.String()),)
		}

		// --- Execute Handler Chain ---
		next.ServeHTTP(wrapped, r)

		// --- RESPONSE LOGGING ---

		// Log the captured status, path, latency, and captured response body
		log.Printf("RESPONSE %d [%s] %s | Latency: %s\n\nBody: %s",
			wrapped.statusCode,
			r.Method,
			fmt.Sprintf("%v%v", r.URL.Path, query.String()),
			time.Since(start),
			wrapped.body.String(), // Access the buffered response body
		)
	})
}
