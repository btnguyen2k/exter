<template>
  <div>
    <CRow>
      <CCol sm="12">
        <CCard accent-color="info">
          <CCardHeader>{{ $t('message.edit_my_app') }}</CCardHeader>
          <CForm @submit.prevent="doSubmit" method="post">
            <CCardBody v-if="foundStatus<0">
              <CAlert color="info">{{ $t('message.wait') }}</CAlert>
            </CCardBody>
            <CCardBody v-if="foundStatus==0">
              <CAlert color="danger">{{ $t('message.error_app_not_exist', {app: this.$route.params.id}) }}</CAlert>
            </CCardBody>
            <CCardBody v-if="foundStatus>0">
              <CAlert v-if="errorMsg" color="danger">{{ errorMsg }}</CAlert>
              <CInput horizontal type="text" v-model="app.id" :label="$t('message.app_id')"
                      :placeholder="$t('message.app_id_placeholder')"
                      readonly="readonly"
              />
              <div class="form-group form-row">
                <CCol :sm="{offset:3,size:9}" class="form-inline">
                  <CInputCheckbox inline :label="$t('message.app_active')" :checked.sync="app.isActive"
                  />
                </CCol>
              </div>
              <div class="form-group form-row">
                <CCol tag="label" sm="3" class="col-form-label">
                  {{ $t('message.app_auth_provider') }}
                </CCol>
                <CCol sm="9" class="form-inline">
                  <CInputCheckbox inline v-for="option in loginChannelList" :label="$t('message.auth_provider_'+option)"
                                  :value="option" :checked.sync="app.idSources[option]"
                  />
                </CCol>
              </div>
              <CInput horizontal type="text" v-model="app.description" :label="$t('message.app_desc')"
                      :placeholder="$t('message.app_desc_placeholder')"
              />
              <CInput horizontal type="text" v-model="app.defaultReturnUrl" :label="$t('message.app_default_return_url')"
                      :placeholder="$t('message.app_default_return_url_placeholder')"
                      :is-valid="validatorUrl"
                      :invalid-feedback="$t('message.app_default_return_url_rule')"
              />
              <CInput horizontal type="text" v-model="app.defaultCancelUrl" :label="$t('message.app_default_cancel_url')"
                      :placeholder="$t('message.app_default_cancel_url_placeholder')"
                      :is-valid="validatorUrl"
                      :invalid-feedback="$t('message.app_default_cancel_url_rule')"
              />
              <CInput horizontal type="text" v-model="app.domains" :label="$t('message.app_domains')"
                      :placeholder="$t('message.app_domains_placeholder')"
              />
              <CInput horizontal type="text" v-model="app.tags" :label="$t('message.app_tags')"
                      :placeholder="$t('message.app_tags_placeholder')"
              />
              <CTextarea horizontal type="text" v-model="app.rsaPublicKey" :label="$t('message.app_rsa_pubkey')"
                         rows="6" :placeholder="$t('message.app_rsa_pubkey_placeholder')"
              />
            </CCardBody>
            <CCardFooter>
              <CButton v-if="foundStatus>0" type="submit" color="primary" style="width: 96px">
                <CIcon name="cil-save"/>
                {{ $t('message.save') }}
              </CButton>
              <CButton type="button" color="info" class="ml-2" style="width: 96px" @click="doCancel">
                <CIcon name="cil-arrow-circle-left"/>
                {{ $t('message.back') }}
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
              vue.loginChannelList = _loginChannels

              const apiUrl = clientUtils.apiMyApp.replaceAll(':app', appId)
              clientUtils.apiDoGet(apiUrl,
                  (apiRes) => {
                    vue.foundStatus = apiRes.status == 200 ? 1 : 0
                    if (vue.foundStatus == 1) {
                      let _app = {
                        id: apiRes.data.id,
                        isActive: apiRes.data.public_attrs.actv,
                        description: apiRes.data.public_attrs.desc,
                        rsaPublicKey: apiRes.data.public_attrs.rpub,
                        defaultReturnUrl: apiRes.data.public_attrs.rurl,
                        defaultCancelUrl: apiRes.data.public_attrs.curl,
                        idSources: apiRes.data.public_attrs.sources,
                        domains: apiRes.data.domains != null ? apiRes.data.domains.join(", ") : "",
                        tags: apiRes.data.public_attrs.tags != null ? apiRes.data.public_attrs.tags.join(", ") : ""
                      }
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
                params: {flashMsg: this.$i18n.t('message.app_updated_successful', {id: vue.app.id})},
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
