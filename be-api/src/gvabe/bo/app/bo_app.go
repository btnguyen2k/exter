// package app contains App business object (BO) and DAO implementations.
package app

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/btnguyen2k/consu/reddo"

	"main/src/gvabe/bo/user"
	"main/src/henge"
)

// NewApp is helper function to create new App bo
func NewApp(appVersion uint64, id, ownerId, desc string) *App {
	app := &App{
		UniversalBo: *henge.NewUniversalBo(id, appVersion),
		ownerId:     strings.TrimSpace(strings.ToLower(ownerId)),
		attrsPublic: AppAttrsPublic{
			IsActive:    true,
			Description: strings.TrimSpace(desc),
		},
	}
	return app.sync()
}

// NewAppFromUbo is helper function to create new App bo from a universal bo
func NewAppFromUbo(ubo *henge.UniversalBo) *App {
	if ubo == nil {
		return nil
	}
	app := App{}
	if err := json.Unmarshal([]byte(ubo.GetDataJson()), &app); err != nil {
		log.Print(fmt.Sprintf("[WARN] NewAppFromUbo - error unmarshalling JSON data: %e", err))
		log.Print(err)
		return nil
	}
	app.UniversalBo = *ubo.Clone()
	if ownerId, err := app.GetExtraAttrAs(FieldApp_OwnerId, reddo.TypeString); err == nil {
		app.ownerId = ownerId.(string)
	}
	return &app
}

type AppAttrsPublic struct {
	IsActive         bool            `json:"actv"` // is this app active or not
	Description      string          `json:"desc"` // description text
	DefaultReturnUrl string          `json:"rurl"` // default return url after login
	IdentitySources  map[string]bool `json:"isrc"` // sources of identity
	Tags             []string        `json:"tags"` // arbitrary tags
	RsaPublicKey     string          `json:"rpub"` // RSA public key in ASCII-armor format
}

func (apub AppAttrsPublic) clone() AppAttrsPublic {
	clone := AppAttrsPublic{
		IsActive:         apub.IsActive,
		Description:      apub.Description,
		DefaultReturnUrl: apub.DefaultReturnUrl,
		RsaPublicKey:     apub.RsaPublicKey,
	}
	if apub.IdentitySources != nil {
		clone.IdentitySources = make(map[string]bool)
		for k, v := range apub.IdentitySources {
			clone.IdentitySources[k] = v
		}
	}
	if apub.Tags != nil {
		clone.Tags = append([]string{}, apub.Tags...)
	}
	return clone
}

const (
	FieldApp_OwnerId    = "oid"
	AppAttr_PublicAttrs = "apub"
	AppAttr_Ubo         = "_ubo"
)

// App is the business object
//	- App inherits unique id from bo.UniversalBo
type App struct {
	henge.UniversalBo `json:"_ubo"`
	ownerId           string         `json:"oid"`  // user id who owns this app
	attrsPublic       AppAttrsPublic `json:"apub"` // app's public attributes, can be access publicly
}

// MarshalJSON implements json.encode.Marshaler.MarshalJSON
func (app *App) MarshalJSON() ([]byte, error) {
	app.sync()
	m := map[string]interface{}{
		AppAttr_Ubo:         app.UniversalBo.Clone(),
		FieldApp_OwnerId:    app.ownerId,
		AppAttr_PublicAttrs: app.attrsPublic.clone(),
	}
	return json.Marshal(m)
}

// UnmarshalJSON implements json.decode.Unmarshaler.UnmarshalJSON
func (app *App) UnmarshalJSON(data []byte) error {
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	var err error
	if m[AppAttr_Ubo] != nil {
		js, _ := json.Marshal(m[AppAttr_Ubo])
		if err := json.Unmarshal(js, &app.UniversalBo); err != nil {
			return err
		}
	}
	if m[AppAttr_PublicAttrs] != nil {
		js, _ := json.Marshal(m[AppAttr_PublicAttrs])
		if err := json.Unmarshal(js, &app.attrsPublic); err != nil {
			return err
		}
	}
	if app.ownerId, err = reddo.ToString(m[FieldApp_OwnerId]); err != nil {
		return err
	}
	return nil
}

// GetOwnerId returns app's 'owner-id' value
func (app *App) GetOwnerId() string {
	return app.ownerId
}

// GetOwnerId sets app's 'owner-id' value
func (app *App) SetOwnerId(value string) *App {
	app.ownerId = strings.TrimSpace(strings.ToLower(value))
	return app
}

// GetAttrsPublic returns app's public attributes
func (app *App) GetAttrsPublic() AppAttrsPublic {
	return app.attrsPublic.clone()
}

// SetAttrsPublic sets app's public attributes
func (app *App) SetAttrsPublic(apub AppAttrsPublic) *App {
	app.attrsPublic = apub.clone()
	return app
}

func (app *App) sync() *App {
	app.SetDataAttr(AppAttr_PublicAttrs, app.attrsPublic)
	app.SetExtraAttr(FieldApp_OwnerId, app.ownerId)
	app.UniversalBo.Sync()
	return app
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

	// GetUserApps retrieves all apps belong to a specific user
	GetUserApps(u *user.User) ([]*App, error)

	// Update modifies an existing business object
	Update(bo *App) (bool, error)
}
