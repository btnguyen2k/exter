//#GovueAdmin-Customized
import Vue from 'vue'
import VueI18n from 'vue-i18n'

const messages = {
    en: {
        _name: 'English',
        _flag: 'cif-gb',
        _demo_msg: 'This is instance is for demo purpose only. Data might be reset without notice.' +
            '<br/>Login with default account <strong>admin@local/s3cr3t</strong> ' +
            'or using your <u>social account</u> via "Login with social account" link ' +
            '(the application <u>will not</u> know or store your social account credential).' +
            '<br/><br/>You can also access the frontend via <a href="/doc/" style="color: yellowgreen">this link</a>.',

        message: {
            actions: 'Actions',
            action_create: 'Create',
            action_save: 'Save',
            action_back: 'Back',
            action_edit: 'Edit',
            action_delete: 'Delete',
            action_move_down: 'Move down',
            action_move_up: 'Move up',

            error: 'Error',
            info: 'Information',
            settings: 'Settings',
            account: 'Account',
            language: 'Language',

            error_field_mandatory: 'Field value is mandatory.',

            yes: 'Yes',
            no: 'No',
            ok: 'Ok',
            cancel: 'Cancel',
            close: 'Close',

            logout: 'Logout',
            login: 'Login',
            login_info: 'Please sign in to continue',
            login_social: 'Login with social account',
            wait: 'Please wait...',
            error_parse_login_token: 'Error parsing login-token',
            home: 'Home',
            dashboard: 'Dashboard',

            icons: "Icons",
            icon_icon: "Icon",
            icon_name: "Name",

            products: "Products",
            topics: "Topics",
            pages: "Pages",

            product_info: "Information",
            product_id: "Id",
            product_id_msg: "Must be unique, will be generated if empty",
            product_is_published: "Published",
            product_is_published_msg: "Product's documents are visible only when published",
            product_name: "Name",
            product_name_msg: "Display name of the product",
            product_desc: "Description",
            product_desc_msg: "Summary description of the product",
            product_domains: "Domain names",
            product_domains_msg: "Product documents are accessible via these domain names (one domain per line)",

            product_contacts: "Contacts",
            product_website: "Website",
            product_website_msg: "Product website",
            product_email: "Email",
            product_email_msg: "Contact email address",
            product_github: "Github",
            product_github_msg: "Url of GitHub page",
            product_facebook: "Facebook",
            product_facebook_msg: "Url of Facebook page",
            product_linkedin: "LinkedIn",
            product_linkedin_msg: "Url of LinkedIn page",
            product_slack: "Slack",
            product_slack_msg: "Url of Slack channel",
            product_twitter: "Twitter",
            product_twitter_msg: "Url of Twitter page",

            error_product_not_found: 'Product "{id}" not found.',

            add_product: "Add new product",
            product_added_msg: 'Product "{name}" has been created successfully.',

            delete_product: "Delete product",
            delete_product_msg: 'All topics and pages belong to the product will also be deleted! This product currently has {numTopics} topic(s). Are you sure you wish to delete the product?',
            product_deleted_msg: 'Product "{name}" has been deleted successfully.',

            edit_product: "Edit product",
            product_updated_msg: 'Product "{name}" has been updated successfully.',

            product_unmap_domain: "Unmap",
            product_unmap_domain_msg: 'Are you sure you wish to unmap domain name "{domain}"? Product documents are no longer accessible via this domain name once unmapped.',
            product_domain_unmapped_msg: 'Domain "{domain}" has been unmapped successfully.',
            product_map_domain: "Map",
            product_map_domain_msg: "Product documents are accessible via mapped domain names",
            product_domain_mapped_msg: 'Domain "{domain}" has been mapped successfully.',

            error_topic_not_found: 'Topic "{id}" not found.',

            add_topic: "Add new topic",
            topic_added_msg: 'Topic "{name}" has been created successfully.',

            delete_topic: "Delete topic",
            delete_topic_msg: 'All pages belong to the topic will also be deleted! This topic currently has {numPages} page(s). Are you sure you wish to delete the topic?',
            topic_deleted_msg: 'Topic "{name}" has been deleted successfully.',

            edit_topic: "Edit topic",
            topic_updated_msg: 'Topic "{name}" has been updated successfully.',

            topic_id: "Id",
            topic_id_msg: "Must be unique, will be generated if empty",
            topic_icon: "Icon",
            topic_icon_msg: "Pick an icon from the list",
            topic_title: "Title",
            topic_title_msg: "Topic title",
            topic_summary: "Summary",
            topic_summary_msg: "Short summary of the topic",

            add_page: "Add new page",
            page_added_msg: 'Page "{name}" has been created successfully.',

            delete_page: "Delete page",
            delete_page_msg: 'Are you sure you wish to delete this page?',
            page_deleted_msg: 'Page "{name}" has been deleted successfully.',

            edit_page: "Edit page",
            page_updated_msg: 'Page "{name}" has been updated successfully.',

            page_id: "Id",
            page_id_msg: "Must be unique, will be generated if empty",
            page_icon: "Icon",
            page_icon_msg: "Pick an icon from the list",
            page_title: "Title",
            page_title_msg: "Page title",
            page_summary: "Summary",
            page_summary_msg: "Short summary of the page",
            page_content: "Content",
            page_content_msg: "Page content (Markdown supported)",

            content_editor: "Editor",
            content_preview: "Preview",

            users: "Users",
            my_profile: "My profile",
            user_is_admin: "Is administrator",
            user_is_admin_msg: "Administrator has permission to create other administrator accounts",
            user_id: "Login id",
            user_id_msg: "User email address as login id, must be unique",
            user_mask_id: "Mask id",
            user_mask_id_msg: "User mask-id to not expose user email address, must be unique",
            user_display_name: "Display name",
            user_display_name_msg: "Name of user for displaying purpose",
            user_password: "Password",
            user_current_password: "Current password",
            user_current_password_msg: "Change password: enter current password and the new one",
            user_new_password: "New password",
            user_new_password_msg: "Enter new password. If empty, user can only login using social network account",
            user_confirmed_password: "Confirmed password",
            user_confirmed_password_msg: "Enter new password again to confirm",
            error_confirmed_password_mismatch: "Password does not match the confirmed one.",

            edit_user_profile: "Edit user profile",
            user_profile_updated_msg: 'User profile "{id}" has been updated successfully.',
            user_password_updated_msg: 'Password of user "{id}" has been updated successfully.',

            add_user: "Add new user",
            user_added_msg: 'User "{id}" has been added successfully.',

            delete_user: "Delete user account",
            delete_user_msg: 'Are you sure you wish to remove user account "{id}"?',
            user_deleted_msg: 'User "{id}" has been deleted successfully.',
        }
    },
    vi: {
        _name: 'Tiếng Việt',
        _flag: 'cif-vn',
        _demo_msg: 'Bản triển khai này dành do mục đích trải nghiệm. Dữ liệu có thể được xoá và cài đặt lại nguyên gốc bất cứ lúc nào mà không cần báo trước.' +
            '<br/>Đăng nhập với tài khoản <strong>admin@local/s3cr3t</strong>, ' +
            'hoặc sử dụng <i>tài khoản mxh</i> bằng cách nhấn vào đường dẫn "Đăng nhập với tài khoản mxh" ' +
            '(máy chủ sẽ không biết và cũng không lưu trữ thông tin đăng nhập tài khoản mxh của bạn).' +
            '<br/><br/>Bạn có thể truy cập vào trang frontend bằng <a href="/doc/" style="color: yellowgreen">đường dẫn này</a>.',

        message: {
            actions: 'Hành động',
            action_create: 'Tạo',
            action_save: 'Lưu',
            action_back: 'Quay lại',
            action_edit: 'Sửa',
            action_delete: 'Xoá',
            action_move_down: 'Chuyển xuống',
            action_move_up: 'Chuyển lên',

            error: 'Có lỗi',
            info: 'Thông tin',
            settings: 'Cài đặt',
            account: 'Tài khoản',
            language: 'Ngôn ngữ',

            yes: 'Có',
            no: 'Không',
            ok: 'Đồng ý',
            cancel: 'Huỷ',
            close: 'Đóng',

            error_field_mandatory: 'Trường dữ liệu không được bỏ trống, vui lòng nhập thông tin.',

            logout: 'Đăng xuất',
            login: 'Đăng nhập',
            login_info: 'Đăng nhập để tiếp tục',
            login_social: 'Đăng nhập với tài khoản mxh',
            wait: 'Vui lòng giờ giây lát...',
            error_parse_login_token: 'Có lỗi khi xử lý login-token',
            home: 'Trang nhà',
            dashboard: 'Tổng hợp',

            icons: "Biểu tượng",
            icon_icon: "Biểu tượng",
            icon_name: "Tên",

            products: "Sản phẩm",
            topics: "Chủ đề",
            pages: "Trang tài liệu",

            product_info: "Thông tin chung",
            product_id: "Id",
            product_id_msg: "Id phải là duy nhất, sẽ được tự động tạo nếu để rỗng",
            product_is_published: "Đăng tải",
            product_is_published_msg: "Tài liệu của sản phẩm chỉ xem được khi trạng thái là 'Đăng tải'",
            product_name: "Tên",
            product_name_msg: "Tên hiển thị của sản phẩm",
            product_desc: "Mô tả",
            product_desc_msg: "Mô tả ngắn về sản phẩm",
            product_domains: "Tên miền",
            product_domains_msg: "Tài liệu của sản phẩm truy cập được từ các tên miền này (mỗi tên miền 1 dòng)",

            product_contacts: "Thông tin liên hệ",
            product_website: "Website",
            product_website_msg: "Địa chỉ trang web của sản phẩm",
            product_email: "Email",
            product_email_msg: "Địa chỉ email liên hệ",
            product_github: "Github",
            product_github_msg: "Trang GitHub của sản phẩm",
            product_facebook: "Facebook",
            product_facebook_msg: "Trang Facebook của sản phẩm",
            product_linkedin: "LinkedIn",
            product_linkedin_msg: "Trang LinkedIn của sản phẩm",
            product_slack: "Slack",
            product_slack_msg: "Nhóm Slack chat của sản phẩm",
            product_twitter: "Twitter",
            product_twitter_msg: "Trang Twitter của sản phẩm",

            error_product_not_found: 'Không tìm thấy sản phẩm "{id}".',

            add_product: "Thêm sản phẩm",
            product_added_msg: 'Sản phẩm "{name}" đã được tạo thành công.',

            delete_product: "Xoá sản phẩm",
            delete_product_msg: 'Xoá sản phẩm sẽ xoá các chủ đề và trang tài liệu của sản phẩm! Sản phẩm này hiện có {numTopics} chủ đề. Bạn có chắc muốn xoá sản phẩm này?',
            product_deleted_msg: 'Sản phẩm "{name}" đã được xoá thành công.',

            edit_product: "Chỉnh sửa sản phẩm",
            product_updated_msg: 'Sản phẩm "{name}" đã được cập nhật thành công.',

            product_unmap_domain: "Bỏ kết nối",
            product_unmap_domain_msg: 'Bạn có chắc bỏ kết nối tên miền "{domain}"? Tài liệu của sản phẩm sẽ không còn truy cập được qua tên miền này sau khi bỏ kết nối.',
            product_domain_unmapped_msg: 'Kết nối với tên miền "{domain}" đã được bỏ thành công.',
            product_map_domain: "Kết nối",
            product_map_domain_msg: "Tài liệu của sản phẩm truy cập được từ các tên miền sau khi được kết nối",
            product_domain_mapped_msg: 'Tên miền "{domain}" đã được kết nối thành công.',

            error_topic_not_found: 'Không tìm thấy chủ để "{id}".',

            add_topic: "Thêm chủ đề",
            topic_added_msg: 'Chủ đề "{name}" đã được tạo thành công.',

            delete_topic: "Xoá chủ đề",
            delete_topic_msg: 'Xoá chủ đề sẽ xoá các trang tài liệu nằm trong chủ đề! Chủ đề này hiện có {numPages} tragn tài liệu. Bạn có chắc muốn xoá chủ đề này?',
            topic_deleted_msg: 'Chủ đề "{name}" đã được xoá thành công.',

            edit_topic: "Chỉnh sửa chủ đề",
            topic_updated_msg: 'Chủ đề "{name}" đã được cập nhật thành công.',

            topic_id: "Id",
            topic_id_msg: "Id phải là duy nhất, sẽ được tự động tạo nếu để rỗng",
            topic_icon: "Biểu tượng",
            topic_icon_msg: "Chọn 1 biểu tượng cho chủ đề trong danh sách",
            topic_title: "Tên",
            topic_title_msg: "Tên hiển thị của chủ đề",
            topic_summary: "Tóm tắt",
            topic_summary_msg: "Phần tóm tắt ngắn về chủ đề",

            add_page: "Thêm trang tài liệu",
            page_added_msg: 'Trang tài liệu "{name}" đã được tạo thành công.',

            delete_page: "Xoá trang tài liệu",
            delete_page_msg: 'Bạn có chắc muốn xoá trang tài liệu này?',
            page_deleted_msg: 'Trang tài liệu "{name}" đã được xoá thành công.',

            edit_page: "Chỉnh sửa trang tài liệu",
            page_updated_msg: 'Trang tài liệu "{name}" đã được cập nhật thành công.',

            page_id: "Id",
            page_id_msg: "Id phải là duy nhất, sẽ được tự động tạo nếu để rỗng",
            page_icon: "Biểu tượng",
            page_icon_msg: "Chọn 1 biểu tượng cho trang tài liệu trong danh sách",
            page_title: "Tên",
            page_title_msg: "Tên hiển thị của trang tài liệu",
            page_summary: "Tóm tắt",
            page_summary_msg: "Phần tóm tắt ngắn về nội dung của trang tài liệu",
            page_content: "Nội dung",
            page_content_msg: "Nội dung của trang tài liệu (hỗ trợ Markdown)",

            content_editor: "Soạn thảo",
            content_preview: "Xem trước",

            users: "Người dùng",
            my_profile: "Thông tin cá nhân",
            user_is_admin: "Quản trị viên",
            user_is_admin_msg: "Quản trị viên sẽ được quyền tạo thêm tài khoản quản trị viên khác",
            user_id: "Tên đăng nhập",
            user_id_msg: "Sử dụng địa chỉ email làm Tên đăng nhập, phải là duy nhất trên hệ thống",
            user_mask_id: "Mask id",
            user_mask_id_msg: "Mask-id sẽ được sử dụng để tránh hiển thị địa chỉ email, phải là duy nhất trên hệ thống",
            user_display_name: "Tên hiển thị",
            user_display_name_msg: "Tên của người dùng",
            user_password: "Mật mã",
            user_current_password: "Mật mã hiện tại",
            user_current_password_msg: "Để đổi mật mã: nhập mật mã hiện tại và mật mã mới",
            user_new_password: "Mật mã mới",
            user_new_password_msg: "Nhập mật mã mới. Nếu rỗng, người dùng chỉ có thể đăng nhập thông qua tài khoản mxh",
            user_confirmed_password: "Xác nhận lại mật mã",
            user_confirmed_password_msg: "Nhập mật mã mới lần nữa để xác nhận",
            error_confirmed_password_mismatch: "Mật mã không khớp nhau",

            edit_user_profile: "Thay đổi thông tin",
            user_profile_updated_msg: 'Thông tin người dùng "{id}" đã được cập nhật thành công.',
            user_password_updated_msg: 'Mật mã người dùng "{id}" đã được cập nhật thành công.',

            add_user: "Thêm tài khoản người dùng",
            user_added_msg: 'Người dùng "{id}" đã được thêm vào hệ thống.',

            delete_user: "Xoá tài khoản người dùng",
            delete_user_msg: 'Bạn có chắc muốn xoá tài khoản "{id}" khỏi hệ thống?',
            user_deleted_msg: 'Tài khoản người dùng "{id}" đã được xoá khỏi hệ thống.',
        }
    }
}

Vue.use(VueI18n)

const i18n = new VueI18n({
    locale: 'en',
    messages: messages
})

export default i18n
