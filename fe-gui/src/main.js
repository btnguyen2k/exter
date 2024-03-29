import Vue from 'vue'
import App from './App'
import i18n from './i18n'
import router from './router'
import CoreuiVue from '@coreui/vue'
import {iconsSet as icons} from './assets/icons/icons.js'
import LoadScript from 'vue-plugin-load-script'

Vue.config.performance = true
Vue.use(CoreuiVue)
Vue.use(LoadScript)

new Vue({
    el: '#app',
    router,
    icons,
    template: '<App/>',
    components: {
        App
    },
    i18n
})
