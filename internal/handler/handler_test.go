// Package handler_test tests all the functionality present in handler package
package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/harshitrajsinha/obj-rest/internal/handler"
	"github.com/harshitrajsinha/obj-rest/internal/middleware"
	"github.com/harshitrajsinha/obj-rest/internal/models"
	"github.com/harshitrajsinha/obj-rest/internal/store"
)

// MockStore embeds and implements ObjectDataAccessor interface
type MockStore struct {
	store.ObjectDataAccessor
}

// GetAllObjs returns a mock list of objects
func (m MockStore) GetAllObjects(ctx context.Context) ([]models.ObjDataFromResponse, error) {
	var objects []models.ObjDataFromResponse
	var testObject models.ObjDataFromResponse = models.ObjDataFromResponse{
		ID:   "123",
		Name: "Test Object",
		Data: map[string]string{
			"Price": "519.99",
		},
	}
	objects = append(objects, testObject)
	return objects, nil
}

// GetObjsByIDs returns a mock list of objects based on requested ID
func (m MockStore) GetObjectsByIDs(ctx context.Context, IDs ...string) ([]models.ObjDataFromResponse, error) {

	var objects []models.ObjDataFromResponse
	testObjectOne := models.ObjDataFromResponse{
		ID:   "1",
		Name: "Test Object One",
		Data: map[string]string{
			"Price": "1",
		},
	}

	testObjectTwo := models.ObjDataFromResponse{
		ID:   "2",
		Name: "Test Object Two",
		Data: map[string]string{
			"Price": "2",
		},
	}

	testObjectThree := models.ObjDataFromResponse{
		ID:   "3",
		Name: "Test Object Three",
		Data: map[string]string{
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

// GetObjByID returns a mock object based on requested ID
func (m MockStore) GetObjectByID(ctx context.Context, ID string) (models.ObjDataFromResponse, error) {
	testObjectOne := models.ObjDataFromResponse{
		ID:   "1",
		Name: "Test Object One",
		Data: map[string]string{
			"Price": "1",
		},
	}
	if ID == "1" {
		return testObjectOne, nil
	}
	return models.ObjDataFromResponse{}, nil
}

// // GetObjByID returns a mock object based on requested ID
// func (m MockStore) CreateNewObject(ctx context.Context, payload models.ObjDataPayload) (models.NewObj, error) {
// 	return models.NewObj{}, nil
// }

// // GetObjByID returns a mock object based on requested ID
// func (m MockStore) UpdateObject(ctx context.Context, objID string, payload models.ObjDataPayload) (models.NewObj, error) {
// 	return models.NewObj{}, nil
// }

// // GetObjByID returns a mock object based on requested ID
// func (m MockStore) UpdateObjectPartially(ctx context.Context, objID string, payload models.ObjDataPayload) (models.NewObj, error) {
// 	return models.NewObj{}, nil
// }

// // GetObjByID returns a mock object based on requested ID
// func (m MockStore) DeleteObject(ctx context.Context, objID string) (map[string]string, error) {
// 	return nil, nil
// }

// TestGetAllObj tests GetAllObj handler
func TestGetAllObj(t *testing.T) {
	t.Run("without authentication", func(t *testing.T) {
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

		message := testResponse["message"].(string)
		if message != "Recognized but you are not allowed to perform this operation" {
			t.Errorf("unexpected unauthorized message, got %s", message)
		}
	})

	t.Run("wrong user role", func(t *testing.T) {
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

		message := testResponse["message"].(string)
		if message != "Recognized but you are not allowed to perform this operation" {
			t.Errorf("unexpected unauthorized message, got %s", message)
		}
	})

	t.Run("admin access", func(t *testing.T) {
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

		message := testResponse["message"].(string)
		if message != "Successfully retrieved all objects" {
			t.Errorf("expected success message, got %s", message)
		}

		data := testResponse["data"].([]interface{})
		object := data[0].(map[string]interface{})
		ID := object["id"].(string)
		Name := object["name"].(string)

		if ID != "123" {
			t.Errorf("want object ID as 123, got %s", ID)
		}

		if Name != "Test Object" {
			t.Errorf("want object Name as Test Object, got %s", Name)
		}
	})

	t.Run("member access", func(t *testing.T) {
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

		message := testResponse["message"].(string)
		if message != "Successfully retrieved all objects" {
			t.Errorf("expected success message, got %s", message)
		}

		data := testResponse["data"].([]interface{})
		object := data[0].(map[string]interface{})
		ID := object["id"].(string)
		Name := object["name"].(string)

		if ID != "123" {
			t.Errorf("want object ID as 123, got %s", ID)
		}

		if Name != "Test Object" {
			t.Errorf("want object Name as Test Object, got %s", Name)
		}
	})
}
