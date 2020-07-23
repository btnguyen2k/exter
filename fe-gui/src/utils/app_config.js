const APP_ID = "exter"
const APP_CONFIG = require('@/config.json')
if (!APP_CONFIG.api_client.bo_api_base_url) {
    APP_CONFIG.api_client.bo_api_base_url = process.env.VUE_APP_BO_API_BASE_URL
}
export default {
    APP_ID,
    APP_CONFIG
}
