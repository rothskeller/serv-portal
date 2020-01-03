import Vue from 'vue'
import Router from 'vue-router'

import CERTAttendance from './pages/reports/CERTAttendance'
import Event from './pages/Event'
import EventAttendance from './pages/EventAttendance'
import Events from './pages/Events'
import Logout from './pages/Logout'
import People from './pages/People'
import Person from './pages/Person'
import Reports from './pages/Reports'
import Role from './pages/Role'
import Team from './pages/Team'
import Teams from './pages/Teams'

Vue.use(Router)

export default new Router({
  mode: 'history',
  base: process.env.BASE_URL,
  routes: [
    { path: '/', redirect: '/events' },
    { path: '/events', component: Events },
    { path: '/events/:id', component: Event },
    { path: '/events/:id/attendance', component: EventAttendance },
    { path: '/logout', component: Logout },
    { path: '/people', component: People },
    { path: '/people/:id', component: Person },
    { path: '/reports', component: Reports },
    { path: '/reports/cert-attendance', component: CERTAttendance, props: route => route.query },
    { path: '/teams', component: Teams },
    { path: '/teams/:id', component: Team, props: route => ({ parent: route.query.parent }) },
    { path: '/teams/:tid/roles/:rid', component: Role },
  ],
})
