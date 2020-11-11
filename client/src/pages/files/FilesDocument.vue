<!--
FilesDocument displays one document in the folder being viewed.
-->

<template lang="pug">
.files-line(@mouseover='hovering = true', @mouseout='hovering = false')
  .files-icon(@click='onEdit')
    SIcon(v-if='canEdit && hovering', icon='edit')
    SIcon(v-else-if='doc.name.endsWith(".docx") || doc.name.endsWith(".doc")', icon='word')
    SIcon(v-else-if='doc.name.endsWith(".pdf")', icon='pdf')
    SIcon(v-else-if='doc.name.endsWith(".ppt") || doc.name.endsWith(".pptx")', icon='powerpoint')
    SIcon(v-else-if='doc.name.endsWith(".xls") || doc.name.endsWith(".xlsx")', icon='excel')
    SIcon(v-else-if='doc.name.endsWith(".jpg") || doc.name.endsWith(".jpeg")', icon='image')
    SIcon(v-else, icon='file')
  span.files-name
    a(:href='docHRef', download) {{ doc.name }}
    span.files-pending(v-if='doc.needsApproval') [pending approval]
</template>

<script lang="ts">
import { computed, defineComponent, PropType, ref } from 'vue'
import { SIcon } from '../../base'
import type { GetFolderDocument } from '../Files.vue'

export default defineComponent({
  components: { SIcon },
  props: {
    folderID: { type: Number, required: true },
    doc: { type: Object as PropType<GetFolderDocument>, required: true },
    canEdit: { type: Boolean, required: true },
  },
  emits: ['edit'],
  setup(props, { emit }) {
    const hovering = ref(false)
    function onEdit() {
      if (props.canEdit) emit('edit')
    }
    const docHRef = computed(() => `/api/folders/${props.folderID}/${props.doc.id}`)
    return { docHRef, hovering, onEdit }
  },
})
</script>
