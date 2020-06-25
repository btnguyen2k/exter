export default [
    {
        _name: 'CSidebarNav',
        _children: [
            {
                _name: 'CSidebarNavItem',
                name: 'Dashboard',
                to: '/dashboard',
                icon: 'cil-speedometer',
                // badge: {
                //     color: 'primary',
                //     text: 'NEW'
                // }
            },
            {
                _name: 'CSidebarNavItem',
                name: 'My Apps',
                to: '/myapps',
                icon: 'cilApplications',
            },
        ]
    }
]
