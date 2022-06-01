<template>
  <div class="c-app flex-row align-items-center">
    <CContainer>
      <CRow class="justify-content-center">
        <CCol md="8">
          <CCardGroup>
            <CCard class="p-4">
              <CCardBody>
                <h1>Login</h1>
                <CAlert v-if="errorMsg!=''" color="danger">{{ errorMsg }}</CAlert>
                <CAlert v-if="initStatus==0" color="info">{{ $t('message.wait') }}</CAlert>
                <CForm method="post" v-if="initStatus>0">
                  <p v-if="infoMsg!=''" class="text-muted">{{ infoMsg }}</p>
                  <CButton v-if="sources.facebook && waitCounter<0" id="loginFb" type="button" name="facebook"
                           color="facebook" class="mb-1" block
                           @click="doLoginFacebook">
                    <CIcon name="cib-facebook"/>
                    {{ $t('message.login_facebook') }}
                  </CButton>
                  <CButton v-if="sources.github && waitCounter<0" id="loginGithub" type="button" name="github"
                           color="github" class="mb-1" block
                           @click="doLoginGitHub">
                    <CIcon name="cib-github"/>
                    {{ $t('message.login_github') }}
                  </CButton>
                  <CButton v-if="sources.google && waitCounter<0" id="loginGoogle" type="button" name="google"
                           color="light" class="mb-1" block
                           @click="doLoginGoogle">
                    <CIcon name="cib-google"/>
                    {{ $t('message.login_google') }}
                  </CButton>
                  <CButton v-if="sources.linkedin && waitCounter<0" id="loginLinkedIn" type="button" name="linkedin"
                           color="linkedin" class="mb-1" block
                           @click="doLoginLinkedIn">
                    <CIcon name="cib-linkedin"/>
                    {{ $t('message.login_linkedin') }}
                  </CButton>
                  <CRow v-if="cancelUrl!=''">
                    <CCol col="12" class="text-right">
                      <CButton color="link" class="px-0" :href="cancelUrl">{{ $t('message.cancel') }}</CButton>
                    </CCol>
                  </CRow>
                  <CSelect horizontal class="py-2" :label="$t('message.language')" :value.sync="$i18n.locale" :options="languageOptions"/>
                </CForm>
              </CCardBody>
            </CCard>
            <CCard v-if="app!=null && app.public_attrs!=null" color="primary" text-color="white"
                   class="text-center py-5 d-md-down-none" style="width:44%" body-wrapper>
              <h2>{{ app.id }}</h2>
              <p>{{ app.public_attrs.desc }}</p>
            </CCard>
          </CCardGroup>
        </CCol>
      </CRow>
    </CContainer>
  </div>
</template>

<script>
import clientUtils from "@/utils/api_client"
import utils from "@/utils/app_utils"
import appConfig from "@/utils/app_config"

const initStatusExterInfoFetched = 1
const initStatusAppInfoFetched = 2
const initStatusGoogleSDKInited = 4
const initStatusFacebookSDKInited = 8

export default {
  name: 'Login',
  computed: {
    appId() {
      return this.$route.query.app ? this.$route.query.app : appConfig.APP_ID
    },
    returnUrl() {
      let appId = this.$route.query.app ? this.$route.query.app : appConfig.APP_ID
      let urlDashboard = this.$router.resolve({name: 'Dashboard'}).href
      let returnUrl = this.$route.query.returnUrl ? this.$route.query.returnUrl : ''
      return returnUrl != '' ? returnUrl : (this.app.public_attrs ? this.app.public_attrs.rurl : (appId == appConfig.APP_ID ? urlDashboard : ''))
    },
    cancelUrl() {
      let urlCancelUrl = this.$route.query.cancelUrl ? this.$route.query.cancelUrl : ''
      return urlCancelUrl != '' ? urlCancelUrl : (this.app.public_attrs ? this.app.public_attrs.curl : '')
    },
    exterInfoInited() {
      return (this.initStatus > 0) && ((this.initStatus | initStatusExterInfoFetched) != 0)
    },
    googleSDKInited() {
      return (this.initStatus > 0) && ((this.initStatus | initStatusGoogleSDKInited) != 0)
    },
    facebookSDKInited() {
      return (this.initStatus > 0) && ((this.initStatus | initStatusFacebookSDKInited) != 0)
    },
    languageOptions() {
      let result = []
      this.$i18n.availableLocales.forEach(locale => {
        result.push({value: locale, label: this.$i18n.messages[locale]._name})
      })
      return result
    },
  },
  data() {
    return {
      // -1: error
      // 0: nothing done,
      // 1st bit (1): Exter info fetched,
      // 2nd bit (2): app info fetched,
      // 3rd bit (4): Google SDK inited,
      // 4rd bit (8): Facebook SDK inited,
      initStatus: 0,

      githubClientId: '',
      githubAuthScope: 'user:email', // https://docs.github.com/en/developers/apps/scopes-for-oauth-apps

      googleAuthScope: 'email profile openid',
      googleClientId: '',

      facebookAppId: '',

      linkedinAppId: '',
      linkedinAuthScope: 'r_liteprofile,r_emailaddress', // https://docs.microsoft.com/en-us/linkedin/consumer/integrations/self-serve/sign-in-with-linkedin

      errorMsg: '',
      infoMsg: this.$i18n.t('message.login_msg'),

      app: {},
      sources: {},

      waitCounter: -1,
    }
  },
  mounted() {
    this.initStatus = 0
    this._loadExterAndAppInfo()
    this._loadGoogleSDK()
    this._loadFacebookSDK()

    const callbackAction = this.$route.query.cba
    switch (callbackAction) {
      case 'gh':
        this._doLoginGitHubCallback()
        break
      case 'ln':
        this._doLoginLinkedInCallback()
        break
    }
  },
  methods: {
    _loadAppInfo(appId) {
      const vue = this
      clientUtils.apiDoGet(clientUtils.apiApp.replaceAll(':app', appId),
          (apiRes) => {
            if (apiRes.status == 404) {
              vue._resetOnError(vue.$i18n.t('message.error_app_not_exist', {app: appId}))
            } else if (apiRes.status != 200) {
              vue._resetOnError(apiRes.message)
            } else if (!apiRes.data.public_attrs.actv) {
              vue._resetOnError(vue.$i18n.t('message.error_app_not_active', {app: appId}))
            } else {
              vue.app = apiRes.data
              vue.initStatus |= initStatusAppInfoFetched
              vue._loadExterInfo()
            }
          },
          (err) => {
            vue._resetOnError(err)
          }
      )
    },
    _loadExterInfo() {
      const vue = this
      let appISources = vue.app.public_attrs.sources
      let iSources = {}
      clientUtils.apiDoGet(clientUtils.apiInfo,
          (apiRes) => {
            if (apiRes.status != 200) {
              vue._resetOnError(apiRes.message)
              return
            }
            vue.githubClientId = apiRes.data.github_client_id
            vue.googleClientId = apiRes.data.google_client_id
            vue.facebookAppId = apiRes.data.facebook_app_id
            vue.linkedinAppId = apiRes.data.linkedin_client_id

            apiRes.data.login_channels.every(function (e) {
              iSources[e] = appISources[e]
              return true
            })
            vue.sources = iSources
            vue.initStatus |= initStatusExterInfoFetched
          },
          (err) => {
            vue._resetOnError(err)
          })
    },
    _loadExterAndAppInfo() {
      this._loadAppInfo(this.appId)
    },
    _loadGoogleSDK() {
      const vue = this
      const googleScriptSrc = 'https://apis.google.com/js/platform.js'
      vue.$loadScript(googleScriptSrc)
          .then(() => {
            gapi.load('auth2', () => {
              vue.initStatus |= initStatusGoogleSDKInited
            })
          })
          .catch(() => {
            vue.$unloadScript(googleScriptSrc)
            const msg = vue.$i18n.t('message.error_loading_gapisdk')
            this._resetOnError(msg)
            console.error(msg)
          })
    },
    _loadFacebookSDK() {
      const vue = this
      window.fbAsyncInit = function () {
        if (!vue.exterInfoInited) {
          // console.log("[DEBUG] Waiting for Exter info...")
          setTimeout(() => {
            window.fbAsyncInit()
          }, 500)
          return
        }
        // console.log("[DEBUG] Initializing Fb SDK..."+vue.initStatus)
        FB.init({
          appId: vue.facebookAppId,
          cookie: true,
          xfbml: false,
          version: 'v8.0'
        });
        vue.initStatus |= initStatusFacebookSDKInited
        FB.AppEvents.logPageView()
      }
      const facebookScriptSrc = 'https://connect.facebook.net/en_US/sdk.js'
      vue.$loadScript(facebookScriptSrc)
          .then(() => {
            // console.log("[DEBUG] Fb SDK loaded.")
          })
          .catch(() => {
            vue.$unloadScript(facebookScriptSrc)
            const msg = vue.$i18n.t('message.error_loading_fbsdk')
            vue._resetOnError(msg)
            console.error(msg)
          })
    },
    _doLoginGitHubCallback() {
      const savedState = utils.localStorageGet('ghoa_state')
      utils.localStorageSet("ghoa_state", null)
      const savedApp = utils.localStorageGet('ghoa_app')
      utils.localStorageSet("ghoa_app", null)
      const savedReturnUrl = utils.localStorageGet('ghoa_returnUrl')
      utils.localStorageSet("ghoa_returnUrl", null)
      const urlState = this.$route.query.state
      const code = this.$route.query.code
      if (savedState == "" || savedState == null || urlState != savedState || code == "" || code == null) {
        //login failed
        this._resetOnError(this.$i18n.t('message.error_login_failed_github'), true)
        if (this.$route.query.app == "" || this.$route.query.app == null || this.$route.query.returnUrl == "" || this.$route.query.returnUrl == null) {
          this.$router.push({
            name: "Login",
            query: {returnUrl: savedReturnUrl, app: savedApp, cba: "gh"}
          })
        }
      } else {
        const data = {
          app: this.$route.query.app,
          source: 'github',
          code: code,
          return_url: this.returnUrl,
        }
        this._doLogin(data)
      }
    },
    doLoginGitHub(e) {
      e.preventDefault()
      const state = utils.crc32("" + Math.random())
      utils.localStorageSet('ghoa_state', state)
      /*
      if user rejects the authorization request, GitHub does _not_ redirect user back to redirect_uri,
      which breaks the login flow by losing 'app' and 'returnUrl' parameters.
      Hence we need to save those params first.
      */
      utils.localStorageSet('ghoa_app', this.app.id)
      utils.localStorageSet('ghoa_returnUrl', this.returnUrl)
      const redirectUrl = window.location.origin + this.$router.resolve({
        name: "Login",
        query: {returnUrl: this.returnUrl, app: this.app.id, cba: "gh"}
      }).href
      let githubLoginUrl = "https://github.com/login/oauth/authorize?login=&state=" + state + "&client_id=" + this.githubClientId + "&scope=" + this.githubAuthScope + "&redirect_uri=" + encodeURIComponent(redirectUrl)
      window.location.href = githubLoginUrl
    },
    _doLoginLinkedInCallback() {
      const savedState = utils.localStorageGet('lnoa_state')
      utils.localStorageSet("lnoa_state", null)
      const savedApp = utils.localStorageGet('lnoa_app')
      utils.localStorageSet("lnoa_app", null)
      const savedReturnUrl = utils.localStorageGet('lnoa_returnUrl')
      utils.localStorageSet("lnoa_returnUrl", null)
      const urlState = this.$route.query.state
      const code = this.$route.query.code
      if (savedState == '' || savedState == null || urlState != savedState || code == '' || code == null) {
        //login failed
        this._resetOnError(this.$i18n.t('message.error_login_failed_linkedin'), true)
        if (this.$route.query.app == '' || this.$route.query.app == null || this.$route.query.returnUrl == '' || this.$route.query.returnUrl == null) {
          this.$router.push({
            name: "Login",
            query: {returnUrl: savedReturnUrl, app: savedApp, cba: "ln"}
          })
        }
      } else if (this.$route.query.app == '' || this.$route.query.app == null) {
          //redirect back to the normal login flow
          utils.localStorageSet('lnoa_state', savedState)
          utils.localStorageSet('lnoa_app', savedApp)
          utils.localStorageSet('lnoa_returnUrl', savedReturnUrl)
          window.location.href = this.$router.resolve({
            name: 'Login',
            query: {returnUrl: savedReturnUrl, app: savedApp, cba: 'ln', code: code, state: savedState}
          }).href
        } else {
          const data = {
            app: this.$route.query.app,
            source: 'linkedin',
            code: code,
            return_url: this.returnUrl,
          }
          this._doLogin(data)
        }
    },
    doLoginLinkedIn(e) {
      e.preventDefault()
      const state = utils.crc32("" + Math.random())
      utils.localStorageSet('lnoa_state', state)
      utils.localStorageSet('lnoa_app', this.app.id)
      utils.localStorageSet('lnoa_returnUrl', this.returnUrl)
      const redirectUrl = window.location.origin + this.$router.resolve({
        name: "Login",
        query: {cba: "ln"}
      }).href
      let linkedinLoginUrl = "https://www.linkedin.com/oauth/v2/authorization?response_type=code&state=" + state + "&client_id=" + this.linkedinAppId + "&scope=" + this.linkedinAuthScope + "&redirect_uri=" + encodeURIComponent(redirectUrl)
      window.location.href = linkedinLoginUrl
    },
    doLoginFacebook(e) {
      e.preventDefault()
      this._resetOnError('', true) //since FB popups auth window, clear any existing error message
      if (!this.facebookSDKInited) {
        alert(this.$i18n.t('message.wait_fbsdk'))
      } else {
        const vue = this
        FB.login(function (response) {
          if (response.status == 'connected') {
            const data = {
              app: vue.app.id,
              source: 'facebook',
              code: response.authResponse.accessToken,
              return_url: vue.returnUrl,
            }
            vue._doLogin(data)
          }
        }, {scope: 'public_profile,email', return_scopes: true, auth_type: 'rerequest'});
      }
    },
    doLoginGoogle(e) {
      e.preventDefault()
      this._resetOnError('', true) //since Google popups auth window, clear any existing error message
      if (!this.googleSDKInited) {
        alert(this.$i18n.t('message.wait_googlesdk'))
      } else {
        const vue = this
        gapi.auth2.authorize({
          client_id: this.googleClientId,
          scope: this.googleAuthScope,
          response_type: "code",
          prompt: "consent",
        }, (resp) => {
          // this.infoMsg = defaultInfoMsg
          if (!resp.error) {
            const data = {
              app: this.app.id,
              source: 'google',
              code: resp.code,
              return_url: this.returnUrl,
            }
            this._doLogin(data)
          }
        })
      }
    },
    _waitPreLogin(token, returnUrl) {
      clientUtils.apiDoPost(clientUtils.apiVerifyLoginToken, {
            token: token,
            app: this.app.id,
            return_url: returnUrl
          },
          (apiRes) => {
            if (300 <= apiRes.status && apiRes.status <= 399) {
              setTimeout(() => {
                this._waitPreLogin(token, returnUrl)
              }, 2000)
            } else if (apiRes.status != 200) {
              this._resetOnError(apiRes.message)
            } else {
              this._doSaveLoginSessionAndLogin(apiRes.data, apiRes.extras.return_url)
            }
          },
          (err) => {
            const msg = "Session verification error, retry in 2 seconds: " + err
            console.error(msg)
            this.errorMsg = msg
            setTimeout(() => {
              this._waitPreLogin(token, returnUrl)
            }, 2000)
          })
    },
    _doSaveLoginSessionAndLogin(token, returnUrl) {
      // this.waitCounter = -1
      if (returnUrl == null || returnUrl == "" || returnUrl == '#') {
        if (this.app.id != appConfig.APP_ID) {
          this.errorMsg = this.$i18n.t('message.error_invalid_return_url')
          return
        }
        returnUrl = this.$router.resolve({name: 'Dashboard'}).href
      }

      // generate and save session token
      const jwt = utils.parseJwt(token)
      utils.saveLoginSession({uid: jwt.payloadObj.uid, name: jwt.payloadObj.name, token: token})

      // redirect to next url
      window.location.href = returnUrl
    },
    _doWaitMessage() {
      if (this.waitCounter >= 0) {
        this.waitCounter++
        this.infoMsg = this.$i18n.t('message.wait_login', {counter: this.waitCounter})
        setTimeout(() => {
          this._doWaitMessage()
        }, 2000)
      }
    },
    _doLogin(data) {
      this.waitCounter = 0
      this._doWaitMessage()
      clientUtils.apiDoPost(
          clientUtils.apiLogin, data,
          (apiRes) => {
            if (apiRes.status != 200) {
              this._resetOnError(apiRes.status + ": " + apiRes.message, true)
            } else {
              const jwt = utils.parseJwt(apiRes.data)
              if (jwt.payloadObj.type == "pre_login") {
                this._waitPreLogin(apiRes.data, this.returnUrl)
              } else {
                this._doSaveLoginSessionAndLogin(apiRes.data, apiRes.extras.return_url)
              }
            }
          },
          (err) => {
            this._resetOnError(err, true)
          }
      )
    },
    _resetOnError(err, preserveInitStatus) {
      this.initStatus = preserveInitStatus ? this.initStatus : -1
      this.errorMsg = err
      this.infoMsg = this.$i18n.t('message.login_msg')
      this.waitCounter = -1
    },
  }
}
</script>
