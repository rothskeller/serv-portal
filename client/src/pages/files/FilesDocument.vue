<!--
FilesDocument displays a single document in the files hierarchy.
-->

<template lang="pug">
.files-doc(
  :style='indentStyle',
  :draggable='folder.canEdit',
  @dragstart='onDragStart',
  @dragend='onDragEnd',
  @mouseover='onMouseOver',
  @mouseout='onMouseOut'
)
  SIcon.files-doc-icon(:icon='icon', :style='iconStyle', @click='onEdit')
  a.files-doc-link(:href='doc.url', :target='target', draggable='false', v-text='doc.name')
  FilesEditFile(v-if='folder.canEdit', ref='editFileModal')
  FilesEditLink(v-if='folder.canEdit', ref='editLinkModal')
</template>

<script lang="ts">
import { computed, defineComponent, PropType, ref } from 'vue'
import SIcon from '../../base/SIcon.vue'
import type { GetFolderDocument, GetFolderFolder } from '../Files.vue'
import FilesEditFile from './FilesEditFile.vue'
import FilesEditLink from './FilesEditLink.vue'

export default defineComponent({
  components: { FilesEditFile, FilesEditLink, SIcon },
  props: {
    folder: { type: Object as PropType<GetFolderFolder>, required: true },
    doc: { type: Object as PropType<GetFolderDocument>, required: true },
    indent: { type: Number, required: true },
  },
  emits: ['reload', 'showTrash'],
  setup(props, { emit }) {
    const isLink = computed(() => !props.doc.url.startsWith('/'))
    const icon = computed(() =>
      hovering.value && props.folder.canEdit
        ? 'edit'
        : isLink.value
        ? 'link'
        : props.doc.name.endsWith('.pdf')
        ? 'pdf'
        : props.doc.name.endsWith('.docx') || props.doc.name.endsWith('.doc')
        ? 'word'
        : props.doc.name.endsWith('.pptx') || props.doc.name.endsWith('.ppt')
        ? 'powerpoint'
        : props.doc.name.endsWith('.xlsx') || props.doc.name.endsWith('.xls')
        ? 'excel'
        : props.doc.name.endsWith('.jpeg') ||
          props.doc.name.endsWith('.jpg') ||
          props.doc.name.endsWith('.png')
        ? 'image'
        : 'file'
    )
    const iconStyle = computed(() =>
      hovering.value && props.folder.canEdit ? { color: 'blue' } : null
    )
    const indentStyle = computed(() => ({ marginLeft: `${props.indent * 1.5}em` }))
    const target = computed(() => (props.doc.newtab ? '_blank' : null))

    // It's possible to drag documents to other folders, or to the trash.
    function onDragStart(evt: DragEvent) {
      if (!props.folder.canEdit) return
      evt.dataTransfer!.setData(
        'x-serv-document',
        JSON.stringify({
          name: props.doc.name,
          url: `${props.folder.url}/${encodeURIComponent(props.doc.name)}`,
          type: isLink.value ? 'link' : 'file',
        })
      )
      const url = isLink.value ? props.doc.url : `https://sunnyvaleserv.org${props.doc.url}`
      evt.dataTransfer!.setData('text/uri-list', `${url}\r\n# ${props.doc.name}`)
      evt.dataTransfer!.setData('text/plain', url)
      emit('showTrash', true)
    }
    function onDragEnd() {
      emit('showTrash', false)
    }

    // Edit button support.
    const hovering = ref(false)
    const editFileModal = ref(null as any)
    const editLinkModal = ref(null as any)
    function onMouseOver() {
      hovering.value = true
    }
    function onMouseOut() {
      hovering.value = false
    }
    async function onEdit() {
      if (!props.folder.canEdit) return
      if (isLink.value) {
        if (await editLinkModal.value.showEdit(props.folder.url, props.doc.name)) emit('reload')
      } else {
        if (await editFileModal.value.show(props.folder.url, props.doc.name)) emit('reload')
      }
    }

    return {
      editFileModal,
      editLinkModal,
      icon,
      iconStyle,
      indentStyle,
      onDragEnd,
      onDragStart,
      onEdit,
      onMouseOver,
      onMouseOut,
      target,
    }
  },
})
</script>

<style lang="postcss">
.files-doc {
  display: flex;
  align-items: center;
  min-height: 0;
  min-width: 0;
}
.files-doc-icon {
  height: 1rem;
  width: 1.5rem;
  margin-right: 0.25rem; /* see note in FilesFolder.vue */
  flex: none;
}
.files-doc-link {
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
}
</style>
