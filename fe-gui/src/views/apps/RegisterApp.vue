<template>
    <div>
        <CRow>
            <CCol sm="12">
                <CCard>
                    <CCardHeader>Register New App</CCardHeader>
                    <CForm @submit.prevent="doSubmit" method="post">
                        <CCardBody>
                            <p v-if="errorMsg!=''" class="alert alert-danger">{{errorMsg}}</p>
                            <div class="form-group form-row">
                                <CCol :sm="{offset:0,size:9}" class="form-inline">
                                    <CInputCheckbox inline label="Active" :checked.sync="form.isActive"
                                    />
                                </CCol>
                            </div>
                            <div class="form-group form-row">
                                <CCol tag="label" sm="3" class="col-form-label">
                                    Login channels
                                </CCol>
                                <CCol sm="9" class="form-inline">
                                    <CInputCheckbox inline v-for="option in loginChannelList" :label="option"
                                                    :value="option" :checked.sync="form.idSources[option]"
                                    />
                                </CCol>
                            </div>
                            <CInput horizontal type="text" v-model="form.id" label="Id"
                                    placeholder="Application's unique id"
                                    :is-valid="validatorAppId"
                                    invalid-feedback="Enter application's id, format [0-9a-z_]+, must be unique."
                                    valid-feedback=""
                            />
                            <CInput horizontal type="text" v-model="form.description" label="Description"
                                    placeholder="Application's description"
                            />
                            <CInput horizontal type="text" v-model="form.defaultReturnUrl" label="Default return URL"
                                    placeholder="http://..."
                                    :is-valid="validatorReturnUrl"
                                    invalid-feedback="Return url must be a http or https link."
                            />
                            <CInput horizontal type="text" v-model="form.tags" label="Tags"
                                    placeholder="Tags separated by comma"
                            />
                            <CTextarea horizontal type="text" v-model="form.rsaPublicKey" label="RSA public key"
                                       rows="6" placeholder="RSA public key in PEM format"
                            />
                        </CCardBody>
                        <CCardFooter>
                            <CButton type="submit" color="primary" style="width: 96px">
                                <CIcon name="cil-save"/>
                                Save
                            </CButton>
                            <CButton type="button" color="info" class="ml-2" style="width: 96px" @click="doCancel">
                                <CIcon name="cil-arrow-circle-left"/>
                                Back
                            </CButton>
                        </CCardFooter>
                    </CForm>
                </CCard>
            </CCol>
        </CRow>
    </div>
</template>

<script>
    import router from "@/router"
    import clientUtils from "@/utils/api_client";

    let patternAppId = /^[0-9a-z_]+$/
    let patternReturnUrl = /^http(s?):\/\//

    export default {
        name: 'RegisterApp',
        data() {
            let loginChannelList = []
            let form = {
                isActive: true,
                id: "", description: "", rsaPublicKey: "", defaultReturnUrl: "",
                tags: "",
                idSources: {},
            }
            clientUtils.apiDoGet(clientUtils.apiInfo,
                (apiRes) => {
                    if (apiRes.status == 200) {
                        apiRes.data.login_channels.every(function (e) {
                            loginChannelList.push(e)
                            form.idSources[e] = true
                            return true
                        })
                    } else {
                        console.error("Getting info was unsuccessful: " + apiRes)
                    }
                },
                (err) => {
                    console.error("Error getting info list: " + err)
                })
            return {
                form: form,
                errorMsg: "",
                loginChannelList: loginChannelList,
            }
        },
        methods: {
            doCancel() {
                router.push("/apps")
            },
            doSubmit(e) {
                e.preventDefault()
                let data = {
                    is_active: this.form.isActive,
                    id: this.form.id, description: this.form.description,
                    default_return_url: this.form.defaultReturnUrl,
                    rsa_public_key: this.form.rsaPublicKey,
                    tags: this.form.tags,
                    id_sources: this.form.idSources,
                }
                clientUtils.apiDoPost(
                    clientUtils.apiMyAppList, data,
                    (apiRes) => {
                        if (apiRes.status != 200) {
                            this.errorMsg = apiRes.status + ": " + apiRes.message
                        } else {
                            this.$router.push({
                                name: "Apps",
                                params: {flashMsg: "Application [" + this.form.id + "] has been registered successfully."},
                            })
                        }
                    },
                    (err) => {
                        console.error(err)
                        this.errorMsg = err
                    }
                )
            },
            validatorAppId(val) {
                return val ? patternAppId.test(val.toString()) : false
            },
            validatorReturnUrl(val) {
                return val ? patternReturnUrl.test(val.toString().trim()) : true
            },
        }
    }
</script>
