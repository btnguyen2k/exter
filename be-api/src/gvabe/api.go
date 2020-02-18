package gvabe

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/btnguyen2k/consu/reddo"

	"main/src/itineris"
)

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
