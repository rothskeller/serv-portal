<!--
People displays the People module.
-->

<template lang="pug">
TabPage(:tabs='tabs')
  router-view(:onLoadPerson='onLoadPerson')
</template>

<script lang="ts">
import { defineComponent, inject, Ref, computed, ref, watchEffect } from 'vue'
import { LoginData } from '../plugins/login'
import { useRoute } from 'vue-router'
import { TabPage, TabDef } from '../base'
import type { GetPerson } from './people/PersonView.vue'

export default defineComponent({
  components: { TabPage },
  setup() {
    const me = inject<Ref<LoginData>>('me')!
    const route = useRoute()

    const person = ref(null as null | GetPerson)
    function onLoadPerson(ev: GetPerson) {
      person.value = ev
    }
    watchEffect(() => {
      if (!route.params.id) person.value = null
    })

    // List of tabs on the page.
    const tabs = computed(() => {
      const tabs = [] as Array<TabDef>
      if (me.value)
        tabs.push({ to: '/people/list', label: 'List' })
      if (me.value)
        tabs.push({ to: '/people/map', label: 'Map' })
      /* not implemented:
      if (me.value && !route.params.id && me.value.canAddPeople)
        tabs.push({ to: '/people/NEW/edit', label: 'Add Person' })
      */
      if (me.value && route.params.id && route.params.id !== 'NEW')
        tabs.push({ to: `/people/${route.params.id}`, label: 'Details' })
      if (person.value && person.value.canHours)
        tabs.push({ to: `/people/${route.params.id}/activity/current`, label: 'Activity' })
      return tabs
    })

    return { tabs, onLoadPerson }
  },
})
</script>
