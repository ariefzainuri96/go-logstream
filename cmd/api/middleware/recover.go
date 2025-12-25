package middleware

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/ariefzainuri96/go-logstream/cmd/api/utils"
)

// Recoverer captures panics, logs the stack trace, and returns a 500 error.
func Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// 1. Set up the defer function to run after the handler chain finishes (or panics)
		defer func() {
			if rvr := recover(); rvr != nil {

				// 2. A panic occurred! Log the full stack trace.
				log.Printf("PANIC RECOVERED: %v", rvr)

				// Optional: You can get the stack trace here for detailed logging
				stack := debug.Stack()

				utils.RespondError(w, http.StatusInternalServerError, string(stack))
			}
		}()

		// 4. Continue to the next handler
		next.ServeHTTP(w, r)
	})
}
