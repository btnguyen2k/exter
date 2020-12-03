package app

import (
	"main/src/gvabe/bo/user"
)

const (
	TableApp = "exter_app"
)

// AppDao defines API to access App storage.
type AppDao interface {
	// Delete removes the specified business object from storage.
	Delete(bo *App) (bool, error)

	// Create persists a new business object to storage.
	Create(bo *App) (bool, error)

	// Get retrieves a business object from storage.
	Get(id string) (*App, error)

	// // getN retrieves N business objects from storage.
	// getN(fromOffset, maxNumRows int) ([]*App, error)
	//
	// // getAll retrieves all available business objects from storage.
	// getAll() ([]*App, error)

	// GetUserApps retrieves all apps belong to a specific user.
	GetUserApps(u *user.User) ([]*App, error)

	// Update modifies an existing business object.
	Update(bo *App) (bool, error)
}
