package session

const (
	TableSession = "exter_session"
)

// SessionDao defines API to access Session storage.
type SessionDao interface {
	// Delete removes the specified business object from storage.
	Delete(bo *Session) (bool, error)

	// Get retrieves a business object from storage.
	Get(id string) (*Session, error)

	// Save persists a new business object to storage or update an existing one.
	Save(bo *Session) (bool, error)
}
