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
	initLinkedinClientSecret()
	initCaches()
	initDaos()
	initApiHandlers(goapi.ApiRouter)
	initApiFilters(goapi.ApiRouter)
	return nil
}

func initRsaKeys() {
	confKeyRsaPrivKeyFile := "gvabe.keys.rsa_privkey_file"
	confKeyRsaPrivKeyPass := "gvabe.keys.rsa_privkey_passphrase"
	rsaPrivKeyFile := goapi.AppConfig.GetString(confKeyRsaPrivKeyFile)
	if rsaPrivKeyFile == "" {
		log.Println(fmt.Sprintf("[WARN] No RSA private key file configured at [%s], generating one...", confKeyRsaPrivKeyFile))
		privKey, err := genRsaKey(2048)
		if err != nil {
			panic(err)
		}
		rsaPrivKey = privKey
	} else {
		log.Println(fmt.Sprintf("[INFO] Loading RSA private key from [%s]...", rsaPrivKeyFile))
		content, err := ioutil.ReadFile(rsaPrivKeyFile)
		if err != nil {
			panic(err)
		}
		block, _ := pem.Decode(content)
		if block == nil {
			panic(fmt.Sprintf("cannot decode PEM from file [%s]", rsaPrivKeyFile))
		}
		var der []byte
		passphrase := goapi.AppConfig.GetString(confKeyRsaPrivKeyPass)
		if passphrase != "" {
			log.Println("[INFO] RSA private key is pass-phrase protected")
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

	pubBlockPKCS1 := pem.Block{
		Type:    "RSA PUBLIC KEY",
		Headers: nil,
		Bytes:   x509.MarshalPKCS1PublicKey(rsaPubKey),
	}
	rsaPubKeyPemPKCS1 = pem.EncodeToMemory(&pubBlockPKCS1)

	pubPKIX, _ := x509.MarshalPKIXPublicKey(rsaPubKey)
	pubBlockPKIX := pem.Block{
		Type:    "PUBLIC KEY",
		Headers: nil,
		Bytes:   pubPKIX,
	}
	rsaPubKeyPemPKIX = pem.EncodeToMemory(&pubBlockPKIX)

	if DEBUG {
		log.Printf("[DEBUG] Exter public key: {Size: %d / Exponent: %d / Modulus: %x}",
			rsaPubKey.Size()*8, rsaPubKey.E, rsaPubKey.N)
		log.Printf("[DEBUG] Exter public key (PKCS1): %s", string(rsaPubKeyPemPKCS1))
		log.Printf("[DEBUG] Exter public key (PKIX): %s", string(rsaPubKeyPemPKIX))
	}
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
			_clientSecretJson := r.ReplaceAllString(clientSecretJson, `"client_secret": "`+clientSecret+`"`)
			log.Printf("[DEBUG] initGoogleClientSecret: %s", _clientSecretJson)
		}
	}
	var err error
	if googleOAuthConf, err = google.ConfigFromJSON([]byte(clientSecretJson)); err != nil {
		panic(err)
	}
}

// available since v0.5.0
func initLinkedinClientSecret() {
	if !enabledLoginChannels[loginChannelLinkedin] {
		return
	}
	clientId := strings.TrimSpace(goapi.AppConfig.GetString("gvabe.channels.linkedin.client_id"))
	if clientId == "" {
		log.Println("[ERROR] No valid LinkedIn OAuth app client-id defined at [gvabe.channels.linkedin.client_id]")
	}
	clientSecret := strings.TrimSpace(goapi.AppConfig.GetString("gvabe.channels.linkedin.client_secret"))
	if clientSecret == "" {
		log.Println("[ERROR] No valid LinkedIn OAuth app client-secret defined at [gvabe.channels.linkedin.client_secret]")
	}
	redirectUri := strings.TrimSpace(goapi.AppConfig.GetString("gvabe.channels.linkedin.redirect_uri"))
	if redirectUri == "" {
		log.Println("[ERROR] No valid LinkedIn OAuth app redirect-uri defined at [gvabe.channels.linkedin.redirect_uri]")
		redirectUri = exterHomeUrl
	}
	linkedinOAuthConf.ClientID = clientId
	linkedinOAuthConf.ClientSecret = clientSecret
	linkedinOAuthConf.RedirectURL = redirectUri
	if DEBUG && clientId != "" && clientSecret != "" {
		log.Printf("[DEBUG] initLinkedinClientSecret: %s/%s/%s", clientId, "***"+clientSecret[len(clientSecret)-4:], redirectUri)
	}
}
