import i18n from '../i18n'

export default [
    {
        _name: 'CSidebarNav',
        _children: [
            {
                _name: 'CSidebarNavItem',
                name: i18n.t('message.dashboard'),
                to: {name: "Dashboard"},
                icon: 'cil-speedometer',
                // badge: {
                //     color: 'primary',
                //     text: 'NEW'
                // }
            },
            {
                _name: 'CSidebarNavItem',
                name: i18n.t('message.myapps'),
                to: {name: "MyApps"},
                icon: 'cilApplications',
                exact: false, //[extract=false] to make this item "active" on child-action (create/edit/delete)
            },
        ]
    }
]
