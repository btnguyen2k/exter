<template>
  <div>
    <CRow>
      <CCol sm="6" lg="3">
        <CWidgetDropdown color="primary"
                         v-bind:text="systemInfo.cpu.cores+' core(s) / '+systemInfo.cpu.load+' load'"
                         header="CPU">
          <template #footer>
            <CChartLineSimple pointed class="mt-3 mx-3" style="height:70px"
                              :data-points="systemInfo.cpu.history_load"
                              point-hover-background-color="primary"
                              label="Load"
            />
          </template>
        </CWidgetDropdown>
      </CCol>
      <CCol sm="6" lg="3">
        <CWidgetDropdown color="info" v-bind:text="systemInfo.go_routines.num+''" header="Go Routines">
          <template #footer>
            <CChartLineSimple pointed class="mt-3 mx-3" style="height:70px"
                              :data-points="systemInfo.go_routines.history"
                              point-hover-background-color="info"
                              :options="{ elements: { line: { tension: 0.00001 }}}"
                              label="Routines"
            />
          </template>
        </CWidgetDropdown>
      </CCol>
      <CCol sm="6" lg="3">
        <CWidgetDropdown color="danger" v-bind:text="systemInfo.app_memory.usedMb+' Mb'"
                         header="Used AppMemory">
          <template #footer>
            <CChartLineSimple class="mt-3" style="height:70px" background-color="rgba(255,255,255,.2)"
                              :data-points="systemInfo.app_memory.history_usedMb"
                              :options="{ elements: { line: { borderWidth: 2.5 }}}"
                              point-hover-background-color="warning" label="AppMemory"
            />
          </template>
        </CWidgetDropdown>
      </CCol>
      <CCol sm="6" lg="3">
        <CWidgetDropdown color="warning" v-bind:text="systemInfo.memory.freeGb+' Gb'"
                         header="Free SystemMemory">
          <template #footer>
            <CChartLineSimple class="mt-3" style="height:70px" background-color="rgb(250, 152, 152)"
                              :data-points="systemInfo.memory.history_freeGb"
                              :options="{ elements: { line: { borderWidth: 2.5 }}}"
                              point-hover-background-color="danger" label="Free SystemMemory"
            />
          </template>
        </CWidgetDropdown>
      </CCol>
    </CRow>
    <CRow>
      <CCol sm="12">
        <CCard accent-color="info">
          <CCardHeader>
            <strong>Applications ({{ myAppList.data.length }})</strong>
            <div class="card-header-actions">
              <CLink class="card-header-action btn-minimize" @click="clickRegisterApp">
                <CIcon name="cil-library-add"/>
              </CLink>
              <CLink class="card-header-action btn-minimize"
                     @click="isCollapsedMyApps = !isCollapsedMyApps">
                <CIcon :name="`cil-chevron-${isCollapsedMyApps ? 'bottom' : 'top'}`"/>
              </CLink>
            </div>
          </CCardHeader>
          <CCollapse :show="isCollapsedMyApps" :duration="400">
            <CCardBody>
              <CDataTable :items="myAppList.data"
                          :fields="[{label:'',key:'active'},'id','description','sources','tags','actions']">
                <template #active="{item}">
                  <td style="vertical-align: middle">
                    <CIcon name="cil-check" :style="`color: ${item.public_attrs.actv?'green':'grey'}`"/>
                  </td>
                </template>
                <template #description="{item}">
                  <td style="vertical-align: middle">
                    {{ item.public_attrs.desc }}
                  </td>
                </template>
                <template #sources="{item}">
                  <td style="vertical-align: middle">
                    {{ item.public_attrs.sources }}
                  </td>
                </template>
                <template #tags="{item}">
                  <td style="vertical-align: middle">
                    {{ item.public_attrs.tags }}
                  </td>
                </template>
                <template #actions="{item}">
                  <td nowrap="nowrap" style="vertical-align: middle">
                    <CLink @click="clickEditMyApp(item.id)" label="Edit" class="btn-sm btn-primary">
                      <CIcon name="cil-pencil"/>
                    </CLink>
                    &nbsp;
                    <CLink @click="clickDeleteMyApp(item.id)" label="Delete"
                           class="btn-sm btn-danger">
                      <CIcon name="cil-trash"/>
                    </CLink>
                  </td>
                </template>
              </CDataTable>
            </CCardBody>
          </CCollapse>
        </CCard>
      </CCol>
    </CRow>
  </div>
</template>

<script>
import {CChartLineSimple} from './charts/index.js'
import clientUtils from "@/utils/api_client"
import appUtils from "@/utils/app_utils"

var intervalUpdateSystemInfo

export default {
  name: 'Dashboard',
  components: {
    CChartLineSimple,
  },
  mounted: function () {
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

    this._updateSystemInfo()
    this.$nextTick(function () {
      intervalUpdateSystemInfo = window.setInterval(() => this._updateSystemInfo(), 10000);
    })
  },
  destroyed: function () {
    if (intervalUpdateSystemInfo) {
      window.clearInterval(intervalUpdateSystemInfo)
      intervalUpdateSystemInfo = null
    }
  },
  data() {
    return {
      isCollapsedMyApps: true,
      isCollapsedGroups: true,
      isCollapsedUsers: true,
      systemInfo: {
        cpu: {cores: -1, load: -1.0, history_load: []},
        memory: {free: 0, freeGb: 0.0, history_freeGb: []},
        app_memory: {usedMb: 0.0, history_usedMb: []},
        go_routines: {num: 0, history: []},
      },
      myAppList: {data: []},
    }
  },
  methods: {
    _updateSystemInfo() {
      clientUtils.apiDoGet(
          clientUtils.apiSystemInfo,
          (apiRes) => {
            if (apiRes.status == 200) {
              this.$data.systemInfo.cpu = apiRes.data.cpu
              this.$data.systemInfo.go_routines = apiRes.data.go_routines
              this.$data.systemInfo.app_memory = apiRes.data.app_memory
              this.$data.systemInfo.memory = apiRes.data.memory
            } else {
              console.error("Getting system info was unsuccessful: " + JSON.stringify(apiRes))
            }
          },
          (err) => {
            console.error("Error getting system info: " + err)
          }
      )
    },
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
