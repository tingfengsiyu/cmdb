import Vue from 'vue'
import axios from 'axios'

let Url = 'http://172.16.2.6:5000/api/v1/'

axios.defaults.baseURL = Url

axios.interceptors.request.use(config => {
    config.headers.Authorization = `Bearer ${localStorage.getItem('token')}`
    return config
})

Vue.prototype.$http = axios

export { Url }
