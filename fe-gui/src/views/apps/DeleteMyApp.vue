<template>
  <div>
    <CRow>
      <CCol sm="12">
        <CCard>
          <CCardHeader>Delete Application</CCardHeader>
          <CForm @submit.prevent="doSubmit" method="post">
            <CCardBody>
              <p v-if="!found" class="alert alert-danger">Application [{{ this.$route.params.id }}] not found</p>
              <p v-if="errorMsg!=''" class="alert alert-danger">{{ errorMsg }}</p>
              <CInput v-if="found" horizontal type="text" v-model="app.id" label="Id"
                      placeholder="Application's unique id"
                      disabled="disabled"
              />
              <div v-if="found" class="form-group form-row">
                <CCol :sm="{offset:3,size:9}" class="form-inline">
                  <CInputCheckbox inline label="Active" :checked.sync="app.isActive"
                                  disabled="disabled"
                  />
                </CCol>
              </div>
              <div v-if="found" class="form-group form-row">
                <CCol tag="label" sm="3" class="col-form-label">
                  Login channels
                </CCol>
                <CCol sm="9" class="form-inline">
                  <CInputCheckbox inline v-for="option in loginChannelList" :label="option"
                                  :value="option" :checked.sync="app.idSources[option]"
                                  disabled="disabled"
                  />
                </CCol>
              </div>
              <CInput v-if="found" horizontal type="text" v-model="app.description" label="Description"
                      placeholder="Application's description"
                      disabled="disabled"
              />
              <CInput v-if="found" horizontal type="text" v-model="app.defaultReturnUrl" label="Default return url"
                      placeholder="http://..."
                      disabled="disabled"
              />
              <CInput v-if="found" horizontal type="text" v-model="app.defaultCancelUrl" label="Default cancel url"
                      placeholder="http://..."
                      disabled="disabled"
              />
              <CInput v-if="found" horizontal type="text" v-model="app.tags" label="Tags"
                      placeholder="Tags separated by comma"
                      disabled="disabled"
              />
              <CTextarea v-if="found" horizontal type="text" v-model="app.rsaPublicKey" label="RSA public key"
                         rows="6" placeholder="RSA public key in PEM format"
                         disabled="disabled"
              />
            </CCardBody>
            <CCardFooter>
              <CButton v-if="found" type="submit" color="danger" style="width: 96px">
                <CIcon name="cil-trash"/>
                Delete
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

export default {
  name: 'DeleteMyApp',
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
            app.isActive = apiRes.data.public_attrs.actv
            app.description = apiRes.data.public_attrs.desc
            app.rsaPublicKey = apiRes.data.public_attrs.rpub
            app.defaultReturnUrl = apiRes.data.public_attrs.rurl
            app.defaultCancelUrl = apiRes.data.public_attrs.curl
            app.idSources = apiRes.data.public_attrs.sources
            app.tags = apiRes.data.public_attrs.tags != null ? apiRes.data.public_attrs.tags.join(", ") : ""

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
      router.push("/myapps")
    },
    doSubmit(e) {
      e.preventDefault()
      clientUtils.apiDoDelete(
          clientUtils.apiMyApp + "/" + this.$route.params.id,
          (apiRes) => {
            if (apiRes.status != 200) {
              this.errorMsg = apiRes.status + ": " + apiRes.message
            } else {
              this.$router.push({
                name: "MyApps",
                params: {flashMsg: "Application [" + this.app.id + "] has been deleted successfully."},
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
