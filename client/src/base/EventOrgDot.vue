<!--
EventOrgDot displays an appropriately colored dot for an event organization.
-->

<template lang="pug">
span.dot(:class='`dot-${org}`', :title='orgNames[org]')
</template>

<script lang="ts">
import { defineComponent } from 'vue'

const orgNames = {
  admin: 'Admin',
  'cert-d': 'CERT Deployment',
  'cert-t': 'CERT Training',
  listos: 'Listos',
  sares: 'SARES',
  snap: 'SNAP',
}

export default defineComponent({
  props: {
    org: String,
  },
  setup() {
    return { orgNames }
  },
})
</script>

<style lang="postcss">
.dot {
  display: inline-block;
  margin: 0 0 0 2px;
  width: 10px;
  height: 10px;
  border-radius: 50%;
  &:first-child {
    margin-left: 0;
  }
}
.dot-admin {
  background-color: #a9a9a9; /* gray */
}
.dot-cert-t {
  background-color: #3cb44b; /* green */
}
.dot-cert-d {
  border: 2px solid #3cb44b; /* green */
  background-color: white;
}
.dot-listos {
  background-color: #f58231; /* orange */
}
.dot-sares {
  background-color: #ffe119; /* yellow */
}
.dot-snap {
  background-color: #4363d8; /* blue */
}
/* colors taken from https://sashat.me/2017/01/11/list-of-20-simple-distinct-colors/ */
</style>
