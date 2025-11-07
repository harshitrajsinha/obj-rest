// Package handler defines the handler for registered routes
package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/harshitrajsinha/obj-rest/internal/middleware"
	"github.com/harshitrajsinha/obj-rest/internal/models"
	"github.com/harshitrajsinha/obj-rest/internal/store"
)

// ObjHandler contains the store reference that will be used to fetch data
type ObjHandler struct {
	store store.ObjectDataAccessor
}

// NewObjHandler initializes and returns a new Handler instance with the provided store.UserStore dependency
func NewObjHandler(store store.ObjectDataAccessor) *ObjHandler {
	return &ObjHandler{
		store: store,
	}
}

// CreateNewObj creates a new object and add to the reserved list of objects
func (h *ObjHandler) CreateNewObj(w http.ResponseWriter, r *http.Request) {

	role := r.Context().Value(middleware.UserRole)
	if role != "admin" {
		if err := models.SendResponse(w, http.StatusForbidden, "Recognized but you are not allowed to perform this operation", nil); err != nil {
			log.Println(err)
		}
		return
	}

	ctxWithTimeout, cancel := context.WithTimeout(r.Context(), 8*time.Second)
	defer cancel()

	var payload models.ObjDataPayload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		if err := models.SendResponse(w, http.StatusBadRequest, "Could not create object, invalid payload provided", nil); err != nil {
			log.Println(err)
		}
		return
	}
	defer r.Body.Close()

	if payload.Name == "" {
		if err := models.SendResponse(w, http.StatusBadRequest, "Could not create object, invalid payload provided", nil); err != nil {
			log.Println(err)
		}
		return
	}

	// objPayload := models.ObjDataPayload{
	// 	Name: "Apple MacBook Pro 16",
	// 	Data: map[string]interface{}{
	// 		"year":           2019,
	// 		"price":          1849.99,
	// 		"CPU model":      "Intel Core i9",
	// 		"Hard disk size": "1 TB",
	// 	},
	// }

	responseData, err := h.store.CreateNewObject(ctxWithTimeout, payload)
	if err != nil {
		log.Println(err)
		if err := models.SendResponse(w, http.StatusInternalServerError, "error creating object, try again later", nil); err != nil {
			log.Println(err)
		}
		return
	}

	if err := models.SendResponse(w, http.StatusOK, "Successfully created the object", responseData); err != nil {
		log.Println(err)
		return
	}

}

// GetAllObj gets a list objects that are reserved or newly added
func (h *ObjHandler) GetAllObj(w http.ResponseWriter, r *http.Request) {

	var objsList []models.ObjDataFromResponse

	role := r.Context().Value(middleware.UserRole)
	if role != "admin" && role != "member" {
		if err := models.SendResponse(w, http.StatusForbidden, "Recognized but you are not allowed to perform this operation", nil); err != nil {
			log.Println(err)
		}
		return
	}

	ctxWithTimeout, cancel := context.WithTimeout(r.Context(), 8*time.Second)
	defer cancel()

	objsList, err := h.store.GetAllObjects(ctxWithTimeout)
	if err != nil {
		log.Println(err)
		return
	}

	if err := models.SendResponse(w, http.StatusOK, "Successfully retrieved all objects", objsList); err != nil {
		log.Println(err)
		return
	}

}

// GetObjByID getse an object from the object list based on ID
func (h *ObjHandler) GetObjByID(w http.ResponseWriter, r *http.Request) {

	var objData models.ObjDataFromResponse
	role := r.Context().Value(middleware.UserRole)
	if role != "admin" && role != "member" {
		if err := models.SendResponse(w, http.StatusForbidden, "Recognized but not allowed to perform action", nil); err != nil {
			log.Println(err)
		}
		return
	}

	id := r.PathValue("id")
	if id == "" {
		if err := models.SendResponse(w, http.StatusBadRequest, "object ID is missing", nil); err != nil {
			log.Println(err)
		}
		return
	}

	ctxWithTimeout, cancel := context.WithTimeout(r.Context(), 8*time.Second)
	defer cancel()
	objData, err := h.store.GetObjectByID(ctxWithTimeout, id)
	if err != nil {
		if strings.Contains(err.Error(), "error - no data retrieved in response") {
			if err := models.SendResponse(w, http.StatusBadRequest, "Object with given ID not available", nil); err != nil {
				log.Println(err)
			}
			return
		}
		if err := models.SendResponse(w, http.StatusInternalServerError, "could not retrieve requested object. Try again later", nil); err != nil {
			log.Println(err)
		}
		return
	}

	if err := models.SendResponse(w, http.StatusOK, "Successfully retrieved object", objData); err != nil {
		log.Println(err)
	}

}
