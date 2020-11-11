import axios, { AxiosResponse } from 'axios'
import router from '../router'
import { clearLoginData } from './login'

// CSRF token provided by server when we logged in.
let csrf = ''

// Function called by login when we received a CSRF token from the server.
export function setAxiosCSRF(to: string) { csrf = to }

// Create and configure an axios instance.
const _axios = axios.create({})
_axios.interceptors.request.use(
    config => {
        if (!config.headers) config.headers = {}
        if (config.method === 'get' || config.method === 'GET') {
            config.headers['Cache-Control'] = 'no-cache'
            config.headers['Pragma'] = 'no-cache'
        } else if (csrf) {
            config.headers['X-CSRF-Token'] = csrf
        }
        return config
    },
    error => Promise.reject(error)
)
_axios.interceptors.response.use(
    response => response,
    error => {
        if (
            error.response &&
            error.response.status === 401 &&
            router &&
            router.currentRoute.value &&
            !router.currentRoute.value.matched.some(record => record.meta.public)
        ) {
            clearLoginData()
            console.log(error.response && error.response.status, router && router.currentRoute, error.request)
            router.replace({ path: '/login', query: { redirect: router.currentRoute.value.fullPath } })
        }
        return Promise.reject(error)
    }
)
export default _axios

export type { AxiosResponse }
