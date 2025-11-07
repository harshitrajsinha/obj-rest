// Package store serves as data layer for the application
// The following functions calls external API to retrieve and then serve data to the application
package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/harshitrajsinha/obj-rest/internal/models"
)

// ObjectStore implements ObjectDataAccessor to define methods to fetch data from external API
type ObjectStore struct {
	APIURL string
}

// NewStore acts as a constructor method for dependency injection
func NewStore(apiURL string) ObjectDataAccessor {
	return &ObjectStore{
		APIURL: apiURL,
	}
}

// GetAllObjects fetches all objects from the external API
func (s ObjectStore) GetAllObjects(ctx context.Context) ([]models.ObjDataFromResponse, error) {

	apiURL := s.APIURL + "/objects"

	// create a new request with context
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request to fetch all objects, %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// create new client
	newClient := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:    100,
			IdleConnTimeout: 90 * time.Second,
		},
	}

	// send request and get response
	resp, err := newClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending response to fetch all objects, %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	// parse response
	var objectsList []models.ObjDataFromResponse
	err = json.NewDecoder(resp.Body).Decode(&objectsList)
	if err != nil {
		return nil, fmt.Errorf("error parsing response of all objects, %w", err)
	}

	if len(objectsList) == 0 || (len(objectsList) != 0 && objectsList[0].ID == "") {
		return nil, errors.New("error - no data retrieved in response")
	}

	return objectsList, nil
}

// GetObjectsByIDs fetches list of all objects based on different IDs input as query params
func (s ObjectStore) GetObjectsByIDs(ctx context.Context, IDs ...string) ([]models.ObjDataFromResponse, error) {

	var str strings.Builder
	str.WriteString("?")
	for _, id := range IDs {
		str.WriteString(fmt.Sprintf("id=%s&", id))
	}
	apiURL := s.APIURL + "/objects" + str.String()

	// create a new request with context
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request to fetch objects based on IDs, %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// create new client
	newClient := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:    100,
			IdleConnTimeout: 90 * time.Second,
		},
	}

	// send request and get response
	resp, err := newClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending response to fetch objects based on IDs, %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	// parse response
	var objectsList []models.ObjDataFromResponse
	err = json.NewDecoder(resp.Body).Decode(&objectsList)
	if err != nil {
		return nil, fmt.Errorf("error parsing response of objects based on IDs, %w", err)
	}

	if len(objectsList) == 0 || (len(objectsList) != 0 && objectsList[0].ID == "") {
		return nil, errors.New("error - no data retrieved in response")
	}

	return objectsList, nil
}

// GetObjectByID fetches a single object based on param value
func (s ObjectStore) GetObjectByID(ctx context.Context, ID string) (models.ObjDataFromResponse, error) {

	var objectData models.ObjDataFromResponse
	apiURL := s.APIURL + "/objects/" + ID

	// create a new request with context
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return objectData, fmt.Errorf("error creating request to fetch object based on ID, %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// create new client
	newClient := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:    100,
			IdleConnTimeout: 90 * time.Second,
		},
	}

	// send request and get response
	resp, err := newClient.Do(req)
	if err != nil {
		return objectData, fmt.Errorf("error sending response to fetch object based on ID, %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return objectData, fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	// parse response

	err = json.NewDecoder(resp.Body).Decode(&objectData)
	if err != nil {
		return objectData, fmt.Errorf("error parsing response of object based on ID, %w", err)
	}

	if objectData.ID == "" {
		return objectData, errors.New("error - no data retrieved in response")
	}

	return objectData, nil
}

// CreateNewObject sends a request to generate new object and add it to the list
func (s ObjectStore) CreateNewObject(ctx context.Context, objPayload models.ObjDataPayload) (models.NewObj, error) {

	var objectData models.NewObj
	apiURL := s.APIURL + "/objects"

	// encode into JSON string
	payloadToSend, err := json.Marshal(objPayload)
	if err != nil {
		return objectData, fmt.Errorf("error creating payload to send to create new object, %w", err)
	}

	// create a new request with context
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, strings.NewReader(string(payloadToSend)))
	if err != nil {
		return objectData, fmt.Errorf("error creating request to create new object, %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// create new client
	newClient := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:    100,
			IdleConnTimeout: 90 * time.Second,
		},
	}

	// send request and get response
	resp, err := newClient.Do(req)
	if err != nil {
		return objectData, fmt.Errorf("error sending response to create new object, %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return objectData, fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	// parse response

	err = json.NewDecoder(resp.Body).Decode(&objectData)
	if err != nil {
		return objectData, fmt.Errorf("error parsing response of newly created object, %w", err)
	}

	if objectData.ID == "" {
		return objectData, errors.New("error - no data retrieved in response")
	}

	return objectData, nil
}

// UpdateObject updates all the fields of an object that has been created
func (s ObjectStore) UpdateObject(ctx context.Context, objID string, objPayload models.ObjDataPayload) (models.NewObj, error) {

	var objectData models.NewObj
	apiURL := s.APIURL + "/objects/" + objID

	// encode into JSON string
	payloadToSend, err := json.Marshal(objPayload)
	if err != nil {
		return objectData, fmt.Errorf("error creating payload to send to update object, %w", err)
	}

	// create a new request with context
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, apiURL, strings.NewReader(string(payloadToSend)))
	if err != nil {
		return objectData, fmt.Errorf("error creating request to update object, %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// create new client
	newClient := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:    100,
			IdleConnTimeout: 90 * time.Second,
		},
	}

	// send request and get response
	resp, err := newClient.Do(req)
	if err != nil {
		return objectData, fmt.Errorf("error sending response to update object, %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return objectData, fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	// parse response

	err = json.NewDecoder(resp.Body).Decode(&objectData)
	if err != nil {
		return objectData, fmt.Errorf("error parsing response of updated object, %w", err)
	}

	if objectData.ID == "" {
		return objectData, errors.New("error - no data retrieved in response")
	}

	return objectData, nil
}

// UpdateObjectPartially updates one or more fields of an object that has been created
func (s ObjectStore) UpdateObjectPartially(ctx context.Context, objID string, objPayload models.ObjDataPayload) (models.NewObj, error) {

	var objectData models.NewObj
	apiURL := s.APIURL + "/objects/" + objID

	// encode into JSON string
	payloadToSend, err := json.Marshal(objPayload)
	if err != nil {
		return objectData, fmt.Errorf("error creating payload to send to partially update object, %w", err)
	}

	// create a new request with context
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, apiURL, strings.NewReader(string(payloadToSend)))
	if err != nil {
		return objectData, fmt.Errorf("error creating request to partially update object, %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// create new client
	newClient := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:    100,
			IdleConnTimeout: 90 * time.Second,
		},
	}

	// send request and get response
	resp, err := newClient.Do(req)
	if err != nil {
		return objectData, fmt.Errorf("error sending response to partially update object, %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return objectData, fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	// parse response

	err = json.NewDecoder(resp.Body).Decode(&objectData)
	if err != nil {
		return objectData, fmt.Errorf("error parsing response of updated object, %w", err)
	}

	if objectData.ID == "" {
		return objectData, errors.New("error - no data retrieved in response")
	}

	return objectData, nil
}

// DeleteObject deletes a created object
func (s ObjectStore) DeleteObject(ctx context.Context, objID string) (map[string]string, error) {

	var response map[string]string

	apiURL := s.APIURL + "/objects/" + objID

	// create a new request with context
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request to delete object, %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// create new client
	newClient := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:    100,
			IdleConnTimeout: 90 * time.Second,
		},
	}

	// send request and get response
	resp, err := newClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending response to delete object, %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	// parse response
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("error parsing response of updated object, %w", err)
	}

	value, ok := response["message"]
	if !ok || !strings.Contains(value, "has been deleted") {
		return nil, errors.New("the response of delete object API is not as expected")
	}

	return response, nil
}
