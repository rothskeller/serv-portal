import Vue from 'vue'
import Router from 'vue-router'
import store from './store'

import CERTAttendance from './pages/reports/CERTAttendance'
import Event from './pages/Event'
import EventAttendance from './pages/EventAttendance'
import Events from './pages/Events'
import Login from './pages/Login'
import Logout from './pages/Logout'
import People from './pages/People'
import Person from './pages/Person'
import PWReset from './pages/PWReset'
import PWResetToken from './pages/PWResetToken'
import Reports from './pages/Reports'
import Role from './pages/Role'
import Roles from './pages/Roles'

Vue.use(Router)

const router = new Router({
  mode: 'history',
  base: process.env.BASE_URL,
  routes: [
    { path: '/', redirect: '/events' },
    { path: '/events', component: Events },
    { path: '/events/:id', component: Event },
    { path: '/events/:id/attendance', component: EventAttendance },
    { path: '/login', component: Login, meta: { publicAccess: true } },
    { path: '/logout', component: Logout },
    { path: '/password-reset', component: PWReset, meta: { publicAccess: true } },
    { path: '/password-reset/:token', component: PWResetToken, meta: { publicAccess: true } },
    { path: '/people', component: People },
    { path: '/people/:id', component: Person },
    { path: '/reports', component: Reports },
    { path: '/reports/cert-attendance', component: CERTAttendance, props: route => route.query },
    { path: '/roles', component: Roles },
    { path: '/roles/:id', component: Role },
  ],
})

router.beforeEach((to, from, next) => {
  if (!store.state.me && !to.matched.some(record => record.meta.publicAccess))
    next({ path: '/login', query: { redirect: to.fullPath } })
  else
    next()
})

export default router
