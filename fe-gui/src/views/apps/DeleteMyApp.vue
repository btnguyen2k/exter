<template>
  <div>
    <CRow>
      <CCol sm="12">
        <CCard>
          <CCardHeader>{{ $t('message.delete_my_app') }}</CCardHeader>
          <CForm @submit.prevent="doSubmit" method="post">
            <CCardBody>
              <p v-if="!found" class="alert alert-danger">{{ $t('message.error_app_not_exist', {app: this.$route.params.id}) }}</p>
              <p v-if="errorMsg!=''" class="alert alert-danger">{{ errorMsg }}</p>
              <CInput v-if="found" horizontal type="text" v-model="app.id" :label="$t('message.app_id')"
                      :placeholder="$t('message.app_id_placeholder')"
                      disabled="disabled"
              />
              <div v-if="found" class="form-group form-row">
                <CCol :sm="{offset:3,size:9}" class="form-inline">
                </CCol>
                <CInputCheckbox inline :label="$t('message.app_active')" :checked.sync="app.isActive" disabled="disabled"
                />
              </div>
              <div v-if="found" class="form-group form-row">
                <CCol tag="label" sm="3" class="col-form-label">
                  {{ $t('message.app_auth_provider') }}
                </CCol>
                <CCol sm="9" class="form-inline">
                  <CInputCheckbox inline v-for="option in loginChannelList" :label="$t('message.auth_provider_'+option)"
                                  :value="option" :checked.sync="app.idSources[option]"
                                  disabled="disabled"
                  />
                </CCol>
              </div>
              <CInput v-if="found" horizontal type="text" v-model="app.description" :label="$t('message.app_desc')"
                      :placeholder="$t('message.app_desc_placeholder')"
                      disabled="disabled"
              />
              <CInput v-if="found" horizontal type="text" v-model="app.defaultReturnUrl" :label="$t('message.app_default_return_url')"
                      :placeholder="$t('message.app_default_return_url_placeholder')"
                      disabled="disabled"
              />
              <CInput v-if="found" horizontal type="text" v-model="app.defaultCancelUrl" :label="$t('message.app_default_cancel_url')"
                      :placeholder="$t('message.app_default_cancel_url_placeholder')"
                      disabled="disabled"
              />
              <CInput v-if="found" horizontal type="text" v-model="app.domains" :label="$t('message.app_domains')"
                      :placeholder="$t('message.app_domains_placeholder')"
                      disabled="disabled"
              />
              <CInput v-if="found" horizontal type="text" v-model="app.tags" :label="$t('message.app_tags')"
                      :placeholder="$t('message.app_tags_placeholder')"
                      disabled="disabled"
              />
              <CTextarea v-if="found" horizontal type="text" v-model="app.rsaPublicKey" :label="$t('message.app_rsa_pubkey')"
                         rows="6" :placeholder="$t('message.app_rsa_pubkey_placeholder')"
                         disabled="disabled"
              />
            </CCardBody>
            <CCardFooter>
              <CButton v-if="found" type="submit" color="danger" style="width: 96px">
                <CIcon name="cil-trash"/>
                {{ $t('message.delete') }}
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

export default {
  name: 'DeleteMyApp',
  data() {
    let loginChannelList = []
    let app = {
      isActive: true,
      id: "", description: "", rsaPublicKey: "", defaultReturnUrl: "", defaultCancelUrl: "",
      domains: "",
      tags: "",
      idSources: {},
    }
    clientUtils.apiDoGet(clientUtils.apiMyApp.replaceAll(':app', this.$route.params.id),
        (apiRes) => {
          this.found = apiRes.status == 200
          if (apiRes.status == 200) {
            app.id = apiRes.data.id
            app.isActive = apiRes.data.public_attrs.actv
            app.description = apiRes.data.public_attrs.desc
            app.rsaPublicKey = apiRes.data.public_attrs.rpub
            app.defaultReturnUrl = apiRes.data.public_attrs.rurl
            app.defaultCancelUrl = apiRes.data.public_attrs.curl
            app.idSources = apiRes.data.public_attrs.sources
            // app.tags = apiRes.data.public_attrs.tags != null ? apiRes.data.public_attrs.tags.join(", ") : ""
            app.tags = JSON.stringify(apiRes.data.public_attrs.tags)
            app.domains = JSON.stringify(apiRes.data.domains)
            clientUtils.apiDoGet(clientUtils.apiInfo,
                (apiRes) => {
                  if (apiRes.status == 200) {
                    apiRes.data.login_channels.every(function (e) {
                      loginChannelList.push(e)
                      return true
                    })
                  } else {
                    console.error("Getting info was unsuccessful: " + apiRes)
                  }
                },
                (err) => {
                  console.error("Error getting info list: " + err)
                })
          }
        },
        (err) => {
          this.errorMsg = err
        })
    return {
      app: app,
      errorMsg: "",
      loginChannelList: loginChannelList,
      found: true,
    }
  },
  methods: {
    doCancel() {
      router.push(router.resolve({name: "MyApps"}).location)
    },
    doSubmit(e) {
      e.preventDefault()
      clientUtils.apiDoDelete(
          clientUtils.apiMyApp.replaceAll(':app', this.$route.params.id),
          (apiRes) => {
            if (apiRes.status != 200) {
              this.errorMsg = apiRes.status + ": " + apiRes.message
            } else {
              this.$router.push({
                name: "MyApps",
                params: {flashMsg: this.$i18n.t('message.app_deleted_successful', {id: this.app.id})},
              })
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
