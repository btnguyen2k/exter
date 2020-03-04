package gvabe

import (
	"context"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"strings"
	"time"

	"github.com/btnguyen2k/consu/reddo"
	"golang.org/x/oauth2"

	"main/src/goapi"
	"main/src/itineris"
)

var (
	// those APIs will not need authentication.
	// "false" means client, however, needs to sends app-id along with the API call
	// "true" means the API is free for public call
	publicApis = map[string]bool{
		"login":           false,
		"info":            true,
		"getApp":          false,
		"checkLoginToken": true,
	}
)

// API handler "info"
func apiInfo(_ *itineris.ApiContext, auth *itineris.ApiAuth, params *itineris.ApiParams) *itineris.ApiResult {
	appInfo := map[string]interface{}{
		"name":        goapi.AppConfig.GetString("app.name"),
		"shortname":   goapi.AppConfig.GetString("app.shortname"),
		"version":     goapi.AppConfig.GetString("app.version"),
		"description": goapi.AppConfig.GetString("app.desc"),
	}

	var publicPEM []byte
	if pubDER, err := x509.MarshalPKIXPublicKey(rsaPubKey); err == nil {
		pubBlock := pem.Block{
			Type:    "PUBLIC KEY",
			Headers: nil,
			Bytes:   pubDER,
		}
		publicPEM = pem.EncodeToMemory(&pubBlock)
	} else {
		publicPEM = []byte(err.Error())
	}

	result := map[string]interface{}{
		"app":              appInfo,
		"rsa_public_key":   string(publicPEM),
		"google_client_id": gConfig.ClientID,
	}

	return itineris.NewApiResult(itineris.StatusOk).SetData(result)
}

func _doLoginGoogle(_ *itineris.ApiContext, aauth *itineris.ApiAuth, authCode string) *itineris.ApiResult {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	if token, err := gConfig.Exchange(ctx, authCode, oauth2.AccessTypeOnline); err != nil {
		return itineris.NewApiResult(itineris.StatusNoPermission).SetMessage(err.Error())
	} else if token == nil {
		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage("Error: exchanged token is nil")
	} else {
		js, _ := json.Marshal(token)
		now := time.Now()
		session := Session{
			ClientId:  aauth.GetAppId(),
			Channel:   loginChannelGoogle,
			CreatedAt: now,
			ExpiredAt: now.Add(2 * 60 * time.Second), // 2 mins
			Data:      js,                            // JSON-serialization of oauth2.Token
		}
		preLoginToken, err := genPreLoginToken(session, "")
		if err != nil {
			return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(err.Error())
		}
		goFetchGoogleProfile(preLoginToken)
		return itineris.NewApiResult(itineris.StatusOk).SetData(preLoginToken)
	}
}

/*
apiLogin handles API call "login".
Upon login successfully, this API returns the login token as JWT.
*/
func apiLogin(ctx *itineris.ApiContext, auth *itineris.ApiAuth, params *itineris.ApiParams) *itineris.ApiResult {
	source, _ := params.GetParamAsType("source", reddo.TypeString)
	if source == nil {
		source = ""
	}
	switch strings.ToLower(source.(string)) {
	case loginChannelGoogle:
		authCode, _ := params.GetParamAsType("code", reddo.TypeString)
		if authCode == nil {
			authCode = ""
		}
		return _doLoginGoogle(ctx, auth, authCode.(string))
	}
	return itineris.NewApiResult(itineris.StatusErrorClient).SetMessage(fmt.Sprintf("Login source is not supported: %s", source))
}

/*
apiCheckLoginToken handles API call "checkLoginToken".
This API expects an input map:

	{
		"token": login token (returned by apiLogin/apiCheckLoginToken),
	}

Upon successful, this API returns the login-token.
*/
func apiCheckLoginToken(_ *itineris.ApiContext, _ *itineris.ApiAuth, params *itineris.ApiParams) *itineris.ApiResult {
	token, _ := params.GetParamAsType("token", reddo.TypeString)
	if token == nil || token == "" {
		return itineris.NewApiResult(itineris.StatusErrorClient).SetMessage("empty token")
	}
	claim, err := parseLoginToken(token.(string))
	if err != nil {
		return itineris.NewApiResult(itineris.StatusNoPermission).SetMessage(err.Error())
	}
	if claim.isExpired() {
		return itineris.NewApiResult(itineris.StatusNoPermission).SetMessage(errorExpiredJwt.Error())
	}
	if claim.Type == sessionTypePreLogin {
		session, err := loadPreLoginSessionFromCache(claim.CacheKey)
		if err != nil {
			return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(err.Error())
		}
		if session == nil {
			return itineris.NewApiResult(itineris.StatusNoPermission).SetMessage("session not found or expired")
		}
		if session.UserId == "" {
			return itineris.NewApiResult(302).SetMessage("please try again after a moment")
		}
		token, err = genLoginToken(*session, "")
		if err != nil {
			return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(err.Error())
		}
	}
	return itineris.NewApiResult(itineris.StatusOk).SetData(token.(string))
}

// API handler "getApp"
//	- This API return only app's public info
func apiGetApp(_ *itineris.ApiContext, _ *itineris.ApiAuth, params *itineris.ApiParams) *itineris.ApiResult {
	id, _ := params.GetParamAsType("id", reddo.TypeString)
	if id == nil || strings.TrimSpace(id.(string)) == "" {
		return itineris.NewApiResult(itineris.StatusNotFound).SetMessage(fmt.Sprintf("App [%s] not found", id))
	}
	if app, err := appDao.Get(id.(string)); err != nil {
		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(err.Error())
	} else if app == nil {
		return itineris.NewApiResult(itineris.StatusNotFound).SetMessage(fmt.Sprintf("App [%s] not found", id))
	} else {
		config := make(map[string]interface{})
		if app.Config != nil {
			js, _ := json.Marshal(app.Config)
			json.Unmarshal(js, &config)
			loginChannels := make(map[string]bool)
			for s, _ := range app.Config.IdentitySources {
				if app.Config.IdentitySources[s] && enabledLoginChannels[s] {
					loginChannels[s] = true
				}
			}
			config["sources"] = loginChannels
		}
		return itineris.NewApiResult(itineris.StatusOk).SetData(map[string]interface{}{"id": app.Id, "config": config})
	}
}
