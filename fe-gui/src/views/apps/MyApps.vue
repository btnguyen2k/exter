<template>
    <CRow>
        <CCol sm="12">
            <CCard accent-color="info">
                <CCardHeader>
                    <strong>Applications ({{myAppList.data.length}})</strong>
                    <div class="card-header-actions">
                        <CButton class="btn-sm btn-primary" @click="clickRegisterApp">
                            <CIcon name="cil-playlist-add"/>
                            Register New App
                        </CButton>
                    </div>
                </CCardHeader>
                <CCardBody>
                    <p v-if="flashMsg" class="alert alert-success">{{flashMsg}}</p>
                    <CDataTable :items="myAppList.data" :fields="[{label:'',key:'active'},'id','description','sources','tags','actions']">
                        <template #active="{item}">
                            <td>
                                <CIcon :name="`${item.config.actv?'cil-check':'cil-check-alt'}`" :style="`color: ${item.config.actv?'green':'grey'}`"/>
                            </td>
                        </template>
                        <template #description="{item}">
                            <td>
                                {{item.config.desc}}
                            </td>
                        </template>
                        <template #sources="{item}">
                            <td>
                                {{item.config.sources}}
                            </td>
                        </template>
                        <template #tags="{item}">
                            <td>
                                {{item.config.tags}}
                            </td>
                        </template>
                        <template #actions="{item}">
                            <td>
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
                            //console.log(myAppList.data)
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
