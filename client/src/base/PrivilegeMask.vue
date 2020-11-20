<!--
PrivilegeMask displays the privilege choices for a role acting on a group.
-->

<template lang="pug">
.privileges
  PrivilegeMaskButton(label='M', v-model='modelValue.member', style='margin-right:0.5rem')
  PrivilegeMaskButton(label='R', v-model='modelValue.roster', variant='left')
  PrivilegeMaskButton(label='C', v-model='modelValue.contact', variant='inner')
  PrivilegeMaskButton(label='A', v-model='modelValue.admin', variant='inner')
  PrivilegeMaskButton(label='E', v-model='modelValue.events', variant='inner')
  PrivilegeMaskButton(label='F', v-model='modelValue.folders', variant='inner')
  PrivilegeMaskButton(label='T', v-model='modelValue.texts', variant='inner')
  PrivilegeMaskButton(label='@', v-model='modelValue.emails', variant='inner')
  PrivilegeMaskButton(label='B', v-model='modelValue.bcc', variant='right')
</template>

<script lang="ts">
import { defineComponent, PropType, watch } from 'vue'
import type { Privileges } from './privileges'
import PrivilegeMaskButton from './PrivilegeMaskButton.vue'

export default defineComponent({
  components: { PrivilegeMaskButton },
  props: {
    modelValue: { type: Object as PropType<Privileges>, required: true },
  },
  emits: ['update:modelValue'],
  setup(props, { emit }) {
    watch(props.modelValue, () => {
      emit('update:modelValue', props.modelValue)
    })
    watch(
      () => props.modelValue.roster,
      () => {
        if (!props.modelValue.roster)
          props.modelValue.contact = props.modelValue.admin = props.modelValue.events = false
      }
    )
    watch(
      () => props.modelValue.contact,
      () => {
        if (props.modelValue.contact) props.modelValue.roster = true
        else props.modelValue.texts = false
      }
    )
    watch(
      () => props.modelValue.admin,
      () => {
        if (props.modelValue.admin) props.modelValue.roster = true
      }
    )
    watch(
      () => props.modelValue.events,
      () => {
        if (props.modelValue.events) props.modelValue.roster = true
      }
    )
    watch(
      () => props.modelValue.texts,
      () => {
        if (props.modelValue.texts) props.modelValue.contact = true
      }
    )
  },
})
</script>

<style lang="postcss">
.privileges {
  margin: 0.125rem 0;
  white-space: nowrap;
}
</style>
