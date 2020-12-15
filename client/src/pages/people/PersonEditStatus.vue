<!--
PersonEditStatus is the dialog box for editing a person's volunteer status.
-->

<template lang="pug">
Modal(ref='modal')
  SForm(
    dialog,
    variant='primary',
    title='Edit Volunteer Status',
    submitLabel='Save',
    :disabled='submitting',
    @submit='onSubmit',
    @cancel='onCancel'
  )
    SSpinner(v-if='loading')
    template(v-else)
      SFInput#person-volgistics(
        label='Volgistics ID',
        v-model='volgistics',
        autofocus,
        :errorFn='volgisticsError',
        :restrictFn='digitsOnly'
      )
      SFInput(
        :id='`person-dsw-cert`',
        type='date',
        :label='`DSW CERT`',
        v-model='person.dswCERT.registered',
        :errorFn='dswCERTError',
        :help='`Date when CERT DSW registration form was signed.`'
      )
      SFInput(
        :id='`person-dsw-comm`',
        type='date',
        :label='`DSW Communications`',
        v-model='person.dswComm.registered',
        :errorFn='dswCommError',
        :help='`Date when Communications DSW registration form was signed.`'
      )
      SFInput#person-background(
        label='BG Check',
        trim,
        v-model='person.backgroundCheck.cleared',
        :errorFn='backgroundError',
        help='Date when background check cleared, or “TRUE” if clearance confirmed but date unknown',
        style='text-transform: uppercase'
      )
      SFCheckGroup#person-identification(
        label='IDs Issued',
        v-model='identification',
        :options='identTypes'
      )
</template>

<script lang="ts">
import { defineComponent, nextTick, ref, watch } from 'vue'
import axios from '../../plugins/axios'
import { Modal, SForm, SFCheckGroup, SFInput, SSpinner } from '../../base'

interface GetPersonStatus {
  id: number
  volgistics: {
    needed: boolean
    id: number
  }
  dswCERT: {
    needed: boolean
    registered?: string
    expires?: string
    expired?: true
  }
  dswComm: {
    needed: boolean
    registered?: string
    expires?: string
    expired?: true
  }
  backgroundCheck: {
    needed: boolean
    cleared?: string
  }
  identification: Array<{
    type: string
    held: boolean
  }>
}

function digitsOnly(s: string): string {
  return s.replace(/[^0-9]/g, '')
}

export default defineComponent({
  components: { Modal, SForm, SFCheckGroup, SFInput, SSpinner },
  props: {
    pid: { type: Number, required: true },
  },
  setup(props) {
    const modal = ref(null as any)
    function show() {
      loadData()
      return modal.value.show()
    }

    // Load the form data.
    const person = ref({} as GetPersonStatus)
    const volgistics = ref('')
    const identification = ref(new Set() as Set<string>)
    const identTypes = ref([] as any)
    const loading = ref(true)
    async function loadData() {
      loading.value = true
      person.value = (await axios.get<GetPersonStatus>(`/api/people/${props.pid}/status`)).data
      volgistics.value = person.value.volgistics ? person.value.volgistics.toString() : ''
      identTypes.value = person.value.identification.map((id) => ({
        label: id.type,
        value: id.type,
      }))
      identification.value = new Set(
        person.value.identification.filter((id) => id.held).map((id) => id.type)
      )
      loading.value = false
    }

    // Field validation.
    function volgisticsError(lostFocus: boolean) {
      if (!lostFocus || !volgistics.value) return ''
      if (parseInt(volgistics.value) < 1) return 'This is not a valid Volgistics ID number.'
      return ''
    }
    function dswCERTError(lostFocus: boolean) {
      if (!lostFocus || !person.value.dswCERT.registered) return ''
      if (!person.value.dswCERT.registered.match(/^20\d\d-\d\d-\d\d$/))
        return 'This is not a valid YYYY-MM-DD date.'
      return ''
    }
    function dswCommError(lostFocus: boolean) {
      if (!lostFocus || !person.value.dswComm.registered) return ''
      if (!person.value.dswComm.registered.match(/^20\d\d-\d\d-\d\d$/))
        return 'This is not a valid YYYY-MM-DD date.'
      return ''
    }
    function backgroundError(lostFocus: boolean) {
      if (
        !lostFocus ||
        !person.value?.backgroundCheck.cleared ||
        person.value?.backgroundCheck.cleared.toUpperCase() === 'TRUE'
      )
        return ''
      if (!person.value.backgroundCheck.cleared.match(/^20\d\d-\d\d-\d\d$/))
        return 'This is not a valid YYYY-MM-DD date.'
      return ''
    }

    // Save and close.
    const submitting = ref(false)
    async function onSubmit() {
      var body = new FormData()
      body.append('volgistics', volgistics.value)
      body.append('backgroundCheck', (person.value.backgroundCheck.cleared || '').toLowerCase())
      body.append('dswCERT', person.value.dswCERT.registered || '')
      body.append('dswComm', person.value.dswComm.registered || '')
      identification.value.forEach((t) => {
        body.append('identification', t)
      })
      submitting.value = true
      await axios.post(`/api/people/${props.pid}/status`, body)
      submitting.value = false
      modal.value.close(true)
    }
    function onCancel() {
      modal.value.close(false)
    }

    return {
      backgroundError,
      digitsOnly,
      dswCERTError,
      dswCommError,
      identification,
      identTypes,
      loading,
      modal,
      onCancel,
      onSubmit,
      person,
      show,
      submitting,
      volgistics,
      volgisticsError,
    }
  },
})
</script>

<style lang="postcss">
</style>
