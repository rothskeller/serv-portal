<!--
NotFound displays a page-not-found error.
-->

<template lang="pug">
#notfound-top
  #notfound-banner Page Not Found
  #notfound-intro
    | The page you attempted to reach does not exist.
    #notfound-button
      SButton(to='/', variant='primary', v-text='buttonLabel')
</template>

<script lang="ts">
import { defineComponent, Ref, computed, inject } from 'vue'
import { LoginData } from '../plugins/login'
import setPage from '../plugins/page'
import { SButton } from '../base'

export default defineComponent({
  components: { SButton },
  setup() {
    setPage({ title: '' })

    const me = inject<Ref<LoginData>>('me')!
    const buttonLabel = computed(() => (me.value ? 'Go to Home Page' : 'Go to Login Page'))
    return { buttonLabel }
  },
})
</script>

<style lang="postcss">
#notfound-top {
  margin: 0 auto;
  padding: 0 0.75rem;
  max-width: calc(0.75rem + 10rem + 20rem + 0.75rem + 1.5rem);
}
#notfound-banner {
  margin-top: 1rem;
  text-align: center;
  font-weight: bold;
  font-size: 1.5rem;
}
#notfound-intro {
  margin-top: 0.5rem;
  text-align: center;
  font-size: 0.9rem;
  line-height: 1.2;
}
#notfound-button {
  margin-top: 1rem;
  text-align: center;
}
</style>
