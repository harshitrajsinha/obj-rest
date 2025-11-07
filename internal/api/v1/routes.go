// Package v1 defines the v1 api routes
package v1

import (
	"net/http"

	"github.com/harshitrajsinha/obj-rest/internal/handler"
	"github.com/harshitrajsinha/obj-rest/internal/middleware"
	"github.com/harshitrajsinha/obj-rest/internal/store"
)

// RegisterV1Routes registers all the routes for api version v1
func RegisterV1Routes(mux *http.ServeMux, storeClient store.ObjectDataAccessor, authSecretKey string) {

	mux.HandleFunc("GET /login", handler.Login)

	objHandler := handler.NewObjHandler(storeClient)
	mux.HandleFunc("POST /api/v1/objects", middleware.AuthMiddleware((objHandler.CreateNewObj), authSecretKey))
	mux.HandleFunc("GET /api/v1/objects", middleware.AuthMiddleware((objHandler.GetAllObj), authSecretKey))

}
