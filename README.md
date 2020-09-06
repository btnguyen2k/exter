![Exter icon](docs/icons/exter_icon.png)xter is an identity gateway that allows application to authenticate users using various identity sources.

Supported identity sources:

- [x] Facebook
- [x] GitHub
- [x] Google
- [ ] Linkedin

Latest release [`v0.4.0`](RELEASE-NOTES.md).

You can [deploy Exter](BuildAndRun.md) on your own infrastructure, on-premises or cloud. Or leverage the [pre-hosted Exter](https://btnguyen2k.github.io/exter/).

## How Exter works

The authentication (login) flow is as the following:

![Exter Integration Flow](docs/Exter_flow_1.png)

When the user accesses the application website, the application calls Exter's API `verifyLoginToken` to check if the user's login-token is valid. Note: if the user curently does not have a login-token, API `verifyLoginToken` will return `invalid login-token` result.

- If the login-token is invalid, the application should redirect user to Exter's login page.
    - Then Exter authenticates user with an available identity source (*).
    - Upon successful authentication, Exter generates a login-token and passes it to the application via the callback url.
    - The application, then, should store the login-token for the user and, again, verify if the login-token is valid.
- If the login-token is valid, the application lets user access its content.

> **(*) Authentication will be done on the identity source site. User will _not_ be asked to enter any credential information on Exter.**

## Read more

- [Setup an Exter instance](BuildAndRun.md)
- [Integrate with Exter](Integration.md)
