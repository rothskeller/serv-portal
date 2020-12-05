<!--
OrgBadge displays the badge or logo for an org.  It is always displayed with
2rem height, and variable width.
-->

<template lang="pug">
img.org-badge(:src='img')
</template>

<script lang="ts">
import { computed, defineComponent, watch } from 'vue'
import certBadge from './orgBadges/cert.png'
import dpsBadge from './orgBadges/dps.png'
import listosBadge from './orgBadges/listos.png'
import saresBadge from './orgBadges/sares.png'
import servBadge from './orgBadges/serv.png'
import snapBadge from './orgBadges/snap.png'

export default defineComponent({
  props: {
    org: { type: String, required: true },
  },
  setup(props) {
    const img = computed(() => {
      switch (props.org) {
        case 'dps':
          return dpsBadge
        case 'cert':
          return certBadge
        case 'listos':
          return listosBadge
        case 'sares':
          return saresBadge
        case 'serv':
          return servBadge
        case 'snap':
          return snapBadge
        default:
          return null
      }
    })
    return { img }
  },
})
</script>

<style lang="postcss">
.org-badge {
  margin-right: 0.5rem;
  width: 32px;
}
</style>
