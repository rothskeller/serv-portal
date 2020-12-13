<!--
PersonViewSection is the framework for one section on the PersonView page.
-->

<template lang="pug">
.person-view-section(:class='{ "person-view-section-editable": editable }')
  .person-view-section-header
    .person-view-section-header-text(v-text='title')
    .person-view-section-edit(v-if='editable')
      SButton.person-view-section-ebtn(variant='primary', small, @click='$emit("edit")') Edit
  slot
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { SButton } from '../../base'

export default defineComponent({
  components: { SButton },
  props: {
    title: { type: String, required: true },
    editable: { type: Boolean, required: true },
  },
  emits: ['edit'],
})
</script>

<style lang="postcss">
.person-view-section {
  margin-bottom: 1.5rem;
}
.person-view-section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background-color: #eee;
}
.person-view-section-header-text {
  color: black;
  font-size: 1.25rem;
  padding-left: 0.75rem;
}
.person-view-section-edit {
  flex: none;
  padding-left: 0.25rem;
  background-color: white;
}
.person-view-section-ebtn {
}
</style>
