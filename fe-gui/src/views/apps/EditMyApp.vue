<template>
  <div>
    <CRow>
      <CCol sm="12">
        <CCard accent-color="info">
          <CCardHeader>Edit Application</CCardHeader>
          <CForm @submit.prevent="doSubmit" method="post">
            <CCardBody v-if="foundStatus<0">
              <CAlert color="info">Please wait...</CAlert>
            </CCardBody>
            <CCardBody v-if="foundStatus==0">
              <CAlert color="danger">Application [{{ this.$route.params.id }}] not found</CAlert>
            </CCardBody>
            <CCardBody v-if="foundStatus>0">
              <CAlert v-if="errorMsg" color="danger">{{ errorMsg }}</CAlert>
              <CInput horizontal type="text" v-model="app.id" label="Id"
                      placeholder="Application's unique id"
                      readonly="readonly"
              />
              <div class="form-group form-row">
                <CCol :sm="{offset:3,size:9}" class="form-inline">
                  <CInputCheckbox inline label="Active" :checked.sync="app.isActive"
                  />
                </CCol>
              </div>
              <div class="form-group form-row">
                <CCol tag="label" sm="3" class="col-form-label">
                  Login channels
                </CCol>
                <CCol sm="9" class="form-inline">
                  <CInputCheckbox inline v-for="option in loginChannelList" :label="option"
                                  :value="option" :checked.sync="app.idSources[option]"
                  />
                </CCol>
              </div>
              <CInput horizontal type="text" v-model="app.description" label="Description"
                      placeholder="Application's description"
              />
              <CInput horizontal type="text" v-model="app.defaultReturnUrl" label="Default return URL"
                      placeholder="http://..."
                      :is-valid="validatorUrl"
                      invalid-feedback="Return url must be a http or https link."
              />
              <CInput horizontal type="text" v-model="app.defaultCancelUrl" label="Default cancel URL"
                      placeholder="http://..."
                      :is-valid="validatorUrl"
                      invalid-feedback="Cancel url must be a http or https link."
              />
              <CInput horizontal type="text" v-model="app.domains" label="Whitelist domains"
                      placeholder="Exter redirects users to only these whitelist domains. Domains separated by spaces, commas or semi-colons"
              />
              <CInput horizontal type="text" v-model="app.tags" label="Tags"
                      placeholder="Tags separated by comma"
              />
              <CTextarea horizontal type="text" v-model="app.rsaPublicKey" label="RSA public key"
                         rows="6" placeholder="RSA public key in PEM format"
              />
            </CCardBody>
            <CCardFooter>
              <CButton v-if="foundStatus>0" type="submit" color="primary" style="width: 96px">
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
import clientUtils from "@/utils/api_client"

let patternUrl = /^http(s?):\/\//

// const sleepSync = (ms) => {
//   const end = new Date().getTime() + ms;
//   while (new Date().getTime() < end) { /* do nothing */ }
// }

export default {
  name: 'EditMyApp',
  mounted() {
    this.loadApp(this.$route.params.id)
  },
  data() {
    return {
      app: {},
      errorMsg: "",
      loginChannelList: [],
      foundStatus: -1,
    }
  },
  methods: {
    loadApp(appId) {
      this.foundStatus = -1
      const vue = this
      clientUtils.apiDoGet(clientUtils.apiInfo,
          (apiRes) => {
            if (apiRes.status == 200) {
              let _loginChannels = []
              apiRes.data.login_channels.every(function (e) {
                _loginChannels.push(e)
                return true
              })
              vue.loadLoginChannels = _loginChannels

              const apiUrl = clientUtils.apiMyApp.replaceAll(':app', appId)
              clientUtils.apiDoGet(apiUrl,
                  (apiRes) => {
                    vue.foundStatus = apiRes.status == 200 ? 1 : 0
                    if (vue.foundStatus == 1) {
                      let _app = {}
                      _app.id = apiRes.data.id
                      _app.isActive = apiRes.data.public_attrs.actv
                      _app.description = apiRes.data.public_attrs.desc
                      _app.rsaPublicKey = apiRes.data.public_attrs.rpub
                      _app.defaultReturnUrl = apiRes.data.public_attrs.rurl
                      _app.defaultCancelUrl = apiRes.data.public_attrs.curl
                      _app.idSources = apiRes.data.public_attrs.sources
                      _app.domains = apiRes.data.domains != null ? apiRes.data.domains.join(", ") : ""
                      _app.tags = apiRes.data.public_attrs.tags != null ? apiRes.data.public_attrs.tags.join(", ") : ""
                      vue.app = _app
                    }
                  },
                  (err) => {
                    vue.errorMsg = "Error calling API getting application info: " + err
                  })
            } else {
              vue.errorMsg = apiRes.status + ": " + apiRes.message
            }
          },
          (err) => {
            vue.errorMsg = "Error calling API getting Exter info: " + err
          })
    },
    doCancel() {
      router.push(router.resolve({name: "MyApps"}).location)
    },
    doSubmit(e) {
      e.preventDefault()
      this.foundStatus = -1
      const vue = this
      let data = {
        is_active: vue.app.isActive,
        id: vue.app.id, description: vue.app.description,
        default_return_url: vue.app.defaultReturnUrl,
        default_cancel_url: vue.app.defaultCancelUrl,
        rsa_public_key: vue.app.rsaPublicKey,
        domains: vue.app.domains,
        tags: vue.app.tags,
        id_sources: vue.app.idSources,
      }
      const apiUrl = clientUtils.apiMyApp.replaceAll(':app', this.$route.params.id)
      clientUtils.apiDoPut(apiUrl, data,
          (apiRes) => {
            if (apiRes.status != 200) {
              vue.errorMsg = apiRes.status + ": " + apiRes.message
              vue.foundStatus = 1
            } else {
              vue.$router.push({
                name: "MyApps",
                params: {flashMsg: "Application [" + vue.app.id + "] has been updated successfully."},
              })
            }
          },
          (err) => {
            vue.errorMsg = err
            vue.foundStatus = 1
          }
      )
    },
    validatorUrl(val) {
      return val ? patternUrl.test(val.toString().trim()) : true
    },
  }
}
</script>
