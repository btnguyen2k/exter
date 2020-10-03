# Exter release notes

## 2020-10-02 - v0.4.2

- AB#27: (bugfix) Redirect to application only with 'login' token.
- AB#28: (Support) `cancel url`.

## 2020-09-28 - v0.4.1

- AB#26: return RSA public key in different formats

## 2020-09-06 - v0.4.0

- AB#15: Fetch display name from Facebook (AB#18), GitHub (AB#17), Google (AB#16).
- User's display name is returned in session JWT.
- Fix incorrect redirection when logging into Exter backend.

## 2020-09-05 - v0.3.1

- Leverage [Gravatar](https://gravatar.com/) for user avatar (used in Exter backend)
- Fix: access to `/favicon.ico`

## 2020-09-04 - v0.3.0

- Support login with `Facebook` (EPIC#2).
- Config `gvabe.channels.google.app_domains` is deprecated, replaced by `gvabe.exter_home_url`

## 2020-08-31 - v0.2.0

- Support login with `GitHub` (EPIC#4).
- Breaking change: change Vue router mode to `history`:
  - Page `/app/#/xlogin` is now `/app/xlogin`.
  - Page `/app/#/xcheck` is now `/app/xcheck`.
- Fix bug: Login page broken on mobile view (AB#1).
- Other fixes and enhancements.

## 2020-07-24 - v0.1.0

First release:

- Login using [Google account](https://www.google.com/account/about/).
