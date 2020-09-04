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
	"os"
	"regexp"
	"strings"

	"golang.org/x/oauth2/google"

	"main/src/goapi"
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
	if os.Getenv("DEBUG") != "" {
		DEBUG = true
	}
	go startUpdateSystemInfo()

	initRsaKeys()
	initLoginChannels()
	initExterHomeUrl()
	initFacebookAppSecret()
	initGithubClientSecret()
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

// available since v0.3.0
func initExterHomeUrl() {
	exterHomeUrl = strings.TrimSpace(goapi.AppConfig.GetString("gvabe.exter_home_url"))
	if DEBUG {
		log.Printf("[DEBUG] Exter home url: %s", exterHomeUrl)
	}
	if exterHomeUrl == "" {
		panic("no valid Exter home-url defined at [gvabe.exter_home_url]")
	}
}

// available since v0.3.0
func initFacebookAppSecret() {
	if !enabledLoginChannels[loginChannelFacebook] {
		return
	}
	appId := strings.TrimSpace(goapi.AppConfig.GetString("gvabe.channels.facebook.app_id"))
	if appId == "" {
		log.Println("[ERROR] No valid Facebook app-id defined at [gvabe.channels.facebook.app_id]")
	}
	appSecret := strings.TrimSpace(goapi.AppConfig.GetString("gvabe.channels.facebook.app_secret"))
	if appSecret == "" {
		log.Println("[ERROR] No valid Facebook app-secret defined at [gvabe.channels.facebook.app_secret]")
	}
	fbOAuthConf.ClientID = appId
	fbOAuthConf.ClientSecret = appSecret
	fbOAuthConf.RedirectURL = exterHomeUrl
	initFacebookApp(appId, appSecret)
	if DEBUG && appId != "" && appSecret != "" {
		log.Printf("[DEBUG] initFacebookAppSecret: %s/%s", appId, "***"+appSecret[len(appSecret)-4:])
	}
}

// available since v0.2.0
func initGithubClientSecret() {
	if !enabledLoginChannels[loginChannelGithub] {
		return
	}
	clientId := strings.TrimSpace(goapi.AppConfig.GetString("gvabe.channels.github.client_id"))
	if clientId == "" {
		log.Println("[ERROR] No valid Github OAuth app client-id defined at [gvabe.channels.github.client_id]")
	}
	clientSecret := strings.TrimSpace(goapi.AppConfig.GetString("gvabe.channels.github.client_secret"))
	if clientSecret == "" {
		log.Println("[ERROR] No valid Github OAuth app client-secret defined at [gvabe.channels.github.client_secret]")
	}
	githubOAuthConf.ClientID = clientId
	githubOAuthConf.ClientSecret = clientSecret
	// githubOAuthConf.RedirectURL = exterHomeUrl //[btnguyen2k-20200904]: do NOT set RedirectURL, or else we encounter error "oauth2: server response missing access_token"
	if DEBUG && clientId != "" && clientSecret != "" {
		log.Printf("[DEBUG] initGithubClientSecret: %s/%s", clientId, "***"+clientSecret[len(clientSecret)-4:])
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
			log.Println("[ERROR] No valid GoogleAPI project-id defined at [gvabe.channels.google.project_id]")
		}
		clientId := strings.TrimSpace(goapi.AppConfig.GetString("gvabe.channels.google.client_id"))
		if clientId == "" {
			log.Println("[ERROR] No valid GoogleAPI client-id defined at [gvabe.channels.google.client_id]")
		}
		clientSecret := strings.TrimSpace(goapi.AppConfig.GetString("gvabe.channels.google.client_secret"))
		if clientSecret == "" {
			log.Println("[ERROR] No valid GoogleAPI client-secret defined at [gvabe.channels.google.client_secret]")
		}
		// appDomainsStr := strings.TrimSpace(goapi.AppConfig.GetString("gvabe.channels.google.app_domains"))
		// if appDomainsStr == "" {
		// 	log.Println("[ERROR] No valid GoogleAPI app-domains defined at [gvabe.channels.google.app_domains]")
		// }
		// appDomains := make([]string, 0)
		// for _, s := range regexp.MustCompile("[,; ]+").Split(appDomainsStr, -1) {
		// 	if strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://") {
		// 		appDomains = append(appDomains, s)
		// 	} else {
		// 		appDomains = append(append(appDomains, "http://"+s), "https://"+s)
		// 	}
		// }
		// appDomainsJs, _ := json.Marshal(appDomains)
		appDomainsJs, _ := json.Marshal([]string{exterHomeUrl})

		clientSecretJson = fmt.Sprintf(`{
		  "type":"authorized_user",
		  "web": {
			"project_id": "%s",
			"client_id": "%s",
			"client_secret": "%s",
			"auth_uri": "https://accounts.google.com/o/oauth2/auth",
			"token_uri": "https://oauth2.googleapis.com/token",
			"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
			"redirect_uris": %s,
			"javascript_origins": %s,
			"access_type": "offline"
		  }
		}`, projectId, clientId, clientSecret, appDomainsJs, appDomainsJs)
	}
	if DEBUG {
		r := regexp.MustCompile(`(?s)"client_secret":\s*"(.*?)"`)
		f := r.FindSubmatch([]byte(clientSecretJson))
		if len(f) > 1 {
			clientSecret := string(f[1])
			clientSecret = "***" + clientSecret[len(clientSecret)-4:]
			clientSecretJson = r.ReplaceAllString(clientSecretJson, `"client_secret": "`+clientSecret+`"`)
			log.Printf("[DEBUG] initGoogleClientSecret: %s", clientSecretJson)
		}
	}
	var err error
	if googleOAuthConf, err = google.ConfigFromJSON([]byte(clientSecretJson)); err != nil {
		panic(err)
	}
}
