![Exter icon](icons/exter_icon.png) is an identity gateway that allows application to authenticate users using various identity sources.

Supported identity sources:

- [ ] Facebook
- [x] GitHub
- [x] Google
- [ ] Linkedin

## ExterOSS

`ExterOSS` is a pre-hosted `Exter` instance which is free to use for personal and open source projects.

> **`ExterOSS` does not provide any SLA or guarantee. You are welcome to [host `Exter` on your own infrastructure](https://github.com/btnguyen2k/exter/blob/master/BuildAndRun.md).**

### Privacy Policy

- Email is used as user id. Upon successful authentication, the `ExterOSS` will store your email address in its database and login session.
- Your email address is only used to uniquely identify yourself and will not be used for any other purpose. 
- `ExterOSS` will remove its session data periodically without notice.

### Technical Info

- Exter version: `latest build`
- URL Home: `http://exteross.gpvcloud.com/`
- URL Login: `http://exteross.gpvcloud.com/app/xlogin`
- URL Check login: `http://exteross.gpvcloud.com/app/xcheck`
- API info: `GET http://exteross.gpvcloud.com/info`
- API verifyLoginToken: `POST http://exteross.gpvcloud.com/api/verifyLoginToken`

> See Exter's [GitHub repository](https://github.com/btnguyen2k/exter/) for documentations.
