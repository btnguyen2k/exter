package gvabe

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/btnguyen2k/consu/gjrc"
	"golang.org/x/oauth2"
	linkedinoauth "golang.org/x/oauth2/linkedin"
)

var (
	linkedinOAuthConf = &oauth2.Config{
		ClientID:     "",
		ClientSecret: "",
		Scopes:       []string{"r_liteprofile", "r_emailaddress"},
		Endpoint:     linkedinoauth.Endpoint,
	}
)

// routine to fetch LinkedIn profile in background
func goFetchLinkedInProfile(sessId string) {
	if bo, err := sessionDao.Get(sessId); err != nil {
		log.Println(fmt.Sprintf("[ERROR] goFetchLinkedInProfile(%s) - error loading session data: %e", sessId, err))
	} else if bo == nil {
		log.Println(fmt.Sprintf("[WARN] goFetchLinkedInProfile(%s) - session does not exist", sessId))
	} else if bo.IsExpired() {
		log.Println(fmt.Sprintf("[WARN] goFetchLinkedInProfile(%s) - session expired", sessId))
	} else if claims, err := parseLoginToken(bo.GetSessionData()); err != nil {
		log.Println(fmt.Sprintf("[ERROR] goFetchLinkedInProfile(%s) - cannot parse JWT token: %e", sessId, err))
	} else if claims.Type != sessionTypePreLogin || claims.isExpired() {
		log.Println(fmt.Sprintf("[WARN] goFetchLinkedInProfile(%s) - invalid claims type of JWT expired", sessId))
	} else {
		sess := &Session{}
		if err := json.Unmarshal(claims.Data, &sess); err != nil {
			log.Println(fmt.Sprintf("[ERROR] goFetchLinkedInProfile(%s) - error decoding session: %e", sessId, err))
			return
		}
		if sess.Channel != loginChannelLinkedin {
			log.Println(fmt.Sprintf("[WARN] goFetchLinkedInProfile(%s) - invalid login channel: %s", sessId, sess.Channel))
			return
		}
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		oauth2Token := &oauth2.Token{}
		if err := json.Unmarshal(sess.Data, &oauth2Token); err != nil {
			log.Println(fmt.Sprintf("[ERROR] goFetchLinkedInProfile - error unmarshalling oauth2.Token: %e", err))
		} else if httpClient := linkedinOAuthConf.Client(ctx, oauth2Token); httpClient == nil {
			log.Println(fmt.Sprintf("[ERROR] goFetchLinkedInProfile - error creating new LinkedIn API httpClient: nill"))
		} else {
			if u, err := createUserAccountFromLinkedInProfile(gjrc.NewGjrc(httpClient, 0)); err != nil {
				log.Println(fmt.Sprintf("[ERROR] goFetchLinkedInProfile - error creating user account from LinkedIn profile: %e", err))
			} else {
				js, _ := json.Marshal(oauth2Token)
				sess.UserId = u.GetId()
				sess.DisplayName = u.GetDisplayName()
				sess.ExpiredAt = oauth2Token.Expiry
				sess.Data = js
				claims, err := genLoginClaims(sessId, sess)
				if err != nil {
					log.Println(fmt.Sprintf("[ERROR] goFetchLinkedInProfile(%s) - error generating login token: %e", sessId, err))
				}
				_, _, err = saveSession(claims)
				if err != nil {
					log.Println(fmt.Sprintf("[ERROR] goFetchLinkedInProfile(%s) - error saving login token: %e", sessId, err))
				}
			}
		}
	}
}
