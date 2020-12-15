<!--
FilesEditLink is the dialog box for creating and editing link documents.
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
    SSpinner(v-if='loading')
    template(v-else)
      SFInput#files-editl-name(
        label='Name',
        trim,
        autofocus,
        v-model='link.name',
        :errorFn='nameError'
      )
      SFInput#files-editl-url(label='URL', trim, v-model='link.url', :errorFn='urlError')
</template>

<script lang="ts">
import { computed, defineComponent, nextTick, PropType, ref, watch } from 'vue'
import axios from '../../plugins/axios'
import { Modal, SForm, SFInput, SSpinner } from '../../base'

interface GetDocumentEdit {
  name: string
  url: string
}

export default defineComponent({
  components: { Modal, SForm, SFInput, SSpinner },
  setup() {
    const modal = ref(null as any)
    const folderURL = ref('')
    const linkName = ref('')
    function showAdd(folder: string) {
      folderURL.value = folder
      link.value = { name: '', url: '' }
      loading.value = false
      return modal.value.show()
    }
    function showEdit(folder: string, ln: string) {
      folderURL.value = folder
      linkName.value = ln
      loadData()
      return modal.value.show()
    }

    // Load the form data.
    const link = ref({} as GetDocumentEdit)
    const loading = ref(true)
    async function loadData() {
      loading.value = true
      link.value = (
        await axios.get<GetDocumentEdit>(
          `/api/document${folderURL.value}/${encodeURIComponent(linkName.value)}?op=edit`
        )
      ).data
      loading.value = false
    }

    // The name field.
    const duplicateName = ref('')
    function nameError(lostFocus: boolean): string {
      if (!lostFocus) return ''
      if (!link.value.name) return 'The link name is required.'
      if (link.value.name === duplicateName.value) return 'This name is already in use.'
      return ''
    }

    // The URL field.
    function urlError(lostFocus: boolean): string {
      if (!lostFocus) return ''
      if (!link.value.url) return 'The link URL is required.'
      if (!link.value.url.startsWith('http://') && !link.value.url.startsWith(`https://`))
        return 'The link URL must start with http:// or https://.'
      return ''
    }

    // Labels.
    const title = computed(() => (linkName.value ? 'Edit Link' : 'Add Link'))
    const submitLabel = computed(() => (linkName.value ? 'Save' : 'Add'))

    // Save and close.
    const submitting = ref(false)
    async function onSubmit() {
      var body = new FormData()
      body.append('name', link.value.name)
      body.append('url', link.value.url)
      submitting.value = true
      try {
        if (linkName.value)
          await axios.post(
            `/api/document${folderURL.value}/${encodeURIComponent(linkName.value)}?op=editLink`,
            body
          )
        else await axios.post(`/api/folders${folderURL.value}?op=newLink`, body)
        modal.value.close(true)
      } catch (err) {
        if (!err.response || err.response.status !== 409) throw err
        duplicateName.value = link.value.name
      } finally {
        submitting.value = false
      }
    }
    function onCancel() {
      modal.value.close(false)
    }

    return {
      link,
      loading,
      modal,
      nameError,
      onCancel,
      onSubmit,
      showAdd,
      showEdit,
      submitLabel,
      submitting,
      title,
      urlError,
    }
  },
})
</script>

<style lang="postcss">
</style>
