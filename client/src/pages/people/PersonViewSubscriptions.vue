<!--
PersonViewSubscriptions displays a person's list subscriptions, as well as their
global unsubscribes if any.
-->

<template lang="pug">
PersonViewSection(
  v-if='person.noEmail || person.noText || person.lists.length || person.canEditLists',
  title='Subscriptions',
  :editable='person.canEditLists',
  @edit='onEditSubscriptions'
)
  #person-view-lists-list
    div(v-if='person.noEmail', style='color: red') Unsubcribed from all email.
    div(v-if='person.noText', style='color: red') Unsubscribed from all text messaging.
    div(v-for='list in person.lists', v-text='list')
  PersonEditSubscriptions(
    v-if='person.canEditLists',
    ref='editSubscriptionsModal',
    :pid='person.id'
  )
</template>

<script lang="ts">
import { computed, defineComponent, PropType, ref } from 'vue'
import { GetPerson } from './PersonView.vue'
import PersonEditSubscriptions from './PersonEditSubscriptions.vue'
import PersonViewSection from './PersonViewSection.vue'

export default defineComponent({
  components: { PersonEditSubscriptions, PersonViewSection },
  props: {
    person: { type: Object as PropType<GetPerson>, required: true },
  },
  emits: ['reload'],
  setup(props, { emit }) {
    const editSubscriptionsModal = ref(null as any)
    async function onEditSubscriptions() {
      if (!(await editSubscriptionsModal.value.show())) return
      emit('reload')
    }
    return {
      editSubscriptionsModal,
      onEditSubscriptions,
    }
  },
})
</script>

<style lang="postcss">
#person-view-lists-list {
  margin-top: 0.75rem;
  line-height: 1.2;
  margin-left: 2em;
  text-indent: -2em;
}
</style>
