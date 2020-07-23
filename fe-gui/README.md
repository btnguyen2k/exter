# Exter Frontend

Frontend component of Exter - based on [govueadmin.g8](https://github.com/btnguyen2k/govueadmin.g8).

Template by [CoreUI for Vue.js](https://coreui.io/vue/docs/introduction/).

##Getting Started

**Project setup**
```
npm install
```

**Compiles and hot-reloads for development**
```
npm run serve
```

**Compiles and minifies for production**
```
npm run build
```

**Lints and fixes files**
```
npm run lint
```

##Application Configurations

`src/config.json`: application's main configuration file. Important config keys:
- `api_client.bo_api_base_url`: point to backend's base URL
- `api_client.app_id`: application id in order to authenticate with backend. _Must match between frontend and backend._
- `api_client.header_app_id` and `api_client.header_access_token`: name of HTTP headers passed along every API call for authentication. _Must match between frontend and backend._

> `api_client.bo_api_base_url` is empty by default. This value can be overridden by environment variable `VUE_APP_BO_API_BASE_URL`.
> E.g. on development machine, create a file `.env` at frontend's root directory with the following content:
>
> ```
> VUE_APP_BO_API_BASE_URL="http://localhost:3000"
> ```
>
> Where `http://localhost:3000` is the base url of the backend.

## References

- [VueJS Configuration Reference](https://cli.vuejs.org/config/)
- [CoreUI for Vue.JS](https://coreui.io/vue/docs/introduction/)

## LICENSE & COPYRIGHT

See [LICENSE.md](../LICENSE.md).
