<!--
Main is the top application component for the client.
-->

<template lang="pug">
#page-top(:class='topClasses')
  #page-heading
    #page-menu-trigger-box(v-if='me', @click='onMenuTrigger')
      SIcon#page-menu-trigger(icon='menu')
    #page-menu-spacer(v-else)
    #page-titlebox
      #page-title(v-text='page.title || "Sunnyvale SERV"')
    #page-search-box(v-if='me', @click='onSearch')
      SIcon#page-search(icon='search')
    #page-menu-spacer(v-else)
  #page-menu(v-if='me')
    #page-menu-welcome
      | Welcome
      br
      b(v-text='me.informalName')
    ul#page-nav
      template(v-for='mi in menuItems')
        li(v-if='mi.show')
          router-link.page-nav-item(
            :to='mi.to',
            v-text='mi.label',
            :class='{ "page-nav-active": mi.active }',
            @click='onMenuClick'
          )
    router-link#page-policies(to='/policies', @click='onMenuClick') Policies/Legal
  Content(:contentClasses='contentClasses')
#modal-port
</template>

<script lang="ts">
import { defineComponent, provide, reactive, ref, computed, h } from 'vue'
import { me, LoginData } from './plugins/login'
import { PageData } from './plugins/page'
import provideSize from './plugins/size'
import { touch } from './plugins/touch'
import { useRoute, useRouter, RouterView } from 'vue-router'
import SIcon from './base/SIcon.vue'

const Content = defineComponent({
  props: {
    contentClasses: Object,
  },
  setup(props) {
    provideSize()
    return () =>
      h(
        'div',
        {
          id: 'page-content',
          class: props.contentClasses,
        },
        [h(RouterView)]
      )
  },
})

export default defineComponent({
  components: { Content, SIcon },
  setup() {
    provide('me', me)
    provide('touch', touch)

    // Set up for pages to provide their title, subtitle, etc.
    const page: PageData = reactive({
      title: '',
      browserTitle: '',
      subtitle: '',
      padding: false,
      menuItem: '',
    })
    function setPage(data: PageData) {
      page.title = data.title
      page.browserTitle = data.browserTitle || data.title
      document.title = page.browserTitle
        ? `${page.browserTitle} | Sunnyvale SERV`
        : 'Sunnyvale SERV'
      page.subtitle = data.subtitle || ''
      page.padding = data.padding !== undefined ? data.padding : !!page.subtitle
      page.menuItem = data.menuItem || ''
    }
    provide('setPage', setPage)

    // Flag indicating that the menu has been opened (on small devices where it
    // isn't permanently open).
    const menuOpen = ref(false)
    function onMenuTrigger() {
      menuOpen.value = !menuOpen.value
    }
    function onMenuClick() {
      menuOpen.value = false
    }

    // Search button.
    const router = useRouter()
    function onSearch() {
      router.push('/search')
    }

    // Top-level classes for the page.
    const topClasses = reactive({
      mouse: computed(() => !touch.value),
      touch: touch,
      'page-menu-open': menuOpen,
    })
    const contentClasses = reactive({
      'page-no-padding': computed(() => !page.padding),
      'page-no-menu': computed(() => !me.value),
    })

    // Menu contents.
    const route = useRoute()
    const menuItems = computed(() => [
      {
        label: 'Events',
        to: '/events',
        active: page.menuItem === 'events',
        show: true,
      },
      {
        label: 'People',
        to: '/people',
        active: page.menuItem === 'people' && !route.path.startsWith(`/people/${me.value?.id}/`),
        show: me.value?.canViewRosters,
      },
      {
        label: 'Files',
        to: '/files',
        active: page.menuItem === 'files',
        show: true,
      },
      {
        label: 'Reports',
        to: '/reports',
        active: page.menuItem === 'reports',
        show: me.value?.canViewReports,
      },
      {
        label: 'Texts',
        to: '/texts',
        active: page.menuItem === 'texts',
        show: me.value?.canSendTextMessages,
      },
      {
        label: 'Admin',
        to: '/admin',
        active: page.menuItem === 'admin',
        show: me.value?.webmaster,
      },
      {
        label: 'Profile',
        to: `/people/${me.value?.id}`,
        active: route.path == `/people/${me.value?.id}`,
        show: !!me.value?.id,
      },
      {
        label: 'Logout',
        to: '/logout',
        active: false,
        show: true,
      },
    ])

    return { me, page, topClasses, contentClasses, onMenuTrigger, onMenuClick, menuItems, onSearch }
  },
})
</script>

<style lang="postcss">
:root {
  --titlebarHeight: 40px;
  --sidebarWidth: 7em;
}
#page-top {
  display: grid;
  height: 100vh;
  grid: 'header header' var(--titlebarHeight) 'menu content' 1fr / var(--sidebarWidth) 1fr;
  @media print {
    display: block;
    height: auto;
  }
}
#page-heading {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background-color: #006600;
  color: #fff;
  grid-area: header;
  @media print {
    display: none;
  }
}
#page-menu-trigger-box,
#page-search-box {
  flex: none;
  margin: 0 6px;
  width: 2rem;
  cursor: pointer;
  user-select: none;
}
#page-menu-trigger {
  @media (min-width: 576px) {
    display: none;
  }
}
#page-search {
  width: 1.5rem;
}
#page-titlebox {
  display: flex;
  flex: auto;
  flex-direction: column;
  text-align: center;
}
#page-title {
  overflow: hidden;
  width: 100%;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-size: 1.5rem;
}
#page-menu {
  display: none;
  flex-direction: column;
  overflow: visible;
  border-right: 1px solid #888;
  background-color: #ccc;
  grid-area: menu;
  .page-menu-open & {
    z-index: 1;
    display: flex;
    flex-direction: column;
  }
  @media (min-width: 576px) {
    display: flex;
    flex-direction: column;
  }
  @media print {
    display: none;
  }
}
#page-menu-welcome {
  margin-bottom: 0.5rem;
  padding: 0.75rem;
  border-bottom: 1px solid white;
  text-align: center;
  white-space: nowrap;
  font-size: 0.75rem;
}
#page-nav {
  flex-direction: column;
  display: flex;
  flex-wrap: wrap;
  margin-bottom: 0;
  list-style: none;
  padding: 0 0.5rem;
  font-size: 1.25rem;
}
.page-nav-item {
  padding: 0.125rem 0.5rem;
  color: black;
  border-radius: 0.25rem;
  display: block;
  &:hover {
    text-decoration: none;
  }
}
.page-nav-active {
  color: #fff;
  background-color: #007bff;
}
#page-policies {
  display: flex;
  flex: 1 1 auto;
  flex-direction: column;
  justify-content: flex-end;
  align-self: center;
  margin-bottom: 0.5rem;
  font-size: 0.75rem;
}
#page-content {
  overflow: auto;
  padding: 1.5rem 0.75rem;
  grid-area: 2 / 1 / 3 / 3;
  @media (min-width: 576px) {
    grid-area: content;
    &.page-no-menu {
      grid-area: 2 / 1 / 3 / 3;
    }
  }
  &.page-no-padding {
    padding: 0;
  }
}
#page-subtitle {
  font-size: 1.5rem;
}
#modal-port {
  position: absolute;
  z-index: 1040;
}
</style>
