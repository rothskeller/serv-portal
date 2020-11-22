<!--
ListsList displays the list of lists.
-->

<template lang="pug">
#lists-list
  SSpinner(v-if='loading')
  template(v-else)
    #lists-list-table(v-if='lists.length')
      .lists-list-name.lists-list-heading List
      .lists-list-num.lists-list-heading Sub
      .lists-list-num.lists-list-heading Unsub
      .lists-list-num.lists-list-heading Send
      template(v-for='l in lists')
        .lists-list-name
          span(v-if='l.type === "sms"', v-text='"SMS: "')
          router-link(
            :to='`/admin/lists/${l.id}`',
            v-text='l.type === "email" ? `${l.name}@SunnyvaleSERV.org` : l.name'
          )
        .lists-list-num(v-text='l.subscribed')
        .lists-list-num(v-text='l.unsubscribed')
        .lists-list-num(v-text='l.senders')
    div(v-else) No lists currently defined.
    #lists-list-buttons
      SButton(variant='primary', @click='onAdd') Add List
</template>

<script lang="ts">
import { defineComponent, ref } from 'vue'
import axios from '../../plugins/axios'
import setPage from '../../plugins/page'
import { SButton, SSpinner } from '../../base'
import { useRouter } from 'vue-router'

type GetListsList = {
  id: number
  type: string
  name: string
  subscribed: number
  unsubscribed: number
  senders: number
}

export default defineComponent({
  components: { SButton, SSpinner },
  setup() {
    const router = useRouter()
    setPage({ title: 'Lists' })
    const loading = ref(true)
    const lists = ref([] as Array<GetListsList>)
    axios.get<Array<GetListsList>>(`/api/lists`).then((resp) => {
      lists.value = resp.data
      loading.value = false
    })
    function fmtType(t: string) {
      return t === 'email' ? 'Email' : 'SMS'
    }
    function onAdd() {
      router.push('/admin/lists/NEW')
    }
    return { fmtType, lists, loading, onAdd }
  },
})
</script>

<style lang="postcss">
#lists-list {
  padding: 1.5rem 0.75rem;
}
#lists-list-table {
  display: grid;
  grid: auto / 1fr;
  column-gap: 1rem;
  .touch & {
    row-gap: calc(40px - 1.5rem);
  }
  @media (min-width: 576px) {
    grid: auto / max-content min-content min-content min-content;
  }
}
.lists-list-heading {
  display: none;
  font-weight: bold;
  @media (min-width: 576px) {
    display: block;
  }
}
.lists-list-name {
  margin-right: 1rem;
}
.lists-list-num {
  display: none;
  text-align: right;
  @media (min-width: 576px) {
    display: block;
  }
}
#lists-list-buttons {
  margin-top: 1.5rem;
}
</style>
