package store

import (
	"context"

	"github.com/harshitrajsinha/obj-rest/internal/models"
)

// ObjectDataAccessor creates an interface for various operations that will be in sync to the endpoints provided by external API
// The method signature declared will then be implemented to call the external API and for testing.
type ObjectDataAccessor interface {
	GetAllObjs(ctx context.Context) ([]models.ObjDataFromResponse, error)
	GetObjsByIDs(ctx context.Context, IDs ...string) ([]models.ObjDataFromResponse, error)
	GetObjByID(ctx context.Context, ID string) (models.ObjDataFromResponse, error)
	CreateNewObj(ctx context.Context, payload models.ObjDataPayload) (models.NewObj, error)
	UpdateObj(ctx context.Context, objID string, payload models.ObjDataPayload) (models.NewObj, error)
	UpdateObjPartially(ctx context.Context, objID string, payload models.ObjDataPayload) (models.NewObj, error)
	DeleteObj(ctx context.Context, objID string) (map[string]string, error)
}
