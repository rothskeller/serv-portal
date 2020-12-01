<!--
FilesEditFile is the dialog box for editing file documents.
-->

<template lang="pug">
Modal(ref='modal')
  SForm(
    dialog,
    variant='primary',
    title='Edit File',
    submitLabel='Save',
    :disabled='submitting',
    @submit='onSubmit',
    @cancel='onCancel'
  )
    SFInput#files-editf-name(
      ref='nameRef',
      label='Name',
      trim,
      v-model='newName',
      :errorFn='nameError'
    )
    SFFile#files-editf-file(label='File', v-model='file')
    SProgress#files-editf-progress.form-item(v-if='progress', :value='progress')
</template>

<script lang="ts">
import { defineComponent, nextTick, ref } from 'vue'
import axios from '../../plugins/axios'
import { Modal, SForm, SFFile, SFInput, SProgress } from '../../base'

export default defineComponent({
  components: { Modal, SForm, SFFile, SFInput, SProgress },
  setup() {
    const modal = ref(null as any)
    let folderURL = ''
    let oldName = ''
    const newName = ref('')
    function show(folder: string, fn: string) {
      folderURL = folder
      oldName = newName.value = fn
      console.log('newName', fn)
      nextTick(() => {
        nextTick(() => {
          nameRef.value.focus()
        })
      })
      return modal.value.show()
    }

    // The name field.
    const nameRef = ref(null as any)
    const duplicateName = ref('')
    function nameError(lostFocus: boolean): string {
      if (!lostFocus) return ''
      if (!newName.value) return 'The file name is required.'
      if (newName.value[0] === '.' || newName.value.includes(':') || newName.value.includes('/'))
        return `This is not a legal file name.`
      if (newName.value === duplicateName.value) return 'This name is already in use.'
      return ''
    }

    // The file field.
    const file = ref(null as null | FileList)

    // Save and close.
    const submitting = ref(false)
    async function onSubmit() {
      const body = new FormData()
      submitting.value = true
      try {
        body.append('name', newName.value)
        if (file.value && file.value.length) body.append('file', file.value.item(0)!)
        await axios.post(
          `/api/document${folderURL}/${encodeURIComponent(oldName)}?op=editFile`,
          body,
          { onUploadProgress }
        )
        modal.value.close(true)
      } catch (err) {
        if (!err.response || err.response.status !== 409) throw err
        duplicateName.value = newName.value
      } finally {
        submitting.value = false
        progress.value = 0
      }
    }
    function onCancel() {
      modal.value.close(false)
    }

    // Display upload progress.
    const progress = ref(0)
    function onUploadProgress(evt: ProgressEvent) {
      if (evt.lengthComputable && evt.total) progress.value = evt.loaded / evt.total
    }

    return {
      file,
      modal,
      nameError,
      nameRef,
      newName,
      onCancel,
      onSubmit,
      progress,
      show,
      submitting,
    }
  },
})
</script>

<style lang="postcss">
#files-editf-progress {
  margin: 0 0.75rem;
}
</style>
