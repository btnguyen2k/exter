<template>
    <CContainer class="d-flex align-items-center min-vh-100">
        <CRow class="justify-content-center">
            <CCol md="8">
                <CCardGroup>
                    <CCard class="p-4">
                        <CCardBody>
                            <h1>Login</h1>
                            <p v-if="errorMsg!=''" class="alert alert-danger">{{errorMsg}}</p>
                            <CForm v-if="errorMsg==''" @submit.prevent="doSubmit" method="post">
                                <p class="text-muted">Please sign in to continue</p>
                                <CButton v-if="sources.facebook" id="loginFb" type="button" name="facebook"
                                         color="facebook" class="mb-1" style="width: 100%" @click="funcNotImplemented">
                                    <CIcon name="cib-facebook"/>
                                    Login with Facebook
                                </CButton>
                                <CButton v-if="sources.google" id="loginGoogle" type="button" name="google"
                                         color="light" class="mb-1" style="width: 100%">
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

    export default {
        name: 'Login',
        data() {
            clientUtils.apiDoGet(clientUtils.apiApp + "/" + this.$route.query.app,
                (apiRes) => {
                    if (apiRes.status != 200) {
                        this.errorMsg = apiRes.message
                    } else {
                        this.sources = apiRes.data.config.sources
                    }
                },
                (err) => {
                    this.errorMsg = err
                })
            return {
                returnUrl: this.$route.query.returnUrl ? this.$route.query.returnUrl : "/",
                app: this.$route.query.app ? this.$route.query.app : appConfig.APP_NAME,
                sources: {},
                pageRegister: Register,
                form: {username: "", password: ""},
                errorMsg: "",
            }
        },
        mounted() {
            const scriptSrc = 'https://apis.google.com/js/platform.js'
            this.$loadScript(scriptSrc)
                .then(() => {
                    gapi.load('auth2', () => {
                        gapi.auth2.init({
                            client_id: '334322862548-9o5rr6edh0fi64vf1km0i2omtpfno1ph.apps.googleusercontent.com',
                            scope: 'email profile openid',
                            cookiepolicy: 'single_host_origin',
                        }).then((auth) => {
                                document.addEventListener('click', (e) => {
                                    let el = document.getElementById('loginGoogle')
                                    if (el != null && e.target == el) {
                                        auth.signIn({
                                            scope: 'email profile openid',
                                            prompt: 'select_account',
                                        }).then((gu) => {
                                            let data = {
                                                token: gu.getAuthResponse().id_token,
                                                email: gu.getBasicProfile().getEmail(),
                                            }
                                            console.log(gu)
                                            console.log(data)
                                        }).catch((err) => {
                                            //console.error(err)
                                        })
                                    }
                                })
                            },
                            () => {
                                this.$unloadScript(scriptSrc)
                                this.errorMsg = 'Error initializing GoogleApi auth2'
                            }
                        )
                    })
                })
                .catch(() => {
                    this.$unloadScript(scriptSrc)
                    this.errorMsg = 'Error loading GoogleApi Javascript'
                })
        },
        methods: {
            funcNotImplemented() {
                console.log(loginGoogle)
            },
            doSubmit(e) {
                e.preventDefault()
                let data = {username: this.form.username, password: this.form.password}
                clientUtils.apiDoPost(
                    clientUtils.apiLogin, data,
                    (apiRes) => {
                        if (apiRes.status != 200) {
                            this.errorMsg = apiRes.status + ": " + apiRes.message
                        } else {
                            utils.saveLoginSession({
                                uid: apiRes.data.uid,
                                token: apiRes.data.token,
                                expiry: apiRes.data.expiry,
                            })
                            this.$router.push(this.returnUrl != "" ? this.returnUrl : "/")
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
