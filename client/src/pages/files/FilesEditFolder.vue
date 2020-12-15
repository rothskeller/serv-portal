<!--
FilesEditFolder is the dialog box for editing a folder (or adding a new one).
-->

<template lang="pug">
Modal(ref='modal')
  SForm(
    dialog,
    variant='primary',
    :title='title',
    :submitLabel='submitLabel',
    :disabled='submitting',
    @submit='onSubmit',
    @cancel='onCancel'
  )
    SSpinner(v-if='!visibilities.length')
    template(v-else)
      SFInput#files-editf-name(
        label='Name',
        trim,
        autofocus,
        v-model='folder.name',
        :errorFn='nameError'
      )
      SFSelect#files-editf-vis(
        label='Visibility',
        :options='visibilities',
        v-model='folder.visibility'
      )
</template>

<script lang="ts">
import { computed, defineComponent, nextTick, PropType, ref, watch } from 'vue'
import axios from '../../plugins/axios'
import { Modal, SForm, SFInput, SFSelect, SSpinner } from '../../base'

interface GetFolderEdit {
  name: string
  visibility: string
  allowedVisibilities: Array<string>
}

const fmtVis: Record<string, string> = {
  public: 'Public',
  serv: 'SERV Volunteers',
  admin: 'SERV Leads',
  'cert-d': 'CERT Deployment Teams',
  'cert-t': 'CERT Training Committee',
  listos: 'Listos Team',
  sares: 'SARES Members',
  snap: 'SNAP Team',
}

export default defineComponent({
  components: { Modal, SForm, SFInput, SFSelect, SSpinner },
  setup() {
    const modal = ref(null as any)
    let url: string
    let op: string
    function showAdd(parent: string) {
      url = parent
      op = 'newFolder'
      loadData()
      return modal.value.show()
    }
    function showEdit(folder: string) {
      url = folder
      op = 'editFolder'
      loadData()
      return modal.value.show()
    }

    // Load the form data.
    const folder = ref({} as GetFolderEdit)
    const visibilities = ref([] as any)
    async function loadData() {
      folder.value = (await axios.get<GetFolderEdit>(`/api/folders${url}?op=${op}`)).data
      visibilities.value = folder.value.allowedVisibilities.map((v) => ({
        value: v,
        label: fmtVis[v],
      }))
    }

    // The name field.
    const duplicateName = ref('')
    function nameError(lostFocus: boolean): string {
      if (!lostFocus) return ''
      if (!folder.value.name) return 'The folder name is required.'
      if (folder.value.name === duplicateName.value) return 'This name is already in use.'
      return ''
    }

    // Labels.
    const title = computed(() => (op === 'newFolder' ? 'Add Folder' : 'Edit Folder'))
    const submitLabel = computed(() => (op === 'newFolder' ? 'Add' : 'Save'))

    // Save and close.
    const submitting = ref(false)
    async function onSubmit() {
      submitting.value = true
      const body = new FormData()
      body.append('name', folder.value.name)
      body.append('visibility', folder.value.visibility)
      try {
        const newURL = (await axios.post<string>(`/api/folders${url}?op=${op}`, body)).data
        modal.value.close(newURL || true)
      } catch (err) {
        if (err.response && err.response.status == 409) duplicateName.value = folder.value.name
        else throw err
      } finally {
        submitting.value = false
      }
    }
    function onCancel() {
      modal.value.close(false)
    }

    return {
      folder,
      modal,
      nameError,
      onCancel,
      onSubmit,
      showAdd,
      showEdit,
      submitLabel,
      submitting,
      title,
      visibilities,
    }
  },
})
</script>

<style lang="postcss">
</style>
