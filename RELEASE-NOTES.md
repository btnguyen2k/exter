# Exter release notes

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
