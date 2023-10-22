<!--
PersonEditNames is the dialog box for editing a person's names.
-->

<template lang="pug">
Modal(ref='modal')
  SForm(
    dialog,
    variant='primary',
    title='Edit Names',
    submitLabel='Save',
    :disabled='submitting',
    @submit='onSubmit',
    @cancel='onCancel'
  )
    SSpinner(v-if='loading')
    template(v-else)
      SFInput#person-informalName(
        label='Name',
        trim,
        autofocus,
        v-model='person.informalName',
        :errorFn='informalNameError',
        help='What you like to be called, e.g. “Joe Banks”'
      )
      SFInput#person-formalName(
        label='Formal name',
        trim,
        v-model='person.formalName',
        :errorFn='formalNameError',
        help='For formal documents, e.g. “Joseph A. Banks, Jr.”'
      )
      SFInput#person-sortName(
        label='Sort name',
        trim,
        v-model='person.sortName',
        :errorFn='sortNameError',
        help='For appearance in sorted lists, e.g. “Banks, Joe”'
      )
      SFInput#person-callSign(
        label='Call sign',
        trim,
        v-model='person.callSign',
        :errorFn='callSignError',
        help='FCC amateur radio license',
        style='text-transform: uppercase'
      )
      template(v-if='person.birthdate || person.birthdate === ""')
        SFInput#person-birthdate(
          label='Birthdate',
          type='date',
          v-model='person.birthdate',
          :errorFn='birthdateError'
        )
</template>

<script lang="ts">
import { defineComponent, nextTick, ref, watch } from 'vue'
import axios from '../../plugins/axios'
import { Modal, SForm, SFInput, SSpinner } from '../../base'

interface GetPersonNames {
  id: number
  informalName: string
  formalName: string
  sortName: string
  callSign: string
  birthdate?: string
}

export default defineComponent({
  components: { Modal, SForm, SFInput, SSpinner },
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
    const person = ref({} as GetPersonNames)
    const loading = ref(true)
    async function loadData() {
      loading.value = true
      person.value = (await axios.get<GetPersonNames>(`/api/people/${props.pid}/names`)).data
      loading.value = false
    }

    // Field validation.
    const duplicateCallSign = ref('')
    const duplicateSortName = ref('')
    function informalNameError(lostFocus: boolean) {
      if (!lostFocus) return ''
      if (!person.value.informalName) return 'A name is required.'
      return ''
    }
    function formalNameError(lostFocus: boolean) {
      if (!lostFocus) return ''
      if (!person.value.formalName) return 'A name is required.'
      return ''
    }
    function sortNameError(lostFocus: boolean) {
      if (!person.value.sortName) return lostFocus ? 'A name is required.' : ''
      if (duplicateSortName.value === person.value?.sortName)
        return 'A different person has this name.'
      return ''
    }
    function callSignError(lostFocus: boolean) {
      if (!person.value.callSign) return ''
      if (lostFocus && !person.value.callSign.match(/^[AKNW][A-Z]?[0-9][A-Z]{1,3}$/i))
        return 'This is not a valid call sign.'
      if (duplicateCallSign.value === person.value?.callSign)
        return 'A different person has this call sign.'
      return ''
    }
    function birthdateError(lostFocus: boolean) {
      if (!person.value.birthdate) return ''
      if (lostFocus && !person.value.birthdate.match(/^(?:19|20)\d\d-\d\d-\d\d$/)) {
        return 'This is not a valid YYYY-MM-DD date.'
      }
      return ''
    }

    // When the informal name is changed, we may update the formal and/or
    // sort name.
    function informalToSort(n: string): string {
      if (!n) return n
      const parts = n.split(/\s+/, 2)
      return parts.length > 1 ? `${parts[1]}, ${parts[0]}` : n
    }
    watch(
      () => person.value.informalName,
      (n, o) => {
        if (person.value.formalName === o) person.value.formalName = n
        if (person.value.sortName === informalToSort(o)) person.value.sortName = informalToSort(n)
      }
    )

    // Save and close.
    const submitting = ref(false)
    async function onSubmit() {
      var body = new FormData()
      body.append('informalName', person.value.informalName)
      body.append('formalName', person.value.formalName)
      body.append('sortName', person.value.sortName)
      body.append('callSign', person.value.callSign)
      body.append('birthdate', person.value.birthdate!)
      submitting.value = true
      try {
        await axios.post(`/api/people/${props.pid}/names`, body)
        modal.value.close(true)
      } catch (err) {
        if (!err.response || err.response.status !== 409) throw err
        switch (err.response.data) {
          case 'sortName':
            duplicateSortName.value = person.value.sortName
            break
          case 'callSign':
            duplicateCallSign.value = person.value.callSign
            break
        }
      } finally {
        submitting.value = false
      }
    }
    function onCancel() {
      modal.value.close(false)
    }

    return {
      birthdateError,
      callSignError,
      formalNameError,
      informalNameError,
      loading,
      modal,
      onCancel,
      onSubmit,
      person,
      show,
      sortNameError,
      submitting,
    }
  },
})
</script>

<style lang="postcss">
</style>
