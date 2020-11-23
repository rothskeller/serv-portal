<!--
PersonEditLists is the dialog box for editing a person's list subscriptions.
-->

<template lang="pug">
SForm(
  :dialog='inModal',
  :title='inModal ? "Edit List Subscriptions" : null',
  submitLabel='Save',
  :cancelLabel='inModal ? "Cancel" : ""',
  :disabled='submitting',
  @submit='onSubmit',
  @cancel='onCancel'
)
  SSpinner(v-if='!lists.length')
  template(v-else)
    SFCheckGroup#person-lists(
      label='Subscriptions',
      :options='lists',
      valueKey='id',
      labelKey='name',
      v-model='subscribed'
    )
    #subscriptions-warning.form-item(v-if='unsubscribeWarnings', v-text='unsubscribeWarnings')
</template>

<script lang="ts">
import { computed, defineComponent, onMounted, ref, watch } from 'vue'
import axios from '../../plugins/axios'
import { Modal, SForm, SFCheckGroup, SSpinner } from '../../base'

interface GetPersonListsList {
  id: number
  name: string
  subscribed: boolean
  subWarn: Array<string>
}

export default defineComponent({
  components: { SFCheckGroup, SForm, SSpinner },
  props: {
    inModal: { type: Boolean, default: false },
    pid: { type: [Number, String], required: true },
    email: { type: String },
  },
  emits: ['done'],
  setup(props, { emit }) {
    // Load the form data.
    const lists = ref([] as Array<GetPersonListsList>)
    const subscribed = ref(new Set<number>())
    onMounted(async () => {
      lists.value = []
      subscribed.value.clear()
      lists.value = (
        await axios.get<Array<GetPersonListsList>>(`/api/people/${props.pid}/lists`)
      ).data
      subscribed.value.clear()
      lists.value
        .filter(
          (l) => l.subscribed && (!props.email || l.name !== `${props.email}@SunnyvaleSERV.org`)
        )
        .forEach((l) => {
          subscribed.value.add(l.id)
        })
    })

    // Unsubscribe warnings.
    const unsubscribeWarnings = computed(() =>
      lists.value
        .filter((l) => !subscribed.value.has(l.id) && l.subWarn.length)
        .map((l) => {
          switch (l.subWarn.length) {
            case 1:
              return `Messages sent to ${l.name} are considered required for the ${l.subWarn[0]} role.  Unsubscribing from it may cause you to lose that role.`
              break
            case 2:
              return `Messages sent to ${l.name} are considered required for the ${l.subWarn[0]} and ${l.subWarn[1]} roles.  Unsubscribing from it may cause you to lose those roles.`
              break
            default:
              return `Messages sent to ${l.name} are considered required for the ${l.subWarn
                .slice(0, -1)
                .join(', ')}, and ${
                l.subWarn[l.subWarn.length - 1]
              } roles.  Unsubscribing from it may cause you to lose those roles.`
          }
        })
        .join('\n\n')
    )

    // Save and close.
    const submitting = ref(false)
    async function onSubmit() {
      var body = new FormData()
      subscribed.value.forEach((lid) => {
        body.append('list', lid.toString())
      })
      submitting.value = true
      await axios.post(`/api/people/${props.pid}/lists`, body)
      submitting.value = false
      emit('done', true)
    }
    function onCancel() {
      emit('done', false)
    }

    return { lists, onCancel, onSubmit, submitting, subscribed, unsubscribeWarnings }
  },
})
</script>

<style lang="postcss">
#subscriptions-warning {
  margin: 0 0.75rem;
  color: red;
  line-height: 1.2;
  white-space: pre-wrap;
}
</style>
