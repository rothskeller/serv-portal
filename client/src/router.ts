import Cookies from 'js-cookie'
import { createRouter, createWebHistory } from 'vue-router'
import { me } from './plugins/login'
import Admin from './pages/Admin.vue'
import Events from './pages/Events.vue'
import People from './pages/People.vue'
import Texts from './pages/Texts.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '',
      component: () => import('./pages/Public.vue'),
      meta: { public: true },
      beforeEnter: (to, from, next) => {
        if (me.value)
          next({ path: '/home' })
        else
          next()
      }
    },
    {
      path: '/admin',
      component: Admin,
      children: [
        {
          path: '',
          redirect: '/admin/groups'
        },
        {
          path: 'groups',
          component: () => import('./pages/admin/GroupsList.vue'),
        },
        {
          name: 'groups-gid',
          path: 'groups/:gid',
          component: () => import('./pages/admin/GroupEdit.vue'),
        },
        {
          path: 'roles',
          component: () => import('./pages/admin/RolesList.vue'),
        },
        {
          name: 'roles-rid',
          path: 'roles/:rid',
          component: () => import('./pages/admin/RoleEdit.vue'),
        },
      ],
    },
    {
      path: '/events',
      component: Events,
      children: [
        {
          path: '',
          redirect: () => {
            const page = Cookies.get('serv-events-page')
            return page ? `/events/${page}` : '/events/calendar'
          },
        },
        {
          path: 'calendar',
          component: () => import('./pages/events/EventsCalendar.vue'),
        },
        {
          path: 'list',
          component: () => import('./pages/events/EventsList.vue'),
        },
        {
          path: ':id',
          component: () => import('./pages/events/EventView.vue'),
        },
        {
          path: ':id/edit',
          component: () => import('./pages/events/EventEdit.vue'),
        },
        {
          path: ':id/attendance',
          component: () => import('./pages/events/EventAttendance.vue'),
        },
      ]
    },
    {
      path: '/files',
      redirect: '/files/0',
    },
    {
      path: '/files/:id',
      component: () => import('./pages/Files.vue'),
    },
    {
      path: '/home',
      component: () => import('./pages/Home.vue'),
    },
    {
      path: '/login',
      component: () => import('./pages/Login.vue'),
      meta: { public: true },
    },
    {
      path: '/logout',
      component: () => import('./pages/Logout.vue'),
    },
    {
      path: '/password-reset',
      component: () => import('./pages/PWReset.vue'),
      meta: { public: true },
    },
    {
      path: '/password-reset/:token',
      component: () => import('./pages/PWResetToken.vue'),
      meta: { public: true },
    },
    {
      path: '/people',
      component: People,
      children: [
        {
          path: '',
          redirect: '/people/list'
        },
        {
          path: 'list',
          component: () => import('./pages/people/PeopleList.vue'),
        },
        {
          path: 'map',
          component: () => import('./pages/people/PeopleMap.vue'),
        },
        {
          path: ':id',
          component: () => import('./pages/people/PersonView.vue'),
        },
        {
          path: ':id/edit',
          component: () => import('./pages/people/PersonEdit.vue'),
        },
        {
          path: ':id/hours',
          component: () => import('./pages/people/PersonHours.vue'),
        },
      ]
    },
    {
      path: '/policies',
      component: () => import('./pages/static/Policies.vue'),
      meta: { public: true },
    },
    {
      path: '/search',
      component: () => import('./pages/Search.vue'),
    },
    {
      path: '/static/email-lists',
      component: () => import('./pages/static/EmailLists.vue'),
    },
    {
      path: '/static/subscribe-calendar',
      component: () => import('./pages/static/SubscribeCalendar.vue'),
    },
    {
      path: '/texts',
      component: Texts,
      children: [
        {
          path: '',
          component: () => import('./pages/texts/TextsList.vue'),
        },
        {
          path: 'send',
          component: () => import('./pages/texts/TextsSend.vue'),
        },
        {
          path: ':id',
          component: () => import('./pages/texts/TextsView.vue'),
        }
      ]
    },
    {
      path: '/unsubscribe/:token/:email?',
      component: () => import('./pages/Unsubscribe.vue'),
    },
    {
      path: '/volunteer-hours/:id',
      component: () => import('./pages/people/PersonHours.vue'),
    },
    {
      path: '/:path(disaster-preparedness|preparacion-para-desastres|neighborhood-organization|disaster-response|emergency-communications)',
      component: () => import('./pages/public/PublicFolder.vue'),
      meta: { public: true },
    },
    {
      path: '/:path(disaster-preparedness|preparacion-para-desastres|neighborhood-organization|disaster-response|emergency-communications)/:rest(.*)',
      component: () => import('./pages/public/PublicFolder.vue'),
      meta: { public: true },
    },
    {
      path: '/:catchAll(.*)',
      component: () => import('./pages/NotFound.vue'),
      meta: { public: true },
    },
  ]
})
router.beforeEach((to, from, next) => {
  if (me.value || to.meta.public) next()
  else next({ path: '/login', query: { redirect: to.fullPath } })
})

export default router
