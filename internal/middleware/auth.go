// Package middleware defines different middlewares around request-response cycle
package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/harshitrajsinha/obj-rest/internal/models"
)

type contextKey string

// UserRole is a constant to be used as context key
const UserRole contextKey = "role"

// AuthMiddleware authenticate the user before accessing protected API routes
func AuthMiddleware(next http.Handler, authSecretKey string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authToken := strings.TrimSpace(r.Header.Get("Authorization"))

		if authToken == "" {
			unauthorized(w, "Missing Authorization header")
			return
		}

		if !strings.HasPrefix(authToken, "Bearer ") {
			unauthorized(w, "Authentication required")
			return
		}

		token := strings.TrimPrefix(authToken, "Bearer ")
		if token == "" {
			unauthorized(w, "Authentication required")
			return
		}

		userrole, err := models.VerifyAuthToken(token, authSecretKey)
		if err != nil {
			log.Println(err)
			unauthorized(w, "Invalid or expired token")
			return
		}

		ctx := context.WithValue(r.Context(), UserRole, userrole)
		r = r.WithContext(ctx)

		log.Println("successfully authenticated")
		next.ServeHTTP(w, r)

	})

}

func unauthorized(w http.ResponseWriter, message string) {
	if err := models.SendResponse(w, http.StatusUnauthorized, message, nil); err != nil {
		log.Println(err)
	}
}
