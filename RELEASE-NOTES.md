# Exter release notes

## 2021-02-18 - v0.6.0

- AB#31: support AWS DynamoDB.
- AB#32: support MongoDB.
- Fix AB#34: Fix bug `403: Return url [] is not allowed for app [exter]`.
- Fix AB#35: Display favicon on browser.
- AB#36: support Azure CosmosDB.

## 2020-11-02 - v0.5.0

- EP#3: support login with LinkedIn (AB#9, AB#10).
- EP#15: fetch display name from LinkedIn (AB#19).

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
