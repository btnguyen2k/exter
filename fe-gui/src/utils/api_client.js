/*
Client to make call to API server using Axios.

@author Thanh Nguyen <btnguyen2k@gmail.com>
@since template-v0.1.0
*/
import Axios from "axios"
import appConfig from "./app_config"
import utils from "@/utils/app_utils"

const apiClient = Axios.create({
    baseURL: appConfig.APP_CONFIG.api_client.bo_api_base_url,
    timeout: 30000,
});

const headerAppId = appConfig.APP_CONFIG.api_client.header_app_id
const headerAccessToken = appConfig.APP_CONFIG.api_client.header_access_token
let appId = appConfig.APP_CONFIG.api_client.app_id + ":" + Math.random()

let apiInfo = "/info"
let apiLogin = "/api/login"
let apiCheckLoginToken = "/api/checkLoginToken"
let apiSystemInfo = "/api/systemInfo"
let apiApp = "/api/app"
let apiMyAppList = "/api/myapps"
let apiGroupList = "/api/groups"
let apiGroup = "/api/group"
let apiUserList = "/api/users"
let apiUser = "/api/user"

function _apiOnSuccess(resp, apiUri, callbackSuccessful) {
    if (apiUri != apiLogin && apiUri != apiCheckLoginToken && resp.hasOwnProperty("data") && resp.data.status == 403) {
        console.error("Error 403 from API [" + apiUri + "], redirecting to login page...")
        router.push({name: "Login", query: {app: appConfig.APP_NAME, returnUrl: router.currentRoute.fullPath}})
        return
    }
    if (resp.hasOwnProperty("data") && resp.data.hasOwnProperty("extras") && resp.data.extras.hasOwnProperty("_access_token_")) {
        console.log("Update new access token from API [" + apiUri + "]")
        let jwt = utils.parseJwt(resp.data.extras._access_token_)
        utils.saveLoginSession({uid: jwt.payloadObj.uid, token: resp.data.extras._access_token_})
    }
    if (callbackSuccessful != null) {
        callbackSuccessful(resp.data)
    }
}

function _apiOnError(err, apiUri, callbackError) {
    console.error("Error calling api [" + apiUri + "]: " + err)
    if (callbackError != null) {
        callbackError(err)
    }
}

function apiDoGet(apiUri, callbackSuccessful, callbackError) {
    let session = utils.loadLoginSession()
    const headers = {}
    headers[headerAppId] = appId
    headers[headerAccessToken] = session != null ? session.token : ""
    return apiClient.get(apiUri, {
        headers: headers
    }).then(res => _apiOnSuccess(res, apiUri, callbackSuccessful)).catch(err => _apiOnError(err, apiUri, callbackError))
}

function apiDoPost(apiUri, data, callbackSuccessful, callbackError) {
    let session = utils.loadLoginSession()
    const headers = {}
    headers[headerAppId] = appId
    headers[headerAccessToken] = session != null ? session.token : ""
    apiClient.post(apiUri, data, {
        headers: headers
    }).then(res => _apiOnSuccess(res, apiUri, callbackSuccessful)).catch(err => _apiOnError(err, apiUri, callbackError))
}

function apiDoPut(apiUri, data, callbackSuccessful, callbackError) {
    let session = utils.loadLoginSession()
    const headers = {}
    headers[headerAppId] = appId
    headers[headerAccessToken] = session != null ? session.token : ""
    apiClient.put(apiUri, data, {
        headers: headers
    }).then(res => _apiOnSuccess(res, apiUri, callbackSuccessful)).catch(err => _apiOnError(err, apiUri, callbackError))
}

function apiDoDelete(apiUri, callbackSuccessful, callbackError) {
    let session = utils.loadLoginSession()
    const headers = {}
    headers[headerAppId] = appId
    headers[headerAccessToken] = session != null ? session.token : ""
    apiClient.delete(apiUri, {
        headers: headers
    }).then(res => _apiOnSuccess(res, apiUri, callbackSuccessful)).catch(err => _apiOnError(err, apiUri, callbackError))
}

export default {
    apiInfo,
    apiLogin,
    apiApp,
    apiCheckLoginToken,
    apiSystemInfo,
    apiMyAppList,
    apiGroupList,
    apiGroup,
    apiUserList,
    apiUser,

    apiDoGet,
    apiDoPost,
    apiDoPut,
    apiDoDelete,
}
