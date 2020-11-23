<!--
Unsubscribe displays the Unsubscribe page.
-->

<template lang="pug">
#unsubscribe(v-if='submitted')
  #unsub-head Subscriptions
  #unsub-intro.
    Your selections have been saved.
  #unsub-warn.
    We would appreciate it if you’d drop us a note at
    <a href="mailto:admin@sunnyvaleserv.org">admin@sunnyvaleserv.org</a>
    and let us know why you unsubscribed.
  #unsub-warn.
    If you ever want to get back on the email lists, come back to this page
    and let us know.
#unsubscribe(v-else)
  #unsub-head Your Subscriptions
  div Turn off any that you don't want.
  SubscriptionsForm(:pid='$route.params.token', :email='$route.params.email', @done='onDone')
</template>

<script lang="ts">
import { defineComponent, ref } from 'vue'
import { useRoute } from 'vue-router'
import setPage from '../plugins/page'
import SubscriptionsForm from './people/SubscriptionsForm.vue'

export default defineComponent({
  components: { SubscriptionsForm },
  setup() {
    const route = useRoute()
    setPage({ title: 'SunnyvaleSERV Subscriptions' })
    const submitted = ref(false)
    function onDone(sub: boolean) {
      submitted.value = sub
    }
    return { onDone, submitted }
  },
})
</script>

<style lang="postcss">
#unsubscribe {
  margin: 1.5rem 0.75rem;
}
#unsub-head {
  font-weight: bold;
  font-size: 1.25rem;
}
#unsub-intro {
  margin: 0.5rem 0;
}
#unsub-warn {
  margin: 0.5rem 0 0.75rem;
  max-width: 40rem;
  line-height: 1.2;
}
</style>
