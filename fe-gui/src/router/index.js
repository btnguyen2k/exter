import Vue from 'vue'
import Router from 'vue-router'

// Containers
const TheContainer = () => import('@/containers/TheContainer')

// Views
const Dashboard = () => import('@/views/Dashboard')

// My Apps
const MyApps = () => import('@/views/apps/MyApps')
const RegisterApp = () => import('@/views/apps/RegisterApp')
const EditMyApp = () => import('@/views/apps/EditMyApp')
const DeleteMyApp = () => import('@/views/apps/DeleteMyApp')

// Views - Pages
const Login = () => import('@/views/pages/Login')
const CheckLogin = () => import('@/views/pages/CheckLogin')

Vue.use(Router)

let router = new Router({
    mode: 'history', // https://router.vuejs.org/api/#mode
    linkActiveClass: 'active',
    //scrollBehavior: () => ({y: 0}),
    base: "/app/",
    routes: configRoutes()
})

import appConfig from "@/utils/app_config"
import utils from "@/utils/app_utils"
import clientUtils from "@/utils/api_client"

router.beforeEach((to, from, next) => {
    if (!to.matched.some(record => record.meta.allowGuest)) {
        let session = utils.loadLoginSession()
        if (session == null) {
            //redirect to login page if not logged in
            return next({name: "Login", query: {returnUrl: router.resolve(to, from).href, app: appConfig.APP_ID}})
        }
        let lastUserTokenCheck = utils.localStorageGetAsInt(utils.lskeyLoginSessionLastCheck)
        if (lastUserTokenCheck + 60 < utils.getUnixTimestamp()) {
            lastUserTokenCheck = utils.getUnixTimestamp()
            let token = session.token
            clientUtils.apiDoPost(clientUtils.apiVerifyLoginToken, {app: appConfig.APP_ID, token: token},
                (apiRes) => {
                    if (apiRes.status != 200) {
                        //redirect to login page if session verification failed
                        console.error("Session verification failed: " + JSON.stringify(apiRes))
                        return next({
                            name: "Login",
                            query: {returnUrl: router.resolve(to, from).href, app: appConfig.APP_ID}
                        })
                    } else {
                        utils.localStorageSet(utils.lskeyLoginSessionLastCheck, lastUserTokenCheck)
                        next()
                    }
                },
                (err) => {
                    //redirect to login page if cannot verify session
                    console.error("Session verification error: " + err)
                    return next({
                        name: "Login",
                        query: {returnUrl: router.resolve(to, from).href, app: appConfig.APP_ID}
                    })
                })
        } else {
            next()
        }
    } else {
        next()
    }
})

export default router

import i18n from '../i18n'

function configRoutes() {
    return [
        {
            path: '/',
            redirect: {name: "Dashboard"},
            name: 'Home', meta: {label: i18n.t('message.home')},
            component: TheContainer,
            children: [
                {
                    path: 'dashboard',
                    name: 'Dashboard', meta: {label: i18n.t('message.dashboard')},
                    component: Dashboard
                },
                {
                    path: 'myapps',
                    meta: {label: i18n.t('message.my_apps')},
                    component: {
                        render(c) {
                            return c('router-view')
                        }
                    },
                    children: [
                        {
                            path: '',
                            name: 'MyApps', meta: {label: i18n.t('message.my_app_list')},
                            component: MyApps,
                            props: true, //for passing flash message
                        },
                        {
                            path: '_register',
                            name: 'RegisterApp', meta: {label: i18n.t('message.register_app')},
                            component: RegisterApp,
                        },
                        {
                            path: '_edit/:id',
                            name: 'EditMyApp', meta: {label: i18n.t('message.edit_my_app')},
                            component: EditMyApp,
                        },
                        {
                            path: '_delete/:id',
                            name: 'DeleteMyApp', meta: {label: i18n.t('message.delete_my_app')},
                            component: DeleteMyApp,
                        },
                    ]
                },
            ]
        },
        {
            path: '/xlogin',
            meta: {
                allowGuest: true
            },
            name: 'Login',
            component: Login,
        },
        {
            path: '/xcheck',
            meta: {
                allowGuest: true
            },
            name: 'CheckLogin',
            component: CheckLogin,
        },
        {
            path: '*',
            redirect: '/',
        }
    ]
}
