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
import type { GetPersonPersonBase } from './people/PersonView.vue'

export default defineComponent({
  components: { TabPage },
  setup() {
    const me = inject<Ref<LoginData>>('me')!
    const route = useRoute()

    const person = ref(null as null | GetPersonPersonBase)
    function onLoadPerson(ev: GetPersonPersonBase) {
      person.value = ev
    }
    watchEffect(() => {
      if (!route.params.id) person.value = null
    })

    // List of tabs on the page.
    const tabs = computed(() => {
      const tabs = [] as Array<TabDef>
      tabs.push({ to: '/people/list', label: 'List' })
      tabs.push({ to: '/people/map', label: 'Map' })
      if (!route.params.id && me.value.canAddPeople)
        tabs.push({ to: '/people/NEW/edit', label: 'Add Person' })
      if (route.params.id && route.params.id !== 'NEW')
        tabs.push({ to: `/people/${route.params.id}`, label: 'Details' })
      if (person.value && person.value.canEdit)
        tabs.push({
          to: `/people/${route.params.id}/edit`,
          label: route.params.id === 'NEW' ? 'Add Person' : 'Edit',
        })
      if (person.value && person.value.canHours)
        tabs.push({ to: `/people/${route.params.id}/hours`, label: 'Hours' })
      return tabs
    })

    return { tabs, onLoadPerson }
  },
})
</script>
