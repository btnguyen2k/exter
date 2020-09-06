package gvabe

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
)

var (
	githubOAuthConf = &oauth2.Config{
		ClientID:     "",
		ClientSecret: "",
		Scopes:       []string{"user:email"},
		Endpoint:     githuboauth.Endpoint,
	}
)

// routine to fetch GitHub profile in background
func goFetchGitHubProfile(sessId string) {
	if bo, err := sessionDao.Get(sessId); err != nil {
		log.Println(fmt.Sprintf("[ERROR] goFetchGitHubProfile(%s) - error loading session data: %e", sessId, err))
	} else if bo == nil {
		log.Println(fmt.Sprintf("[WARN] goFetchGitHubProfile(%s) - session does not exist", sessId))
	} else if bo.IsExpired() {
		log.Println(fmt.Sprintf("[WARN] goFetchGitHubProfile(%s) - session expired", sessId))
	} else if claims, err := parseLoginToken(bo.GetSessionData()); err != nil {
		log.Println(fmt.Sprintf("[ERROR] goFetchGitHubProfile(%s) - cannot parse JWT token: %e", sessId, err))
	} else if claims.Type != sessionTypePreLogin || claims.isExpired() {
		log.Println(fmt.Sprintf("[WARN] goFetchGitHubProfile(%s) - invalid claims type of JWT expired", sessId))
	} else {
		sess := &Session{}
		if err := json.Unmarshal(claims.Data, &sess); err != nil {
			log.Println(fmt.Sprintf("[ERROR] goFetchGitHubProfile(%s) - error decoding session: %e", sessId, err))
			return
		}
		if sess.Channel != loginChannelGithub {
			log.Println(fmt.Sprintf("[WARN] goFetchGitHubProfile(%s) - invalid login channel: %s", sessId, sess.Channel))
			return
		}
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		oauth2Token := &oauth2.Token{}
		if err := json.Unmarshal(sess.Data, &oauth2Token); err != nil {
			log.Println(fmt.Sprintf("[ERROR] goFetchGitHubProfile - error unmarshalling oauth2.Token: %e", err))
		} else if githubClient := github.NewClient(githubOAuthConf.Client(ctx, oauth2Token)); githubClient == nil {
			log.Println(fmt.Sprintf("[ERROR] goFetchGitHubProfile - error creating new GitHub API client: nill"))
		} else if userinfo, _, err := githubClient.Users.Get(ctx, ""); err != nil {
			log.Println(fmt.Sprintf("[ERROR] goFetchGitHubProfile - error fetching GitHub userinfo: %e", err))
		} else {
			if u, err := createUserAccountFromGitHubProfile(userinfo); err != nil {
				log.Println(fmt.Sprintf("[ERROR] goFetchGitHubProfile - error creating user account from GitHub userinfo: %e", err))
			} else {
				js, _ := json.Marshal(oauth2Token)
				sess.UserId = u.GetId()
				sess.DisplayName = u.GetDisplayName()
				sess.ExpiredAt = oauth2Token.Expiry
				sess.Data = js
				claims, err := genLoginClaims(sessId, sess)
				if err != nil {
					log.Println(fmt.Sprintf("[ERROR] goFetchGitHubProfile(%s) - error generating login token: %e", sessId, err))
				}
				_, _, err = saveSession(claims)
				if err != nil {
					log.Println(fmt.Sprintf("[ERROR] goFetchGitHubProfile(%s) - error saving login token: %e", sessId, err))
				}
			}
		}
	}
}

func githubFetchUserProfile(ctx context.Context, accessToken string) (map[string]interface{}, error) {
	ghUrl := "https://api.github.com/user"
	if ctx == nil {
		ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	}
	req, _ := http.NewRequestWithContext(ctx, "GET", ghUrl, nil)
	req.Header.Set("Authorization", "token "+accessToken)
	req.Header.Set("Accept", "application/json")
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New("github OAuth response status: " + resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s", body)
	return nil, nil
	// ghresp := make(map[string]interface{})
	// if err := json.Unmarshal(body, &ghresp); err != nil {
	// 	return "", err
	// }
	// accessToken, ok := ghresp["access_token"].(string)
	// if !ok {
	// 	if DEBUG {
	// 		log.Printf("[DEBUG] githubExchangeAccessToken - GitHub OAuth response: %s", body)
	// 	}
	// 	return "", errors.New("invalid GitHub OAuth access token")
	// }
	// return accessToken, nil
}
