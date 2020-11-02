<template>
  <div class="c-app flex-row align-items-center">
    <CContainer>
      <CRow class="justify-content-center">
        <CCol md="8">
          <CCardGroup>
            <CCard class="p-4">
              <CCardBody>
                <h1>Login</h1>
                <p v-if="errorMsg!=''" class="alert alert-danger">{{ errorMsg }}</p>
                <CForm method="post">
                  <p v-if="infoMsg!=''" class="text-muted">{{ infoMsg }}</p>
                  <CButton v-if="sources.facebook" id="loginFb" type="button" name="facebook"
                           color="facebook" class="mb-1" block
                           @click="doLoginFacebook">
                    <CIcon name="cib-facebook"/>
                    Login with Facebook
                  </CButton>
                  <CButton v-if="sources.github" id="loginGithub" type="button" name="github"
                           color="github" class="mb-1" block
                           @click="doLoginGitHub">
                    <CIcon name="cib-github"/>
                    Login with GitHub
                  </CButton>
                  <CButton v-if="sources.google" id="loginGoogle" type="button" name="google"
                           color="light" class="mb-1" block
                           @click="doLoginGoogle">
                    <CIcon name="cib-google"/>
                    Login with Google
                  </CButton>
                  <CButton v-if="sources.linkedin" id="loginLinkedIn" type="button" name="linkedin"
                           color="linkedin" class="mb-1" block
                           @click="doLoginLinkedIn">
                    <CIcon name="cib-linkedin"/>
                    Login with LinkedIn
                  </CButton>
                  <CRow v-if="cancelUrl!=''">
                    <CCol col="12" class="text-right">
                      <CButton color="link" class="px-0" :href="cancelUrl">Cancel</CButton>
                    </CCol>
                  </CRow>
                </CForm>
              </CCardBody>
            </CCard>
            <CCard v-if="app!=null && app.config!=null" color="primary" text-color="white"
                   class="text-center py-5 d-md-down-none" style="width:44%" body-wrapper>
              <h2>{{ app.id }}</h2>
              <p>{{ app.config.desc }}</p>
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
import router from "@/router";

const defaultInfoMsg = "Please sign in to continue"
const waitInfoMsg = "Please wait..."
const waitLoginInfoMsg = "Logging in, please wait..."
const invalidReturnUrlErrMsg = "Error: invalid return url"

export default {
  name: 'Login',
  computed: {
    returnUrl() {
      let appId = this.$route.query.app ? this.$route.query.app : appConfig.APP_ID
      let urlDashboard = this.$router.resolve({name: 'Dashboard'}).href
      let returnUrl = this.$route.query.returnUrl ? this.$route.query.returnUrl : ''
      return returnUrl != '' ? returnUrl : (this.app.config ? this.app.config.rurl : (appId == appConfig.APP_ID ? urlDashboard : ''))
    },
    cancelUrl() {
      let urlCancelUrl = this.$route.query.cancelUrl ? this.$route.query.cancelUrl : ''
      return urlCancelUrl != '' ? urlCancelUrl : (this.app.config ? this.app.config.curl : '')
    },
  },
  data() {
    this.infoMsg = waitInfoMsg
    let appId = this.$route.query.app ? this.$route.query.app : appConfig.APP_ID
    clientUtils.apiDoGet(clientUtils.apiApp + "/" + appId,
        (apiRes) => {
          if (apiRes.status != 200) {
            this._resetOnError(apiRes.message)
            return
          }
          this.app = apiRes.data
          this.appInited = true
          if (!this.app.config.actv) {
            this._resetOnError("App [" + appId + "] is not active")
            return
          }
          let appISources = this.app.config.sources
          let iSources = {}
          clientUtils.apiDoGet(clientUtils.apiInfo,
              (apiRes) => {
                if (apiRes.status != 200) {
                  this._resetOnError(apiRes.message)
                  return
                }
                this.githubClientId = apiRes.data.github_client_id
                this.googleClientId = apiRes.data.google_client_id
                this.facebookAppId = apiRes.data.facebook_app_id
                this.linkedinAppId = apiRes.data.linkedin_client_id

                apiRes.data.login_channels.every(function (e) {
                  iSources[e] = appISources[e]
                  return true
                })
                this.sources = iSources
                this.exterInfoInited = true
                this.infoMsg = defaultInfoMsg
              },
              (err) => {
                this.errorMsg = err
              })
        },
        (err) => {
          this.errorMsg = err
        })
    return {
      exterInfoInited: false,
      appInited: false,

      githubClientId: '',
      //https://docs.github.com/en/developers/apps/scopes-for-oauth-apps
      githubAuthScope: 'user:email',

      googleInited: false,
      googleAuthScope: 'email profile openid',
      googleClientId: '',

      facebookInited: false,
      facebookAppId: '',

      linkedinAppId: '',
      //https://docs.microsoft.com/en-us/linkedin/consumer/integrations/self-serve/sign-in-with-linkedin
      linkedinAuthScope: 'r_liteprofile,r_emailaddress',

      errorMsg: '',
      infoMsg: defaultInfoMsg,

      app: {},
      sources: {},

      waitCounter: -1,
    }
  },
  mounted() {
    const callbackAction = this.$route.query.cba
    switch (callbackAction) {
      case 'gh':
        this._doLoginGitHubCallback()
      case 'ln':
        this._doLoginLinkedInCallback()
    }

    const vue = this

    const googleScriptSrc = 'https://apis.google.com/js/platform.js'
    vue.$loadScript(googleScriptSrc)
        .then(() => {
          gapi.load('auth2', () => {
            vue.googleInited = true
          })
        })
        .catch(() => {
          vue.$unloadScript(googleScriptSrc)
          const msg = 'Error loading GoogleApi SDK'
          vue.errorMsg += '<br>' + msg
          console.error(msg)
        })

    window.fbAsyncInit = function () {
      if (!vue.exterInfoInited) {
        setTimeout(() => {
          window.fbAsyncInit()
        }, 1000)
        return
      }
      FB.init({
        appId: vue.facebookAppId,
        cookie: true,
        xfbml: false,
        version: 'v8.0'
      });
      vue.facebookInited = true
      FB.AppEvents.logPageView();
    }
    const facebookScriptSrc = 'https://connect.facebook.net/en_US/sdk.js'
    vue.$loadScript(facebookScriptSrc)
        .then(() => {
        })
        .catch(() => {
          vue.$unloadScript(facebookScriptSrc)
          const msg = 'Error loading Facebook SDK'
          vue.errorMsg += '<br>' + msg
          console.error(msg)
        })
  },
  methods: {
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
        this._resetOnError('GitHub login failed.')
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
        this._resetOnError('LinkedIn login failed.')
        if (this.$route.query.app == '' || this.$route.query.app == null || this.$route.query.returnUrl == '' || this.$route.query.returnUrl == null) {
          this.$router.push({
            name: "Login",
            query: {returnUrl: savedReturnUrl, app: savedApp, cba: "ln"}
          })
        }
      } else {
        if (this.$route.query.app == '' || this.$route.query.app == null) {
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
      if (!this.facebookInited) {
        alert('Please wait, Facebook SDK is being loaded.')
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
      if (!this.googleInited) {
        alert('Please wait, Google SDK is being loaded.')
      } else {
        this.infoMsg = waitInfoMsg
        gapi.auth2.authorize({
          client_id: this.googleClientId,
          scope: this.googleAuthScope,
          response_type: "code",
          prompt: "consent",
        }, (resp) => {
          this.infoMsg = defaultInfoMsg
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
      this.waitCounter = -1
      if (returnUrl == null || returnUrl == "" || returnUrl == '#') {
        if (this.app.id != appConfig.APP_ID) {
          this.errorMsg = invalidReturnUrlErrMsg
          return
        } else {
          returnUrl = this.$router.resolve({name: 'Dashboard'}).href
        }
      }
      const jwt = utils.parseJwt(token)
      utils.saveLoginSession({uid: jwt.payloadObj.uid, name: jwt.payloadObj.name, token: token})
      window.location.href = returnUrl
    },
    _doWaitMessage() {
      if (this.waitCounter >= 0) {
        this.waitCounter++
        this.infoMsg = waitLoginInfoMsg + " " + this.waitCounter
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
              this._resetOnError(apiRes.status + ": " + apiRes.message)
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
            this._resetOnError(err)
          }
      )
    },
    _resetOnError(err) {
      this.errorMsg = err
      this.infoMsg = defaultInfoMsg
      this.waitCounter = -1
    },
  }
}
</script>
