package gvabe

import (
	"context"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/btnguyen2k/consu/reddo"
	"golang.org/x/oauth2"

	"main/src/goapi"
	"main/src/gvabe/bo/app"
	"main/src/gvabe/bo/user"
	"main/src/itineris"
)

/*
Setup API handlers: application register its api-handlers by calling router.SetHandler(apiName, apiHandlerFunc)

    - api-handler function must has the following signature: func (itineris.ApiContext, itineris.ApiAuth, itineris.ApiParams) *itineris.ApiResult
*/
func initApiHandlers(router *itineris.ApiRouter) {
	router.SetHandler("info", apiInfo)
	router.SetHandler("login", apiLogin)
	router.SetHandler("verifyLoginToken", apiVerifyLoginToken)
	router.SetHandler("systemInfo", apiSystemInfo)

	router.SetHandler("getApp", apiGetApp)
	router.SetHandler("myAppList", apiMyAppList)
	router.SetHandler("getMyApp", apiGetMyApp)
	router.SetHandler("registerApp", apiRegisterApp)
	router.SetHandler("updateMyApp", apiUpdateMyApp)
	router.SetHandler("deleteMyApp", apiDeleteMyApp)
}

/*------------------------------ shared variables and functions ------------------------------*/

var (
	// those APIs will not need authentication.
	// "false" means client, however, needs to sends app-id along with the API call
	// "true" means the API is free for public call
	publicApis = map[string]bool{
		"login":            false,
		"info":             true,
		"getApp":           false,
		"verifyLoginToken": true,
		"loginChannelList": true,
	}
)

func _parseLoginTokenFromApi(_token interface{}) (*itineris.ApiResult, *SessionClaims, *user.User) {
	stoken, ok := _token.(string)
	if !ok || stoken == "" {
		return itineris.NewApiResult(itineris.StatusNoPermission).SetMessage("empty token"), nil, nil
	}

	var claim *SessionClaims
	var user *user.User
	var err error
	if claim, err = parseLoginToken(stoken); err != nil {
		return itineris.NewApiResult(itineris.StatusNoPermission).SetMessage(err.Error()), nil, nil
	} else if claim.isExpired() {
		return itineris.NewApiResult(itineris.StatusNoPermission).SetMessage(errorExpiredJwt.Error()), nil, nil
	} else if claim.Type != sessionTypeLogin {
		return itineris.NewApiResult(itineris.StatusNoPermission).SetMessage("invalid session type"), nil, nil
	}
	if user, err = userDao.Get(claim.UserId); err != nil {
		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(err.Error()), nil, nil
	} else if user == nil {
		return itineris.NewApiResult(itineris.StatusNoPermission).SetMessage("session user not found"), nil, nil
	}
	return nil, claim, user
}

func extractAppAttrsPublic(myApp *app.App) map[string]interface{} {
	result := make(map[string]interface{})
	attrsPublic := myApp.GetAttrsPublic()
	js, _ := json.Marshal(attrsPublic)
	json.Unmarshal(js, &result)
	loginChannels := make(map[string]bool)
	for s, _ := range attrsPublic.IdentitySources {
		if attrsPublic.IdentitySources[s] && enabledLoginChannels[s] {
			loginChannels[s] = true
		}
	}
	result["sources"] = loginChannels
	return result
}

/*------------------------------ public APIs ------------------------------*/

// API handler "info"
func apiInfo(_ *itineris.ApiContext, _ *itineris.ApiAuth, _ *itineris.ApiParams) *itineris.ApiResult {
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

	loginChannels := make([]string, 0)
	for channel, _ := range enabledLoginChannels {
		loginChannels = append(loginChannels, channel)
	}
	sort.Strings(loginChannels)
	result := map[string]interface{}{
		"app":              appInfo,
		"login_channels":   loginChannels,
		"rsa_public_key":   string(publicPEM),
		"google_client_id": googleOAuthConf.ClientID,
		"github_client_id": githubOAuthConf.ClientID,
		"facebook_app_id":  fbOAuthConf.ClientID,
	}

	return itineris.NewApiResult(itineris.StatusOk).SetData(result)
}

// API handler "systemInfo"
func apiSystemInfo(_ *itineris.ApiContext, _ *itineris.ApiAuth, _ *itineris.ApiParams) *itineris.ApiResult {
	data := lastSystemInfo()
	return itineris.NewApiResult(itineris.StatusOk).SetData(data)
}

/*------------------------------ login & session APIs ------------------------------*/

func _doLoginFacebook(_ *itineris.ApiContext, _ *itineris.ApiAuth, accessToken string, app *app.App, returnUrl string) *itineris.ApiResult {
	if DEBUG {
		log.Printf("[DEBUG] START _doLoginFacebook")
		t := time.Now().UnixNano()
		defer func() {
			d := time.Now().UnixNano() - t
			log.Printf("[DEBUG] END _doLoginFacebook: %d ms", d/1000000)
		}()
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	// firstly exchange for long-live token
	if token, err := fbExchangeForLongLiveToken(ctx, accessToken); err != nil {
		if DEBUG {
			log.Printf("[DEBUG] ERROR _doLoginFacebook: %s / %s", accessToken[len(accessToken)-4:], err)
		}
		return itineris.NewApiResult(itineris.StatusNoPermission).SetMessage(err.Error())
	} else if token == nil {
		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage("Error: exchanged token is nil")
	} else {
		// secondly embed accessToken into exter's session as a JWT
		js, _ := json.Marshal(token)
		now := time.Now()
		claims, err := genPreLoginClaims(&Session{
			ClientId:  app.GetId(),
			Channel:   loginChannelFacebook,
			CreatedAt: now,
			ExpiredAt: token.Expiry,
			Data:      js, // JSON-serialization of oauth2.Token
		})
		if err != nil {
			return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(err.Error())
		}
		_, jwt, err := saveSession(claims)
		if err != nil {
			return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(err.Error())
		}
		// lastly use accessToken to fetch Facebook profile info
		go goFetchFacebookProfile(claims.Id)
		returnUrl = strings.ReplaceAll(returnUrl, "${token}", jwt)
		return itineris.NewApiResult(itineris.StatusOk).SetData(jwt).SetExtras(map[string]interface{}{apiResultExtraReturnUrl: returnUrl})
	}
}

func _doLoginGitHub(_ *itineris.ApiContext, _ *itineris.ApiAuth, authCode string, app *app.App, returnUrl string) *itineris.ApiResult {
	if DEBUG {
		log.Printf("[DEBUG] START _doLoginGitHub")
		t := time.Now().UnixNano()
		defer func() {
			d := time.Now().UnixNano() - t
			log.Printf("[DEBUG] END _doLoginGitHub: %d ms", d/1000000)
		}()
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	// firstly exchange authCode for accessToken
	if token, err := githubOAuthConf.Exchange(ctx, authCode /*, oauth2.AccessTypeOnline*/); err != nil {
		if DEBUG {
			log.Printf("[DEBUG] ERROR _doLoginGithub: %s / %s", authCode[len(authCode)-4:], err)
		}
		return itineris.NewApiResult(itineris.StatusNoPermission).SetMessage(err.Error())
	} else if token == nil {
		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage("Error: exchanged token is nil")
	} else {
		now := time.Now()
		/*
			GitHub's OAuth token does not set expiry, which causes the token to be expired immediately.
			Hence we need to explicitly set expiry.
		*/
		token.Expiry = now.Add(1 * time.Hour)
		// secondly embed accessToken into exter's session as a JWT
		js, _ := json.Marshal(token)
		claims, err := genPreLoginClaims(&Session{
			ClientId:  app.GetId(),
			Channel:   loginChannelGithub,
			CreatedAt: now,
			ExpiredAt: token.Expiry,
			Data:      js, // JSON-serialization of oauth2.Token
		})
		if err != nil {
			return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(err.Error())
		}
		_, jwt, err := saveSession(claims)
		if err != nil {
			return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(err.Error())
		}
		// lastly use accessToken to fetch GitHub profile info
		go goFetchGitHubProfile(claims.Id)
		returnUrl = strings.ReplaceAll(returnUrl, "${token}", jwt)
		return itineris.NewApiResult(itineris.StatusOk).SetData(jwt).SetExtras(map[string]interface{}{apiResultExtraReturnUrl: returnUrl})
	}
}

func _doLoginGoogle(_ *itineris.ApiContext, _ *itineris.ApiAuth, authCode string, app *app.App, returnUrl string) *itineris.ApiResult {
	if DEBUG {
		log.Printf("[DEBUG] START _doLoginGoogle")
		t := time.Now().UnixNano()
		defer func() {
			d := time.Now().UnixNano() - t
			log.Printf("[DEBUG] END _doLoginGoogle: %d ms", d/1000000)
		}()
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	// firstly exchange authCode for accessToken
	if token, err := googleOAuthConf.Exchange(ctx, authCode, oauth2.AccessTypeOnline); err != nil {
		if DEBUG {
			log.Printf("[DEBUG] ERROR _doLoginGoogle: %s / %s", authCode[len(authCode)-4:], err)
		}
		return itineris.NewApiResult(itineris.StatusNoPermission).SetMessage(err.Error())
	} else if token == nil {
		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage("Error: exchanged token is nil")
	} else {
		// secondly embed accessToken into exter's session as a JWT
		js, _ := json.Marshal(token)
		now := time.Now()
		claims, err := genPreLoginClaims(&Session{
			ClientId:  app.GetId(),
			Channel:   loginChannelGoogle,
			CreatedAt: now,
			ExpiredAt: token.Expiry,
			Data:      js, // JSON-serialization of oauth2.Token
		})
		if err != nil {
			return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(err.Error())
		}
		_, jwt, err := saveSession(claims)
		if err != nil {
			return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(err.Error())
		}
		// lastly use accessToken to fetch Google profile info
		go goFetchGoogleProfile(claims.Id)
		returnUrl = strings.ReplaceAll(returnUrl, "${token}", jwt)
		return itineris.NewApiResult(itineris.StatusOk).SetData(jwt).SetExtras(map[string]interface{}{apiResultExtraReturnUrl: returnUrl})
	}
}

/*
apiLogin handles API call "login".

- Upon login successfully, this API returns the login token as JWT.
*/
func apiLogin(ctx *itineris.ApiContext, auth *itineris.ApiAuth, params *itineris.ApiParams) *itineris.ApiResult {
	appId := _extractParam(params, "app", reddo.TypeString, "", nil)
	app, err := appDao.Get(appId.(string))
	if err != nil {
		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(err.Error())
	} else if app == nil {
		return itineris.NewApiResult(itineris.StatusNoPermission).SetMessage(fmt.Sprintf("App [%s] not found", appId))
	} else if !app.GetAttrsPublic().IsActive {
		return itineris.NewApiResult(itineris.StatusNoPermission).SetMessage(fmt.Sprintf("App [%s] is not active", appId))
	}

	returnUrl := _extractParam(params, "return_url", reddo.TypeString, "", nil)
	if returnUrl != "" {
		if returnUrl = app.GenerateReturnUrl(returnUrl.(string)); returnUrl == "" {
			return itineris.NewApiResult(itineris.StatusNoPermission).SetMessage(fmt.Sprintf("Return url [%s] is not allowed for app [%s]", returnUrl, appId))
		}
	}

	source := _extractParam(params, "source", reddo.TypeString, "", nil)

	switch strings.ToLower(source.(string)) {
	case loginChannelGoogle:
		authCode := _extractParam(params, "code", reddo.TypeString, "", nil)
		return _doLoginGoogle(ctx, auth, authCode.(string), app, returnUrl.(string))
	case loginChannelGithub:
		authCode := _extractParam(params, "code", reddo.TypeString, "", nil)
		return _doLoginGitHub(ctx, auth, authCode.(string), app, returnUrl.(string))
	case loginChannelFacebook:
		authCode := _extractParam(params, "code", reddo.TypeString, "", nil)
		return _doLoginFacebook(ctx, auth, authCode.(string), app, returnUrl.(string))
	}
	return itineris.NewApiResult(itineris.StatusErrorClient).SetMessage(fmt.Sprintf("Login source is not supported: %s", source))
}

/*
apiVerifyLoginToken handles API call "verifyLoginToken".
This API expects an input map:

	{
		"token": login token (returned by apiLogin/apiVerifyLoginToken),
		"app": application's id,
	}

- Upon successful, this API returns the login-token.
*/
func apiVerifyLoginToken(_ *itineris.ApiContext, _ *itineris.ApiAuth, params *itineris.ApiParams) *itineris.ApiResult {
	// firstly extract JWT token from request and convert it into claims
	token := _extractParam(params, "token", reddo.TypeString, "", nil)
	if token == "" {
		return itineris.NewApiResult(itineris.StatusErrorClient).SetMessage("empty token")
	}
	claims, err := parseLoginToken(token.(string))
	if err != nil {
		return itineris.NewApiResult(itineris.StatusNoPermission).SetMessage(err.Error())
	}
	if claims.isExpired() {
		return itineris.NewApiResult(itineris.StatusNoPermission).SetMessage(errorExpiredJwt.Error())
	}

	// secondly verify the client app
	appId := _extractParam(params, "app", reddo.TypeString, "", nil)
	app, err := appDao.Get(appId.(string))
	if err != nil {
		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(err.Error())
	} else if app == nil || !app.GetAttrsPublic().IsActive {
		return itineris.NewApiResult(itineris.StatusNoPermission).SetMessage("invalid app")
	}

	// also verify 'return-url'
	returnUrl := _extractParam(params, "return_url", reddo.TypeString, "", nil)
	if returnUrl != "" {
		if returnUrl = app.GenerateReturnUrl(returnUrl.(string)); returnUrl == "" {
			return itineris.NewApiResult(itineris.StatusNoPermission).SetMessage(fmt.Sprintf("Return url [%s] is not allowed for app [%s]", returnUrl, appId))
		}
	}

	// thirdly verify the session
	sess, err := sessionDao.Get(claims.Id)
	if err != nil {
		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(err.Error())
	}
	if sess == nil || sess.IsExpired() {
		return itineris.NewApiResult(itineris.StatusNoPermission).SetMessage(fmt.Sprintf("Session not exists not expired"))
	}

	// lastly return the session encoded as JWT
	if sess.GetSessionType() == sessionTypePreLogin {
		return itineris.NewApiResult(302).SetMessage("please try again after a moment")
	} else {
		returnUrl = strings.ReplaceAll(returnUrl.(string), "${token}", sess.GetSessionData())
	}
	return itineris.NewApiResult(itineris.StatusOk).SetData(sess.GetSessionData()).SetExtras(map[string]interface{}{apiResultExtraReturnUrl: returnUrl})
}

/* app APIs */

/*
API handler "myAppList"
This API expects an input map:

	{
		"token": login token (returned by apiLogin/apiVerifyLoginToken),
	}
*/
func apiMyAppList(_ *itineris.ApiContext, _ *itineris.ApiAuth, params *itineris.ApiParams) *itineris.ApiResult {
	token, _ := params.GetParamAsType("token", reddo.TypeString)
	errResult, _, user := _parseLoginTokenFromApi(token)
	if errResult != nil {
		return errResult
	}
	appList, err := appDao.GetUserApps(user)
	if err != nil {
		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(err.Error())
	}
	result := make([]map[string]interface{}, 0)
	for _, myApp := range appList {
		attrsPublic := extractAppAttrsPublic(myApp)
		appInfo := map[string]interface{}{"id": myApp.GetId(), "config": attrsPublic}
		result = append(result, appInfo)
	}
	return itineris.NewApiResult(itineris.StatusOk).SetData(result)
}

/*
API handler "getMyApp"

Notes:
	- This API return only app's public info
*/
func apiGetMyApp(ctx *itineris.ApiContext, _ *itineris.ApiAuth, params *itineris.ApiParams) *itineris.ApiResult {
	id, _ := params.GetParamAsType("id", reddo.TypeString)
	if id == nil || strings.TrimSpace(id.(string)) == "" {
		return itineris.NewApiResult(itineris.StatusNotFound).SetMessage(fmt.Sprintf("App [%s] not found", id))
	}
	if myApp, err := appDao.Get(id.(string)); err != nil {
		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(err.Error())
	} else if myApp == nil {
		return itineris.NewApiResult(itineris.StatusNotFound).SetMessage(fmt.Sprintf("App [%s] not found", id))
	} else {
		sessionClaim, ok := ctx.GetContextValue(ctxFieldSession).(*SessionClaims)
		if !ok || sessionClaim == nil {
			return itineris.NewApiResult(itineris.StatusNoPermission).SetMessage("Cannot obtain current logged in user info")
		}
		if myApp.GetOwnerId() != sessionClaim.UserId {
			// purposely return "not found" error
			return itineris.NewApiResult(itineris.StatusNotFound).SetMessage(fmt.Sprintf("App [%s] not found", id))
		}
		attrsPublic := extractAppAttrsPublic(myApp)
		return itineris.NewApiResult(itineris.StatusOk).SetData(map[string]interface{}{"id": myApp.GetId(), "config": attrsPublic})
	}
}

/*
API handler "getApp"

Notes:
	- This API return only app's public info
*/
func apiGetApp(_ *itineris.ApiContext, _ *itineris.ApiAuth, params *itineris.ApiParams) *itineris.ApiResult {
	id, _ := params.GetParamAsType("id", reddo.TypeString)
	if id == nil || strings.TrimSpace(id.(string)) == "" {
		return itineris.NewApiResult(itineris.StatusNotFound).SetMessage(fmt.Sprintf("App [%s] not found", id))
	}
	if myApp, err := appDao.Get(id.(string)); err != nil {
		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(err.Error())
	} else if myApp == nil {
		return itineris.NewApiResult(itineris.StatusNotFound).SetMessage(fmt.Sprintf("App [%s] not found", id))
	} else {
		attrsPublic := extractAppAttrsPublic(myApp)
		return itineris.NewApiResult(itineris.StatusOk).SetData(map[string]interface{}{"id": myApp.GetId(), "config": attrsPublic})
	}
}

func _extractParam(params *itineris.ApiParams, paramName string, typ reflect.Type, defValue interface{}, regexp *regexp.Regexp) interface{} {
	v, _ := params.GetParamAsType(paramName, typ)
	if v == nil {
		v = defValue
	}
	if v != nil {
		if _, ok := v.(string); ok {
			v = strings.TrimSpace(v.(string))
			if regexp != nil && !regexp.Match([]byte(v.(string))) {
				return nil
			}
		}
	}
	return v
}

func _extractAppParams(ctx *itineris.ApiContext, params *itineris.ApiParams) (*app.App, *itineris.ApiResult) {
	id := _extractParam(params, "id", reddo.TypeString, nil, regexp.MustCompile("^[0-9A-Za-z_]+$"))
	if id == nil {
		return nil, itineris.NewApiResult(itineris.StatusErrorClient).SetMessage("Missing or invalid value for parameter [id]")
	} else {
		id = strings.ToLower(id.(string))
	}
	isActive := _extractParam(params, "is_active", reddo.TypeBool, false, nil)
	desc := _extractParam(params, "description", reddo.TypeString, "", nil)
	defaultReturnUrl := _extractParam(params, "default_return_url", reddo.TypeString, "", nil)
	if defaultReturnUrl != "" && !regexp.MustCompile("^(?i)https?://.*$").Match([]byte(defaultReturnUrl.(string))) {
		return nil, itineris.NewApiResult(itineris.StatusErrorClient).SetMessage("Invalid value for parameter [default_return_url]")
	}
	tagsStr := _extractParam(params, "tags", reddo.TypeString, "", nil)
	tags := regexp.MustCompile("[,;]+").Split(tagsStr.(string), -1)
	for i, tag := range tags {
		tags[i] = strings.TrimSpace(tag)
	}
	idSources := _extractParam(params, "id_sources", reflect.TypeOf(map[string]bool{}), make(map[string]bool), nil)
	rsaPubicKeyPem := _extractParam(params, "rsa_public_key", reddo.TypeString, "", nil)
	if rsaPubicKeyPem != "" {
		_, err := parseRsaPublicKeyFromPem(rsaPubicKeyPem.(string))
		if err != nil {
			return nil, itineris.NewApiResult(itineris.StatusErrorClient).SetMessage(err.Error())
		}
	}

	sessionClaim, ok := ctx.GetContextValue(ctxFieldSession).(*SessionClaims)
	if !ok || sessionClaim == nil {
		return nil, itineris.NewApiResult(itineris.StatusNoPermission).SetMessage("Cannot obtain current logged in user info")
	}
	ownerId := sessionClaim.UserId

	boApp := app.NewApp(goapi.AppVersionNumber, id.(string), ownerId, desc.(string))
	boApp.SetAttrsPublic(app.AppAttrsPublic{
		IsActive:         isActive.(bool),
		Description:      desc.(string),
		DefaultReturnUrl: defaultReturnUrl.(string),
		IdentitySources:  idSources.(map[string]bool),
		Tags:             tags,
		RsaPublicKey:     rsaPubicKeyPem.(string),
	})

	return boApp, nil
}

// API handler "registerApp"
func apiRegisterApp(ctx *itineris.ApiContext, _ *itineris.ApiAuth, params *itineris.ApiParams) *itineris.ApiResult {
	newApp, apiResult := _extractAppParams(ctx, params)
	if apiResult != nil {
		return apiResult
	}

	if existingApp, err := appDao.Get(newApp.GetId()); err != nil {
		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(err.Error())
	} else if existingApp != nil {
		return itineris.NewApiResult(itineris.StatusErrorClient).SetMessage(fmt.Sprintf("App [%s] already exist", newApp.GetId()))
	}

	if ok, err := appDao.Create(newApp); err != nil {
		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(err.Error())
	} else if !ok {
		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(fmt.Sprintf("Unknown error while registering app [%s]", newApp.GetId()))
	}
	return itineris.NewApiResult(itineris.StatusOk).SetMessage(fmt.Sprintf("App [%s] has been registered successfully", newApp.GetId()))
}

// API handler "updateMyApp"
func apiUpdateMyApp(ctx *itineris.ApiContext, _ *itineris.ApiAuth, params *itineris.ApiParams) *itineris.ApiResult {
	submitApp, apiResult := _extractAppParams(ctx, params)
	if apiResult != nil {
		return apiResult
	}

	if existingApp, err := appDao.Get(submitApp.GetId()); err != nil {
		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(err.Error())
	} else if existingApp == nil {
		return itineris.NewApiResult(itineris.StatusErrorClient).SetMessage(fmt.Sprintf("App [%s] does not exist", submitApp.GetId()))
	} else if existingApp.GetOwnerId() != submitApp.GetOwnerId() {
		return itineris.NewApiResult(itineris.StatusNoPermission).SetMessage(fmt.Sprintf("App [%s] does not belong to user", submitApp.GetId()))
	}

	if ok, err := appDao.Update(submitApp); err != nil {
		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(err.Error())
	} else if !ok {
		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(fmt.Sprintf("Unknown error while updating app [%s]", submitApp.GetId()))
	}
	return itineris.NewApiResult(itineris.StatusOk).SetMessage(fmt.Sprintf("App [%s] has been updated successfully", submitApp.GetId()))
}

// API handler "deleteMyApp"
func apiDeleteMyApp(ctx *itineris.ApiContext, _ *itineris.ApiAuth, params *itineris.ApiParams) *itineris.ApiResult {
	submitApp, apiResult := _extractAppParams(ctx, params)
	if apiResult != nil {
		return apiResult
	}

	if existingApp, err := appDao.Get(submitApp.GetId()); err != nil {
		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(err.Error())
	} else if existingApp == nil {
		return itineris.NewApiResult(itineris.StatusErrorClient).SetMessage(fmt.Sprintf("App [%s] does not exist", submitApp.GetId()))
	} else if existingApp.GetOwnerId() != submitApp.GetOwnerId() {
		return itineris.NewApiResult(itineris.StatusNoPermission).SetMessage(fmt.Sprintf("App [%s] does not belong to user", submitApp.GetId()))
	}

	if submitApp.GetId() == systemAppId {
		return itineris.NewApiResult(itineris.StatusNoPermission).SetMessage(fmt.Sprintf("App [%s] can not be deleted", submitApp.GetId()))
	}

	if ok, err := appDao.Delete(submitApp); err != nil {
		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(err.Error())
	} else if !ok {
		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(fmt.Sprintf("Unknown error while deleting app [%s]", submitApp.GetId()))
	}
	return itineris.NewApiResult(itineris.StatusOk).SetMessage(fmt.Sprintf("App [%s] has been deleted successfully", submitApp.GetId()))
}
