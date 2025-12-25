package middleware

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/ariefzainuri96/go-logstream/cmd/api/utils"
	"github.com/golang-jwt/jwt/v5"
)

// Key to store user information in context
type contextKey string

const UserContextKey contextKey = "user"

// Function to get user data from request
func GetUserFromContext(r *http.Request) (map[string]any, bool) {
	userData, ok := r.Context().Value(UserContextKey).(map[string]any)
	return userData, ok
}

/*
This Authentication middleware usage is for route

mux.Handle("/v1/product/", middleware.Authentication(http.StripPrefix("/v1/product", app.ProductRouter())))
*/
func Authentication(next http.Handler) http.Handler {
	jwtSecret := os.Getenv("SECRET_KEY")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			utils.RespondError(w, http.StatusUnauthorized, "Missing Authorization Header")			
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			utils.RespondError(w, http.StatusUnauthorized, "Invalid Token")
			return
		}

		// Extract claims from the token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			utils.RespondError(w, http.StatusUnauthorized, "Invalid Token Claims")			
			return
		}

		// Extract data from token
		userId, _ := claims["user_id"].(float64) // go standart store json numbers as float64
		email, _ := claims["email"].(string)
		isAdmin, _ := claims["is_admin"].(bool)

		// Store the user data in the request context
		ctx := context.WithValue(r.Context(), UserContextKey, map[string]any{
			"user_id":  int(userId),
			"email":    email,
			"is_admin": isAdmin,
		})

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

/*
This AdminHandler usage is for each endpoint

authRouter.HandleFunc("POST /login", middleware.AdminHandler(app.login))
*/
func AdminHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r)

		if !ok {
			utils.RespondError(w, http.StatusUnauthorized, "Unauthorized, please re login!")
			return
		}

		isAdmin := user["is_admin"].(bool)

		if !isAdmin {
			utils.RespondError(w, http.StatusForbidden, "You are not authorized to perform this action!")
			return
		}

		next(w, r)
	}
}

/*
This UserHandler usage is for each endpoint

authRouter.HandleFunc("POST /login", middleware.UserHandler(app.login))
*/
func UserHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r)

		if !ok {
			utils.RespondError(w, http.StatusUnauthorized, "Unauthorized, please re login!")
			return
		}

		isAdmin := user["is_admin"].(bool)

		if isAdmin {
			utils.RespondError(w, http.StatusForbidden, "You are not authorized to perform this action!")
			return
		}

		next(w, r)
	}
}
