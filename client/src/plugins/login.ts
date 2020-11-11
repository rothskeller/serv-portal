import { ref } from "vue"
import axios, { AxiosResponse, setAxiosCSRF } from './axios'

export interface LoginData {
    id: number
    informalName: string
    webmaster: boolean
    canAddEvents: boolean
    canAddPeople: boolean
    canSendTextMessages: boolean
    canViewReports: boolean
    canViewRosters: boolean
    csrf: string
}

export const me = ref(null as null | LoginData)

// On startup, check whether already logged in via cookie.
export async function checkLogin() {
    try {
        const resp: AxiosResponse<LoginData> = await axios.get('/api/login')
        me.value = resp.data
        setAxiosCSRF(resp.data.csrf)
    } catch (err) { }
}

// Clear the login data (called when we get a 401).
export function clearLoginData() {
    me.value = null
    setAxiosCSRF('')
}

// Send a login request.  Returns true for success, false for failure.
export async function login(email: string, password: string, remember: boolean): Promise<boolean> {
    const body = new (FormData)
    body.append('username', email)
    body.append('password', password)
    if (remember) body.append('remember', 'true')
    try {
        const resp: AxiosResponse<LoginData> = await axios.post('/api/login', body)
        me.value = resp.data
        setAxiosCSRF(resp.data.csrf)
        return true
    } catch (err) {
        console.error(err)
        return false
    }
}

// Send a password reset request.  Returns true for success, false for failure.
export async function passwordReset(token: string, password: string): Promise<boolean> {
    const body = new (FormData)
    body.append('password', password)
    try {
        const resp: AxiosResponse<LoginData> = await axios.post(`/api/password-reset/${token}`, body)
        me.value = resp.data
        setAxiosCSRF(resp.data.csrf)
        return true
    } catch (err) {
        console.error(err)
        return false
    }
}
