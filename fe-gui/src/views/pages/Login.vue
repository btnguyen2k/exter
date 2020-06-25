<template>
    <CContainer class="d-flex align-items-center min-vh-100">
        <CRow class="justify-content-center">
            <CCol md="8">
                <CCardGroup>
                    <CCard class="p-4">
                        <CCardBody>
                            <h1>Login</h1>
                            <p v-if="errorMsg!=''" class="alert alert-danger">{{errorMsg}}</p>
                            <CForm method="post">
                                <p v-if="infoMsg!=''" class="text-muted">{{infoMsg}}</p>
                                <CButton v-if="sources.facebook" id="loginFb" type="button" name="facebook"
                                         color="facebook" class="mb-1" style="width: 100%"
                                         @click="doLoginFacebook">
                                    <CIcon name="cib-facebook"/>
                                    Login with Facebook
                                </CButton>
                                <CButton v-if="sources.google" id="loginGoogle" type="button" name="google"
                                         color="light" class="mb-1" style="width: 100%"
                                         @click="doLoginGoogle">
                                    <CIcon name="cib-google"/>
                                    Login with Google
                                </CButton>

                                <!--                                <CInput placeholder="Username" autocomplete="username email" name="username"-->
                                <!--                                        id="username" v-model="form.username">-->
                                <!--                                    <template #prepend-content>-->
                                <!--                                        <CIcon name="cil-user"/>-->
                                <!--                                    </template>-->
                                <!--                                </CInput>-->
                                <!--                                <CInput placeholder="Password" type="password" autocomplete="current-password"-->
                                <!--                                        name="password" id="password" v-model="form.password">-->
                                <!--                                    <template #prepend-content>-->
                                <!--                                        <CIcon name="cil-lock-locked"/>-->
                                <!--                                    </template>-->
                                <!--                                </CInput>-->
                                <!--                                <CRow>-->
                                <!--                                    <CCol col="6">-->
                                <!--                                        <CButton color="primary" class="px-4" type="submit">-->
                                <!--                                            Login-->
                                <!--                                        </CButton>-->
                                <!--                                    </CCol>-->
                                <!--                                    <CCol col="6" class="text-right">-->
                                <!--                                        <CButton color="link" class="px-0" @click="funcNotImplemented">Forgot password?-->
                                <!--                                        </CButton>-->
                                <!--                                    </CCol>-->
                                <!--                                </CRow>-->
                            </CForm>
                        </CCardBody>
                    </CCard>
                    <CCard color="primary" text-color="white" class="text-center py-5 d-md-down-none" style="width:44%"
                           body-wrapper>
                        <h2>Demo</h2>
                        <p>This is instance is for demo purpose only. Login with administrator account <strong>admin/s3cr3t</strong>.
                            You can create/edit/delete other user group or user account. This special admin account,
                            however, can not be modified or deleted.</p>
                        <!--
                        <CButton color="primary" class="active mt-3" :to="pageRegister">
                            Register Now!
                        </CButton>
                        -->
                    </CCard>
                </CCardGroup>
            </CCol>
        </CRow>
    </CContainer>
</template>

<script>
    import Register from "@/views/pages/Register"
    import clientUtils from "@/utils/api_client"
    import utils from "@/utils/app_utils"
    import appConfig from "@/utils/app_config"

    const defaultInfoMsg = "Please sign in to continue"
    const waitInfoMsg = "Please wait..."
    const waitLoginInfoMsg = "Logging in, please wait..."

    export default {
        name: 'Login',
        data() {
            this.infoMsg = waitInfoMsg
            clientUtils.apiDoGet(clientUtils.apiApp + "/" + this.$route.query.app,
                (apiRes) => {
                    if (apiRes.status != 200) {
                        this.errorMsg = apiRes.message
                    } else {
                        this.sources = apiRes.data.config.sources
                        this.infoMsg = defaultInfoMsg
                    }
                },
                (err) => {
                    this.errorMsg = err
                })
            clientUtils.apiDoGet("/info",
                (apiRes) => {
                    if (apiRes.status != 200) {
                        this.errorMsg = apiRes.message
                    } else {
                        this.googleClientId = apiRes.data.google_client_id
                        this.rsaPublicKeyPEM = apiRes.data.rsa_public_key
                    }
                },
                (err) => {
                    this.errorMsg = err
                })
            return {
                googleInited: false,
                googleAuthScope: 'email profile openid',
                googleClientId: '',

                rsaPublicKeyPEM: '',

                returnUrl: this.$route.query.returnUrl ? this.$route.query.returnUrl : "/",
                app: this.$route.query.app ? this.$route.query.app : appConfig.APP_NAME,
                sources: {},
                pageRegister: Register,
                errorMsg: "",
                infoMsg: defaultInfoMsg,
            }
        },
        mounted() {
            const vue = this
            const scriptSrc = 'https://apis.google.com/js/platform.js'
            vue.$loadScript(scriptSrc)
                .then(() => {
                    gapi.load('auth2', () => {
                        //use gapi.auth2.init together with gapi.auth2.getAuthInstance().signIn or gapi.auth2.getAuthInstance().grantOfflineAccess
                        //if to use gapi.auth2.authorize, do NOT call gapi.auth2.init
                        // gapi.auth2.init({
                        //     client_id: vue.googleClientId,
                        //     scope: vue.googleAuthScope,
                        // }).then(
                        //     () => {
                        //         vue.googleInited = true
                        //     },
                        //     () => {
                        //         vue.$unloadScript(scriptSrc)
                        //         vue.errorMsg = "Error while initializing Google SDK"
                        //     }
                        // )
                        vue.googleInited = true
                    })
                })
                .catch(() => {
                    vue.$unloadScript(scriptSrc)
                    vue.errorMsg = 'Error loading GoogleApi SDK'
                })
        },
        methods: {
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
                    // gapi.auth2.getAuthInstance().grantOfflineAccess({
                    //     prompt: "consent",
                    //     scope: this.googleAuthScope,
                    // }).then(
                    //     (resp) => {
                    //         if (!resp.error) {
                    //             // const data = {
                    //             //     source: 'google',
                    //             //     code: resp.code,
                    //             // }
                    //             // this._doLogin(data)
                    //         }
                    //     }
                    // )
                    gapi.auth2.authorize({
                        client_id: this.googleClientId,
                        scope: this.googleAuthScope,
                        response_type: "code",
                        prompt: "consent",
                    }, (resp) => {
                        this.infoMsg = defaultInfoMsg
                        if (!resp.error) {
                            const data = {
                                source: 'google',
                                code: resp.code,
                            }
                            this._doLogin(data)
                        }
                    })
                    // this.googleAuth.signIn({scope: this.googleAuthScope, prompt: 'consent'})
                    //     .then((auth) => {
                    //         const data = {
                    //             source: 'google',
                    //             id_token: auth.getAuthResponse().id_token,
                    //             access_token: auth.getAuthResponse().access_token,
                    //             email: auth.getBasicProfile().getEmail(),
                    //         }
                    //         this._doLogin(data)
                    //     })
                    //     .catch((err) => {
                    //         //this.errorMsg = 'Error while logging with Google account: ' + err
                    //     })
                }
            },
            _waitPreLogin(token) {
                this.infoMsg = waitInfoMsg
                clientUtils.apiDoPost(clientUtils.apiCheckLoginToken, {token: token},
                    (apiRes) => {
                        if (300 <= apiRes.status && apiRes.status <= 399) {
                            // console.log("Server is creating login session: " + JSON.stringify(apiRes))
                            this.infoMsg = waitLoginInfoMsg
                            setTimeout(() => {
                                this._waitPreLogin(token)
                            }, 2000)
                        } else if (apiRes.status != 200) {
                            this.errorMsg = apiRes.message
                        } else {
                            this._doSaveLoginSessionAndLogin(apiRes.data)
                        }
                    },
                    (err) => {
                        const msg = "Session verification error, retry in 2 seconds: " + err
                        console.error(msg)
                        this.errorMsg = msg
                        setTimeout(() => {
                            this._waitPreLogin(token)
                        }, 2000)
                    })
            },
            _doSaveLoginSessionAndLogin(token) {
                const jwt = utils.parseJwt(token)
                utils.saveLoginSession({uid: jwt.payloadObj.uid, token: token})
                this.$router.push(this.returnUrl != "" ? this.returnUrl : "/")
            },
            _doLogin(data) {
                this.infoMsg = waitInfoMsg
                clientUtils.apiDoPost(
                    clientUtils.apiLogin, data,
                    (apiRes) => {
                        if (apiRes.status != 200) {
                            this.errorMsg = apiRes.status + ": " + apiRes.message
                        } else {
                            const jwt = utils.parseJwt(apiRes.data)
                            if (jwt.payloadObj.type == "pre_login") {
                                this._waitPreLogin(apiRes.data)
                            } else {
                                this._doSaveLoginSessionAndLogin(apiRes.data)
                            }
                        }
                    },
                    (err) => {
                        this.errorMsg = err
                    }
                )
            },
        }
    }
</script>
