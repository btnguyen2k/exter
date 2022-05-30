<template>
  <div>
    <CRow>
      <CCol sm="12">
        <CCard>
          <CCardHeader>{{ $t('message.register_app') }}</CCardHeader>
          <CForm @submit.prevent="doSubmit" method="post">
            <CCardBody>
              <p v-if="errorMsg!=''" class="alert alert-danger">{{ errorMsg }}</p>
              <CInput horizontal type="text" v-model="form.id" :label="$t('message.app_id')"
                      :placeholder="$t('message.app_id_placeholder')"
                      :is-valid="validatorAppId"
                      :invalid-feedback="$t('message.app_id_rule')"
                      valid-feedback=""
              />
              <div class="form-group form-row">
                <CCol :sm="{offset:3,size:9}" class="form-inline">
                  <CInputCheckbox inline :label="$t('message.app_active')" :checked.sync="form.isActive"/>
                </CCol>
              </div>
              <div class="form-group form-row">
                <CCol tag="label" sm="3" class="col-form-label">
                  {{ $t('message.app_auth_provider') }}
                </CCol>
                <CCol sm="9" class="form-inline">
                  <CInputCheckbox inline v-for="(option, _) in loginChannelList" :label="$t('message.auth_provider_'+option)"
                                  :value="option" :checked.sync="form.idSources[option]"
                  />
                </CCol>
              </div>
              <CInput horizontal type="text" v-model="form.description" :label="$t('message.app_desc')"
                      :placeholder="$t('message.app_desc_placeholder')"
              />
              <CInput horizontal type="text" v-model="form.defaultReturnUrl" :label="$t('message.app_default_return_url')"
                      :placeholder="$t('message.app_default_return_url_placeholder')"
                      :is-valid="validatorUrl"
                      :invalid-feedback="$t('message.app_default_return_url_rule')"
              />
              <CInput horizontal type="text" v-model="form.defaultCancelUrl" :label="$t('message.app_default_cancel_url')"
                      :placeholder="$t('message.app_default_cancel_url_placeholder')"
                      :is-valid="validatorUrl"
                      :invalid-feedback="$t('message.app_default_cancel_url_rule')"
              />
              <CInput horizontal type="text" v-model="form.domains" :label="$t('message.app_domains')"
                      :placeholder="$t('message.app_domains_placeholder')"
              />
              <CInput horizontal type="text" v-model="form.tags" :label="$t('message.app_tags')"
                      :placeholder="$t('message.app_tags_placeholder')"
              />
              <CTextarea horizontal type="text" v-model="form.rsaPublicKey" :label="$t('message.app_rsa_pubkey')"
                         rows="6" :placeholder="$t('message.app_rsa_pubkey_placeholder')"
              />
            </CCardBody>
            <CCardFooter>
              <CButton type="submit" color="primary" style="width: 96px">
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
import clientUtils from "@/utils/api_client";

let patternAppId = /^[0-9a-z_]+$/
let patternUrl = /^http(s?):\/\//

export default {
  name: 'RegisterApp',
  mounted() {
    const vue = this
    let loginChannelList = []
    clientUtils.apiDoGet(clientUtils.apiInfo,
        (apiRes) => {
          if (apiRes.status == 200) {
            apiRes.data.login_channels.every(function (e) {
              loginChannelList.push(e)
              vue.form.idSources[e] = true
              return true
            })
            vue.loginChannelList = loginChannelList
          } else {
            console.error("Calling api "+clientUtils.apiInfo+" was unsuccessful: " + apiRes)
          }
        },
        (err) => {
          console.error("Error calling api "+clientUtils.apiInfo+": " + err)
        })
  },
  data() {
    return {
      form: {
        isActive: true,
        id: "", description: "", rsaPublicKey: "", defaultReturnUrl: "", defaultCancelUrl: "",
        domains: "",
        tags: "",
        idSources: {},
      },
      errorMsg: "",
      loginChannelList: [],
    }
  },
  methods: {
    doCancel() {
      router.push(router.resolve({name: "MyApps"}).location)
    },
    doSubmit(e) {
      e.preventDefault()
      let data = {
        is_active: this.form.isActive,
        id: this.form.id, description: this.form.description,
        default_return_url: this.form.defaultReturnUrl,
        default_cancel_url: this.form.defaultCancelUrl,
        rsa_public_key: this.form.rsaPublicKey,
        domains: this.form.domains,
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
                name: "MyApps",
                params: {flashMsg: this.$i18n.t('message.app_registered_successful', {id: this.form.id})},
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
    validatorUrl(val) {
      return val ? patternUrl.test(val.toString().trim()) : true
    },
  }
}
</script>
