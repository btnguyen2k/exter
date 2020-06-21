<template>
    <CSidebar fixed :minimize="minimize" :show.sync="show">
        <div class="c-sidebar-brand">
            <a href="/">
                <span class="c-sidebar-brand-full" style="color: #fff; font-weight: bolder; font-size: x-large">{{appName}}</span>
                <span class="c-sidebar-brand-minimized"
                      style="color: #fff; font-weight: bolder; font-size: large">{{appInitial}}</span>
            </a>
        </div>
        <CRenderFunction flat :content-to-render="nav"/>
        <CSidebarMinimizer class="d-md-down-none" @click.native="minimize = !minimize"/>
    </CSidebar>
</template>

<script>
    import nav from './_nav'
    import cfg from '@/utils/app_config'

    export default {
        name: 'TheSidebar',
        data() {
            return {
                minimize: false,
                nav,
                show: 'responsive',
                appName: cfg.APP_CONFIG.app.name,
                appInitial: cfg.APP_CONFIG.app.initial,
            }
        },
        mounted() {
            this.$root.$on('toggle-sidebar', () => {
                const sidebarOpened = this.show === true || this.show === 'responsive'
                this.show = sidebarOpened ? false : 'responsive'
            })
            this.$root.$on('toggle-sidebar-mobile', () => {
                const sidebarClosed = this.show === 'responsive' || this.show === false
                this.show = sidebarClosed ? true : 'responsive'
            })
        }
    }
</script>
