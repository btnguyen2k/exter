// package app contains business object (BO) and data access object (DAO) implementations for Application.
package app

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/btnguyen2k/consu/reddo"
	"github.com/btnguyen2k/henge"
)

// NewApp is helper function to create new App bo.
func NewApp(tagVersion uint64, id, ownerId, desc string) *App {
	app := &App{
		UniversalBo: henge.NewUniversalBo(id, tagVersion),
		ownerId:     strings.TrimSpace(strings.ToLower(ownerId)),
		attrsPublic: AppAttrsPublic{
			IsActive:    true,
			Description: strings.TrimSpace(desc),
		},
	}
	return app.sync()
}

// NewAppFromUbo is helper function to create new App bo from a universal bo.
func NewAppFromUbo(ubo *henge.UniversalBo) *App {
	if ubo == nil {
		return nil
	}
	app := App{UniversalBo: &henge.UniversalBo{}}
	if err := json.Unmarshal([]byte(ubo.GetDataJson()), &app); err != nil {
		log.Print(fmt.Sprintf("[WARN] NewAppFromUbo - error unmarshalling JSON data: %e", err))
		log.Print(err)
		return nil
	}
	app.UniversalBo = ubo.Clone()
	if ownerId, err := app.GetExtraAttrAs(FieldApp_OwnerId, reddo.TypeString); err == nil {
		app.ownerId = ownerId.(string)
	}
	return &app
}

// AppAttrsPublic holds application's public attributes.
type AppAttrsPublic struct {
	IsActive         bool            `json:"actv"` // is this app active or not
	Description      string          `json:"desc"` // description text
	DefaultReturnUrl string          `json:"rurl"` // default return url after login
	DefaultCancelUrl string          `json:"curl"` // default cancel url after login
	IdentitySources  map[string]bool `json:"isrc"` // sources of identity
	Tags             []string        `json:"tags"` // arbitrary tags
	RsaPublicKey     string          `json:"rpub"` // RSA public key in ASCII-armor format
}

func (apub AppAttrsPublic) clone() AppAttrsPublic {
	clone := AppAttrsPublic{
		IsActive:         apub.IsActive,
		Description:      apub.Description,
		DefaultReturnUrl: apub.DefaultReturnUrl,
		DefaultCancelUrl: apub.DefaultCancelUrl,
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
	FieldApp_OwnerId = "oid"

	AttrApp_PublicAttrs = "apub"
	AttrApp_Ubo         = "_ubo"
)

// App is the business object.
// App inherits unique id from bo.UniversalBo.
type App struct {
	*henge.UniversalBo `json:"_ubo"`
	ownerId            string         `json:"oid"`  // user id who owns this app
	attrsPublic        AppAttrsPublic `json:"apub"` // app's public attributes, can be access publicly
}

// GenerateReturnUrl validates 'preferredReturnUrl' and builds "return url" for the app.
//
// - if 'preferredReturnUrl' is invalid, this function returns empty string
func (app *App) GenerateReturnUrl(preferredReturnUrl string) string {
	preferredReturnUrl = strings.TrimSpace(preferredReturnUrl)
	if preferredReturnUrl == "" {
		return app.attrsPublic.DefaultReturnUrl
	}
	urlPreferredReturnUrl, err := url.Parse(preferredReturnUrl)
	if err != nil {
		log.Println("[WARN] Preferred return url is invalid: " + preferredReturnUrl)
		return ""
	}
	urlDefaultReturnUrl, err := url.Parse(strings.TrimSpace(app.attrsPublic.DefaultReturnUrl))
	if err != nil {
		log.Println("[WARN] Default return url is invalid: " + app.attrsPublic.DefaultReturnUrl)
		return ""
	}
	if !urlDefaultReturnUrl.IsAbs() {
		if urlPreferredReturnUrl.IsAbs() {
			log.Printf("[WARN] Preferred return url [%s] is not valid against app's default one [%s]", preferredReturnUrl, app.attrsPublic.DefaultReturnUrl)
			return ""
		}
		return preferredReturnUrl
	}
	if !urlPreferredReturnUrl.IsAbs() {
		return urlDefaultReturnUrl.Scheme + "://" + urlDefaultReturnUrl.Host + "/" + strings.TrimPrefix(preferredReturnUrl, "/")
	}
	if urlDefaultReturnUrl.Host != urlPreferredReturnUrl.Host {
		log.Printf("[WARN] Preferred return url [%s] is not valid against app's default one [%s]", preferredReturnUrl, app.attrsPublic.DefaultReturnUrl)
		return ""
	}
	return preferredReturnUrl
}

// GenerateCancelUrl validates 'preferredCancelUrl' and builds "cancel url" for the app.
//
//	- if 'preferredCancelUrl' is invalid, this function returns empty string
func (app *App) GenerateCancelUrl(preferredCancelUrl string) string {
	preferredCancelUrl = strings.TrimSpace(preferredCancelUrl)
	if preferredCancelUrl == "" {
		return app.attrsPublic.DefaultCancelUrl
	}
	urlPreferredCancelUrl, err := url.Parse(preferredCancelUrl)
	if err != nil {
		log.Println("[WARN] Preferred return url is invalid: " + preferredCancelUrl)
		return ""
	}
	urlDefaultCancelUrl, err := url.Parse(strings.TrimSpace(app.attrsPublic.DefaultCancelUrl))
	if err != nil {
		log.Println("[WARN] Default cancel url is invalid: " + app.attrsPublic.DefaultCancelUrl)
		return ""
	}
	if !urlDefaultCancelUrl.IsAbs() {
		if urlPreferredCancelUrl.IsAbs() {
			log.Printf("[WARN] Preferred cancel url [%s] is not valid against app's default one [%s]", preferredCancelUrl, app.attrsPublic.DefaultCancelUrl)
			return ""
		}
		return preferredCancelUrl
	}
	if !urlPreferredCancelUrl.IsAbs() {
		return urlDefaultCancelUrl.Scheme + "://" + urlDefaultCancelUrl.Host + "/" + strings.TrimPrefix(preferredCancelUrl, "/")
	}
	if urlDefaultCancelUrl.Host != urlPreferredCancelUrl.Host {
		log.Printf("[WARN] Preferred cancel url [%s] is not valid against app's default one [%s]", preferredCancelUrl, app.attrsPublic.DefaultCancelUrl)
		return ""
	}
	return preferredCancelUrl
}

// MarshalJSON implements json.encode.Marshaler.MarshalJSON.
//	TODO: lock for read?
func (app *App) MarshalJSON() ([]byte, error) {
	app.sync()
	m := map[string]interface{}{
		AttrApp_Ubo:         app.UniversalBo.Clone(),
		FieldApp_OwnerId:    app.ownerId,
		AttrApp_PublicAttrs: app.attrsPublic.clone(),
	}
	return json.Marshal(m)
}

// UnmarshalJSON implements json.decode.Unmarshaler.UnmarshalJSON.
//	TODO: lock for write?
func (app *App) UnmarshalJSON(data []byte) error {
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	var err error
	if m[AttrApp_Ubo] != nil {
		js, _ := json.Marshal(m[AttrApp_Ubo])
		if err := json.Unmarshal(js, &app.UniversalBo); err != nil {
			return err
		}
	}
	if m[AttrApp_PublicAttrs] != nil {
		js, _ := json.Marshal(m[AttrApp_PublicAttrs])
		if err := json.Unmarshal(js, &app.attrsPublic); err != nil {
			return err
		}
	}
	if app.ownerId, err = reddo.ToString(m[FieldApp_OwnerId]); err != nil {
		return err
	}
	app.sync()
	return nil
}

// GetOwnerId returns app's 'owner-id' value.
func (app *App) GetOwnerId() string {
	return app.ownerId
}

// GetOwnerId sets app's 'owner-id' value.
func (app *App) SetOwnerId(value string) *App {
	app.ownerId = strings.TrimSpace(strings.ToLower(value))
	return app
}

// GetAttrsPublic returns app's public attributes.
func (app *App) GetAttrsPublic() AppAttrsPublic {
	return app.attrsPublic.clone()
}

// SetAttrsPublic sets app's public attributes.
func (app *App) SetAttrsPublic(apub AppAttrsPublic) *App {
	app.attrsPublic = apub.clone()
	return app
}

func (app *App) sync() *App {
	app.SetDataAttr(AttrApp_PublicAttrs, app.attrsPublic)
	app.SetExtraAttr(FieldApp_OwnerId, app.ownerId)
	app.UniversalBo.Sync()
	return app
}
