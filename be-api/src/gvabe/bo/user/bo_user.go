// Package user contains business object (BO) and data access object (DAO) implementations for User.
package user

import (
	"encoding/json"
	"strings"

	"github.com/btnguyen2k/consu/reddo"
	"main/src/gvabe/bo"

	"github.com/btnguyen2k/henge"

	"main/src/utils"
)

// NewUser is helper function to create new User bo.
func NewUser(appVersion uint64, id string) *User {
	user := &User{
		UniversalBo: henge.NewUniversalBo(id, appVersion, henge.UboOpt{TimeLayout: bo.UboTimeLayout, TimestampRounding: bo.UboTimestampRouding}),
	}
	user.
		SetDisplayName(id).
		SetAesKey(utils.RandomString(16))
	return user.sync()
}

// NewUserFromUbo is helper function to create User App bo from a universal bo.
func NewUserFromUbo(ubo *henge.UniversalBo) *User {
	if ubo == nil {
		return nil
	}
	ubo = ubo.Clone()
	user := &User{UniversalBo: ubo}
	if v, err := ubo.GetDataAttrAs(AttrUserAesKey, reddo.TypeString); err != nil {
		return nil
	} else if v != nil {
		user.SetAesKey(v.(string))
	}
	if v, err := ubo.GetDataAttrAs(AttrUserDisplayName, reddo.TypeString); err != nil {
		return nil
	} else if v != nil {
		user.SetDisplayName(v.(string))
	}
	return user.sync()
}

const (
	AttrUserUbo         = "_ubo"
	AttrUserAesKey      = "aes"
	AttrUserDisplayName = "dname"
)

// User is the business object.
// User inherits unique id from bo.UniversalBo. Email address is used to uniquely identify user (e.g. user-id is email address).
type User struct {
	*henge.UniversalBo `json:"_ubo"`
	aesKey             string `json:"aes"`
	displayName        string `json:"dname"`
}

// MarshalJSON implements json.encode.Marshaler.MarshalJSON.
//	TODO: lock for read?
func (u *User) MarshalJSON() ([]byte, error) {
	u.sync()
	m := map[string]interface{}{
		AttrUserUbo: u.UniversalBo.Clone(),
		bo.SerKeyAttrs: map[string]interface{}{
			AttrUserAesKey:      u.GetAesKey(),
			AttrUserDisplayName: u.GetDisplayName(),
		},
	}
	return json.Marshal(m)
}

// UnmarshalJSON implements json.decode.Unmarshaler.UnmarshalJSON.
//	TODO: lock for write?
func (u *User) UnmarshalJSON(data []byte) error {
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	var err error
	if m[AttrUserUbo] != nil {
		js, _ := json.Marshal(m[AttrUserUbo])
		if err = json.Unmarshal(js, &u.UniversalBo); err != nil {
			return err
		}
	}
	if _attrs, ok := m[bo.SerKeyAttrs].(map[string]interface{}); ok {
		if v, err := reddo.ToString(_attrs[AttrUserAesKey]); err != nil {
			return err
		} else {
			u.SetAesKey(v)
		}
		if v, err := reddo.ToString(_attrs[AttrUserDisplayName]); err != nil {
			return err
		} else {
			u.SetDisplayName(v)
		}
	}
	u.sync()
	return nil
}

// GetAesKey returns value of user's 'aes-key' attribute.
func (u *User) GetAesKey() string {
	return u.aesKey
}

// SetAesKey sets value of user's 'aes-key' attribute.
func (u *User) SetAesKey(v string) *User {
	u.aesKey = strings.TrimSpace(v)
	return u
}

// GetDisplayName returns value of user's 'display-name' attribute.
// available since v0.4.0
func (u *User) GetDisplayName() string {
	return u.displayName
}

// SetDisplayName sets value of user's 'display-name' attribute.
// available since v0.4.0
func (u *User) SetDisplayName(v string) *User {
	u.displayName = strings.TrimSpace(v)
	return u
}

func (u *User) sync() *User {
	u.SetDataAttr(AttrUserAesKey, u.aesKey)
	u.SetDataAttr(AttrUserDisplayName, u.displayName)
	u.UniversalBo.Sync()
	return u
}
