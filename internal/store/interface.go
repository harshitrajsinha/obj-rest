package store

import (
	"context"

	"github.com/harshitrajsinha/obj-rest/internal/models"
)

// ObjectDataAccessor creates an interface for various operations that will be in sync to the endpoints provided by external API
// The method signature declared will then be implemented to call the external API and for testing.
type ObjectDataAccessor interface {
	GetAllObjects(ctx context.Context) ([]models.ObjDataFromResponse, error)
	GetObjectsByIDs(ctx context.Context, IDs ...string) ([]models.ObjDataFromResponse, error)
	GetObjectByID(ctx context.Context, ID string) (models.ObjDataFromResponse, error)
	CreateNewObject(ctx context.Context, payload models.ObjDataPayload) (models.NewObj, error)
	UpdateObject(ctx context.Context, objID string, payload models.ObjDataPayload) (models.NewObj, error)
	UpdateObjectPartially(ctx context.Context, objID string, payload models.ObjDataPayload) (models.NewObj, error)
	DeleteObject(ctx context.Context, objID string) (map[string]string, error)
}
