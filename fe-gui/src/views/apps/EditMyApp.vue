<template>
  <div>
    <CRow v-if="!found">
      <CCol sm="12">
        <CCard>
          <CCardHeader>Edit Application</CCardHeader>
          <CCardBody>
            <p class="alert alert-danger">Application [{{ this.$route.params.id }}] not found</p>
          </CCardBody>
          <CCardFooter>
            <CButton type="button" color="info" class="ml-2" style="width: 96px" @click="doCancel">
              <CIcon name="cil-arrow-circle-left"/>
              Back
            </CButton>
          </CCardFooter>
        </CCard>
      </CCol>
    </CRow>
    <CRow v-if="found">
      <CCol sm="12">
        <CCard>
          <CCardHeader>Edit User</CCardHeader>
          <CForm @submit.prevent="doSubmit" method="post">
            <CCardBody>
              <p v-if="errorMsg!=''" class="alert alert-danger">{{ errorMsg }}</p>
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
              <CInput horizontal type="text" v-model="app.tags" label="Tags"
                      placeholder="Tags separated by comma"
              />
              <CTextarea horizontal type="text" v-model="app.rsaPublicKey" label="RSA public key"
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
import clientUtils from "@/utils/api_client"

let patternUrl = /^http(s?):\/\//

export default {
  name: 'EditMyApp',
  data() {
    let loginChannelList = []
    let app = {
      isActive: true,
      id: "", description: "", rsaPublicKey: "", defaultReturnUrl: "", defaultCancelUrl: "",
      tags: "",
      idSources: {},
    }
    clientUtils.apiDoGet(clientUtils.apiMyApp + "/" + this.$route.params.id,
        (apiRes) => {
          this.found = apiRes.status == 200
          if (apiRes.status == 200) {
            app.id = apiRes.data.id
            app.isActive = apiRes.data.config.actv
            app.description = apiRes.data.config.desc
            app.rsaPublicKey = apiRes.data.config.rpub
            app.defaultReturnUrl = apiRes.data.config.rurl
            app.defaultCancelUrl = apiRes.data.config.curl
            app.idSources = apiRes.data.config.sources
            app.tags = apiRes.data.config.tags != null ? apiRes.data.config.tags.join(", ") : ""

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
      let data = {
        is_active: this.app.isActive,
        id: this.app.id, description: this.app.description,
        default_return_url: this.app.defaultReturnUrl,
        default_cancel_url: this.app.defaultCancelUrl,
        rsa_public_key: this.app.rsaPublicKey,
        tags: this.app.tags,
        id_sources: this.app.idSources,
      }
      clientUtils.apiDoPut(
          clientUtils.apiMyApp + "/" + this.$route.params.id, data,
          (apiRes) => {
            if (apiRes.status != 200) {
              this.errorMsg = apiRes.status + ": " + apiRes.message
            } else {
              this.$router.push({
                name: "MyApps",
                params: {flashMsg: "Application [" + this.app.id + "] has been updated successfully."},
              })
            }
          },
          (err) => {
            this.errorMsg = err
          }
      )
    },
    validatorUrl(val) {
      return val ? patternUrl.test(val.toString().trim()) : true
    },
  }
}
</script>
