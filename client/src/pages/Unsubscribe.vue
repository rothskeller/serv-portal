<!--
Unsubscribe displays the Unsubscribe page.
-->

<template lang="pug">
#unsubscribe(v-if='loading')
  SSpinner
#unsubscribe(v-else-if='submitted')
  #unsub-head Unsubscribe
  #unsub-intro(v-if='noEmail').
    You have been removed from all of our email lists.
  #unsub-intro(v-else).
    You have been removed from the selected email lists.
  #unsub-warn.
    We would appreciate it if you’d drop us a note at
    <a href="mailto:admin@sunnyvaleserv.org">admin@sunnyvaleserv.org</a>
    and let us know why you unsubscribed.
  #unsub-warn.
    If you ever want to get back on the email lists, come back to this page
    and let us know.
#unsubscribe(v-else)
  #unsub-head Unsubscribe
  #unsub-intro Which email list(s) do you want to unsubscribe from?
  div(v-for='g in groups')
    SCheck(:id='`unsub-${g.id}`', :label='`${g.email}@SunnyvaleSERV.org`', v-model='g.unsub')
  div(style='margin-top: 0.5rem')
    SCheck#unsub-all(label='All SunnyvaleSERV email lists', v-model='noEmail')
  #unsub-warn.
    Please note that, if you unsubscribe from a critical mailing list for one
    of our volunteer groups, you may no longer be able to participate in that
    group.
  #unsub-buttons
    SButton(variant='primary', @click='onSubmit') Unsubscribe
</template>

<script lang="ts">
import { defineComponent, ref } from 'vue'
import { useRoute } from 'vue-router'
import axios from '../plugins/axios'
import setPage from '../plugins/page'
import { SButton, SCheck, SSpinner } from '../base'

interface GetUnsubscribeGroup {
  id: number
  email: string
  unsub: boolean
}
interface GetUnsubscribe {
  noEmail: boolean
  groups: Array<GetUnsubscribeGroup>
}

export default defineComponent({
  components: { SButton, SCheck, SSpinner },
  setup() {
    const route = useRoute()
    setPage({ title: 'SunnyvaleSERV Unsubscribe' })

    const loading = ref(true)
    const noEmail = ref(false)
    const groups = ref([] as Array<GetUnsubscribeGroup>)
    axios.get<GetUnsubscribe>(`/api/unsubscribe/${route.params.token}`).then((resp) => {
      noEmail.value = resp.data.noEmail
      groups.value = resp.data.groups
      loading.value = false
      if (route.params.email) {
        groups.value.forEach((g) => {
          if (g.email === route.params.email) g.unsub = true
        })
      }
    })

    const submitted = ref(false)
    async function onSubmit() {
      const body = new FormData()
      body.append('noEmail', noEmail.value.toString())
      groups.value.forEach((g) => {
        body.append(`unsub:${g.id}`, g.unsub.toString())
      })
      await axios.post(`/api/unsubscribe/${route.params.token}`, body)
      submitted.value = true
    }

    return { groups, loading, noEmail, onSubmit, submitted }
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
