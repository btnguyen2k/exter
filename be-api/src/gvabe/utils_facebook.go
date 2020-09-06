package gvabe

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/btnguyen2k/consu/reddo"
	"github.com/btnguyen2k/consu/semita"
	fbv2 "github.com/huandu/facebook/v2"
	"golang.org/x/oauth2"
	facebookoauth "golang.org/x/oauth2/facebook"
)

var (
	fbOAuthConf = &oauth2.Config{
		ClientID:     "",
		ClientSecret: "",
		Scopes:       []string{"public_profile,email"},
		Endpoint:     facebookoauth.Endpoint,
	}
	fbApp *fbv2.App
)

type FbToken struct {
	Token     string `json:"token,access_token"`
	Type      string `json:"token_type"`
	ExpiresIn int    `json:"expires_in"`
}

func initFacebookApp(appId, appSecret string) {
	fbApp = fbv2.New(appId, appSecret)
}

// exchange a short-live token for long-live one. Can be used to renew token?
func fbExchangeForLongLiveToken(ctx context.Context, accessToken string) (*oauth2.Token, error) {
	if ctx == nil {
		ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	}
	resp, err := fbApp.Session(accessToken).WithContext(ctx).Get(
		"/oauth/access_token",
		fbv2.Params{
			"grant_type":        "fb_exchange_token",
			"client_id":         fbOAuthConf.ClientID,
			"client_secret":     fbOAuthConf.ClientSecret,
			"fb_exchange_token": accessToken,
		},
	)
	if err != nil {
		return nil, err
	}
	s := semita.NewSemita(resp)
	token := &oauth2.Token{
	}
	if value, err := s.GetValueOfType("access_token", reddo.TypeString); err == nil {
		token.AccessToken = value.(string)
		token.RefreshToken = value.(string)
	} else {
		return nil, err
	}
	if value, err := s.GetValueOfType("token_type", reddo.TypeString); err == nil {
		token.TokenType = value.(string)
	} else {
		return nil, err
	}
	if value, err := s.GetValueOfType("expires_in", reddo.TypeInt); err == nil {
		token.Expiry = time.Now().Add(time.Duration(value.(int64)) * time.Second)
	} else {
		return nil, err
	}
	return token, nil
}

func fbGetProfile(ctx context.Context, accessToken string) (map[string]interface{}, error) {
	if ctx == nil {
		ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	}
	return fbApp.Session(accessToken).WithContext(ctx).Get(
		"/me",
		fbv2.Params{"access_token": accessToken, "fields": "email,name"},
	)
}

// routine to fetch Facebook profile in background
func goFetchFacebookProfile(sessId string) {
	if bo, err := sessionDao.Get(sessId); err != nil {
		log.Println(fmt.Sprintf("[ERROR] goFetchFacebookProfile(%s) - error loading session data: %e", sessId, err))
	} else if bo == nil {
		log.Println(fmt.Sprintf("[WARN] goFetchFacebookProfile(%s) - session does not exist", sessId))
	} else if bo.IsExpired() {
		log.Println(fmt.Sprintf("[WARN] goFetchFacebookProfile(%s) - session expired", sessId))
	} else if claims, err := parseLoginToken(bo.GetSessionData()); err != nil {
		log.Println(fmt.Sprintf("[ERROR] goFetchFacebookProfile(%s) - cannot parse JWT token: %e", sessId, err))
	} else if claims.Type != sessionTypePreLogin || claims.isExpired() {
		log.Println(fmt.Sprintf("[WARN] goFetchFacebookProfile(%s) - invalid claims type of JWT expired", sessId))
	} else {
		sess := &Session{}
		if err := json.Unmarshal(claims.Data, &sess); err != nil {
			log.Println(fmt.Sprintf("[ERROR] goFetchFacebookProfile(%s) - error decoding session: %e", sessId, err))
			return
		}
		if sess.Channel != loginChannelFacebook {
			log.Println(fmt.Sprintf("[WARN] goFetchFacebookProfile(%s) - invalid login channel: %s", sessId, sess.Channel))
			return
		}
		oauth2Token := &oauth2.Token{}
		if err := json.Unmarshal(sess.Data, &oauth2Token); err != nil {
			log.Println(fmt.Sprintf("[ERROR] goFetchFacebookProfile - error unmarshalling oauth2.Token: %e", err))
		} else {
			ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
			if profile, err := fbGetProfile(ctx, oauth2Token.AccessToken); err != nil {
				log.Println(fmt.Sprintf("[ERROR] goFetchFacebookProfile - error fetching Facebook userinfo: %e", err))
			} else {
				if u, err := createUserAccountFromFacebookProfile(profile); err != nil {
					log.Println(fmt.Sprintf("[ERROR] goFetchFacebookProfile - error creating user account from Facebook userinfo: %e", err))
				} else {
					js, _ := json.Marshal(oauth2Token)
					sess.UserId = u.GetId()
					sess.DisplayName = u.GetDisplayName()
					sess.ExpiredAt = oauth2Token.Expiry
					sess.Data = js // JSON-serialization of oauth2.Token
					claims, err := genLoginClaims(sessId, sess)
					if err != nil {
						log.Println(fmt.Sprintf("[ERROR] goFetchFacebookProfile(%s) - error generating login token: %e", sessId, err))
					}
					_, _, err = saveSession(claims)
					if err != nil {
						log.Println(fmt.Sprintf("[ERROR] goFetchFacebookProfile(%s) - error saving login token: %e", sessId, err))
					}
				}
			}
		}
	}
}
