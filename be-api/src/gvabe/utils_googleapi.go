package gvabe

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/btnguyen2k/consu/reddo"
	"github.com/btnguyen2k/consu/semita"
	"golang.org/x/oauth2"
)

const (
	urlGoogleServiceTokenInfo = "https://oauth2.googleapis.com/tokeninfo"
)

var (
	gClientSecretJson   []byte
	gClientSecretData   map[string]interface{}
	gConfig             *oauth2.Config
	sGoogleClientSecret *semita.Semita
)

type IdTokenVerified struct {
	Data      map[string]interface{}
	IssuedAt  time.Time
	ExpiredAt time.Time
	s         *semita.Semita
}

type AccessTokenVerified struct {
	Data      map[string]interface{}
	ExpiredAt time.Time
	s         *semita.Semita
}

var (
	cachedIdTokens         = make(map[string]*IdTokenVerified)
	mutexCacheIdTokens     sync.Mutex
	cachedAccessTokens     = make(map[string]*AccessTokenVerified)
	mutexCacheAccessTokens sync.Mutex
)

// parseAndVerifyGoogleIdToken calls Google's tokeninfo service to parse and verify id_token.
func parseAndVerifyGoogleIdToken(idToken string) (*IdTokenVerified, error) {
	mutexCacheIdTokens.Lock()
	defer mutexCacheIdTokens.Unlock()
	sha := sha1.Sum([]byte(idToken))
	shaHash := hex.EncodeToString(sha[:])
	vIdToken := cachedIdTokens[shaHash]
	if vIdToken == nil || time.Now().After(vIdToken.ExpiredAt) {
		data := make(map[string]interface{})
		url := urlGoogleServiceTokenInfo + "?id_token=" + idToken
		if resp, err := httpClient.Get(url); err != nil {
			return nil, err
		} else {
			defer resp.Body.Close()
			if resp.StatusCode != 200 {
				return nil, errors.New("error while validating Google id_token: " + resp.Status)
			}
			if body, err := ioutil.ReadAll(resp.Body); err != nil {
				return nil, err
			} else if err := json.Unmarshal(body, &data); err != nil {
				return nil, errors.New("error while decoding json data from Google tokeninfo service")
			}
		}
		s := semita.NewSemita(data)
		issueTime, _ := s.GetValueOfType("iat", reddo.TypeInt)
		expiryTime, _ := s.GetValueOfType("exp", reddo.TypeInt)
		expiry := time.Now().Add(time.Minute * 15)
		if expiry.Unix() > expiryTime.(int64) {
			expiry = time.Unix(expiryTime.(int64), 0)
		}
		vIdToken = &IdTokenVerified{
			Data:      data,
			IssuedAt:  time.Unix(issueTime.(int64), 0),
			ExpiredAt: expiry,
			s:         s,
		}
		cachedIdTokens[shaHash] = vIdToken
	}
	now := time.Now()
	for h, token := range cachedIdTokens {
		if now.After(token.ExpiredAt) {
			cachedIdTokens[h] = nil
			delete(cachedIdTokens, h)
		}
	}
	return vIdToken, nil
}

// parseAndVerifyGoogleAccessToken calls Google's tokeninfo service to parse and verify access_token.
func parseAndVerifyGoogleAccessToken(accessToken string) (*AccessTokenVerified, error) {
	mutexCacheAccessTokens.Lock()
	defer mutexCacheAccessTokens.Unlock()
	sha := sha1.Sum([]byte(accessToken))
	shaHash := hex.EncodeToString(sha[:])
	vAccessToken := cachedAccessTokens[shaHash]
	if vAccessToken == nil || time.Now().After(vAccessToken.ExpiredAt) {
		data := make(map[string]interface{})
		url := urlGoogleServiceTokenInfo + "?access_token=" + accessToken
		if resp, err := httpClient.Get(url); err != nil {
			return nil, err
		} else {
			defer resp.Body.Close()
			if resp.StatusCode != 200 {
				return nil, errors.New("error while validating Google access_token: " + resp.Status)
			}
			if body, err := ioutil.ReadAll(resp.Body); err != nil {
				return nil, err
			} else if err := json.Unmarshal(body, &data); err != nil {
				return nil, errors.New("error while decoding json data from Google tokeninfo service")
			}
		}
		s := semita.NewSemita(data)
		expiryTime, _ := s.GetValueOfType("exp", reddo.TypeInt)
		expiry := time.Now().Add(time.Minute * 15)
		if expiry.Unix() > expiryTime.(int64) {
			expiry = time.Unix(expiryTime.(int64), 0)
		}
		vAccessToken = &AccessTokenVerified{
			Data:      data,
			ExpiredAt: expiry,
			s:         s,
		}
		cachedAccessTokens[shaHash] = vAccessToken
	}
	now := time.Now()
	for h, token := range cachedAccessTokens {
		if now.After(token.ExpiredAt) {
			cachedAccessTokens[h] = nil
			delete(cachedAccessTokens, h)
		}
	}
	return vAccessToken, nil
}

func urlsafeB64decode(str string) []byte {
	if m := len(str) % 4; m != 0 {
		str += strings.Repeat("=", 4-m)
	}
	bt, _ := base64.URLEncoding.DecodeString(str)
	return bt
}

func btrToInt(a io.Reader) int {
	var e uint64
	binary.Read(a, binary.BigEndian, &e)
	return int(e)
}

func byteToInt(bt []byte) *big.Int {
	a := big.NewInt(0)
	a.SetBytes(bt)
	return a
}

func byteToBtr(bt0 []byte) *bytes.Reader {
	var bt1 []byte
	if len(bt0) < 8 {
		bt1 = make([]byte, 8-len(bt0), 8)
		bt1 = append(bt1, bt0...)
	} else {
		bt1 = bt0
	}
	return bytes.NewReader(bt1)
}
