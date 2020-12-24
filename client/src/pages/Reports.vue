<!--
Reports displays the available reports.
-->

<template lang="pug">
TabPage(:tabs='tabs')
  router-view
</template>

<script lang="ts">
import { defineComponent, inject, Ref, computed } from 'vue'
import { LoginData } from '../plugins/login'
import { TabPage, TabDef } from '../base'

export default defineComponent({
  components: { TabPage },
  setup() {
    const me = inject<Ref<LoginData>>('me')!

    // List of tabs on the page.
    const tabs = computed(() => {
      const tabs = [] as Array<TabDef>
      if (me.value.canViewReports) tabs.push({ to: '/reports/attendance', label: 'Attendance' })
      if (me.value.canViewReports) tabs.push({ to: '/reports/clearance', label: 'Clearance' })
      return tabs
    })

    return { tabs }
  },
})
</script>
