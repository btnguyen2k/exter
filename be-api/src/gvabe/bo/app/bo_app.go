// package app contains App business object (BO) and DAO implementations.
package app

import (
	"encoding/json"
	"log"
	"strings"

	"main/src/gvabe/bo"
)

type appConfig struct {
	DefaultReturnUrl string          `json:"return_url"`
	IdentitySources  map[string]bool `json:"sources"`
	Tags             []string        `json:"tags"`
}

// App is the business object
//	- App inherits unique id from bo.UniversalBo
type App struct {
	*bo.UniversalBo `json:"-"`
	OwnerId         string     `json:"owner_id"` // user id who owns this app
	Description     string     `json:"desc"`
	IsActive        bool       `json:"active"` // is this app active or not
	Config          *appConfig `json:"config"`
}

func (app *App) sync() *App {
	js, _ := json.Marshal(app)
	app.UniversalBo.DataJson = string(js)
	return app
}

// NewAppFromUniversal is helper function to create new App bo from a universal bo
func NewAppFromUniversal(ubo *bo.UniversalBo) *App {
	if ubo == nil {
		return nil
	}
	js := []byte(ubo.DataJson)
	app := App{}
	if err := json.Unmarshal(js, &app); err != nil {
		log.Print(err)
		return nil
	}
	app.UniversalBo = ubo.Clone()
	if app.Config == nil {
		app.Config = &appConfig{}
	}
	return &app
}

// NewApp is helper function to create new App bo
func NewApp(appVersion uint64, id, ownerId, desc string) *App {
	app := &App{
		UniversalBo: bo.NewUniversalBo(id, appVersion),
		OwnerId:     strings.TrimSpace(strings.ToLower(ownerId)),
		Description: strings.TrimSpace(desc),
		IsActive:    true,
		Config:      &appConfig{},
	}
	return app.sync()
}

// AppDao defines API to access App storage
type AppDao interface {
	// Delete removes the specified business object from storage
	Delete(bo *App) (bool, error)

	// Create persists a new business object to storage
	Create(bo *App) (bool, error)

	// Get retrieves a business object from storage
	Get(id string) (*App, error)

	// GetN retrieves N business objects from storage
	GetN(fromOffset, maxNumRows int) ([]*App, error)

	// GetAll retrieves all available business objects from storage
	GetAll() ([]*App, error)

	// Update modifies an existing business object
	Update(bo *App) (bool, error)
}
