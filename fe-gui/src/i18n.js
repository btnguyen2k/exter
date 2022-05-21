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
            myapps: 'My Apps',
            settings: 'Settings',
            profile: 'Profile',

            actions: 'Actions',
            cancel: 'Cancel',

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

            app_id: 'Id',
            app_desc: 'Descriptions',
            app_id_src: 'Id sources',
            app_tags: 'Tags',
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
            myapps: 'Ứng dụng',
            settings: 'Cài đặt',
            profile: 'Cá nhân',

            actions: 'Hành động',
            cancel: 'Huỷ bỏ',

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

            app_id: 'Id',
            app_desc: 'Mô tả',
            app_id_src: 'Nguồn đăng nhập',
            app_tags: 'Thẻ',
        }
    }
}

Vue.use(VueI18n)

const i18n = new VueI18n({
    locale: 'en',
    messages: messages
})

export default i18n
