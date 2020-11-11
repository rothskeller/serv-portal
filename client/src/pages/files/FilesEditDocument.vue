<!--
FilesEditDocument displays the document editor.
-->

<template lang="pug">
Modal(ref='modal')
  SForm(
    v-if='doc',
    dialog,
    title='Edit File',
    submitLabel='OK',
    :disabled='confirmingDelete',
    @submit='onSubmit',
    @cancel='onCancel'
  )
    SFInput#files-edit-doc-name(
      ref='nameRef',
      label='Name',
      trim,
      v-model='name',
      :errorFn='nameError'
    )
    SFSelect#files-edit-doc-folder(
      label='In Folder',
      :options='allowedParents',
      v-model='folder',
      valueKey='id',
      labelKey='name'
    )
    template(#extraButtons)
      SButton(variant='danger', :disabled='confirmingDelete', @click='onDelete') Delete
    template(v-if='confirmingDelete', #feedback)
      #files-edit-doc-confirm-label.form-item.
        Are you sure you want to delete the document "{{doc.name}}"?
        It cannot be restored.
      #files-edit-doc-confirm-buttons.form-item
        SButton(variant='secondary', @click='cancelDelete') Keep
        SButton(variant='danger', @click='confirmDelete') Delete
</template>

<script lang="ts">
import { defineComponent, onMounted, PropType, ref, watchEffect } from 'vue'
import axios from '../../plugins/axios'
import { Modal, SButton, SFInput, SForm, SFSelect } from '../../base'
import type { GetFolderEditAllowedParent, GetFolderDocument, GetFolderEdit } from '../Files.vue'

function extension(path: string): string {
  const dot = path.lastIndexOf('.')
  return dot < 0 ? '' : path.substr(dot)
}

export default defineComponent({
  components: {
    Modal,
    SButton,
    SFInput,
    SForm,
    SFSelect,
  },
  props: {
    allowedParents: { type: Array as PropType<Array<GetFolderEditAllowedParent>>, required: true },
    doc: Object as PropType<GetFolderDocument>,
    parentID: { type: Number, required: true },
    siblings: { type: Array as PropType<Array<GetFolderDocument>>, required: true },
  },
  setup(props) {
    // Show the modal on request.
    const modal = ref(null as any)
    function show() {
      confirmingDelete.value = false
      return modal.value.show()
    }

    // Copy information about the folder being edited into local refs.
    const name = ref('')
    const folder = ref(0)
    watchEffect(() => {
      name.value = props.doc?.name || ''
      folder.value = props.parentID
    })

    // Focus on the name field when mounted.
    const nameRef = ref(null as null | HTMLInputElement)
    onMounted(() => {
      if (nameRef.value) nameRef.value.focus()
    })

    // Validate the file name.
    function nameError(lostFocus: boolean): string {
      if (!lostFocus) return ''
      if (!name.value) return 'The file name is required.'
      if (!extension(name.value)) return 'The file must have an extension indicating its type.'
      if (extension(name.value) !== extension(props.doc!.name))
        return `The new file name must end with "${extension(props.doc!.name)}".`
      if (props.siblings.find((f) => f !== props.doc && f.name === name.value))
        return 'This file name is already in use.'
      return ''
    }

    // Submit the document edit.
    async function onSubmit() {
      const body = new FormData()
      body.append('name', name.value)
      body.append('folder', folder.value.toString())
      const response = (
        await axios.post<GetFolderEdit>(`/api/folders/${props.parentID}/${props.doc!.id}`, body)
      ).data
      modal.value.close(response)
    }
    function onCancel() {
      modal.value.close(null)
    }

    // Delete.
    const confirmingDelete = ref(false)
    function onDelete() {
      confirmingDelete.value = true
    }
    function cancelDelete() {
      confirmingDelete.value = false
    }
    async function confirmDelete() {
      const response = (
        await axios.delete<GetFolderEdit>(`/api/folders/${props.parentID}/${props.doc!.id}`)
      ).data
      modal.value.close(response)
    }

    return {
      cancelDelete,
      confirmDelete,
      confirmingDelete,
      folder,
      modal,
      name,
      nameError,
      nameRef,
      onCancel,
      onDelete,
      onSubmit,
      show,
    }
  },
})
</script>

<style lang="postcss">
#files-edit-doc-confirm-label {
  border-top: 1px solid #dee2e6;
  padding: 0.75rem;
}
#files-edit-doc-confirm-buttons {
  padding: 0 0.75rem 0.75rem;
  text-align: right;
  & .sbtn {
    margin-left: 0.5rem;
  }
}
</style>
