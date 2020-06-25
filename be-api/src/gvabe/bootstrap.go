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
