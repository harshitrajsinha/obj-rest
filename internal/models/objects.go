// Package models defines data structures and functions that are used across the application
package models

// DataField represents data structure of the `data` field of an object
type DataField struct {
	Color        *string  `json:"color,omitempty"`
	Capacity     *string  `json:"capacity,omitempty"`
	Price        *float32 `json:"price,omitempty"`
	Generation   *string  `json:"generation,omitempty"`
	CPUModel     *string  `json:"CPU model,omitempty"`
	HardDiskSize *string  `json:"Hard disk size,omitempty"`
	StrapColor   *string  `json:"Strap Color,omitempty"`
	CaseSize     *string  `json:"Case size,omitempty"`
	ScreenSize   *string  `json:"Screen size,omitempty"`
	Description  *string  `json:"Description,omitempty"`
}

// ObjDataFromResponse represents strucutre of an object that will be received from response
type ObjDataFromResponse struct {
	ID   string                 `json:"id"`
	Name string                 `json:"name"`
	Data map[string]interface{} `json:"data,omitempty"`
}

// ObjDataPayload represents strucutre of an object that will be send as a payload to create or update object
type ObjDataPayload struct {
	Name string                 `json:"name"`
	Data map[string]interface{} `json:"data,omitempty"`
}

// NewObj represents strucutre of an object that has newly been created
type NewObj struct {
	ID        string                 `json:"id"`
	CreatedAt string                 `json:"createdAt"`
	Name      string                 `json:"name"`
	Data      map[string]interface{} `json:"data,omitempty"`
}
