<!--
TabPage displays a page with a tab bar across the top and a content area
underneath.
-->

<template lang="pug">
#tabpage
  nav#tabpage-bar
    ul#tabpage-tabs
      li.tabpage-tab(v-for='tab in tabs')
        router-link.tabpage-link(
          :to='tab.to',
          exact-active-class='tabpage-active',
          v-text='tab.label'
        )
  #tabpage-content
    slot
</template>

<script lang="ts">
import { defineComponent, PropType, Ref } from 'vue'
import type { TabDef } from './tabdef'

export default defineComponent({
  props: {
    tabs: { type: Array as PropType<Array<TabDef>>, required: true },
  },
})
</script>

<style lang="postcss">
#tabpage {
  height: 100%;
  border: none;
  position: relative;
  display: flex;
  flex-direction: column;
  min-width: 0;
}
#tabpage-bar {
  padding: 0.75rem 1.25rem;
  margin-bottom: 0;
  background-color: rgba(0, 0, 0, 0.03);
  border-bottom: 1px solid rgba(0, 0, 0, 0.125);
}
#tabpage-tabs {
  margin: 0 -0.625rem -0.75rem;
  border-bottom: 0;
  display: flex;
  flex-wrap: wrap;
  padding-left: 0;
  list-style: none;
}
.tabpage-tab {
  margin-bottom: -1px;
}
.tabpage-link {
  border: 1px solid transparent;
  border-top-left-radius: 0.25rem;
  border-top-right-radius: 0.25rem;
  display: block;
  padding: 0.5rem 1rem;
}
.tabpage-active {
  color: #495057;
  background-color: #fff;
  border-color: #dee2e6 #dee2e6 #fff;
}
#tabpage-content {
  flex: auto;
  overflow-x: hidden;
  overflow-y: auto;
}
</style>
