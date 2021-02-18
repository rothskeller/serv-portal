<!--
Events displays the list of events.
-->

<template lang="pug">
TabPage(:tabs='tabs')
  router-view(:onLoadEvent='onLoadEvent')
</template>

<script lang="ts">
import { defineComponent, inject, Ref, computed, ref, watchEffect } from 'vue'
import { LoginData } from '../plugins/login'
import { useRoute } from 'vue-router'
import { TabPage, TabDef } from '../base'
import { GetEventEvent } from './events/EventView.vue'

export default defineComponent({
  components: { TabPage },
  setup() {
    const me = inject<Ref<LoginData>>('me')!
    const route = useRoute()

    const event = ref(null as null | GetEventEvent)
    function onLoadEvent(ev: GetEventEvent) {
      event.value = ev
    }
    watchEffect(() => {
      if (!route.params.id) event.value = null
    })

    // List of tabs on the page.
    const tabs = computed(() => {
      const tabs = [] as Array<TabDef>
      tabs.push({ to: '/events/calendar', label: 'Calendar' })
      tabs.push({ to: '/events/list', label: 'List' })
      tabs.push({ to: '/events/signups', label: 'Signups' })
      if (!route.params.id && me.value.canAddEvents)
        tabs.push({ to: '/events/NEW/edit', label: 'Add Event' })
      if (route.params.id && route.params.id !== 'NEW')
        tabs.push({ to: `/events/${route.params.id}`, label: 'Details' })
      if (event.value && event.value.canEdit)
        tabs.push({
          to: `/events/${route.params.id}/edit`,
          label: route.params.id === 'NEW' ? 'Add Event' : 'Edit',
        })
      if (route.path === `/events/${route.params.id}/timesheet`)
        tabs.push({ to: `/events/${route.params.id}/timesheet`, label: 'Timesheet' })
      else if (event.value && event.value.canViewAttendance)
        tabs.push({ to: `/events/${route.params.id}/attendance`, label: 'Attendance' })
      return tabs
    })

    return { tabs, onLoadEvent }
  },
})
</script>
