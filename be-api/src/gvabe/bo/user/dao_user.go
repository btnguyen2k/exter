package user

const (
	TableUser = "exter_user"
)

// UserDao defines API to access User storage.
type UserDao interface {
	// Delete removes the specified business object from storage.
	Delete(bo *User) (bool, error)

	// Create persists a new business object to storage.
	Create(bo *User) (bool, error)

	// Get retrieves a business object from storage.
	Get(username string) (*User, error)

	// GetN retrieves N business objects from storage.
	GetN(fromOffset, maxNumRows int) ([]*User, error)

	// GetAll retrieves all available business objects from storage.
	GetAll() ([]*User, error)

	// Update modifies an existing business object.
	Update(bo *User) (bool, error)
}
