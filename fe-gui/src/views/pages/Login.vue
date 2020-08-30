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
                           color="facebook" class="mb-1" style="width: 100%"
                           @click="doLoginFacebook">
                    <CIcon name="cib-facebook"/>
                    Login with Facebook
                  </CButton>
                  <CButton v-if="sources.github" id="loginGithub" type="button" name="github"
                           color="light" class="mb-1" style="width: 100%"
                           @click="doLoginGitHub">
                    <CIcon name="cib-github"/>
                    Login with GitHub
                  </CButton>
                  <CButton v-if="sources.google" id="loginGoogle" type="button" name="google"
                           color="light" class="mb-1" style="width: 100%"
                           @click="doLoginGoogle">
                    <CIcon name="cib-google"/>
                    Login with Google
                  </CButton>
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

export default {
  name: 'Login',
  data() {
    this.infoMsg = waitInfoMsg
    let appId = this.$route.query.app ? this.$route.query.app : appConfig.APP_ID
    clientUtils.apiDoGet(clientUtils.apiApp + "/" + appId,
        (apiRes) => {
          if (apiRes.status != 200) {
            this.errorMsg = apiRes.message
            this.infoMsg = ""
          } else {
            this.app = apiRes.data
            if (!this.app.config.actv) {
              this.errorMsg = "App [" + appId + "] is not active"
              this.infoMsg = ""
              return
            }
            this.infoMsg = defaultInfoMsg
            let appISources = this.app.config.sources
            let iSources = {}
            clientUtils.apiDoGet(clientUtils.apiInfo,
                (apiRes) => {
                  if (apiRes.status != 200) {
                    this.errorMsg = apiRes.message
                  } else {
                    this.githubClientId = apiRes.data.github_client_id
                    this.googleClientId = apiRes.data.google_client_id

                    apiRes.data.login_channels.every(function (e) {
                      iSources[e] = appISources[e]
                      return true
                    })
                    this.sources = iSources
                  }
                },
                (err) => {
                  this.errorMsg = err
                })
          }
        },
        (err) => {
          this.errorMsg = err
        })
    return {
      githubClientId: '',
      //https://docs.github.com/en/developers/apps/scopes-for-oauth-apps
      githubAuthScope: 'user:email',

      googleInited: false,
      googleAuthScope: 'email profile openid',
      googleClientId: '',

      returnUrl: this.$route.query.returnUrl ? this.$route.query.returnUrl : "/",
      errorMsg: "",
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
        return this._doLoginGitHubCallback()
    }

    const vue = this
    const scriptSrc = 'https://apis.google.com/js/platform.js'
    vue.$loadScript(scriptSrc)
        .then(() => {
          gapi.load('auth2', () => {
            vue.googleInited = true
          })
        })
        .catch(() => {
          vue.$unloadScript(scriptSrc)
          vue.errorMsg = 'Error loading GoogleApi SDK'
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
        this.errorMsg = 'GitHub login failed.'
        this.infoMsg = defaultInfoMsg
        this.waitCounter = -1
        if (this.$route.query.app == "" || this.$route.query.app == null || this.$route.query.returnUrl == "" || this.$route.query.returnUrl == null) {
          this.$router.push({
            name: "Login",
            query: {returnUrl: savedReturnUrl, app: savedApp, cba: "gh"}
          })
        }
        return false
      }
      const data = {
        app: this.$route.query.app,
        source: 'github',
        code: code,
        return_url: this.returnUrl,
      }
      this._doLogin(data)
      return true
    },
    doLoginGitHub(e) {
      e.preventDefault()
      const state = utils.crc32("" + Math.random())
      utils.localStorageSet('ghoa_state', state)
      /*
      if user rejects the authorization request, GitHug does _not_ redirect user back to redirect_uri,
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
      //console.log(githubLoginUrl)
    },
    doLoginFacebook(e) {
      e.preventDefault()
      alert('Please wait, Facebook SDK is being loaded.')
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
              this.errorMsg = apiRes.message
              this.infoMsg = defaultInfoMsg
              this.waitCounter = -1
            } else {
              let returnUrl = apiRes.extras.return_url
              this._doSaveLoginSessionAndLogin(apiRes.data, returnUrl)
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
      const jwt = utils.parseJwt(token)
      utils.saveLoginSession({uid: jwt.payloadObj.uid, token: token})
      // console.log(returnUrl)
      window.location.href = returnUrl != "" ? returnUrl : "/"
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
              this.errorMsg = apiRes.status + ": " + apiRes.message
              this.infoMsg = defaultInfoMsg
              this.waitCounter = -1
            } else {
              let returnUrl = apiRes.extras.return_url
              const jwt = utils.parseJwt(apiRes.data)
              if (jwt.payloadObj.type == "pre_login") {
                this._waitPreLogin(apiRes.data, returnUrl)
              } else {
                this._doSaveLoginSessionAndLogin(apiRes.data, returnUrl)
              }
            }
          },
          (err) => {
            this.errorMsg = err
            this.infoMsg = defaultInfoMsg
            this.waitCounter = -1
          }
      )
    },
  }
}
</script>
