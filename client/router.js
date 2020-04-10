import Vue from 'vue'
import Router from 'vue-router'
import store from './store'
import MainMenu from './MainMenu'
import Admin from './pages/Admin'
import Events from './pages/Events'
import People from './pages/People'
import Texts from './pages/Texts'

Vue.use(Router)

const router = new Router({
  mode: 'history',
  base: process.env.BASE_URL,
  routes: [
    {
      path: '/login',
      component: () => import(/* webpackChunkName: "Login" */ './pages/Login'),
      meta: { allow401: true },
    },
    {
      path: '/password-reset',
      component: () => import(/* webpackChunkName: "PWReset" */ './pages/PWReset'),
    },
    {
      path: '/password-reset/:token',
      component: () => import(/* webpackChunkName: "PWResetToken" */ './pages/PWResetToken'),
    },
    {
      path: '/policies',
      component: () => import(/* webpackChunkName: "Policies" */ './pages/Policies'),
    },
    {
      path: '/unsubscribe/:token/:email?',
      component: () => import(/* webpackChunkName: "Unsubscribe" */ './pages/Unsubscribe'),
    },
    {
      path: '*',
      component: MainMenu,
      children: [
        {
          path: '/admin',
          component: Admin,
          meta: { menuItem: 'admin', tabbed: true },
          children: [
            {
              path: '',
              redirect: '/admin/groups'
            },
            {
              path: 'emails',
              component: () => import(/* webpackChunkName: "EmailsList" */ './pages/admin/EmailsList'),
            },
            {
              path: 'groups',
              component: () => import(/* webpackChunkName: "GroupsList" */ './pages/admin/GroupsList'),
            },
            {
              name: 'groups-gid',
              path: 'groups/:gid',
              component: () => import(/* webpackChunkName: "GroupEdit" */ './pages/admin/GroupEdit'),
            },
            {
              path: 'roles',
              component: () => import(/* webpackChunkName: "RolesList" */ './pages/admin/RolesList'),
            },
            {
              name: 'roles-rid',
              path: 'roles/:rid',
              component: () => import(/* webpackChunkName: "RoleEdit" */ './pages/admin/RoleEdit'),
            },
          ],
        },
        {
          path: '/events',
          component: Events,
          meta: { menuItem: 'events', tabbed: true },
          children: [
            {
              path: '',
              redirect: '/events/calendar'
            },
            {
              path: 'calendar',
              component: () => import(/* webpackChunkName: "EventsCalendar" */ './pages/events/EventsCalendar'),
            },
            {
              path: 'list',
              component: () => import(/* webpackChunkName: "EventsList" */ './pages/events/EventsList'),
            },
            {
              path: ':id',
              component: () => import(/* webpackChunkName: "EventView" */ './pages/event/EventView'),
            },
            {
              path: ':id/edit',
              component: () => import(/* webpackChunkName: "EventEdit" */ './pages/event/EventEdit'),
            },
            {
              path: ':id/attendance',
              component: () => import(/* webpackChunkName: "EventAttendance" */ './pages/event/EventAttendance'),
            },
          ]
        },
        {
          path: '/files',
          redirect: '/files/0',
        },
        {
          path: '/files/:id',
          component: () => import(/* webpackChunkName: "Files" */ './pages/Files'),
          meta: { menuItem: 'files', tabbed: true },
        },
        {
          path: '/help/email-lists',
          component: () => import(/* webpackChunkName: "Help" */ './pages/help/EmailLists'),
        },
        {
          path: '/help/subscribe-calendar',
          component: () => import(/* webpackChunkName: "Help" */ './pages/help/SubscribeCalendar'),
        },
        {
          path: '/logout',
          component: () => import(/* webpackChunkName: "Logout" */ './pages/Logout'),
        },
        {
          path: '/policies-m',
          component: () => import(/* webpackChunkName: "Policies" */ './base/Policies'),
        },
        {
          path: '/people',
          component: People,
          meta: { menuItem: 'people', tabbed: true },
          children: [
            {
              path: '',
              redirect: '/people/list'
            },
            {
              path: 'list',
              component: () => import(/* webpackChunkName: "PeopleList" */ './pages/people/PeopleList'),
            },
            {
              path: 'map',
              component: () => import(/* webpackChunkName: "PeopleMap" */ './pages/people/PeopleMap'),
            },
            {
              path: ':id',
              component: () => import(/* webpackChunkName: "PersonView" */ './pages/person/PersonView'),
            },
            {
              path: ':id/edit',
              component: () => import(/* webpackChunkName: "PersonEdit" */ './pages/person/PersonEdit'),
            },
          ]
        },
        {
          path: '/reports',
          component: () => import(/* webpackChunkName: "Reports" */ './pages/Reports'),
          meta: { menuItem: 'reports' },
        },
        {
          path: '/reports/cert-attendance',
          component: () => import(/* webpackChunkName: "CERTAttendance" */ './pages/reports/CERTAttendance'),
          props: route => route.query,
          meta: { menuItem: 'reports' },
        },
        {
          path: '/search',
          component: () => import(/* webpackChunkName: "Search" */ './pages/Search'),
        },
        {
          path: '/texts',
          component: Texts,
          meta: { menuItem: 'texts', tabbed: true },
          children: [
            {
              path: '',
              component: () => import(/* webpackChunkName: "TextsList" */ './pages/texts/TextsList'),
            },
            {
              path: 'send',
              component: () => import(/* webpackChunkName: "TextsSend" */ './pages/texts/TextsSend'),
            },
            {
              path: ':id',
              component: () => import(/* webpackChunkName: "TextsView" */ './pages/texts/TextsView'),
            }
          ]
        },
        {
          path: '',
          component: () => import(/* webpackChunkName: "Home" */ './pages/Home'),
          meta: { tabbed: true }, // OK not really but this inhibits padding
        },
      ],
      beforeEnter: (to, from, next) => {
        if (!store.state.me)
          next({ path: '/login', query: { redirect: to.fullPath } })
        else
          next()
      }
    }
  ],
})


export default router
