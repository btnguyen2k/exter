// package user contains User business object (BO) and DAO implementations.
package user

import (
	"encoding/json"
	"fmt"
	"log"

	"main/src/gvabe/bo"
	"main/src/utils"
)

// User is the business object
//	- User inherits unique id from bo.UniversalBo
//	- Email address is used to uniquely identify user (e.g. user-id is email address)
type User struct {
	*bo.UniversalBo `json:"-"`
	AesKey          string `json:"aes_key"`
}

func (user *User) sync() *User {
	js, _ := json.Marshal(user)
	user.UniversalBo.DataJson = string(js)
	return user
}

// NewUserFromUniversal is helper function to create User App bo from a universal bo
func NewUserFromUniversal(ubo *bo.UniversalBo) *User {
	if ubo == nil {
		return nil
	}
	js := []byte(ubo.DataJson)
	app := User{}
	if err := json.Unmarshal(js, &app); err != nil {
		log.Print(fmt.Sprintf("[WARN] NewUserFromUniversal - error unmarshalling JSON data: %e", err))
		return nil
	}
	app.UniversalBo = ubo.Clone()
	return &app
}

// NewUser is helper function to create new User bo
func NewUser(appVersion uint64, id string) *User {
	user := &User{
		UniversalBo: bo.NewUniversalBo(id, appVersion),
		AesKey:      utils.RandomString(16),
	}
	return user.sync()
}

// UserDao defines API to access User storage
type UserDao interface {
	// Delete removes the specified business object from storage
	Delete(bo *User) (bool, error)

	// Create persists a new business object to storage
	Create(bo *User) (bool, error)

	// Get retrieves a business object from storage
	Get(username string) (*User, error)

	// GetN retrieves N business objects from storage
	GetN(fromOffset, maxNumRows int) ([]*User, error)

	// GetAll retrieves all available business objects from storage
	GetAll() ([]*User, error)

	// Update modifies an existing business object
	Update(bo *User) (bool, error)
}
