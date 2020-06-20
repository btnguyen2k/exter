package gvabe

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/btnguyen2k/consu/reddo"

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
	router.SetHandler("checkLoginToken", apiCheckLoginToken)
	router.SetHandler("systemInfo", apiSystemInfo)

	router.SetHandler("myAppList", apiMyAppList)
	router.SetHandler("getApp", apiGetApp)

	router.SetHandler("groupList", apiGroupList)
	router.SetHandler("getGroup", apiGetGroup)
	router.SetHandler("createGroup", apiCreateGroup)
	router.SetHandler("deleteGroup", apiDeleteGroup)
	router.SetHandler("updateGroup", apiUpdateGroup)

	router.SetHandler("userList", apiUserList)
	router.SetHandler("getUser", apiGetUser)
	router.SetHandler("createUser", apiCreateUser)
	router.SetHandler("deleteUser", apiDeleteUser)
	router.SetHandler("updateUser", apiUpdateUser)
}

/*----------------------------------------------------------------------*/

func _parseLoginTokenFromApi(_token interface{}) (*itineris.ApiResult, *SessionClaim, *user.User) {
	stoken, ok := _token.(string)
	if !ok || stoken == "" {
		return itineris.NewApiResult(itineris.StatusNoPermission).SetMessage("empty token"), nil, nil
	}

	var claim *SessionClaim
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

/*
API handler "myAppList"
This API expects an input map:

	{
		"token": login token (returned by apiLogin/apiCheckLoginToken),
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
	for _, app := range appList {
		config := make(map[string]interface{})
		if app.config != nil {
			js, _ := json.Marshal(app.config)
			json.Unmarshal(js, &config)
			loginChannels := make(map[string]bool)
			for s, _ := range app.config.IdentitySources {
				if app.config.IdentitySources[s] && enabledLoginChannels[s] {
					loginChannels[s] = true
				}
			}
			config["sources"] = loginChannels
		}
		result = append(result, config)
	}
	return itineris.NewApiResult(itineris.StatusOk).SetData(result)
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
	if app, err := appDao.Get(id.(string)); err != nil {
		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(err.Error())
	} else if app == nil {
		return itineris.NewApiResult(itineris.StatusNotFound).SetMessage(fmt.Sprintf("App [%s] not found", id))
	} else {
		config := make(map[string]interface{})
		if app.config != nil {
			js, _ := json.Marshal(app.config)
			json.Unmarshal(js, &config)
			loginChannels := make(map[string]bool)
			for s, _ := range app.config.IdentitySources {
				if app.config.IdentitySources[s] && enabledLoginChannels[s] {
					loginChannels[s] = true
				}
			}
			config["sources"] = loginChannels
		}
		return itineris.NewApiResult(itineris.StatusOk).SetData(map[string]interface{}{"id": app.id, "config": config})
	}
}
