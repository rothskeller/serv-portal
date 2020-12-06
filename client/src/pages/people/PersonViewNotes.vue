<!--
PersonViewNotes displays the notes on a person's record.
-->

<template lang="pug">
PersonViewSection(
  v-if='person.notes.length || person.canEditNotes',
  title='Notes',
  :editable='person.canEditNotes',
  @edit='onEditNotes'
)
  #person-view-notes
    .person-view-note(v-for='note in person.notes')
      .person-view-note-date(v-text='note.date')
      .person-view-note-text(v-text='note.note')
    div(v-if='!person.notes.length') No notes on file.
  PersonEditNotes(v-if='person.canEdit', ref='editNotesModal', :pid='person.id')
</template>

<script lang="ts">
import { computed, defineComponent, PropType, ref } from 'vue'
import { GetPerson } from './PersonView.vue'
import PersonEditNotes from './PersonEditNotes.vue'
import PersonViewSection from './PersonViewSection.vue'

export default defineComponent({
  components: { PersonEditNotes, PersonViewSection },
  props: {
    person: { type: Object as PropType<GetPerson>, required: true },
  },
  emits: ['reload'],
  setup(props, { emit }) {
    const editNotesModal = ref(null as any)
    async function onEditNotes() {
      if (await editNotesModal.value.show()) emit('reload')
    }
    return {
      editNotesModal,
      onEditNotes,
    }
  },
})
</script>

<style lang="postcss">
#person-view-notes {
  margin-top: 0.75rem;
}
.person-view-note {
  line-height: 1.2;
  margin-left: 2rem;
  text-indent: -2rem;
}
.person-view-note-date {
  display: inline;
  margin-right: 0.75rem;
  color: #666;
  font-variant: tabular-nums;
}
.person-view-note-text {
  display: inline;
  margin-left: 0.5rem;
}
</style>
