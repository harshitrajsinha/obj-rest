// Package handler_test tests all the functionality present in handler package
package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/harshitrajsinha/obj-rest/internal/handler"
	"github.com/harshitrajsinha/obj-rest/internal/middleware"
	"github.com/harshitrajsinha/obj-rest/internal/models"
	"github.com/harshitrajsinha/obj-rest/internal/store"

	"github.com/google/uuid"
)

// MockStore embeds and implements ObjectDataAccessor interface
type MockStore struct {
	store.ObjectDataAccessor
}

// GetAllObjects returns a mock list of objects
func (m MockStore) GetAllObjects(_ context.Context) ([]models.ObjDataFromResponse, error) {

	var objects []models.ObjDataFromResponse
	var testObject models.ObjDataFromResponse = models.ObjDataFromResponse{
		ID:   "123",
		Name: "Test Object",
		Data: map[string]interface{}{
			"Price": "519.99",
		},
	}
	objects = append(objects, testObject)
	return objects, nil
}

// GetObjectsByIDs returns a mock list of objects based on requested ID
func (m MockStore) GetObjectsByIDs(_ context.Context, IDs ...string) ([]models.ObjDataFromResponse, error) {

	var objects []models.ObjDataFromResponse
	testObjectOne := models.ObjDataFromResponse{
		ID:   "1",
		Name: "Test Object One",
		Data: map[string]interface{}{
			"Price": "1",
		},
	}

	testObjectTwo := models.ObjDataFromResponse{
		ID:   "2",
		Name: "Test Object Two",
		Data: map[string]interface{}{
			"Price": "2",
		},
	}

	testObjectThree := models.ObjDataFromResponse{
		ID:   "3",
		Name: "Test Object Three",
		Data: map[string]interface{}{
			"Price": "3",
		},
	}

	for _, id := range IDs {
		if id == "1" {
			objects = append(objects, testObjectOne)
		}
		if id == "2" {
			objects = append(objects, testObjectTwo)
		}
		if id == "3" {
			objects = append(objects, testObjectThree)
		}
	}

	return objects, nil

}

// GetObjectByID returns a mock object based on requested ID
func (m MockStore) GetObjectByID(_ context.Context, ID string) (models.ObjDataFromResponse, error) {
	testObjectOne := models.ObjDataFromResponse{
		ID:   "1",
		Name: "Test Object One",
		Data: map[string]interface{}{
			"Price": "1",
		},
	}
	if ID == "1" {
		return testObjectOne, nil
	}
	return models.ObjDataFromResponse{}, errors.New("error - no data retrieved in response")
}

// CreateNewObject returns a mock object based on requested ID
func (m MockStore) CreateNewObject(_ context.Context, payload models.ObjDataPayload) (models.NewObj, error) {

	newID, err := uuid.NewUUID()
	if err != nil {
		return models.NewObj{}, errors.New("unable to create uuid for mock on POST /objects")
	}

	newObject := models.NewObj{
		ID:        newID.String(),
		Name:      payload.Name,
		Data:      payload.Data,
		CreatedAt: time.Now().UTC().GoString(),
	}
	return newObject, nil
}

// // UpdateObject returns a mock object based on requested ID
// func (m MockStore) UpdateObject(ctx context.Context, objID string, payload models.ObjDataPayload) (models.NewObj, error) {
// 	return models.NewObj{}, nil
// }

// // UpdateObjectPartially returns a mock object based on requested ID
// func (m MockStore) UpdateObjectPartially(ctx context.Context, objID string, payload models.ObjDataPayload) (models.NewObj, error) {
// 	return models.NewObj{}, nil
// }

// // DeleteObject returns a mock object based on requested ID
// func (m MockStore) DeleteObject(ctx context.Context, objID string) (map[string]string, error) {
// 	return nil, nil
// }

// TestGetAllObj tests GetAllObj handler
func TestGetAllObj(t *testing.T) {
	t.Run("unauthorized access", func(t *testing.T) {
		var mockStore MockStore

		// create request
		req := httptest.NewRequest(http.MethodGet, "/api/v1/objects", nil)

		// create request recorder
		rec := httptest.NewRecorder()

		objHandlerForTest := handler.NewObjHandler(mockStore)
		objHandlerForTest.GetAllObj(rec, req)

		// check status code
		if rec.Result().StatusCode != http.StatusForbidden {
			t.Errorf("expected status code 403, but got %d", rec.Result().StatusCode)
		}

		// check content-type header
		if rec.Result().Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type header `application/json`, but got %s", rec.Result().Header.Get("Content-Type"))
		}

		var testResponse map[string]interface{}
		if err := json.NewDecoder(rec.Body).Decode(&testResponse); err != nil {
			t.Fatalf("unexpected error occured %v", err)
		}

		message, ok := testResponse["message"].(string)
		if !ok {
			t.Fatalf("message key missing or not a string: %v", testResponse)
		}
		if message != "Recognized but you are not allowed to perform this operation" {
			t.Errorf("unexpected unauthorized message, got %s", message)
		}
	})

	t.Run("incorrect user role", func(t *testing.T) {
		var mockStore MockStore

		// create request
		req := httptest.NewRequest(http.MethodGet, "/api/v1/objects", nil)

		ctxWithValue := context.WithValue(req.Context(), middleware.UserRole, "user")
		req = req.WithContext(ctxWithValue)

		// create request recorder
		rec := httptest.NewRecorder()

		objHandlerForTest := handler.NewObjHandler(mockStore)
		objHandlerForTest.GetAllObj(rec, req)

		// check status code
		if rec.Result().StatusCode != http.StatusForbidden {
			t.Errorf("expected status code 403, but got %d", rec.Result().StatusCode)
		}

		// check content-type header
		if rec.Result().Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type header `application/json`, but got %s", rec.Result().Header.Get("Content-Type"))
		}

		var testResponse map[string]interface{}
		if err := json.NewDecoder(rec.Body).Decode(&testResponse); err != nil {
			t.Fatalf("unexpected error occured %v", err)
		}

		message, ok := testResponse["message"].(string)
		if !ok {
			t.Fatalf("message key missing or not a string: %v", testResponse)
		}
		if message != "Recognized but you are not allowed to perform this operation" {
			t.Errorf("unexpected unauthorized message, got %s", message)
		}
	})

	t.Run("admin role access", func(t *testing.T) {
		var mockStore MockStore

		// create request
		req := httptest.NewRequest(http.MethodGet, "/api/v1/objects", nil)

		ctxWithValue := context.WithValue(req.Context(), middleware.UserRole, "admin")
		req = req.WithContext(ctxWithValue)

		// create request recorder
		rec := httptest.NewRecorder()

		objHandlerForTest := handler.NewObjHandler(mockStore)
		objHandlerForTest.GetAllObj(rec, req)

		// check status code
		if rec.Result().StatusCode != http.StatusOK {
			t.Errorf("expected status code 200, but got %d", rec.Result().StatusCode)
		}

		// check content-type header
		if rec.Result().Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type header `application/json`, but got %s", rec.Result().Header.Get("Content-Type"))
		}

		var testResponse map[string]interface{}
		if err := json.NewDecoder(rec.Body).Decode(&testResponse); err != nil {
			t.Fatalf("unexpected error occured %v", err)
		}

		message, ok := testResponse["message"].(string)
		if !ok {
			t.Fatalf("message key missing or not a string: %v", testResponse)
		}
		if message != "Successfully retrieved all objects" {
			t.Errorf("expected success message, got %s", message)
		}

		data := testResponse["data"].([]interface{})
		if !ok {
			t.Fatalf("data key missing or not a []interface{}: %v", testResponse)
		}
		object := data[0].(map[string]interface{})
		if !ok {
			t.Fatalf("data[0] is not a map[string]interface{}: %v", testResponse)
		}
		ID := object["id"].(string)
		if !ok {
			t.Fatalf("id key missing or not a string: %v", testResponse)
		}
		Name := object["name"].(string)
		if !ok {
			t.Fatalf("name key missing or not a string: %v", testResponse)
		}

		if ID != "123" {
			t.Errorf("want object ID as 123, got %s", ID)
		}

		if Name != "Test Object" {
			t.Errorf("want object Name as Test Object, got %s", Name)
		}
	})

	t.Run("member role access", func(t *testing.T) {
		var mockStore MockStore

		// create request
		req := httptest.NewRequest(http.MethodGet, "/api/v1/objects", nil)

		ctxWithValue := context.WithValue(req.Context(), middleware.UserRole, "member")
		req = req.WithContext(ctxWithValue)

		// create request recorder
		rec := httptest.NewRecorder()

		objHandlerForTest := handler.NewObjHandler(mockStore)
		objHandlerForTest.GetAllObj(rec, req)

		// check status code
		if rec.Result().StatusCode != http.StatusOK {
			t.Errorf("expected status code 200, but got %d", rec.Result().StatusCode)
		}

		// check content-type header
		if rec.Result().Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type header `application/json`, but got %s", rec.Result().Header.Get("Content-Type"))
		}

		var testResponse map[string]interface{}
		if err := json.NewDecoder(rec.Body).Decode(&testResponse); err != nil {
			t.Fatalf("unexpected error occured %v", err)
		}

		message, ok := testResponse["message"].(string)
		if !ok {
			t.Fatalf("message key missing or not a string: %v", testResponse)
		}
		if message != "Successfully retrieved all objects" {
			t.Errorf("expected success message, got %s", message)
		}

		data := testResponse["data"].([]interface{})
		if !ok {
			t.Fatalf("data key missing or not a []interface{}: %v", testResponse)
		}
		object := data[0].(map[string]interface{})
		if !ok {
			t.Fatalf("data[0] is not a map[string]interface{}: %v", testResponse)
		}
		ID := object["id"].(string)
		if !ok {
			t.Fatalf("id key missing or not a string: %v", testResponse)
		}
		Name := object["name"].(string)
		if !ok {
			t.Fatalf("name key missing or not a string: %v", testResponse)
		}

		if ID != "123" {
			t.Errorf("want object ID as 123, got %s", ID)
		}

		if Name != "Test Object" {
			t.Errorf("want object Name as Test Object, got %s", Name)
		}
	})
}

// TestCreateObj tests CreateObj handler
func TestCreateObj(t *testing.T) {
	t.Run("unauthorized access", func(t *testing.T) {
		var mockStore MockStore

		objPayload := models.ObjDataPayload{
			Name: "Apple MacBook Pro 16",
			Data: map[string]interface{}{
				"year":           2019,
				"price":          1849.99,
				"CPU model":      "Intel Core i9",
				"Hard disk size": "1 TB",
			},
		}

		payloadToSendInReq, _ := json.Marshal(objPayload)

		// create request
		req := httptest.NewRequest(http.MethodPost, "/api/v1/objects", strings.NewReader(string(payloadToSendInReq)))

		// create request recorder
		rec := httptest.NewRecorder()

		objHandlerForTest := handler.NewObjHandler(mockStore)
		objHandlerForTest.CreateNewObj(rec, req)

		// check status code
		if rec.Result().StatusCode != http.StatusForbidden {
			t.Errorf("expected status code 403, but got %d", rec.Result().StatusCode)
		}

		// check content-type header
		if rec.Result().Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type header `application/json`, but got %s", rec.Result().Header.Get("Content-Type"))
		}

		var testResponse map[string]interface{}
		if err := json.NewDecoder(rec.Body).Decode(&testResponse); err != nil {
			t.Fatalf("unexpected error occured %v", err)
		}

		message, ok := testResponse["message"].(string)
		if !ok {
			t.Fatalf("message key missing or not a string: %v", testResponse)
		}
		if message != "Recognized but you are not allowed to perform this operation" {
			t.Errorf("unexpected unauthorized message, got %s", message)
		}
	})

	t.Run("incorrect user role", func(t *testing.T) {
		var mockStore MockStore

		objPayload := models.ObjDataPayload{
			Name: "Apple MacBook Pro 16",
			Data: map[string]interface{}{
				"year":           2019,
				"price":          1849.99,
				"CPU model":      "Intel Core i9",
				"Hard disk size": "1 TB",
			},
		}

		payloadToSendInReq, _ := json.Marshal(objPayload)

		// create request
		req := httptest.NewRequest(http.MethodPost, "/api/v1/objects", strings.NewReader(string(payloadToSendInReq)))

		ctxWithValue := context.WithValue(req.Context(), middleware.UserRole, "user")
		req = req.WithContext(ctxWithValue)

		// create request recorder
		rec := httptest.NewRecorder()

		objHandlerForTest := handler.NewObjHandler(mockStore)
		objHandlerForTest.CreateNewObj(rec, req)

		// check status code
		if rec.Result().StatusCode != http.StatusForbidden {
			t.Errorf("expected status code 403, but got %d", rec.Result().StatusCode)
		}

		// check content-type header
		if rec.Result().Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type header `application/json`, but got %s", rec.Result().Header.Get("Content-Type"))
		}

		var testResponse map[string]interface{}
		if err := json.NewDecoder(rec.Body).Decode(&testResponse); err != nil {
			t.Fatalf("unexpected error occured %v", err)
		}

		message, ok := testResponse["message"].(string)
		if !ok {
			t.Fatalf("message key missing or not a string: %v", testResponse)
		}
		if message != "Recognized but you are not allowed to perform this operation" {
			t.Errorf("unexpected unauthorized message, got %s", message)
		}
	})

	t.Run("member access", func(t *testing.T) {
		var mockStore MockStore

		objPayload := models.ObjDataPayload{
			Name: "Apple MacBook Pro 16",
			Data: map[string]interface{}{
				"year":           2019,
				"price":          1849.99,
				"CPU model":      "Intel Core i9",
				"Hard disk size": "1 TB",
			},
		}

		payloadToSendInReq, _ := json.Marshal(objPayload)

		// create request
		req := httptest.NewRequest(http.MethodPost, "/api/v1/objects", strings.NewReader(string(payloadToSendInReq)))

		ctxWithValue := context.WithValue(req.Context(), middleware.UserRole, "member")
		req = req.WithContext(ctxWithValue)

		// create request recorder
		rec := httptest.NewRecorder()

		objHandlerForTest := handler.NewObjHandler(mockStore)
		objHandlerForTest.CreateNewObj(rec, req)

		// check status code
		if rec.Result().StatusCode != http.StatusForbidden {
			t.Errorf("expected status code 403, but got %d", rec.Result().StatusCode)
		}

		// check content-type header
		if rec.Result().Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type header `application/json`, but got %s", rec.Result().Header.Get("Content-Type"))
		}

		var testResponse map[string]interface{}
		if err := json.NewDecoder(rec.Body).Decode(&testResponse); err != nil {
			t.Fatalf("unexpected error occured %v", err)
		}

		message, ok := testResponse["message"].(string)
		if !ok {
			t.Fatalf("message key missing or not a string: %v", testResponse)
		}
		if message != "Recognized but you are not allowed to perform this operation" {
			t.Errorf("unexpected unauthorized message, got %s", message)
		}
	})

	t.Run("admin access", func(t *testing.T) {
		var mockStore MockStore

		objPayload := models.ObjDataPayload{
			Name: "Apple MacBook Pro 16",
			Data: map[string]interface{}{
				"year":           2019,
				"price":          1849.99,
				"CPU model":      "Intel Core i9",
				"Hard disk size": "1 TB",
			},
		}

		payloadToSendInReq, _ := json.Marshal(objPayload)

		// create request
		req := httptest.NewRequest(http.MethodPost, "/api/v1/objects", strings.NewReader(string(payloadToSendInReq)))

		ctxWithValue := context.WithValue(req.Context(), middleware.UserRole, "admin")
		req = req.WithContext(ctxWithValue)

		// create request recorder
		rec := httptest.NewRecorder()

		objHandlerForTest := handler.NewObjHandler(mockStore)
		objHandlerForTest.CreateNewObj(rec, req)

		// check status code
		if rec.Result().StatusCode != http.StatusOK {
			t.Errorf("expected status code 200, but got - %d", rec.Result().StatusCode)
		}

		// check content-type header
		if rec.Result().Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type header `application/json`, but got %s", rec.Result().Header.Get("Content-Type"))
		}

		var testResponse map[string]interface{}
		if err := json.NewDecoder(rec.Body).Decode(&testResponse); err != nil {
			t.Fatalf("unexpected error occured %v", err)
		}

		message, ok := testResponse["message"].(string)
		if !ok {
			t.Fatalf("message key missing or not a string: %v", testResponse)
		}
		if message != "Successfully created the object" {
			t.Errorf("unexpected unauthorized message, got - %s", message)
		}

		data := testResponse["data"].(map[string]interface{})
		name, ok := data["name"].(string)
		if !ok {
			t.Fatalf("name key missing or not a string: %v", testResponse)
		}
		if name != "Apple MacBook Pro 16" {
			t.Errorf("expected created object name `Apple MacBook Pro 16`, got - %s", name)
		}
	})

	t.Run("invalid payload", func(t *testing.T) {
		var mockStore MockStore

		objPayload := models.ObjDataPayload{}

		payloadToSendInReq, _ := json.Marshal(objPayload)

		// create request
		req := httptest.NewRequest(http.MethodPost, "/api/v1/objects", strings.NewReader(string(payloadToSendInReq)))

		ctxWithValue := context.WithValue(req.Context(), middleware.UserRole, "admin")
		req = req.WithContext(ctxWithValue)

		// create request recorder
		rec := httptest.NewRecorder()

		objHandlerForTest := handler.NewObjHandler(mockStore)
		objHandlerForTest.CreateNewObj(rec, req)

		// check status code
		if rec.Result().StatusCode != http.StatusBadRequest {
			t.Errorf("expected status code 400, but got - %d", rec.Result().StatusCode)
		}

		// check content-type header
		if rec.Result().Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type header `application/json`, but got %s", rec.Result().Header.Get("Content-Type"))
		}

		var testResponse map[string]interface{}
		if err := json.NewDecoder(rec.Body).Decode(&testResponse); err != nil {
			t.Fatalf("unexpected error occured %v", err)
		}

		message, ok := testResponse["message"].(string)
		if !ok {
			t.Fatalf("message key missing or not a string: %v", testResponse)
		}
		if message != "Could not create object, invalid payload provided" {
			t.Errorf("unexpected message, got - %s", message)
		}
	})
}

// TestGetObjByID tests GetObjByID handler
func TestGetObjByID(t *testing.T) {

	t.Run("unauthorized access", func(t *testing.T) {

		var mockStore MockStore
		req := httptest.NewRequest(http.MethodGet, "/api/v1/objects/2", nil)
		// no authentication

		rec := httptest.NewRecorder()
		objHandler := handler.NewObjHandler(mockStore)
		objHandler.GetObjByID(rec, req)

		if rec.Result().StatusCode != http.StatusForbidden {
			t.Errorf("expected status code as 403, got %d", rec.Result().StatusCode)
		}

		if rec.Result().Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Header Content-Type as application/json, got %s", rec.Result().Header.Get("Content-Type"))
		}

		var testResponse map[string]interface{}
		if err := json.NewDecoder(rec.Body).Decode(&testResponse); err != nil {
			t.Fatalf("unexpected error occured %v", err)
		}

		message, ok := testResponse["message"].(string)
		if !ok {
			t.Fatalf("message key missing or not a string: %v", testResponse)
		}
		if message != "Recognized but not allowed to perform action" {
			t.Errorf("unexpected message, got %s", message)
		}

	})

	t.Run("invalid object id", func(t *testing.T) {

		var mockStore MockStore
		req := httptest.NewRequest(http.MethodGet, "/api/v1/objects/2", nil)
		ctxWithValue := context.WithValue(req.Context(), middleware.UserRole, "member")
		req = req.WithContext(ctxWithValue)

		rec := httptest.NewRecorder()
		objHandler := handler.NewObjHandler(mockStore)

		mux := http.NewServeMux()
		mux.HandleFunc("GET /api/v1/objects/{id}", objHandler.GetObjByID)
		mux.ServeHTTP(rec, req)

		if rec.Result().StatusCode != http.StatusBadRequest {
			t.Errorf("expected status code as 400, got %d", rec.Result().StatusCode)
		}

		if rec.Result().Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Header Content-Type as application/json, got %s", rec.Result().Header.Get("Content-Type"))
		}

		var testResponse map[string]interface{}
		if err := json.NewDecoder(rec.Body).Decode(&testResponse); err != nil {
			t.Fatalf("unexpected error occured %v", err)
		}

		message, ok := testResponse["message"].(string)
		if !ok {
			t.Fatalf("message key missing or not a string: %v", testResponse)
		}
		if message != "Object with given ID not available" {
			t.Errorf("unexpected message, got %s", message)
		}

	})

	t.Run("valid request", func(t *testing.T) {

		var mockStore MockStore
		req := httptest.NewRequest(http.MethodGet, "/api/v1/objects/1", nil)
		ctxWithValue := context.WithValue(req.Context(), middleware.UserRole, "member")
		req = req.WithContext(ctxWithValue)

		rec := httptest.NewRecorder()
		objHandler := handler.NewObjHandler(mockStore)

		mux := http.NewServeMux()
		mux.HandleFunc("GET /api/v1/objects/{id}", objHandler.GetObjByID)
		mux.ServeHTTP(rec, req)

		if rec.Result().StatusCode != http.StatusOK {
			t.Errorf("expected status code as 200, got %d", rec.Result().StatusCode)
		}

		if rec.Result().Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Header Content-Type as application/json, got %s", rec.Result().Header.Get("Content-Type"))
		}

		var testResponse map[string]interface{}
		if err := json.NewDecoder(rec.Body).Decode(&testResponse); err != nil {
			t.Fatalf("unexpected error occured %v", err)
		}

		message, ok := testResponse["message"].(string)
		if !ok {
			t.Fatalf("message key missing or not a string: %v", testResponse)
		}
		if message != "Successfully retrieved object" {
			t.Errorf("unexpected message, got %s", message)
		}

		data := testResponse["data"].(map[string]interface{})
		if !ok {
			t.Fatalf("data key missing or not a map[string]interface{}: %v", testResponse)
		}
		ID := data["id"].(string)
		if !ok {
			t.Fatalf("id key missing or not a string: %v", testResponse)
		}
		if ID != "1" {
			t.Errorf("expected ID as 1, got %s", ID)
		}
		name := data["name"].(string)
		if !ok {
			t.Fatalf("name key missing or not a string: %v", testResponse)
		}
		if name != "Test Object One" {
			t.Errorf("expected name as Test Object One, got %s", name)
		}
	})

}
