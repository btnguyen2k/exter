import Vue from 'vue'
import Router from 'vue-router'

// Containers
const TheContainer = () => import('@/containers/TheContainer')

// Views
const Dashboard = () => import('@/views/Dashboard')

// Apps
const Apps = () => import('@/views/apps/Apps')
const RegisterApp = () => import('@/views/apps/RegisterApp')
const EditApp = () => import('@/views/apps/EditApp')
const DeleteApp = () => import('@/views/apps/DeleteApp')

// Groups
const Groups = () => import('@/views/groups/Groups')
const CreateGroup = () => import('@/views/groups/CreateGroup')
const EditGroup = () => import('@/views/groups/EditGroup')
const DeleteGroup = () => import('@/views/groups/DeleteGroup')

// Users
const Users = () => import('@/views/users/Users')
const CreateUser = () => import('@/views/users/CreateUser')
const EditUser = () => import('@/views/users/EditUser')
const DeleteUser = () => import('@/views/users/DeleteUser')

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
                    path: 'apps',
                    meta: {label: 'Apps'},
                    component: {
                        render(c) {
                            return c('router-view')
                        }
                    },
                    children: [
                        {
                            path: '',
                            meta: {label: 'App List'},
                            name: 'Apps',
                            component: Apps,
                            props: true,
                        },
                        {
                            path: '_register',
                            meta: {label: 'Register New App'},
                            name: 'RegisterApp',
                            component: RegisterApp,
                        },
                        {
                            path: '_edit/:app',
                            meta: {label: 'Edit App'},
                            name: 'EditApp',
                            component: EditApp,
                        },
                        {
                            path: '_delete/:app',
                            meta: {label: 'Delete App'},
                            name: 'DeleteApp',
                            component: DeleteApp,
                        },
                    ]
                },
                {
                    path: 'groups',
                    meta: {label: 'Groups'},
                    component: {
                        render(c) {
                            return c('router-view')
                        }
                    },
                    children: [
                        {
                            path: '',
                            meta: {label: 'Group List'},
                            name: 'Groups',
                            component: Groups,
                            props: true,
                        },
                        {
                            path: '_create',
                            meta: {label: 'Create New Group'},
                            name: 'CreateGroup',
                            component: CreateGroup,
                        },
                        {
                            path: '_edit/:id',
                            meta: {label: 'Edit Group'},
                            name: 'EditGroup',
                            component: EditGroup,
                        },
                        {
                            path: '_delete/:id',
                            meta: {label: 'Delete Group'},
                            name: 'DeleteGroup',
                            component: DeleteGroup,
                        },
                    ]
                },
                {
                    path: 'users',
                    meta: {label: 'Users'},
                    component: {
                        render(c) {
                            return c('router-view')
                        }
                    },
                    children: [
                        {
                            path: '',
                            meta: {label: 'User List'},
                            name: 'Users',
                            component: Users,
                            props: true,
                        },
                        {
                            path: '_create',
                            meta: {label: 'Create New User'},
                            name: 'CreateUser',
                            component: RegisterApp,
                        },
                        {
                            path: '_edit/:username',
                            meta: {label: 'Edit User'},
                            name: 'EditUser',
                            component: EditUser,
                        },
                        {
                            path: '_delete/:username',
                            meta: {label: 'Delete User'},
                            name: 'DeleteUser',
                            component: DeleteUser,
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
        }
    ]
}

