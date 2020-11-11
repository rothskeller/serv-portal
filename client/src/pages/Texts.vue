<!--
Texts displays the Texts module.
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
      tabs.push({ to: '/texts', label: 'Sent Messages' })
      tabs.push({ to: '/texts/send', label: 'New Message' })
      if (route.params.id) tabs.push({ to: `/texts/${route.params.id}`, label: 'Delivery Status' })
      return tabs
    })

    return { tabs }
  },
})
</script>
