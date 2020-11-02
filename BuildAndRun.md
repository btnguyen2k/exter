# Exter: Build & Run

`Exter` is composed of 2 components that can be built and deployed separately.
- [Frontend](fe-gui/):
    - Which is a single-page application (SPA) built with [Vue.js](https://vuejs.org/) framework.
    - Frontend interacts with backend via REST APIs.
- [Backend](be-api/):
    - Which is a [Go](https://golang.org/) project.
    - Backend has no GUI and offers only APIs for frontend to interfact with.

## Build Docker image

`Exter` can be built and packaged as an _all-on-one_ Docker image that includes both the frontend end backend components.
Simply run the following command at the project's root directory:

```
docker build --rm -t exter .
```

Build Docker image for the frontend by running the following command at the [frontend's root directory](fe-gui/):

```
docker build --rm -t exter-fe .
```

Build Docker image for the backend by running the following command at the [backend's root directory](be-api/):

```
docker build --rm -t exter-be .
```

> See [Docker documentation](https://docs.docker.com/engine/reference/commandline/build/) for command arguments.

## Run from Docker image

Run the _all-in-one_ Docker image on local machine:

```
docker run -d --name exter -p 8000:8000 \
    -e EXTER_HOME_URL="http://localhost:8000" \
    -e GOOGLE_API_PROJECT_ID=<proj-id> \
    -e GOOGLE_API_CLIENT_ID=<client-id> \
    -e GOOGLE_API_CLIENT_SECRET=<client-secret> \
    exter
```

> See [Docker documentation](https://docs.docker.com/engine/reference/commandline/run/) for command arguments.

See [Environment variables](#environment-variables) section for information on setting envinronment.

## Run on developer machine

On developer environment you can either run Exter as a container from [Docker image](#run-from-docker-image) or build & run Exter directly from source code. Note that the [frontend](fe-gui/) is a Vue.js project and the [backend](be-api/) is a Go project.

[Environment variables](#environment-variables) applies when running Exter from source code.

> You may need to configure the frontend when running from source code so that it knows where the backend is. The easiest way to do that is to create a file named `.env.local` at the frontend's root directory with the following content:
>
> ```
> VUE_APP_BO_API_BASE_URL="http://localhost:3000"
> ```
>
> Where `http://localhost:3000` is where the backend is serving.

## Environment variables

> Environment variables are only used to override the backend configurations.

**Common Application Configurations**

|Env variable                |Description                              |Default value   |
|----------------------------|-----------------------------------------|----------------|
|TIMEZONE                    |Timezone for date/time-related operations|`Asia/Ho_Chi_Minh`|
|HTTP_LISTEN_ADDR            |Listen address for REST API|`0.0.0.0`|
|HTTP_LISTEN_PORT            |Listen port for REST API|`8000`|
|HTTP_HEADER_APP_ID (1)      |Name of HTTP header that holds "application id" info passed from client|`X-App-Id`|
|HTTP_HEADER_ACCESS_TOKEN (1)|Name of HTTP header that holds "access token" info passed from client|`X-Access-Token`|
|API_MAX_REQUEST_SIZE (2)    |Maximum size of a HTTP request that client can send to Exter backend|`64kB`|
|API_REQUEST_TIMEOUT (3)     |Exter backend only waits up to this amount of time to read and parse request from client|`10s`|
|INIT_SYSTEM_OWNER_ID (4)    |User id of system "exter" app's owner||

> - (1) Changing these configurations will affect _all clients_, including Exter frontend. Do not change them unless you have a good reason to.
> - (2) Value of this configuration follows the format in this document https://github.com/lightbend/config/blob/master/HOCON.md#size-in-bytes-format
> - (3) Value of this configuration follows the format in this document https://github.com/lightbend/config/blob/master/HOCON.md#duration-format
> - (4) This is the email address of the user who will be the owner of the system "exter" app.

**Security-related Configuration**

|Env variable                |Description                              |Default value   |
|----------------------------|-----------------------------------------|----------------|
|HTTP_ALLOW_ORIGINS (1)      |CORS: value for "Access-Control-Allow-Origin" response header|`*`|
|RSA_PRIVKEY_FILE (2)        |Path to RSA private key (PEM format)|`./config/keys/exter_priv.pem`|
|RSA_PRIVKEY_PASSPHRASE (2)  |Pass-phrase for RSA private key|`exters3cr3t`|

> - (1) This affects only the Exter frontend. On development env you can use the default value. On production env put your fronend domains here. Domain names are separated by spaces or commas or semi-colons. For example `exteross.gpvcloud.com,exteross.mydomain.com;exteross.mydomain.net`.
> - (2) On production env, do _not_ use the default private key. _Generate and use your own key_.

**Database Backend Configuration**

|Env variable                |Description                              |Default value   |
|----------------------------|-----------------------------------------|----------------|
|DB_TYPE                     |Type of database backend|`sqlite`|
|DB_PGSQL_URL                |Connection string for PostgreSQL||

> - Supported database backend:
>   - `sqlite`: use [SQLite](https://sqlite.org/index.html) as database backend. Not recommended for production use. Directory to store SQLite data is `./data/sqlite`
>   - `pgsql`: use [PostgreSQL](https://www.postgresql.org/) as database backend. Recommended for production use. PostgreSQL connection string is read from environment variable `DB_PGSQL_URL`. An example of the connection string: `postgres://test:test@localhost:5432/test?sslmode=disable&client_encoding=UTF-8&application_name=exter`

**Identity Source Configurations**

|Env variable                     |Description                              |Default value   |
|---------------------------------|-----------------------------------------|----------------|
|LOGIN_CHANNELS (1)               |List of enabled login channels, comma separated|`facebook,github,google,linkedin`|
|EXTER_HOME_URL (2)               |Exter home url, used as "redirect_uri" for OAuth2||
|GOOGLE_API_PROJECT_ID (3)        |Google API's project-id||
|GOOGLE_API_CLIENT_ID (3)         |Google API's client-id||
|GOOGLE_API_CLIENT_SECRET (3)     |Google API's client-secret||
|GOOGLE_API_CLIENT_SECRET_JSON (4)|Full content of client secret file||
|GITHUB_OAUTHAPP_CLIENT_ID (5)    |GitHub OAuth App's Client ID||
|GITHUB_OAUTHAPP_CLIENT_SECRET (5)|GitHub OAuth App's Client Secret||
|FACEBOOK_APP_ID (6)              |Facebook App ID||
|FACEBOOK_APP_SECRET (6)          |Facebook App Secret||
|LINKEDIN_CLIENT_ID (7)           |LinkedIn OAuth App's Client ID||
|LINKEDIN_CLIENT_SECRET (7)       |LinkedIn OAuth App's Client Secret||
|LINKEDIN_REDIRECT_URI (8)        |Redirect uri for LinkedIn OAuth flow||

> - (1) As of version `0.5.0`, supported identity sources are `facebook`, `github`, `google` and `linkedin`.
> - (2) Used as `redirect_uri` for OAuth2 (since `v0.3.0`).
> - (3)(4) Create your Google API project at https://console.developers.google.com/apis/ and generate client secret info on page https://console.developers.google.com/apis/credentials. Either supply full content of the download client secret file in `GOOGLE_API_CLIENT_SECRET_JSON` environment variable; or supply project-id, client-id, client-secret and authorized domains info:
>   - `GOOGLE_API_PROJECT_ID`: your Google API's project id
>   - `GOOGLE_API_CLIENT_ID`: your Google API's client id
>   - `GOOGLE_API_CLIENT_SECRET`: your Google API's client secret
> - (5) Create your GitHub OAuth app at https://github.com/settings/developers
>   - Set app's `Authorization callback URL` to `<exter-url>/app/xlogin?cba=gh`
>   - `GITHUB_OAUTHAPP_CLIENT_ID`: your GitHub OAuth app's `Client ID` value
>   - `GITHUB_OAUTHAPP_CLIENT_SECRET`: your GitHub OAuth app's `Client Secret` value
> - (6) Create your Facebook app at https://developers.facebook.com/apps/
>   - `FACEBOOK_APP_ID`: your Facebook app's `App ID` value
>   - `FACEBOOK_APP_SECRET`: your Facebook app's `App Secret` value
> - (7)(8) Create your LinkedIn app with `Sign In with LinkedIn` product at https://www.linkedin.com/developers/
>   - Set app's `Authorized redirect URL` to `<exter-url>/app/xlogin?cba=ln`
>   - `LINKEDIN_CLIENT_ID`: your LinkedIn OAuth app's `Client ID` value
>   - `LINKEDIN_CLIENT_SECRET`: your LinkedIn OAuth app's `Client Secret` value
>   - `LINKEDIN_REDIRECT_URI`: same as the `Authorized redirect URL` above

## Read more

- [Integrate with Exter](Integration.md)
