<!--
Logout handles a logout request.
-->

<template lang="pug">
div(style='margin: 1.5rem 0.75rem') Logging out...
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { useRouter } from 'vue-router'
import axios from '../plugins/axios'
import { clearLoginData } from '../plugins/login'

export default defineComponent({
  setup() {
    const router = useRouter()
    axios.post('/api/logout').then(() => {
      clearLoginData()
      router.replace('/')
    })
  },
})
</script>
