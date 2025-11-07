// Package handler defines the handler for registered routes
package handler

import (
	"log"
	"net/http"
	"os"

	"github.com/harshitrajsinha/obj-rest/internal/models"
)

// Login verifies user role and create auth token for the user
func Login(w http.ResponseWriter, r *http.Request) {

	requestQuery := r.URL.Query()
	role := requestQuery.Get("role")
	if role != "admin" && role != "member" {
		if err := models.SendResponse(w, http.StatusBadRequest, "invalid role option", nil); err != nil {
			log.Println(err)
		}
		return
	}

	secretKey := os.Getenv("AUTH_SECRET_KEY")
	if secretKey == "" {
		log.Fatalf("Auth key not found")
	}

	token, err := models.GenerateAuthToken(role, secretKey)
	if err != nil {
		if err := models.SendResponse(w, http.StatusInternalServerError, "could not authenticate. Try again later", nil); err != nil {
			log.Println(err)
		}
		return
	}

	tokenData := map[string]string{
		"token": token,
	}

	if err := models.SendResponse(w, http.StatusCreated, "Successfully authenticated", tokenData); err != nil {
		log.Println(err)
	}
}
