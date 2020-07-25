# ExterOSS

![Exter icon](icons/exter_icon.png) is an identity gateway that allows application to authenticate users using various identity sources.

Supported identity sources:
- [ ] Facebook
- [ ] GitHub
- [x] Google
- [ ] Linkedin

See Exter's [GitHub repository](https://github.com/btnguyen2k/exter/) for documentations.

## ExterOSS

`ExterOSS` is a pre-hosted `Exter` instance which is free to use for personal and open source projects.

> **`ExterOSS` does not provide any SLA or guarantee. You are encouraged to [host `Exter` on your own infrastructure](README.md#buil--run).**

### Privacy Policy

- Email is used as user id. Upon successful authentication, the `ExterOSS` will store your email address in its database and login session.
- Your email address is only used to uniquely identify yourself and will not be used for any other purpose. 
- `ExterOSS` will remove its session data periodically without notice.

### Technical Info

- Exter version: `latest build`
- URL Home: `https://exteross.herokuapp.com/`
- URL Login: `https://exteross.herokuapp.com/app/#/xlogin`
- URL Check login: `https://exteross.herokuapp.com/app/#/xcheck`
- API info: `GET https://exteross.herokuapp.com/info`
- API verifyLoginToken: `POST https://exteross.herokuapp.com/api/verifyLoginToken`
