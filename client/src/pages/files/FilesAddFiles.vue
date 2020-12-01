<!--
FilesAddFiles is the dialog box for creating file documents.
-->

<template lang="pug">
Modal(ref='modal')
  SForm(
    dialog,
    variant='primary',
    title='Add Files',
    submitLabel='Add',
    :disabled='submitting',
    @submit='onSubmit',
    @cancel='onCancel'
  )
    SFFile#files-addf-files(label='File(s)', multiple, v-model='files', :errorFn='filesError')
    SProgress#files-addf-progress.form-item(v-if='progress', :value='progress')
</template>

<script lang="ts">
import { defineComponent, ref, watch } from 'vue'
import axios from '../../plugins/axios'
import { Modal, SForm, SFFile, SProgress } from '../../base'

export default defineComponent({
  components: { Modal, SForm, SFFile, SProgress },
  setup() {
    const modal = ref(null as any)
    let folderURL = ''
    function show(folder: string) {
      folderURL = folder
      return modal.value.show()
    }

    // The files field.
    const files = ref(null as null | FileList)
    const filesConflict = ref('')
    function filesError(lostFocus: boolean, submitted: boolean): string {
      if (!submitted) return ''
      if (!files.value || !files.value.length) return 'Select the file(s) to be added.'
      for (let i = 0; i < files.value.length; i++) {
        const file = files.value.item(i)!
        if (
          !file.name ||
          file.name[0] === '.' ||
          file.name.includes(':') ||
          file.name.includes('/')
        )
          return `“${file.name}” is not a legal file name.`
      }
      if (filesConflict.value) return filesConflict.value
      return ''
    }
    watch(files, () => {
      filesConflict.value = ''
    })

    // Save and close.
    const submitting = ref(false)
    async function onSubmit() {
      const body = new FormData()
      submitting.value = true
      try {
        for (let i = 0; i < files.value!.length; i++)
          body.append('name', files.value!.item(i)!.name)
        await axios.post(`/api/folders${folderURL}?op=checkNames`, body)
        body.delete('name')
        for (let i = 0; i < files.value!.length; i++) body.append('file', files.value!.item(i)!)
        await axios.post(`/api/folders${folderURL}?op=newFiles`, body, { onUploadProgress })
        modal.value.close(true)
      } catch (err) {
        if (!err.response || err.response.status !== 409) throw err
        filesConflict.value = err.response.data
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

    return { files, filesError, modal, onCancel, onSubmit, progress, show, submitting }
  },
})
</script>

<style lang="postcss">
#files-addf-progress {
  margin: 0 0.75rem;
}
</style>
