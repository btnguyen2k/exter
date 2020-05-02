// package bo defines business object and dao interface
package bo

import (
	"strings"
	"sync"
	"time"

	"github.com/btnguyen2k/godal"
)

const (
	FieldId          = "id"
	FieldData        = "data"
	FieldTimeCreated = "t_created"
	FieldTimeUpdated = "t_updated"
	FieldAppVersion  = "app_version"
)

var (
	allFields = []string{FieldId, FieldData, FieldTimeCreated, FieldTimeUpdated, FieldAppVersion}
)

// NewUniversalBo is helper function to create a new UniversalBo instance
func NewUniversalBo(id string, appVersion uint64) *UniversalBo {
	now := time.Now()
	return &UniversalBo{
		Id:          strings.ToLower(strings.TrimSpace(id)),
		TimeCreated: now,
		TimeUpdated: now,
		AppVersion:  appVersion,
	}
}

// UniversalBo is the "universal" business object
type UniversalBo struct {
	Id          string                 `json:"id"`          // bo's unique identifier
	DataJson    string                 `json:"data"`        // bo's attributes encoded as JSON string
	TimeCreated time.Time              `json:"t_created"`   // bo's creation timestamp
	TimeUpdated time.Time              `json:"t_updated"`   // bo's last update timestamp
	AppVersion  uint64                 `json:"app_version"` // for internal use
	extraFields map[string]interface{} `json:"-"`
	extraMutex  sync.Mutex
}

func (ubo *UniversalBo) SetExtraField(field string, value interface{}) *UniversalBo {
	ubo.extraMutex.Lock()
	defer ubo.extraMutex.Unlock()
	if ubo.extraFields == nil {
		ubo.extraFields = make(map[string]interface{})
	}
	ubo.extraFields[field] = value
	return ubo
}

// Clone creates a cloned copy of the "universal" business object
func (ubo *UniversalBo) Clone() *UniversalBo {
	return &UniversalBo{
		Id:          ubo.Id,
		DataJson:    ubo.DataJson,
		TimeCreated: ubo.TimeCreated,
		TimeUpdated: ubo.TimeUpdated,
		AppVersion:  ubo.AppVersion,
		extraFields: cloneMap(ubo.extraFields),
	}
}

// UniversalDao defines API to access UniversalBo storage
type UniversalDao interface {
	// ToUniversalBo transforms godal.IGenericBo to business object
	ToUniversalBo(gbo godal.IGenericBo) *UniversalBo

	// ToGenericBo transforms business object to godal.IGenericBo
	ToGenericBo(ubo *UniversalBo) godal.IGenericBo

	// Delete removes the specified business object from storage
	Delete(bo *UniversalBo) (bool, error)

	// Create persists a new business object to storage
	Create(bo *UniversalBo) (bool, error)

	// Get retrieves a business object from storage
	Get(id string) (*UniversalBo, error)

	// GetN retrieves N business objects from storage
	GetN(fromOffset, maxNumRows int) ([]*UniversalBo, error)

	// GetAll retrieves all available business objects from storage
	GetAll() ([]*UniversalBo, error)

	// Update modifies an existing business object
	Update(bo *UniversalBo) (bool, error)
}
