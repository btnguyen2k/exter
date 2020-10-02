<template>
    <CHeader fixed with-subheader light>
        <CToggler in-header class="ml-3 d-lg-none" v-c-emit-root-event:toggle-sidebar-mobile/>
        <CToggler in-header class="ml-3 d-md-down-none" v-c-emit-root-event:toggle-sidebar/>
        <a href="/" class="c-header-brand mx-auto d-lg-none" style="font-weight: bolder; font-size: x-large">{{appName}}</a>
        <!--
        <CHeaderBrand
                class="mx-auto d-lg-none"
                src="img/brand/coreui-vue-logo.svg"
                width="190"
                height="46"
                alt="CoreUI Logo"
        />
        -->
        <CHeaderNav class="d-md-down-none mr-auto">
            <CHeaderNavItem class="px-3">
                <CHeaderNavLink :to="{name:'Dashboard'}">
                    Dashboard
                </CHeaderNavLink>
            </CHeaderNavItem>
            <CHeaderNavItem class="px-3">
                <CHeaderNavLink :to="{name:'MyApps'}" exact>
                    My Apps
                </CHeaderNavLink>
            </CHeaderNavItem>
        </CHeaderNav>
        <CHeaderNav class="mr-4">
            <CHeaderNavItem class="d-md-down-none mx-2">
                <CHeaderNavLink>
                    <CIcon name="cil-bell"/>
                </CHeaderNavLink>
            </CHeaderNavItem>
            <CHeaderNavItem class="d-md-down-none mx-2">
                <CHeaderNavLink>
                    <CIcon name="cil-list"/>
                </CHeaderNavLink>
            </CHeaderNavItem>
            <CHeaderNavItem class="d-md-down-none mx-2">
                <CHeaderNavLink>
                    <CIcon name="cil-envelope-open"/>
                </CHeaderNavLink>
            </CHeaderNavItem>
            <TheHeaderDropdownAccnt/>
        </CHeaderNav>
        <CSubheader class="px-3">
            <CBreadcrumbRouter class="border-0"/>
        </CSubheader>
    </CHeader>
</template>

<script>
    import TheHeaderDropdownAccnt from './TheHeaderDropdownAccnt'
    import cfg from '@/utils/app_config'
    export default {
        name: 'TheHeader',
        data() {
            return {
                appName: cfg.APP_CONFIG.app.name,
            }
        },
        components: {
            TheHeaderDropdownAccnt
        }
    }
</script>
