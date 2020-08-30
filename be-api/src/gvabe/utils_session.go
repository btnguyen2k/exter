package gvabe

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/go-github/github"
	goauthv2 "google.golang.org/api/oauth2/v2"

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

func createUserAccountFromGitHubProfile(ui *github.User) (*user.User, error) {
	var u *user.User
	var err error
	if ui.Email == nil {
		return nil, errors.New("github profile does not contain email address")
	}
	if u, err = userDao.Get(*ui.Email); err == nil && u == nil {
		u = user.NewUser(goapi.AppVersionNumber, *ui.Email)
		var ok bool
		if ok, err = userDao.Create(u); err != nil || !ok {
			u = nil
		}
	}
	return u, err
}

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
