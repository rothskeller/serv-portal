<!--
FilesEditDocument displays the document editor.
-->

<template lang="pug">
Modal(ref='modal')
  SForm(
    dialog,
    title='Add File',
    submitLabel='OK',
    :disabled='confirmingReplace || uploadProgress !== null',
    @submit='onSubmit',
    :onCancel='onCancel'
  )
    SFFile#files-edit-doc-files(multiple, label='File(s)', v-model='files', :errorFn='filesError')
    template(v-if='confirmingReplace', #feedback)
      #files-add-doc-confirm-label.form-item.
        This will replace {{ replace.length > 1 ? `${replace.length} existing files with the same names` : `the existing file named "${replace[0]}"` }}. Are you sure?
      #files-add-doc-confirm-buttons.form-item
        SButton(variant='secondary', @click='cancelReplace') Keep
        SButton(variant='warning', @click='confirmReplace') Replace
    template(v-if='uploadProgress !== null', #feedback)
      .form-item
        SProgress(:value='uploadProgress')
</template>

<script lang="ts">
import { defineComponent, onMounted, PropType, ref, watchEffect } from 'vue'
import axios from '../../plugins/axios'
import { Modal, SButton, SFFile, SForm, SProgress } from '../../base'
import type { GetFolderDocument, GetFolderEdit } from '../Files.vue'

export default defineComponent({
  components: { Modal, SButton, SFFile, SForm, SProgress },
  props: {
    parentID: { type: Number, required: true },
    siblings: { type: Array as PropType<Array<GetFolderDocument>>, required: true },
  },
  setup(props) {
    // Show the modal on request.
    const modal = ref(null as any)
    function show() {
      confirmingReplace.value = false
      return modal.value.show()
    }

    // Validate the list of files.
    const files = ref(null as null | FileList)
    function filesError(lostFocus: boolean, submitted: boolean): string {
      if (!submitted) return ''
      if (!files.value || !files.value.length) return 'Select the file(s) to be added.'
      return ''
    }

    // Submit the document edit.
    const replace = ref([] as Array<string>)
    const confirmingReplace = ref(false)
    const uploadProgress = ref(null as null | number)
    function onUploadProgress(evt: ProgressEvent) {
      if (evt.lengthComputable && evt.total) uploadProgress.value = evt.loaded / evt.total
    }
    function onSubmit() {
      for (let idx = 0; idx < files.value!.length; idx++) {
        const file = files.value!.item(idx)!
        if (props.siblings.find((d) => d.name === file.name)) replace.value.push(file.name)
      }
      if (replace.value.length) confirmingReplace.value = true
      else confirmReplace()
    }
    function cancelReplace() {
      confirmingReplace.value = false
    }
    async function confirmReplace() {
      confirmingReplace.value = false
      uploadProgress.value = 0
      const body = new FormData()
      for (let idx = 0; idx < files.value!.length; idx++) {
        const file = files.value!.item(idx)!
        body.append('file', file)
      }
      const updated = (
        await axios.post(`/api/folders/${props.parentID}/NEW`, body, {
          onUploadProgress,
        })
      ).data
      uploadProgress.value = null
      modal.value.close(updated)
    }
    function onCancel() {
      modal.value.close(null)
    }

    return {
      cancelReplace,
      confirmReplace,
      confirmingReplace,
      files,
      filesError,
      modal,
      replace,
      onCancel,
      onSubmit,
      show,
      uploadProgress,
    }
  },
})
</script>

<style lang="postcss">
#files-add-doc-confirm-label {
  border-top: 1px solid #dee2e6;
  padding: 0.75rem;
}
#files-add-doc-confirm-buttons {
  padding: 0 0.75rem 0.75rem;
  text-align: right;
  & .sbtn {
    margin-left: 0.5rem;
  }
}
</style>
