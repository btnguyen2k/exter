// Package app contains business object (BO) and data access object (DAO) implementations for Application.
package app

import (
	"encoding/json"
	"log"
	"net/url"
	"reflect"
	"strings"

	"github.com/btnguyen2k/consu/reddo"
	"github.com/btnguyen2k/henge"
	"main/src/gvabe/bo"
)

// NewApp is helper function to create new App bo.
func NewApp(tagVersion uint64, id, ownerId, desc string) *App {
	app := &App{
		UniversalBo: henge.NewUniversalBo(id, tagVersion, henge.UboOpt{TimeLayout: bo.UboTimeLayout, TimestampRounding: bo.UboTimestampRounding}),
	}
	app.
		SetOwnerId(ownerId).
		SetAttrsPublic(AppAttrsPublic{IsActive: true, Description: strings.TrimSpace(desc)})
	return app.sync()
}

var typMapStrBool = reflect.TypeOf(map[string]bool{})
var typSliceStr = reflect.TypeOf([]string{})

// NewAppFromUbo is helper function to create new App bo from a universal bo.
func NewAppFromUbo(ubo *henge.UniversalBo) *App {
	if ubo == nil {
		return nil
	}
	ubo = ubo.Clone()
	app := &App{UniversalBo: ubo}
	if ownerId, err := app.GetExtraAttrAs(FieldAppOwnerId, reddo.TypeString); err == nil {
		app.SetOwnerId(ownerId.(string))
	}
	if publicAttrsRaw, err := app.GetDataAttr(AttrAppPublicAttrs); err == nil && publicAttrsRaw != nil {
		var publicAttrs AppAttrsPublic
		var v interface{}
		ok := false
		if v, ok = publicAttrsRaw.(AppAttrsPublic); ok {
			publicAttrs = v.(AppAttrsPublic)
		} else if v, ok = publicAttrsRaw.(*AppAttrsPublic); ok {
			publicAttrs = *v.(*AppAttrsPublic)
		}
		if !ok {
			if v, err := app.GetDataAttrAs(AttrAppPublicAttrs+".actv", reddo.TypeBool); err == nil && v != nil {
				publicAttrs.IsActive = v.(bool)
			}
			if v, err := app.GetDataAttrAs(AttrAppPublicAttrs+".desc", reddo.TypeString); err == nil && v != nil {
				publicAttrs.Description = strings.TrimSpace(v.(string))
			}
			if v, err := app.GetDataAttrAs(AttrAppPublicAttrs+".rurl", reddo.TypeString); err == nil && v != nil {
				publicAttrs.DefaultReturnUrl = strings.TrimSpace(v.(string))
			}
			if v, err := app.GetDataAttrAs(AttrAppPublicAttrs+".curl", reddo.TypeString); err == nil && v != nil {
				publicAttrs.DefaultCancelUrl = strings.TrimSpace(v.(string))
			}
			if v, err := app.GetDataAttrAs(AttrAppPublicAttrs+".rpub", reddo.TypeString); err == nil && v != nil {
				publicAttrs.RsaPublicKey = strings.TrimSpace(v.(string))
			}
			if v, err := app.GetDataAttrAs(AttrAppPublicAttrs+".isrc", typMapStrBool); err == nil && v != nil {
				publicAttrs.IdentitySources = v.(map[string]bool)
			}
			if v, err := app.GetDataAttrAs(AttrAppPublicAttrs+".tags", typSliceStr); err == nil && v != nil {
				publicAttrs.Tags = v.([]string)
			}
		}
		app.SetAttrsPublic(publicAttrs)
	}

	return app.sync()
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
	FieldAppOwnerId = "oid"

	AttrAppPublicAttrs = "apub"
	AttrAppUbo         = "_ubo"
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
// - if 'preferredCancelUrl' is invalid, this function returns empty string
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
		AttrAppUbo:         app.UniversalBo.Clone(),
		FieldAppOwnerId:    app.ownerId,
		AttrAppPublicAttrs: app.attrsPublic.clone(),
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
	if m[AttrAppUbo] != nil {
		js, _ := json.Marshal(m[AttrAppUbo])
		if err := json.Unmarshal(js, &app.UniversalBo); err != nil {
			return err
		}
	}
	if m[AttrAppPublicAttrs] != nil {
		js, _ := json.Marshal(m[AttrAppPublicAttrs])
		if err := json.Unmarshal(js, &app.attrsPublic); err != nil {
			return err
		}
	}
	if app.ownerId, err = reddo.ToString(m[FieldAppOwnerId]); err != nil {
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
	app.SetDataAttr(AttrAppPublicAttrs, app.attrsPublic)
	app.SetExtraAttr(FieldAppOwnerId, app.ownerId)
	app.UniversalBo.Sync()
	return app
}
