<!--
Admin displays the Admin module.
-->

<template lang="pug">
TabPage(:tabs='tabs')
  router-view
</template>

<script lang="ts">
import { defineComponent, computed } from 'vue'
import { useRoute } from 'vue-router'
import { TabPage, TabDef } from '../base'

export default defineComponent({
  components: { TabPage },
  setup() {
    const route = useRoute()

    // List of tabs on the page.
    const tabs = computed(() => {
      const tabs = [] as Array<TabDef>
      tabs.push({ to: '/admin/groups', label: 'Groups' })
      if (route.params.gid)
        tabs.push({
          to: `/admin/groups/${route.params.gid}`,
          label: route.params.gid === 'NEW' ? 'Add Group' : 'Edit Group',
        })
      tabs.push({ to: '/admin/roles', label: 'Roles' })
      if (route.params.rid)
        tabs.push({
          to: `/admin/roles/${route.params.rid}`,
          label: route.params.rid === 'NEW' ? 'Add Role' : 'Edit Role',
        })
      return tabs
    })

    return { tabs }
  },
})
</script>
