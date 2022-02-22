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
	if v, err := app.GetExtraAttrAs(FieldAppOwnerId, reddo.TypeString); err == nil && v != nil {
		app.SetOwnerId(v.(string))
	}
	if v, err := app.GetDataAttrAs(AttrAppDomains, typSliceStr); err == nil && v != nil {
		app.SetDomains(v.([]string))
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

	AttrAppDomains     = "domains"
	AttrAppPublicAttrs = "apub"
	AttrAppUbo         = "_ubo"
)

// App is the business object.
// App inherits unique id from bo.UniversalBo.
type App struct {
	*henge.UniversalBo `json:"_ubo"`
	ownerId            string         `json:"oid"`     // user id who owns this app
	domains            []string       `json:"domains"` // app's domain whitelist (must contain domains from AppAttrsPublic.DefaultReturnUrl and AppAttrsPublic.DefaultCancelUrl)
	attrsPublic        AppAttrsPublic `json:"apub"`    // app's public attributes, can be access publicly
}

// _generateUrl validates 'preferred-url' and build the final url.
// If 'preferred-url' is invalid, this function returns empty string.
func _generateUrl(preferredUrl, defaultUrl string, whitelistDomains []string) string {
	preferredUrl = strings.TrimSpace(preferredUrl)
	if preferredUrl == "" {
		return defaultUrl
	}
	uPreferredUrl, err := url.Parse(preferredUrl)
	if err != nil {
		log.Printf("[WARN] Preferred url is invalid: %s", preferredUrl)
		return ""
	}

	if !uPreferredUrl.IsAbs() {
		// if preferred-url is relative, then default-url must not be empty
		if defaultUrl == "" {
			log.Printf("[WARN] Default url is empty")
			return ""
		}

		uDefaultUrl, err := url.Parse(defaultUrl)
		if err != nil {
			log.Printf("[WARN] Default url is invalid: %s", defaultUrl)
			return ""
		}
		if !uDefaultUrl.IsAbs() {
			// preferred-url and default-url are both relative
			return preferredUrl
		}
		// default-url is absolute, complete the url by prepending default-url's scheme and host
		return uDefaultUrl.Scheme + "://" + uDefaultUrl.Host + "/" + strings.TrimPrefix(preferredUrl, "/")
	}

	// if preferred-url is absolute, its host must be in whitelist
	for _, domain := range whitelistDomains {
		if uPreferredUrl.Host == domain {
			return preferredUrl
		}
	}
	log.Printf("[WARN] Preferred url [%s] is not in whitelist.", preferredUrl)
	return ""

	// if !uDefaultUrl.IsAbs() {
	// 	// if default-url is relative, then preferred-url must be relative too
	// 	if uPreferredUrl.IsAbs() {
	// 		log.Printf("[WARN] Preferred url [%s] is not valid against default one [%s]", preferredUrl, defaultUrl)
	// 		return ""
	// 	}
	// 	return preferredUrl
	// }
	// if !uPreferredUrl.IsAbs() {
	// 	// preferred-url is relative, complete the url by prepending default-url's scheme and host
	// 	return uDefaultUrl.Scheme + "://" + uDefaultUrl.Host + "/" + strings.TrimPrefix(preferredUrl, "/")
	// }
	// if uDefaultUrl.Host != uPreferredUrl.Host {
	// 	// preferred-url is absolute, its host must match default-url, or in whitelist
	// 	log.Printf("[WARN] Preferred url [%s] is not valid against default one [%s]", preferredUrl, defaultUrl)
	// 	return ""
	// }
	// return preferredUrl
}

// GenerateReturnUrl validates 'preferredReturnUrl' and builds "return url" for the app.
//
// - if 'preferredReturnUrl' is invalid, this function returns empty string
func (app *App) GenerateReturnUrl(preferredReturnUrl string) string {
	return _generateUrl(preferredReturnUrl, app.attrsPublic.DefaultReturnUrl, app.domains)
}

// GenerateCancelUrl validates 'preferredCancelUrl' and builds "cancel url" for the app.
//
// - if 'preferredCancelUrl' is invalid, this function returns empty string
func (app *App) GenerateCancelUrl(preferredCancelUrl string) string {
	return _generateUrl(preferredCancelUrl, app.attrsPublic.DefaultCancelUrl, app.domains)
}

// MarshalJSON implements json.encode.Marshaler.MarshalJSON.
//	TODO: lock for read?
func (app *App) MarshalJSON() ([]byte, error) {
	app.sync()
	m := map[string]interface{}{
		AttrAppUbo: app.UniversalBo.Clone(),
		bo.SerKeyFields: map[string]interface{}{
			FieldAppOwnerId: app.GetOwnerId(),
		},
		bo.SerKeyAttrs: map[string]interface{}{
			AttrAppDomains:     app.GetDomains(),
			AttrAppPublicAttrs: app.attrsPublic.clone(),
		},
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

	if m[AttrAppUbo] != nil {
		js, _ := json.Marshal(m[AttrAppUbo])
		if err := json.Unmarshal(js, &app.UniversalBo); err != nil {
			return err
		}
	}
	if _cols, ok := m[bo.SerKeyFields].(map[string]interface{}); ok {
		if v, err := reddo.ToString(_cols[FieldAppOwnerId]); err != nil {
			return err
		} else {
			app.SetOwnerId(v)
		}
	}
	if _attrs, ok := m[bo.SerKeyAttrs].(map[string]interface{}); ok {
		if v, err := reddo.ToSlice(_attrs[AttrAppDomains], typSliceStr); err != nil {
			return err
		} else {
			app.SetDomains(v.([]string))
		}
		if _attrs[AttrAppPublicAttrs] != nil {
			js, _ := json.Marshal(_attrs[AttrAppPublicAttrs])
			if err := json.Unmarshal(js, &app.attrsPublic); err != nil {
				return err
			}
		}
	}

	app.sync()
	return nil
}

// GetOwnerId returns app's 'owner-id' value.
func (app *App) GetOwnerId() string {
	return app.ownerId
}

// SetOwnerId sets app's 'owner-id' value.
func (app *App) SetOwnerId(value string) *App {
	app.ownerId = strings.TrimSpace(strings.ToLower(value))
	return app
}

// GetDomains returns app's 'whitelist-domains' value.
//
// Available since v0.7.0
func (app *App) GetDomains() []string {
	domains := make([]string, len(app.domains))
	copy(domains, app.domains)
	return domains
}

// SetDomains sets app's 'whitelist-domains' value.
//
// Available since v0.7.0
func (app *App) SetDomains(value []string) *App {
	domainsMap := make(map[string]bool)
	for _, domain := range value {
		domainsMap[domain] = true
	}
	app.domains = make([]string, len(domainsMap))
	i := 0
	for k := range domainsMap {
		app.domains[i] = k
		i++
	}
	return app
}

// GetAttrsPublic returns app's public attributes.
func (app *App) GetAttrsPublic() AppAttrsPublic {
	return app.attrsPublic.clone()
}

// SetAttrsPublic sets app's public attributes.
func (app *App) SetAttrsPublic(apub AppAttrsPublic) *App {
	app.attrsPublic = apub.clone()
	domains := app.GetDomains()
	if u, e := url.Parse(app.attrsPublic.DefaultReturnUrl); e == nil && u.Host != "" {
		domains = append(domains, u.Host)
	}
	if u, e := url.Parse(app.attrsPublic.DefaultCancelUrl); e == nil && u.Host != "" {
		domains = append(domains, u.Host)
	}
	app.SetDomains(domains)
	return app
}

func (app *App) sync() *App {
	app.SetExtraAttr(FieldAppOwnerId, app.ownerId)
	app.SetDataAttr(AttrAppDomains, app.domains)
	app.SetDataAttr(AttrAppPublicAttrs, app.attrsPublic)
	app.UniversalBo.Sync()
	return app
}
