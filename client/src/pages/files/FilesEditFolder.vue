<!--
FilesEditFolder displays the modal popup for editing a folder.
-->

<template lang="pug">
Modal(ref='modal')
  SForm#files-edit-folder-form(
    dialog,
    :title='editFolder ? "Edit Folder" : "Add Folder"',
    submitLabel='OK',
    :disabled='confirmingDelete',
    @submit='onSubmit',
    :onCancel='onCancel'
  )
    SFInput#files-edit-folder-name(
      ref='editFolderNameRef',
      label='Name',
      trim,
      v-model='editFolderName',
      :errorFn='editFolderNameError'
    )
    SFSelect#files-edit-folder-group(
      label='Group',
      v-model='editFolderGroup',
      :options='allowedGroups',
      valueKey='id',
      labelKey='name'
    )
    SFSelect#files-edit-folder-parent(
      v-if='editFolder',
      label='Parent Folder',
      v-model='editFolderParent',
      :options='allowedParents',
      valueKey='id',
      labelKey='name'
    )
    template(v-if='editFolder', #extraButtons)
      SButton(@click='onDelete', :disabled='confirmingDelete', variant='danger') Delete
    template(v-if='confirmingDelete', #feedback)
      #files-edit-folder-confirm-label.form-item.
        Are you sure you want to delete the folder "{{editFolder.name}}",
        along with all of its files and sub-folders?
        It cannot be restored.
      #files-edit-folder-confirm-buttons.form-item
        SButton(variant='secondary', @click='cancelDelete') Keep
        SButton(variant='danger', @click='confirmDelete') Delete
</template>

<script lang="ts">
import { computed, defineComponent, onMounted, PropType, ref, toRefs, watchEffect } from 'vue'
import axios from '../../plugins/axios'
import { Modal, SButton, SFInput, SForm, SFSelect } from '../../base'
import type {
  GetFolderEditAllowedGroup,
  GetFolderEditAllowedParent,
  GetFolderChild,
  GetFolderEdit,
} from '../Files.vue'

export default defineComponent({
  components: { Modal, SButton, SFInput, SForm, SFSelect },
  props: {
    allowedGroups: { type: Array as PropType<Array<GetFolderEditAllowedGroup>>, required: true },
    allowedParents: { type: Array as PropType<Array<GetFolderEditAllowedParent>>, required: true },
    editFolder: Object as PropType<GetFolderChild>,
    parentID: { type: Number, required: true },
    siblings: { type: Array as PropType<Array<GetFolderChild>>, required: true },
  },
  emits: ['delete'],
  setup(props, { emit }) {
    const { allowedParents, editFolder, parentID } = toRefs(props)

    // Show the modal.
    const modal = ref(null as any)
    function show() {
      confirmingDelete.value = false
      return modal.value.show()
    }

    // Copy information about the folder being edited into local refs.
    const editFolderName = ref('')
    const editFolderGroup = ref(0)
    const editFolderParent = ref(0)
    watchEffect(() => {
      editFolderName.value = editFolder?.value?.name || ''
      editFolderGroup.value = editFolder?.value?.group || 0
      editFolderParent.value = parentID.value
    })

    // Focus on the name field when mounted.
    const editFolderNameRef = ref(null as null | HTMLInputElement)
    onMounted(() => {
      if (editFolderNameRef.value) editFolderNameRef.value.focus()
    })

    // Filter the overall list of allowed parents to remove the folder being
    // edited and its descendants.
    const filteredParents = computed(() => {
      if (!editFolder!.value) return []
      let ignoreAbove: null | number
      return allowedParents.value.filter((p) => {
        if (p.id === editFolder!.value!.id) {
          ignoreAbove = p.indent
          return false
        }
        if (ignoreAbove !== null && p.indent > ignoreAbove) {
          return false
        }
        ignoreAbove = null
        return true
      })
    })

    // Validate the folder name.
    function editFolderNameError(lostFocus: boolean): string {
      if (!lostFocus) return ''
      if (!editFolderName.value) return 'The folder name is required.'
      if (props.siblings.find((f) => f !== editFolder?.value && f.name === editFolderName.value))
        return 'This folder name is already in use.'
      return ''
    }

    // Submit the folder edit.
    async function onSubmit() {
      const body = new FormData()
      body.append('name', editFolderName.value)
      body.append('group', editFolderGroup.value.toString())
      if (editFolder?.value) body.append('parent', editFolderParent.value.toString())
      let response: GetFolderEdit
      if (editFolder?.value)
        response = (await axios.put<GetFolderEdit>(`/api/folders/${editFolder.value.id}`, body))
          .data
      else response = (await axios.post<GetFolderEdit>(`/api/folders/${parentID.value}`, body)).data
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
      const response = (await axios.delete<GetFolderEdit>(`/api/folders/${editFolder!.value!.id}`))
        .data
      modal.value.close(response)
    }

    return {
      allowedParents: filteredParents,
      cancelDelete,
      confirmDelete,
      confirmingDelete,
      editFolderGroup,
      editFolderName,
      editFolderNameError,
      editFolderNameRef,
      editFolderParent,
      modal,
      onCancel,
      onDelete,
      onSubmit,
      show,
    }
  },
})
</script>

<style lang="postcss">
#files-edit-folder-confirm-label {
  border-top: 1px solid #dee2e6;
  padding: 0.75rem;
}
#files-edit-folder-confirm-buttons {
  padding: 0 0.75rem 0.75rem;
  text-align: right;
  & .sbtn {
    margin-left: 0.5rem;
  }
}
</style>
