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
const Page404 = () => import('@/views/pages/Page404')
const Page500 = () => import('@/views/pages/Page500')
const Login = () => import('@/views/pages/Login')
const Register = () => import('@/views/pages/Register')

Vue.use(Router)

let router = new Router({
    mode: 'hash', // https://router.vuejs.org/api/#mode
    linkActiveClass: 'active',
    //scrollBehavior: () => ({y: 0}),
    base: "/app/",
    routes: configRoutes()
})

import appConfig from "@/utils/app_config"
import utils from "@/utils/app_utils"
import api_client from "@/utils/api_client"

router.beforeEach((to, from, next) => {
    if (!to.matched.some(record => record.meta.allowGuest)) {
        let session = utils.loadLoginSession()
        if (session == null) {
            //redirect to login page if not logged in
            return next({name: "Login", query: {returnUrl: to.fullPath, app: appConfig.APP_NAME}})
        }
        let lastUserTokenCheck = utils.localStorageGetAsInt(utils.lskeyLoginSessionLastCheck)
        if (lastUserTokenCheck + 60 < utils.getUnixTimestamp()) {
            lastUserTokenCheck = utils.getUnixTimestamp()
            let uid = session.uid
            let token = session.token
            api_client.apiDoPost(api_client.apiCheckLoginToken, {uid: uid, token: token},
                (apiRes) => {
                    if (apiRes.status != 200) {
                        //redirect to login page if session verification failed
                        console.error("Session verification failed: " + JSON.stringify(apiRes))
                        return next({name: "Login", query: {returnUrl: to.fullPath, app: appConfig.APP_NAME}})
                    } else {
                        utils.localStorageSet(utils.lskeyLoginSessionLastCheck, lastUserTokenCheck)
                        next()
                    }
                },
                (err) => {
                    console.error("Session verification error: " + err)
                    //redirect to login page if cannot verify session
                    return next({name: "Login", query: {returnUrl: to.fullPath, app: appConfig.APP_NAME}})
                })
        } else {
            next()
        }
    } else {
        next()
    }
})

export default router

function configRoutes() {
    return [
        {
            path: '/',
            redirect: '/dashboard',
            name: 'Home',
            component: TheContainer,
            children: [
                {
                    path: 'dashboard',
                    name: 'Dashboard',
                    component: Dashboard
                },
                {
                    path: 'myapps',
                    meta: {label: 'My Apps'},
                    component: {
                        render(c) {
                            return c('router-view')
                        }
                    },
                    children: [
                        {
                            path: '',
                            meta: {label: 'App List'},
                            name: 'MyApps',
                            component: MyApps,
                            props: true,
                        },
                        {
                            path: '_register',
                            meta: {label: 'Register New App'},
                            name: 'RegisterApp',
                            component: RegisterApp,
                        },
                        {
                            path: '_edit/:id',
                            meta: {label: 'Edit My App'},
                            name: 'EditMyApp',
                            component: EditMyApp,
                        },
                        {
                            path: '_delete/:id',
                            meta: {label: 'Delete My App'},
                            name: 'DeleteMyApp',
                            component: DeleteMyApp,
                        },
                    ]
                },
            ]
        },
        {
            path: '/pages',
            redirect: '/pages/404',
            name: 'Pages',
            component: {
                render(c) {
                    return c('router-view')
                }
            },
            meta: {
                allowGuest: true
            },
            children: [
                {
                    path: '404',
                    name: 'Page404',
                    component: Page404
                },
                {
                    path: '500',
                    name: 'Page500',
                    component: Page500
                },
                {
                    path: 'login',
                    name: 'Login',
                    component: Login,
                    //props: (route) => ({returnUrl: route.query.returnUrl, app: route.query.app}),
                    //params: (route) => ({returnUrl: route.query.returnUrl, app: route.query.app}),
                },
                {
                    path: 'register',
                    name: 'Register',
                    component: Register
                }
            ]
        },
        {
            path: '*',
            redirect: '/',
        }
    ]
}
