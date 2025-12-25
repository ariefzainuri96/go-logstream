package middleware

import (
	"bytes"
	"net/http"
)

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
	body       bytes.Buffer
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	// Only set the status code once
    if w.statusCode == http.StatusOK {
        w.statusCode = statusCode
    }
    w.ResponseWriter.WriteHeader(statusCode)
}

func (w *wrappedWriter) Write(p []byte) (int, error) {
    // 1. Write to the buffer (capture for logging)
    w.body.Write(p)
    
    // 2. Write to the underlying ResponseWriter (send to client)
    return w.ResponseWriter.Write(p)
}

type Middleware func(http.Handler) http.Handler

func CreateStack(middlewares ...Middleware) Middleware {
	// Return a function that builds and returns the full http.Handler chain
	return func(next http.Handler) http.Handler {

		// 1. Iterate BACKWARDS over the slice of middlewares
		for i := len(middlewares) - 1; i >= 0; i-- {
			middleware := middlewares[i]

			// 2. Build the chain: The result of the previous iteration
			//    becomes the 'next' handler for the current middleware.
			next = middleware(next)
		}

		// The final 'next' is the complete, correctly ordered chain
		return next
	}
}
