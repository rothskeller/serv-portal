<!--
FilesTrash displays the trash icon and receives documents and folders dragged
to it.
-->

<template lang="pug">
.files-trash(
  :style='dragOverStyle',
  @dragenter='onDragEnter',
  @dragover='onDragOver',
  @dragleave='onDragLeave',
  @drop='onDrop'
)
  SIcon.files-trash-icon(icon='trash')
  span.files-trash-text Drag here to delete
  MessageBox(
    ref='deleteModal',
    :title='deleteTitle',
    cancelLabel='Keep',
    okLabel='Delete',
    variant='danger'
  )
    | {{ deleteText }}
</template>

<script lang="ts">
import { computed, defineComponent, inject, PropType, ref, Ref } from 'vue'
import axios from '../../plugins/axios'
import { MessageBox, SIcon } from '../../base'

export default defineComponent({
  components: { MessageBox, SIcon },
  emits: ['reload', 'showTrash'],
  setup(props, { emit }) {
    const dragOverStyle = computed(() => (draggingOver.value ? { backgroundColor: '#ccc' } : ''))

    // Drop support.
    const draggingOver = ref(0)
    const deleteText = ref('')
    const deleteTitle = ref('')
    const deleteModal = ref(null as any)
    function dropSupported(evt: DragEvent): boolean {
      if (evt.dataTransfer!.types.includes('x-serv-folder')) return true
      if (evt.dataTransfer!.types.includes('x-serv-document')) return true
      return false
    }
    function onDragEnter(evt: DragEvent) {
      if (!dropSupported(evt)) return
      evt.preventDefault()
      draggingOver.value++
    }
    function onDragOver(evt: DragEvent) {
      if (!dropSupported(evt)) return
      evt.preventDefault()
    }
    function onDragLeave(evt: DragEvent) {
      if (!dropSupported(evt)) return
      evt.preventDefault()
      if (draggingOver.value) draggingOver.value--
    }
    async function onDrop(evt: DragEvent) {
      if (!dropSupported(evt)) return
      evt.preventDefault()
      emit('showTrash', true) // prevent ourselves from disappearing
      draggingOver.value = 0
      if (evt.dataTransfer!.types.includes('x-serv-folder')) {
        onDropFolder(evt.dataTransfer!)
      } else {
        onDropDocument(evt.dataTransfer!)
      }
    }
    async function onDropFolder(dt: DataTransfer) {
      const fdd: { name: string; url: string } = JSON.parse(dt.getData('x-serv-folder'))
      const parent = fdd.url.replace(/\/[^/]*$/, '')
      deleteTitle.value = 'Delete Folder'
      deleteText.value = `Are you sure you want to delete the folder “${fdd.name}”, along with all of its files and sub-folders?  It cannot be restored.`
      if (!(await deleteModal.value.show())) {
        emit('showTrash', false)
        return
      }
      await axios.delete(`/api/folders${fdd.url}`)
      emit('showTrash', false)
      emit('reload', function (u: string): string {
        if (u === fdd.url) return parent
        if (u.startsWith(fdd.url + '/')) return parent
        return u
      })
    }
    async function onDropDocument(dt: DataTransfer) {
      const ddd: { name: string; url: string; type: string } = JSON.parse(
        dt.getData('x-serv-document')
      )
      if (ddd.type === 'link') {
        deleteTitle.value = 'Delete Link'
        deleteText.value = `Are you sure you want to delete the link “${ddd.name}”?  It cannot be restored.`
      } else {
        deleteTitle.value = 'Delete File'
        deleteText.value = `Are you sure you want to delete the file “${ddd.name}”?  It cannot be restored.`
      }
      if (!(await deleteModal.value.show())) {
        emit('showTrash', false)
        return
      }
      await axios.delete(`/api/document${ddd.url}`)
      emit('showTrash', false)
      emit('reload')
    }
    return {
      deleteModal,
      deleteText,
      deleteTitle,
      dragOverStyle,
      onDragEnter,
      onDragLeave,
      onDragOver,
      onDrop,
    }
  },
})
</script>

<style lang="postcss">
.files-trash {
  display: flex;
  align-items: center;
  min-height: 0;
  min-width: 0;
  color: red;
  margin-top: 1.5rem;
}
.files-trash-icon {
  flex: none;
  height: 1rem;
  width: 1.5rem;
  margin-right: 0.25rem;
}
.files-trash-text {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
