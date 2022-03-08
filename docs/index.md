![Exter icon](icons/exter_icon.png)(xter) is an identity gateway that allows applications to authenticate users using various identity sources.

Supported identity sources:

- [x] Facebook
- [x] GitHub
- [x] Google
- [x] Linkedin
- [ ] Twitter

## ExterOSS

`ExterOSS` is a pre-hosted `Exter` instance that is free to use for personal and open source projects.

> **`ExterOSS` does not provide any SLA or guarantee. You are welcome to [host `Exter` on your infrastructure](https://github.com/btnguyen2k/exter/blob/master/BuildAndRun.md).**

### Privacy Policy

- `Exter` will collect the following information:
  - User's _email address_: to uniquely identify the user.
  - User's _display name_: for displaying purpose.
- How `Exter` would use the information:
  - Once the authentication is successful, `Exter` forwards the above user information to the _client application_ for further processing.
  - `Exter` has no further use of the user information.
- How `Exter` would store the user information:
  - `Exter` stores the user information in short-live storage (session) with expiry.
  - The _client application_ may also store the user information. It is outside the scope of `Exter`. Should the user do not want the client application to receive the information, simply cancel the authentication flow.
- How users request their information to be removed:
  - `Exter` stores the user information in short-live storage, called session, with expiry. At the expiry time, the session data is automatically removed. There is no further action required from the user.
  - The _client application_ may have a different policy handling user's information. It is outside the scope of `Exter`. Users may need to review the _client application_'s privacy policy for any detail.

### Technical Info

- Exter version: `latest build`
- URL Home: `https://exteross.gpvcloud.com/`
- URL Login: `https://exteross.gpvcloud.com/app/xlogin`
- URL Check login: `https://exteross.gpvcloud.com/app/xcheck`
- API info: `GET https://exteross.gpvcloud.com/info`
- API verifyLoginToken: `POST https://exteross.gpvcloud.com/api/verifyLoginToken`

> See Exter's [GitHub repository](https://github.com/btnguyen2k/exter/) for documentations.
