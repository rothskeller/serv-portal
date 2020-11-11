<!--
FilesFolder displays one sub-folder of the folder being viewed.
-->

<template lang="pug">
.files-line(@mouseover='hovering = true', @mouseout='hovering = false')
  .files-icon(@click='onEdit')
    SIcon(v-if='canEdit && hovering', icon='edit')
    SIcon(v-else, icon='folder')
  .files-name
    router-link(:to='`/files/${folder.id}`') {{ folder.name }}
    span.files-pending(v-if='folder.approvals') [pending approvals]
</template>

<script lang="ts">
import { defineComponent, PropType, ref } from 'vue'
import { SIcon } from '../../base'
import type { GetFolderChild } from '../Files.vue'

export default defineComponent({
  components: { SIcon },
  props: {
    folder: { type: Object as PropType<GetFolderChild>, required: true },
    canEdit: { type: Boolean, required: true },
  },
  emits: ['edit'],
  setup(props, { emit }) {
    // Handle editing.
    const hovering = ref(false)
    function onEdit() {
      if (props.canEdit) emit('edit')
    }

    return { hovering, onEdit }
  },
})
</script>
