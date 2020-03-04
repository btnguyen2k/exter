package gvabe

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/oauth2"
	goauthv2 "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"

	"main/src/goapi"
	"main/src/gvabe/bo/user"
	"main/src/mico"
	"main/src/utils"
)

const (
	sessionTypePreLogin = "pre_login"
	sessionTypeLogin    = "login"
)

var (
	errorInvalidClient = errors.New("invalid client id")
	errorInvalidJwt    = errors.New("cannot decode token")
	errorExpiredJwt    = errors.New("token has expired")
)

type Session struct {
	ClientId  string    `json:"cid"`
	Channel   string    `json:"channel"`
	UserId    string    `json:"uid"`
	CreatedAt time.Time `json:"cat"`
	ExpiredAt time.Time `json:"eat"`
	Data      []byte    `json:"data"`
}

type SessionClaim struct {
	Type     string `json:"type"`
	UserId   string `json:"uid,omitempty"`
	Data     []byte `json:"data,omitempty"`
	CacheKey string `json:"ckey,omitempty"`
	jwt.StandardClaims
}

func (s *SessionClaim) isExpired() bool {
	return s.ExpiresAt > 0 && s.ExpiresAt < time.Now().Unix()
}

func (s *SessionClaim) isGoingExpired(numSec int64) bool {
	return s.ExpiresAt > 0 && s.ExpiresAt-numSec < time.Now().Unix()
}

/*----------------------------------------------------------------------*/

func loadPreLoginSessionFromCache(key string) (*Session, error) {
	var session Session
	if data, err := preLoginSessionCache.Get(key); err != nil {
		return nil, err
	} else if data != nil {
		if err := json.Unmarshal(data, &session); err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}
	return &session, nil
}

func createUserAccountFromGoogleProfile(ui *goauthv2.Userinfoplus) (*user.User, error) {
	var u *user.User
	var err error
	if u, err = userDao.Get(ui.Email); err == nil && u == nil {
		u = user.NewUser(goapi.AppVersionNumber, ui.Email)
		var ok bool
		if ok, err = userDao.Create(u); err != nil || !ok {
			u = nil
		}
	}
	return u, err
}

func goFetchGoogleProfile(jwtToken string) {
	var sessionClaim *SessionClaim
	var err error
	if sessionClaim, err = parseLoginToken(jwtToken); err != nil {
		log.Println(fmt.Sprintf("[WARN] goFetchGoogleProfile - cannot parse JWT token: %e", err))
	} else if sessionClaim.Type == sessionTypePreLogin {
		var session *Session
		if session, err = loadPreLoginSessionFromCache(sessionClaim.CacheKey); err != nil {
			log.Println(fmt.Sprintf("[WARN] goFetchGoogleProfile - error fetching session from storage: %e", err))
		} else if session != nil {
			var oauth2Token oauth2.Token
			if err := json.Unmarshal(session.Data, &oauth2Token); err != nil {
				log.Println(fmt.Sprintf("[WARN] goFetchGoogleProfile - error unmarshalling oauth2.Token: %e", err))
			} else if oauth2Service, err := goauthv2.NewService(context.Background(), option.WithTokenSource(gConfig.TokenSource(context.Background(), &oauth2Token))); err != nil {
				log.Println(fmt.Sprintf("[WARN] goFetchGoogleProfile - error creating new Google Service: %e", err))
			} else if userinfo, err := oauth2Service.Userinfo.V2.Me.Get().Do(); err != nil {
				log.Println(fmt.Sprintf("[WARN] goFetchGoogleProfile - error fetching Google userinfo: %e", err))
			} else {
				if u, err := createUserAccountFromGoogleProfile(userinfo); err != nil {
					log.Println(fmt.Sprintf("[WARN] goFetchGoogleProfile - error creating user account from Google userinfo: %e", err))
				} else {
					js, _ := json.Marshal(oauth2Token)
					session.UserId = u.Id
					session.ExpiredAt = time.Now().Add(3600 * time.Second)
					session.Data = js
					if _, err = serializeAndCacheSession(preLoginSessionCache, *session, sessionClaim.CacheKey); err != nil {
						log.Println(fmt.Sprintf("[WARN] goFetchGoogleProfile - error caching user session: %e", err))
					}
				}
			}
		}
	}
}

func genJws(claim SessionClaim) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claim)
	return token.SignedString(rsaPrivKey)
}

func serializeAndCacheSession(cache mico.ICache, session Session, cacheKey string) ([]byte, error) {
	serSesion, err := json.Marshal(session)
	if cache != nil && cacheKey != "" {
		err = cache.Set(cacheKey, serSesion)
	}
	return serSesion, err
}

// genLoginToken generates a login token in JWT format:
//   - the supplied session is serialized and stored in preLoginSessionCache
//   - a SessionClaim is created with type=login and populated with data from supplied session
//   - the session claim is used to created JWT, the JWT is then signed with rsaPrivKey
func genLoginToken(session Session, cacheKey string) (string, error) {
	u, err := userDao.Get(session.UserId)
	if err != nil {
		return "", err
	}
	if u == nil {
		return "", errors.New(fmt.Sprintf("user [%s] not found", session.UserId))
	}
	// if strings.TrimSpace(cacheKey) == "" {
	// 	cacheKey = utils.UniqueId()
	// }
	serSesion, err := serializeAndCacheSession(preLoginSessionCache, session, cacheKey)
	if err != nil {
		return "", err
	}
	serSesion, err = zipAndEncrypt(serSesion, []byte(u.AesKey))
	if err != nil {
		return "", err
	}
	claim := SessionClaim{
		Type:     sessionTypeLogin,
		UserId:   session.UserId,
		Data:     serSesion,
		CacheKey: cacheKey,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  session.CreatedAt.Unix(),
			ExpiresAt: session.ExpiredAt.Unix(),
		},
	}
	return genJws(claim)
}

// genPreLoginToken generates a pre-login token in JWT format:
//   - the supplied session is serialized and stored in preLoginSessionCache
//   - a SessionClaim is created with type=pre-login and populated with data from supplied session
//   - the session claim is used to created JWT, the JWT is then signed with rsaPrivKey
func genPreLoginToken(session Session, cacheKey string) (string, error) {
	if strings.TrimSpace(cacheKey) == "" {
		cacheKey = utils.UniqueId()
	}
	_, err := serializeAndCacheSession(preLoginSessionCache, session, cacheKey)
	if err != nil {
		return "", err
	}
	claim := SessionClaim{
		Type:     sessionTypePreLogin,
		CacheKey: cacheKey,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  session.CreatedAt.Unix(),
			ExpiresAt: session.ExpiredAt.Unix(),
		},
	}
	return genJws(claim)
}

func parseLoginToken(jwtStr string) (*SessionClaim, error) {
	token, err := jwt.Parse(jwtStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("enexpected signing method: %v", token.Header["alg"])
		}
		return rsaPubKey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		var result SessionClaim
		js, _ := json.Marshal(claims)
		return &result, json.Unmarshal(js, &result)
	} else {
		return nil, errors.New("invalid claim")
	}
}
