// package user contains business object (BO) and data access object (DAO) implementations for User.
package user

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/btnguyen2k/consu/reddo"

	"main/src/henge"
	"main/src/utils"
)

// NewUser is helper function to create new User bo
func NewUser(appVersion uint64, id string) *User {
	user := &User{
		UniversalBo: *henge.NewUniversalBo(id, appVersion),
	}
	user.SetAesKey(utils.RandomString(16))
	return user.sync()
}

// NewUserFromUbo is helper function to create User App bo from a universal bo
func NewUserFromUbo(ubo *henge.UniversalBo) *User {
	if ubo == nil {
		return nil
	}
	user := User{}
	if err := json.Unmarshal([]byte(ubo.GetDataJson()), &user); err != nil {
		log.Print(fmt.Sprintf("[WARN] NewUserFromUbo - error unmarshalling JSON data: %e", err))
		log.Print(err)
		return nil
	}
	user.UniversalBo = *ubo.Clone()
	return &user
}

const (
	AttrUser_AesKey = "aes_key"
	AttrUser_Ubo    = "_ubo"
)

// User is the business object
//	- User inherits unique id from bo.UniversalBo
//	- Email address is used to uniquely identify user (e.g. user-id is email address)
type User struct {
	henge.UniversalBo `json:"_ubo"`
}

// GetAesKey returns value of user's 'aes-key' attribute
func (user *User) GetAesKey() string {
	v, err := user.GetDataAttrAs(AttrUser_AesKey, reddo.TypeString)
	if err != nil || v == nil {
		return ""
	}
	return v.(string)
}

// SetAesKey sets value of user's 'aes-key' attribute
func (user *User) SetAesKey(v string) *User {
	user.SetDataAttr(AttrUser_AesKey, strings.TrimSpace(v))
	return user
}

func (user *User) sync() *User {
	user.UniversalBo.Sync()
	return user
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
