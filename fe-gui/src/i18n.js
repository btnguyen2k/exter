//#GovueAdmin-Customized
import Vue from 'vue'
import VueI18n from 'vue-i18n'

const messages = {
    en: {
        _name: 'English',
        _flag: 'cif-gb',
        // _demo_msg: 'This is instance is for demo purpose only. Data might be reset without notice.' +
        //     '<br/>Login with default account <strong>admin@local/s3cr3t</strong> ' +
        //     'or using your <u>social account</u> via "Login with social account" link ' +
        //     '(the application <u>will not</u> know or store your social account credential).' +
        //     '<br/><br/>You can also access the frontend via <a href="/doc/" style="color: yellowgreen">this link</a>.',

        message: {
            home: 'Home',
            dashboard: 'Dashboard',
            my_apps: 'My Apps',
            my_app_list: 'My app list',
            settings: 'Settings',
            profile: 'Profile',

            actions: 'Actions',
            cancel: 'Cancel',
            save: 'Save',
            back: 'Back',
            delete: 'Delete',

            wait: 'Please wait...',
            wait_login: 'Logging in, please wait...{counter}',
            wait_fbsdk: 'Please wait, Facebook SDK is being loaded.',
            wait_googlesdk: 'Please wait, GoogleAPI SDK is being loaded.',

            login: 'Sign in',
            logout: 'Sign out',
            login_msg: 'Please log in to continue',
            login_facebook: 'Login with Facebook',
            login_github: 'Login with GitHub',
            login_google: 'Login with Google',
            login_linkedin: 'Login with LinkedIn',
            error_login_failed_facebook: 'Facebook login failed.',
            error_login_failed_github: 'GitHub login failed.',
            error_login_failed_google: 'Google login failed.',
            error_login_failed_linkedin: 'LinkedIn login failed.',
            error_invalid_return_url: 'The return URL is invalid',

            error_app_not_exist: 'App "{app}" does not exist.',
            error_app_not_active: 'App "{app}" is not active.',

            error_loading_gapisdk: 'Error while loading GoogleAPI SDK',
            error_loading_fbsdk: 'Error while loading Facebook SDK',

            register_app: 'Register New App',
            app_registered_successful: 'Application "{id}" has been registered successfully.',

            delete_my_app: 'Delete App',
            app_deleted_successful: 'Application "{id}" has been deleted successfully.',

            edit_my_app: 'Edit App',
            app_updated_successful: 'Application "{id}" has been updated successfully.',

            app_id: 'Id',
            app_id_placeholder: "Application's id, must be unique",
            app_id_rule: "Application's id must not be empty, unique and follow format [0-9a-z_]+.",
            app_active: 'Active',
            app_desc: 'Description',
            app_desc_placeholder: "Application's short description",
            app_default_return_url: 'Default return URL',
            app_default_return_url_placeholder: 'User is redirected to this URL after successful authentication.',
            app_default_return_url_rule: 'URL must be in format https://... or http://...',
            app_default_cancel_url: 'Default cancel URL',
            app_default_cancel_url_placeholder: 'User is redirected to this URL when authentication is cancelled.',
            app_default_cancel_url_rule: 'URL must be in format https://... or http://...',
            app_domains: 'Whitelist domains',
            app_domains_placeholder: 'Exter only redirects users to these whitelist domains. Domains separated by spaces, commas or semi-colons.',
            app_auth_provider: 'Authentication provider',
            app_tags: 'Tags',
            app_tags_placeholder: 'Tags separated by commas or semi-colons.',
            app_rsa_pubkey: 'RSA public key',
            app_rsa_pubkey_placeholder: "Application's RSA public key in PEM format.",

            auth_provider_facebook: 'Facebook',
            auth_provider_github: 'GitHub',
            auth_provider_google: 'Google',
            auth_provider_linkedin: 'LinkedIn',
        }
    },
    vi: {
        _name: 'Tiếng Việt',
        _flag: 'cif-vn',
        // _demo_msg: 'Bản triển khai này dành do mục đích trải nghiệm. Dữ liệu có thể được xoá và cài đặt lại nguyên gốc bất cứ lúc nào mà không cần báo trước.' +
        //     '<br/>Đăng nhập với tài khoản <strong>admin@local/s3cr3t</strong>, ' +
        //     'hoặc sử dụng <i>tài khoản mxh</i> bằng cách nhấn vào đường dẫn "Đăng nhập với tài khoản mxh" ' +
        //     '(máy chủ sẽ không biết và cũng không lưu trữ thông tin đăng nhập tài khoản mxh của bạn).' +
        //     '<br/><br/>Bạn có thể truy cập vào trang frontend bằng <a href="/doc/" style="color: yellowgreen">đường dẫn này</a>.',

        message: {
            home: 'Trang nhà',
            dashboard: 'Bảng thông tin',
            my_apps: 'Ứng dụng',
            my_app_list: 'Danh sách ứng dụng',
            settings: 'Cài đặt',
            profile: 'Cá nhân',

            actions: 'Hành động',
            cancel: 'Huỷ bỏ',
            save: 'Lưu',
            back: 'Quay lại',
            delete: 'Xoá',

            wait: 'Vui lòng giờ giây lát...',
            wait_login: 'Đang đăng nhập, vui lòng chờ giây lát...{counter}',
            wait_fbsdk: 'Vui lòng chờ giây lát, đang tải Facebook SDK .',
            wait_googlesdk: 'Vui lòng chờ giây lát, đang tải GoogleAPI SDK.',

            login: 'Đăng nhập',
            logout: 'Đăng xuất',
            login_msg: 'Vui lòng đăng nhập',
            login_facebook: 'Đăng nhập với tài khoản Facebook',
            login_github: 'Đăng nhập với tài khoản GitHub',
            login_google: 'Đăng nhập với tài khoản Google',
            login_linkedin: 'Đăng nhập với tài khoản LinkedIn',
            error_login_failed_facebook: 'Đăng nhập với tài khoản Facebook không thành công.',
            error_login_failed_github: 'Đăng nhập với tài khoản GitHub không thành công.',
            error_login_failed_google: 'Đăng nhập với tài khoản Google không thành công.',
            error_login_failed_linkedin: 'Đăng nhập với tài khoản LinkedIn không thành công.',
            error_invalid_return_url: 'URL chuyển tiếp không hợp lệ',

            error_app_not_exist: 'Ứng dụng "{app}" không tồn tại.',
            error_app_not_active: 'Ứng dụng "{app}" không ở trạng thái "có hiệu lực".',

            error_loading_gapisdk: 'Có lỗi khi tải GoogleAPI SDK',
            error_loading_fbsdk: 'Có lỗi khi tải Facebook SDK',

            register_app: 'Đăng Ký Ứng Dụng',
            app_registered_successful: 'Ứng dụng "{id}" đã được đăng ký thành công.',

            delete_my_app: 'Xoá Ứng Dụng',
            app_deleted_successful: 'Ứng dụng "{id}" đã được xoá thành công.',

            edit_my_app: 'Cập Nhật Ứng Dụng',
            app_edit_successful: 'Ứng dụng "{id}" đã được cập nhật thành công.',

            app_id: 'Định danh',
            app_id_placeholder: "Định danh ứng dụng, không được trùng lắp",
            app_id_rule: "Định danh ứng dụng không được rỗng hoặc trùng lắp, và phải theo định dạng [0-9a-z_]+.",
            app_active: 'Có hiệu lực',
            app_desc: 'Mô tả',
            app_desc_placeholder: 'Thông tin ngắn gọn về ứng dụng',
            app_default_return_url: 'URL xác thực',
            app_default_return_url_placeholder: 'URL được gọi sau khi xác thực thành công.',
            app_default_return_url_rule: 'URL xác thực phải ở dạng https://... hoặc http://...',
            app_default_cancel_url: 'URL huỷ',
            app_default_cancel_url_placeholder: 'URL được gọi khi user huỷ quá trình xác thực.',
            app_default_cancel_url_rule: 'URL huỷ phải ở dạng https://... hoặc http://...',
            app_domains: 'Danh sách tên miền',
            app_domains_placeholder: 'Exter chỉ gọi URL nằm trong danh sách tên miền. Các tên miền phân cách nhau bằng khoảng trắng, dấu phảy (,) hoặc chấm phảy (;).',
            app_auth_provider: 'Nguồn đăng nhập',
            app_tags: 'Thẻ',
            app_tags_placeholder: 'Các thẻ phân cách nhau bằng dấu phảy (,) hoặc chấm phảy (;).',
            app_rsa_pubkey: 'Mã công khai RSA',
            app_rsa_pubkey_placeholder: 'Mã công khai RSA của ứng dụng ở định dạng PEM.',

            auth_provider_facebook: 'Facebook',
            auth_provider_github: 'GitHub',
            auth_provider_google: 'Google',
            auth_provider_linkedin: 'LinkedIn',
        }
    }
}

Vue.use(VueI18n)

const i18n = new VueI18n({
    locale: 'en',
    messages: messages
})

export default i18n
