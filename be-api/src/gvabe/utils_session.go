package gvabe

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/oauth2"
	goauthv2 "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"

	"main/src/goapi"
	"main/src/gvabe/bo/session"
	"main/src/gvabe/bo/user"
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

// Session captures a user-login-session. Session object is to be serialized and embedded into a SessionClaims.
type Session struct {
	ClientId  string    `json:"cid"`  // application's id
	Channel   string    `json:"chan"` // login source/channel (Google, Facebook, etc)
	UserId    string    `json:"uid"`  // id of logged-in user
	CreatedAt time.Time `json:"cat"`  // timestamp when the session is created
	ExpiredAt time.Time `json:"eat"`  // timestamp when the session expires
	Data      []byte    `json:"data"` // session's arbitrary data
}

// SessionClaims is an extended structure of JWT's standard claims
type SessionClaims struct {
	Type   string `json:"type"`           // session type (pre-login or logged-in)
	UserId string `json:"uid,omitempty"`  // id of logged-in user
	Data   []byte `json:"data,omitempty"` // session's arbitrary data
	jwt.StandardClaims
}

func (s *SessionClaims) isExpired() bool {
	return s.ExpiresAt > 0 && s.ExpiresAt < time.Now().Unix()
}

func (s *SessionClaims) isGoingExpired(numSec int64) bool {
	return s.ExpiresAt > 0 && s.ExpiresAt-numSec < time.Now().Unix()
}

/*----------------------------------------------------------------------*/
func saveSession(claims *SessionClaims) (*session.Session, string, error) {
	if claims.Id == "" {
		claims.Id = utils.UniqueId()
	}
	jwt, err := genJws(claims)
	if err != nil {
		return nil, "", err
	}
	expiry := time.Unix(claims.ExpiresAt, 0)
	sess := session.NewSession(goapi.AppVersionNumber, claims.Id, claims.Type, claims.Subject, claims.Audience, claims.UserId, jwt, expiry)
	_, err = sessionDao.Save(sess)
	return sess, jwt, err
}

/*----------------------------------------------------------------------*/

func createUserAccountFromGoogleProfile(ui *goauthv2.Userinfo) (*user.User, error) {
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

func goFetchGoogleProfile(sessId string) {
	t1 := time.Now()
	if bo, err := sessionDao.Get(sessId); err != nil {
		log.Println(fmt.Sprintf("[ERROR] goFetchGoogleProfile(%s) - error loading session data: %e", sessId, err))
	} else if bo == nil {
		log.Println(fmt.Sprintf("[WARN] goFetchGoogleProfile(%s) - session does not exist", sessId))
	} else if bo.IsExpired() {
		log.Println(fmt.Sprintf("[WARN] goFetchGoogleProfile(%s) - session expired", sessId))
	} else if claims, err := parseLoginToken(bo.GetSessionData()); err != nil {
		log.Println(fmt.Sprintf("[ERROR] goFetchGoogleProfile(%s) - cannot parse JWT token: %e", sessId, err))
	} else if claims.Type != sessionTypePreLogin || claims.isExpired() {
		log.Println(fmt.Sprintf("[WARN] goFetchGoogleProfile(%s) - invalid claims type of JWT expired", sessId))
	} else {
		sess := &Session{}
		if err := json.Unmarshal(claims.Data, &sess); err != nil {
			log.Println(fmt.Sprintf("[ERROR] goFetchGoogleProfile(%s) - error decoding session: %e", sessId, err))
			return
		}
		if sess.Channel != loginChannelGoogle {
			log.Println(fmt.Sprintf("[WARN] goFetchGoogleProfile(%s) - invalid login channel: %s", sessId, sess.Channel))
			return
		}
		oauth2Token := &oauth2.Token{}
		if err := json.Unmarshal(sess.Data, &oauth2Token); err != nil {
			log.Println(fmt.Sprintf("[ERROR] goFetchGoogleProfile - error unmarshalling oauth2.Token: %e", err))
		} else if oauth2Service, err := goauthv2.NewService(context.Background(), option.WithTokenSource(gConfig.TokenSource(context.Background(), oauth2Token))); err != nil {
			log.Println(fmt.Sprintf("[ERROR] goFetchGoogleProfile - error creating new Google Service: %e", err))
		} else if userinfo, err := oauth2Service.Userinfo.V2.Me.Get().Do(); err != nil {
			log.Println(fmt.Sprintf("[ERROR] goFetchGoogleProfile - error fetching Google userinfo: %e", err))
		} else {
			if u, err := createUserAccountFromGoogleProfile(userinfo); err != nil {
				log.Println(fmt.Sprintf("[ERROR] goFetchGoogleProfile - error creating user account from Google userinfo: %e", err))
			} else {
				now := time.Now()
				expiry := now.Add(3600 * time.Second)
				js, _ := json.Marshal(oauth2Token)
				sess.UserId = u.GetId()
				sess.ExpiredAt = expiry
				sess.Data = js
				claims, err := genLoginClaims(sessId, sess)
				if err != nil {
					log.Println(fmt.Sprintf("[ERROR] goFetchGoogleProfile(%s) - error generating login token: %e", sessId, err))
				}
				_, _, err = saveSession(claims)
				if err != nil {
					log.Println(fmt.Sprintf("[ERROR] goFetchGoogleProfile(%s) - error saving login token: %e", sessId, err))
				}
				log.Printf("[goFetchGoogleProfile] finished in %d ms", time.Now().Sub(t1).Milliseconds())
			}
		}
	}
}

func genJws(claim *SessionClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claim)
	return token.SignedString(rsaPrivKey)
}

// genLoginClaims generates a login token as SessionClaims:
//   - the SessionClaims is created with type=login and populated with data from supplied session
func genLoginClaims(id string, sess *Session) (*SessionClaims, error) {
	if id == "" {
		id = utils.UniqueId()
	}
	u, err := userDao.Get(sess.UserId)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, errors.New(fmt.Sprintf("user [%s] not found", sess.UserId))
	}
	sessData, err := json.Marshal(sess)
	if err != nil {
		return nil, err
	}
	sessData, err = zipAndEncrypt(sessData, []byte(u.GetAesKey()))
	return &SessionClaims{
		UserId: sess.UserId,
		Type:   sessionTypeLogin,
		Data:   sessData,
		StandardClaims: jwt.StandardClaims{
			Audience:  sess.ClientId,
			ExpiresAt: sess.ExpiredAt.Unix(),
			Id:        id,
			IssuedAt:  sess.CreatedAt.Unix(),
			Subject:   sess.Channel,
		},
	}, err
}

// genLoginToken generates a login token in JWT format:
//   - a SessionClaims is created with type=login and populated with data from supplied session
//   - the session claim is used to created JWT, the JWT is then signed with rsaPrivKey
func genLoginToken(id string, sess *Session) (*SessionClaims, string, error) {
	claims, err := genLoginClaims(id, sess)
	if err != nil {
		return nil, "", err
	}
	jwt, err := genJws(claims)
	return claims, jwt, err
}

// genPreLoginClaims generates a pre-login token as SessionClaims:
//   - the SessionClaims is created with type=pre-login and populated with data from supplied session
func genPreLoginClaims(sess *Session) (*SessionClaims, error) {
	sessData, err := json.Marshal(sess)
	return &SessionClaims{
		Type: sessionTypePreLogin,
		Data: sessData,
		StandardClaims: jwt.StandardClaims{
			Audience:  sess.ClientId,
			ExpiresAt: sess.ExpiredAt.Unix(),
			Id:        utils.UniqueId(),
			IssuedAt:  sess.CreatedAt.Unix(),
			Subject:   sess.Channel,
		},
	}, err
}

// genPreLoginToken generates a pre-login token in JWT format:
//   - a SessionClaims is created with type=pre-login and populated with data from supplied session
//   - the session claim is used to created JWT, the JWT is then signed with rsaPrivKey
func genPreLoginToken(sess *Session) (*SessionClaims, string, error) {
	claims, err := genPreLoginClaims(sess)
	if err != nil {
		return nil, "", err
	}
	jwt, err := genJws(claims)
	return claims, jwt, err
}

func parseLoginToken(jwtStr string) (*SessionClaims, error) {
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
		var result SessionClaims
		js, _ := json.Marshal(claims)
		return &result, json.Unmarshal(js, &result)
	} else {
		return nil, errors.New("invalid claim")
	}
}
