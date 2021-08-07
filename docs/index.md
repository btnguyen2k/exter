![Exter icon](icons/exter_icon.png) is an identity gateway that allows application to authenticate users using various identity sources.

Supported identity sources:

- [x] Facebook
- [x] GitHub
- [x] Google
- [x] Linkedin
- [ ] Twitter

## ExterOSS

`ExterOSS` is a pre-hosted `Exter` instance which is free to use for personal and open source projects.

> **`ExterOSS` does not provide any SLA or guarantee. You are welcome to [host `Exter` on your own infrastructure](https://github.com/btnguyen2k/exter/blob/master/BuildAndRun.md).**

### Privacy Policy

- After authentication, `Exter` forwards user's basic info (_email address_ and _display name_) to the registered application as JWT.
`Exter` does _not_ store user info in a long-term storage. Contact the registered application's owner to learn more about its Terms of Service and Privacy Policy.
- `Exter` stores the JWT in user login session. The login session will expire after 60 minutes and will be automatically invalidated.

### Technical Info

- Exter version: `latest build`
- URL Home: `https://exteross.gpvcloud.com/`
- URL Login: `https://exteross.gpvcloud.com/app/xlogin`
- URL Check login: `https://exteross.gpvcloud.com/app/xcheck`
- API info: `GET https://exteross.gpvcloud.com/info`
- API verifyLoginToken: `POST https://exteross.gpvcloud.com/api/verifyLoginToken`

> See Exter's [GitHub repository](https://github.com/btnguyen2k/exter/) for documentations.
