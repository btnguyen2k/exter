package main

import (
	"fmt"
	"regexp"
)

func main() {
	str := `{
                  "type":"authorized_user",
                  "web": {
                        "project_id": "btnguyen2k",
                        "client_id": "334322862548-9o5rr6edh0fi64vf1km0i2omtpfno1ph.apps.googleusercontent.com",
                        "client_secret": "OLkc-Ki8_Nu8OfYlSfQHtL5c",
                        "auth_uri": "https://accounts.google.com/o/oauth2/auth",
                        "token_uri": "https://oauth2.googleapis.com/token",
                        "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
                        "redirect_uris": ["http://localhost:8080","https://localhost:8080"],
                        "javascript_origins": ["http://localhost:8080","https://localhost:8080"],
                        "access_type": "offline"
                  }
                }`
	r := regexp.MustCompile(`(?s)"client_secret":\s*"(.*?)"`)
	fmt.Printf("Find %s\n", r.FindSubmatch([]byte(str))[1])
	// fmt.Printf("Replace %s\n", r.ReplaceAll([]byte(str), []byte(`"client_secret": "***L5c"`)))
}
