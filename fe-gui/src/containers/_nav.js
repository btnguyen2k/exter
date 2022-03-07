export default [
    {
        _name: 'CSidebarNav',
        _children: [
            {
                _name: 'CSidebarNavItem',
                name: 'Dashboard',
                to: {name: "Dashboard"},
                icon: 'cil-speedometer',
                // badge: {
                //     color: 'primary',
                //     text: 'NEW'
                // }
            },
            {
                _name: 'CSidebarNavItem',
                name: 'My Apps',
                to: {name: "MyApps"},
                icon: 'cilApplications',
                exact: false, //[extract=false] to make this item "active" on child-action (create/edit/delete)
            },
        ]
    }
]
