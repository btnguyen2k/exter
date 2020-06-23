/*
Package gvabe provides backend API for GoVueAdmin Frontend.

@author Thanh Nguyen <btnguyen2k@gmail.com>
@since template-v0.1.0
*/
package gvabe

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"

	"github.com/btnguyen2k/consu/semita"
	"golang.org/x/oauth2/google"

	"main/src/goapi"
	"main/src/itineris"
	"main/src/mico"
)

type MyBootstrapper struct {
	name string
}

var Bootstrapper = &MyBootstrapper{name: "gvabe"}

/*
Bootstrap implements goapi.IBootstrapper.Bootstrap

Bootstrapper usually does:
- register api-handlers with the global ApiRouter
- other initializing work (e.g. creating DAO, initializing database, etc)
*/
func (b *MyBootstrapper) Bootstrap() error {
	go startUpdateSystemInfo()

	initRsaKeys()
	initLoginChannels()
	initGoogleClientSecret()
	initCaches()
	initDaos()
	initApiHandlers(goapi.ApiRouter)
	initApiFilters(goapi.ApiRouter)
	return nil
}

func initRsaKeys() {
	rsaPrivKeyFile := goapi.AppConfig.GetString("gvabe.keys.rsa_privkey_file")
	if rsaPrivKeyFile == "" {
		log.Println("WARN: no RSA private key file configured at [gvabe.keys.rsa_privkey_file], generating one...")
		privKey, err := genRsaKey(2048)
		if err != nil {
			panic(err)
		}
		rsaPrivKey = privKey
	} else {
		log.Println(fmt.Sprintf("INFO: loading RSA private key from [%s]...", rsaPrivKeyFile))
		content, err := ioutil.ReadFile(rsaPrivKeyFile)
		if err != nil {
			panic(err)
		}
		block, _ := pem.Decode(content)
		if block == nil {
			panic(fmt.Sprintf("cannot decode PEM from file [%s]", rsaPrivKeyFile))
		}
		var der []byte
		passphrase := goapi.AppConfig.GetString("gvabe.keys.rsa_privkey_passphrase")
		if passphrase != "" {
			log.Println("INFO: RSA private key is pass-phrase protected")
			if decrypted, err := x509.DecryptPEMBlock(block, []byte(passphrase)); err != nil {
				panic(err)
			} else {
				der = decrypted
			}
		} else {
			der = block.Bytes
		}
		if block.Type == "RSA PRIVATE KEY" {
			if privKey, err := x509.ParsePKCS1PrivateKey(der); err != nil {
				panic(err)
			} else {
				rsaPrivKey = privKey
			}
		} else if block.Type == "PRIVATE KEY" {
			if privKey, err := x509.ParsePKCS8PrivateKey(der); err != nil {
				panic(err)
			} else {
				rsaPrivKey = privKey.(*rsa.PrivateKey)
			}
		}
	}

	rsaPubKey = &rsaPrivKey.PublicKey
	pubDER := x509.MarshalPKCS1PublicKey(rsaPubKey)
	pubBlock := pem.Block{
		Type:    "RSA PUBLIC KEY",
		Headers: nil,
		Bytes:   pubDER,
	}
	publicPEM := pem.EncodeToMemory(&pubBlock)
	log.Println(string(publicPEM))
}

func initCaches() {
	cacheConfig := &mico.CacheConfig{}
	sessionCache = mico.NewMemoryCache(cacheConfig)
	preLoginSessionCache = mico.NewMemoryCache(cacheConfig)
}

func initLoginChannels() {
	loginChannels := regexp.MustCompile("[,;\\s]+").Split(goapi.AppConfig.GetString("gvabe.login_channels"), -1)
	for _, channel := range loginChannels {
		channel = strings.TrimSpace(strings.ToLower(channel))
		enabledLoginChannels[channel] = true
	}
}

func initGoogleClientSecret() {
	if !enabledLoginChannels[loginChannelGoogle] {
		return
	}
	clientSecretJson := strings.TrimSpace(goapi.AppConfig.GetString("gvabe.channels.google.client_secret_json"))
	if clientSecretJson == "" {
		log.Println("[INFO] No valid GoogleAPI client secret defined at [gvabe.channels.google.client_secret_json], falling back to {project_id, client_id, client_secret}")

		projectId := strings.TrimSpace(goapi.AppConfig.GetString("gvabe.channels.google.project_id"))
		if projectId == "" {
			log.Println("[ERROR] No valid GoogleAPI project id defined at [gvabe.channels.google.project_id]")
		}
		clientId := strings.TrimSpace(goapi.AppConfig.GetString("gvabe.channels.google.client_id"))
		if clientId == "" {
			log.Println("[ERROR] No valid GoogleAPI client id defined at [gvabe.channels.google.client_id]")
		}
		clientSecret := strings.TrimSpace(goapi.AppConfig.GetString("gvabe.channels.google.client_secret"))
		if clientSecret == "" {
			log.Println("[ERROR] No valid GoogleAPI client id defined at [gvabe.channels.google.client_secret]")
		}
		clientSecretJson = fmt.Sprintf(`{
		  "type":"authorized_user",
		  "web": {
			"project_id": "%s",
			"client_id": "%s",
			"client_secret": "%s",
			"auth_uri": "https://accounts.google.com/o/oauth2/auth",
			"token_uri": "https://oauth2.googleapis.com/token",
			"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
			"redirect_uris": ["http://localhost:8080"],
			"javascript_origins": ["http://localhost:8080"],
			"access_type": "offline"
		  }
		}`, projectId, clientId, clientSecret)
	}
	gClientSecretJson = []byte(clientSecretJson)
	var err error
	if gConfig, err = google.ConfigFromJSON([]byte(clientSecretJson)); err != nil {
		panic(err)
	}
	if err = json.Unmarshal([]byte(clientSecretJson), &gClientSecretData); err != nil {
		panic(err)
	}
	sGoogleClientSecret = semita.NewSemita(gClientSecretData)
}

/*----------------------------------------------------------------------*/

// API handler "userList"
func apiUserList(_ *itineris.ApiContext, _ *itineris.ApiAuth, _ *itineris.ApiParams) *itineris.ApiResult {
	return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage("not implemented")
	// userList, err := userDao.GetAll()
	// if err != nil {
	// 	return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(err.Error())
	// }
	// data := make([]map[string]interface{}, 0)
	// for _, u := range userList {
	// 	data = append(data, map[string]interface{}{
	// 		"username": u.GetUsername(), "name": u.GetName(), "gid": u.GetGroupId(),
	// 	})
	// }
	// return itineris.NewApiResult(itineris.StatusOk).SetData(data)
}

// API handler "getUser"
func apiGetUser(_ *itineris.ApiContext, _ *itineris.ApiAuth, params *itineris.ApiParams) *itineris.ApiResult {
	return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage("not implemented")
	// username, _ := params.GetParamAsType("username", reddo.TypeString)
	// if username == nil || strings.TrimSpace(username.(string)) == "" {
	// 	return itineris.NewApiResult(itineris.StatusNotFound).SetMessage(fmt.Sprintf("User [%s] not found", username))
	// }
	// if user, err := userDao.Get(username.(string)); err != nil {
	// 	return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(err.Error())
	// } else if user == nil {
	// 	return itineris.NewApiResult(itineris.StatusNotFound).SetMessage(fmt.Sprintf("User [%s] not found", username))
	// } else {
	// 	return itineris.NewApiResult(itineris.StatusOk).SetData(map[string]interface{}{
	// 		"username": user.GetUsername(), "name": user.GetName(), "gid": user.GetGroupId(),
	// 	})
	// }
}

// API handler "updateUser"
func apiUpdateUser(_ *itineris.ApiContext, _ *itineris.ApiAuth, params *itineris.ApiParams) *itineris.ApiResult {
	return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage("not implemented")
	// username, _ := params.GetParamAsType("username", reddo.TypeString)
	// if username == nil || strings.TrimSpace(username.(string)) == "" {
	// 	return itineris.NewApiResult(itineris.StatusNotFound).SetMessage(fmt.Sprintf("User [%s] not found", username))
	// }
	// if user, err := userDao.Get(username.(string)); err != nil {
	// 	return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(err.Error())
	// } else if user == nil {
	// 	return itineris.NewApiResult(itineris.StatusNotFound).SetMessage(fmt.Sprintf("User [%s] not found", username))
	// } else {
	// 	// TODO check current user's permission
	//
	// 	// FIXME this is for demo purpose only!
	// 	if user.GetUsername() == systemAdminUsername {
	// 		return itineris.NewApiResult(itineris.StatusNoPermission).SetMessage(fmt.Sprintf("Cannot edit system admin user [%s]", user.GetUsername()))
	// 	}
	//
	// 	password, _ := params.GetParamAsType("password", reddo.TypeString)
	// 	var newPassword, newPassword2 interface{}
	// 	if password != nil && strings.TrimSpace(password.(string)) != "" {
	// 		password = strings.TrimSpace(password.(string))
	// 		if encryptPassword(user.GetUsername(), password.(string)) != user.GetPassword() {
	// 			return itineris.NewApiResult(itineris.StatusErrorClient).SetMessage("Current password does not match")
	// 		}
	//
	// 		newPassword, _ = params.GetParamAsType("new_password", reddo.TypeString)
	// 		if newPassword == nil || strings.TrimSpace(newPassword.(string)) == "" {
	// 			return itineris.NewApiResult(itineris.StatusErrorClient).SetMessage("Missing or empty parameter [new_password]")
	// 		}
	// 		newPassword = strings.TrimSpace(newPassword.(string))
	// 		newPassword2, _ = params.GetParamAsType("new_password2", reddo.TypeString)
	// 		if newPassword2 == nil || strings.TrimSpace(newPassword2.(string)) != newPassword {
	// 			return itineris.NewApiResult(itineris.StatusErrorClient).SetMessage("New password does not match confirmed one")
	// 		}
	// 	}
	//
	// 	name, _ := params.GetParamAsType("name", reddo.TypeString)
	// 	if name == nil || strings.TrimSpace(name.(string)) == "" {
	// 		return itineris.NewApiResult(itineris.StatusErrorClient).SetMessage("Missing or empty parameter [name]")
	// 	}
	// 	name = strings.TrimSpace(name.(string))
	//
	// 	groupId, _ := params.GetParamAsType("group_id", reddo.TypeString)
	// 	if groupId == nil || strings.TrimSpace(groupId.(string)) == "" {
	// 		return itineris.NewApiResult(itineris.StatusErrorClient).SetMessage("Missing or empty parameter [group_id]")
	// 	}
	// 	groupId = strings.TrimSpace(strings.ToLower(groupId.(string)))
	// 	if group, err := groupDao.Get(groupId.(string)); err != nil {
	// 		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(err.Error())
	// 	} else if group == nil {
	// 		return itineris.NewApiResult(itineris.StatusErrorClient).SetMessage(fmt.Sprintf("Group [%s] does not exist", groupId))
	// 	}
	//
	// 	user.SetName(strings.TrimSpace(name.(string))).
	// 		SetGroupId(groupId.(string))
	// 	if password != nil && strings.TrimSpace(password.(string)) != "" {
	// 		user.SetPassword(encryptPassword(user.GetUsername(), newPassword.(string)))
	// 	}
	//
	// 	if ok, err := userDao.Update(user); err != nil {
	// 		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(err.Error())
	// 	} else if !ok {
	// 		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(fmt.Sprintf("User [%s] has not been updated", username.(string)))
	// 	}
	// 	return itineris.NewApiResult(itineris.StatusOk)
	// }
}

// API handler "deleteUser"
func apiDeleteUser(_ *itineris.ApiContext, _ *itineris.ApiAuth, params *itineris.ApiParams) *itineris.ApiResult {
	return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage("not implemented")
	// username, _ := params.GetParamAsType("username", reddo.TypeString)
	// if username == nil || strings.TrimSpace(username.(string)) == "" {
	// 	return itineris.NewApiResult(itineris.StatusNotFound).SetMessage(fmt.Sprintf("User [%s] not found", username))
	// }
	// if user, err := userDao.Get(username.(string)); err != nil {
	// 	return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(err.Error())
	// } else if user == nil {
	// 	return itineris.NewApiResult(itineris.StatusNotFound).SetMessage(fmt.Sprintf("User [%s] not found", username))
	// } else {
	// 	// TODO check current user's permission
	//
	// 	// FIXME this is for demo purpose only!
	// 	if user.GetUsername() == systemAdminUsername {
	// 		return itineris.NewApiResult(itineris.StatusNoPermission).SetMessage(fmt.Sprintf("Cannot delete system admin user [%s]", user.GetUsername()))
	// 	}
	//
	// 	if ok, err := userDao.Delete(user); err != nil {
	// 		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(err.Error())
	// 	} else if !ok {
	// 		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(fmt.Sprintf("User [%s] has not been deleted", username.(string)))
	// 	}
	// 	return itineris.NewApiResult(itineris.StatusOk)
	// }
}

// API handler "createUser"
func apiCreateUser(_ *itineris.ApiContext, _ *itineris.ApiAuth, params *itineris.ApiParams) *itineris.ApiResult {
	return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage("not implemented")
	// username, _ := params.GetParamAsType("username", reddo.TypeString)
	// if username == nil || strings.TrimSpace(username.(string)) == "" {
	// 	return itineris.NewApiResult(itineris.StatusErrorClient).SetMessage("Missing or empty parameter [username]")
	// }
	// username = strings.TrimSpace(strings.ToLower(username.(string)))
	// if !regexp.MustCompile("^[0-9a-z_]+$").Match([]byte(username.(string))) {
	// 	return itineris.NewApiResult(itineris.StatusErrorClient).SetMessage("Invalid value for parameter [username]")
	// }
	//
	// password, _ := params.GetParamAsType("password", reddo.TypeString)
	// if password == nil || strings.TrimSpace(password.(string)) == "" {
	// 	return itineris.NewApiResult(itineris.StatusErrorClient).SetMessage("Missing or empty parameter [password]")
	// }
	// password = strings.TrimSpace(password.(string))
	// password2, _ := params.GetParamAsType("password2", reddo.TypeString)
	// if password2 == nil || strings.TrimSpace(password2.(string)) != password {
	// 	return itineris.NewApiResult(itineris.StatusErrorClient).SetMessage("Password does not match confirmed one")
	// }
	//
	// name, _ := params.GetParamAsType("name", reddo.TypeString)
	// if name == nil || strings.TrimSpace(name.(string)) == "" {
	// 	return itineris.NewApiResult(itineris.StatusErrorClient).SetMessage("Missing or empty parameter [name]")
	// }
	// name = strings.TrimSpace(name.(string))
	//
	// groupId, _ := params.GetParamAsType("group_id", reddo.TypeString)
	// if groupId == nil || strings.TrimSpace(groupId.(string)) == "" {
	// 	return itineris.NewApiResult(itineris.StatusErrorClient).SetMessage("Missing or empty parameter [group_id]")
	// }
	// groupId = strings.TrimSpace(strings.ToLower(groupId.(string)))
	// if group, err := groupDao.Get(groupId.(string)); err != nil {
	// 	return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(err.Error())
	// } else if group == nil {
	// 	return itineris.NewApiResult(itineris.StatusErrorClient).SetMessage(fmt.Sprintf("Group [%s] does not exist", groupId))
	// }
	//
	// if user, err := userDao.Get(username.(string)); err != nil {
	// 	return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(err.Error())
	// } else if user != nil {
	// 	return itineris.NewApiResult(itineris.StatusErrorClient).SetMessage(fmt.Sprintf("User [%s] already existed", username))
	// }
	// user := user.NewUserBo(username.(string), "").
	// 	SetPassword(encryptPassword(username.(string), password.(string))).
	// 	SetName(name.(string)).
	// 	SetGroupId(groupId.(string)).
	// 	SetAesKey(utils.RandomString(16))
	// if ok, err := userDao.Create(user); err != nil {
	// 	return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(err.Error())
	// } else if !ok {
	// 	return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(fmt.Sprintf("User [%s] has not been created", username))
	// }
	// return itineris.NewApiResult(itineris.StatusOk)
}
