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
        v-model='person.dswCERT',
        :errorFn='dswCERTError',
        :help='`Date when CERT DSW registration form was signed.`'
      )
      SFInput(
        :id='`person-dsw-comm`',
        type='date',
        :label='`DSW Communications`',
        v-model='person.dswComm',
        :errorFn='dswCommError',
        :help='`Date when Communications DSW registration form was signed.`'
      )
      SFCheckGroup#person-identification(
        label='IDs Issued',
        v-model='identification',
        :options='identTypes'
      )
      PersonEditBGChecks(
        ref='bgChecksRef',
        :checks='person.bgChecks',
        :types='person.bgCheckTypes'
      )
    template(#extraButtons)
      SButton(variant='primary', @click.prevent='onAddBGCheck') Add
</template>

<script lang="ts">
import { defineComponent, nextTick, ref, watch } from 'vue'
import axios from '../../plugins/axios'
import { Modal, SButton, SForm, SFCheckGroup, SFInput, SSpinner } from '../../base'
import type { GetPersonStatusBGCheck } from './api'
import PersonEditBGChecks from './PersonEditBGChecks'

interface GetPersonStatus {
  id: number
  volgistics: number
  dswCERT: string
  dswComm: string
  bgChecks: Array<GetPersonStatusBGCheck>
  bgCheckTypes: Array<string>
  identification: Array<{
    type: string
    held: boolean
  }>
}

function digitsOnly(s: string): string {
  return s.replace(/[^0-9]/g, '')
}

export default defineComponent({
  components: { Modal, PersonEditBGChecks, SButton, SForm, SFCheckGroup, SFInput, SSpinner },
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
      if (!lostFocus || !person.value.dswCERT) return ''
      if (!person.value.dswCERT.match(/^20\d\d-\d\d-\d\d$/))
        return 'This is not a valid YYYY-MM-DD date.'
      return ''
    }
    function dswCommError(lostFocus: boolean) {
      if (!lostFocus || !person.value.dswComm) return ''
      if (!person.value.dswComm.match(/^20\d\d-\d\d-\d\d$/))
        return 'This is not a valid YYYY-MM-DD date.'
      return ''
    }

    // Background Checks.
    const bgChecksRef = ref(null as any)
    function onAddBGCheck() {
      console.log('onAddBGCheck called')
      bgChecksRef.value?.startAdd()
    }

    // Save and close.
    const submitting = ref(false)
    async function onSubmit() {
      bgChecksRef.value?.prepareForSave()
      var body = new FormData()
      body.append('volgistics', volgistics.value)
      body.append('dswCERT', person.value.dswCERT || '')
      body.append('dswComm', person.value.dswComm || '')
      identification.value.forEach((t) => {
        body.append('identification', t)
      })
      person.value.bgChecks.forEach(bc => {
        body.append('bgCheck', `${bc.date}:${bc.types.join(',')}:${bc.assumed}`)
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
      bgChecksRef,
      digitsOnly,
      dswCERTError,
      dswCommError,
      identification,
      identTypes,
      loading,
      modal,
      onAddBGCheck,
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
#person-edbg-header {
  margin: 0 0.75rem;
  border-top: 1px solid #ccc;
  padding-top: 0.25rem;
}
#person-edbg-help {
  margin: 0 0.75rem;
  color: #6c757d;
  font-size: 80%;
}
</style>
