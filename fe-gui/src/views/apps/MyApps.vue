<template>
  <CRow>
    <CCol sm="12">
      <CCard accent-color="info">
        <CCardHeader>
          <strong>{{ $t('message.my_apps') }} ({{ myAppList.data.length }})</strong>
          <div class="card-header-actions">
            <CButton class="btn-sm btn-primary" @click="clickRegisterApp">
              <CIcon name="cil-playlist-add"/>
              {{ $t('message.register_app') }}
            </CButton>
          </div>
        </CCardHeader>
        <CCardBody>
          <p v-if="flashMsg" class="alert alert-success">{{ flashMsg }}</p>
          <CDataTable :items="myAppList.data"
                      :fields="[
                              {label:'',key:'active'},
                              {label:$t('message.app_id'),key:'id'},
                              {label:$t('message.app_desc'),key:'description'},
                              {label:$t('message.app_auth_provider'),key:'sources'},
                              {label:$t('message.app_tags'),key:'tags'},
                              {label:$t('message.actions'),key:'actions'}
                          ]">
            <template #active="{item}">
              <td>
                <CIcon :name="`${item.public_attrs.actv?'cil-check':'cil-check-alt'}`"
                       :style="`color: ${item.public_attrs.actv?'green':'grey'}`"/>
              </td>
            </template>
            <template #description="{item}">
              <td>
                {{ item.public_attrs.desc }}
              </td>
            </template>
            <template #sources="{item}">
              <td>
                {{ item.public_attrs.sources }}
              </td>
            </template>
            <template #tags="{item}">
              <td>
                {{ item.public_attrs.tags }}
              </td>
            </template>
            <template #actions="{item}">
              <td>
                <CLink @click="clickEditMyApp(item.id)" class="btn-sm btn-primary">
                  <CIcon name="cil-pencil"/>
                </CLink>
                &nbsp;
                <CLink @click="clickDeleteMyApp(item.id)" class="btn-sm btn-danger">
                  <CIcon name="cil-trash"/>
                </CLink>
              </td>
            </template>
          </CDataTable>
        </CCardBody>
      </CCard>
    </CCol>
  </CRow>
</template>

<script>
import clientUtils from "@/utils/api_client"
import appUtils from "@/utils/app_utils"

export default {
  name: 'MyApps',
  data: () => {
    let myAppList = {data: []}
    let session = appUtils.loadLoginSession()
    if (session != null) {
      clientUtils.apiDoGet(clientUtils.apiMyAppList + "?token=" + session.token,
          (apiRes) => {
            if (apiRes.status == 200) {
              myAppList.data = apiRes.data
            } else {
              console.error("Getting my app list was unsuccessful: " + JSON.stringify(apiRes))
            }
          },
          (err) => {
            console.error("Error getting my app list: " + err)
          })
    }
    return {
      myAppList: myAppList,
    }
  },
  props: ["flashMsg"],
  methods: {
    clickRegisterApp(e) {
      this.$router.push({name: "RegisterApp"})
    },
    clickEditMyApp(id) {
      this.$router.push({name: "EditMyApp", params: {id: id.toString()}})
    },
    clickDeleteMyApp(id) {
      this.$router.push({name: "DeleteMyApp", params: {id: id.toString()}})
    },
  }
}
</script>
