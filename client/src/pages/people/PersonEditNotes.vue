<!--
PersonEditNotes is the dialog box for editing a person's notes.
-->

<template lang="pug">
Modal(ref='modal')
  SForm(
    dialog,
    variant='primary',
    title='Edit Notes',
    submitLabel='Save',
    :disabled='submitting',
    @submit='onSubmit',
    @cancel='onCancel'
  )
    SSpinner(v-if='loading')
    template(v-else)
      #person-edit-note-help.form-item(v-if='notesBefore.length || notesAfter.length')
        | Click on a note to edit it, or click the Add button to add one.
      .person-edit-note.form-item(v-for='note in notesBefore', @click='editNote(note)')
        div Date
        div(v-text='note.date')
        div Note
        div(v-text='note.note')
        div Vis.
        div(v-text='visibilityLabels[note.visibility]')
      template(v-if='noteEditing.date')
        SFInput#person-edit-note-date(label='Date', type='date', v-model='noteEditing.date')
        SFInput#person-edit-note-note(
          ref='noteRef',
          label='Note Text',
          trim,
          v-model='noteEditing.note'
        )
        SFSelect#person-edit-note-visibility(
          label='Visibility',
          :options='visibilityOptions',
          v-model='noteEditing.visibility'
        )
      .person-edit-note.form-item(v-for='note in notesAfter', @click='editNote(note)')
        div Date
        div(v-text='note.date')
        div Note
        div(v-text='note.note')
        div Vis.
        div(v-text='visibilityLabels[note.visibility]')
    template(#extraButtons)
      SButton(variant='primary', @click='onAdd') Add
</template>

<script lang="ts">
import { defineComponent, nextTick, ref, watch } from 'vue'
import moment from 'moment-mini'
import axios from '../../plugins/axios'
import { Modal, SButton, SForm, SFInput, SFSelect, SSpinner } from '../../base'

interface GetPersonNotesNote {
  date: string
  note: string
  visibility: string
}
interface GetPersonNotes {
  notes: Array<GetPersonNotesNote>
  visibilities: Array<string>
}

const visibilityLabels = {
  webmaster: 'Webmasters only',
  admin: 'DPS staff only',
  leader: 'SERV leads only',
  contact: 'Anyone who can see contact info',
  person: 'Anyone',
}
let visibilityOptions = [
  { value: 'webmaster', label: 'Webmasters only' },
  { value: 'admin', label: 'DPS staff only' },
  { value: 'leader', label: 'SERV leads only' },
  { value: 'contact', label: 'Anyone who can see contact info' },
  { value: 'person', label: 'Anyone' },
]

export default defineComponent({
  components: { Modal, SButton, SForm, SFInput, SFSelect, SSpinner },
  props: {
    pid: { type: Number, required: true },
  },
  setup(props) {
    const modal = ref(null as any)
    function show() {
      loadData()
      return modal.value.show()
    }

    // Load the form data.
    const notesBefore = ref([] as Array<GetPersonNotesNote>)
    const notesAfter = ref([] as Array<GetPersonNotesNote>)
    const noteEditing = ref({} as GetPersonNotesNote)
    const loading = ref(true)
    async function loadData() {
      loading.value = true
      const data = (await axios.get<GetPersonNotes>(`/api/people/${props.pid}/notes`)).data
      notesBefore.value = data.notes
      notesBefore.value.sort((a, b) => {
        if (a.date !== b.date) return a.date < b.date ? -1 : 1
        return a.note.localeCompare(b.note)
      })
      notesAfter.value = []
      visibilityOptions = visibilityOptions.filter((v) => data.visibilities.includes(v.value))
      if (notesBefore.value.length) noteEditing.value = { date: '', note: '', visibility: '' }
      else
        noteEditing.value = {
          date: moment().format('YYYY-MM-DD'),
          note: '',
          visibility: visibilityOptions[0].value,
        }
      loading.value = false
    }

    // Editing an existing note.
    const noteRef = ref(null as any)
    function editNote(note: GetPersonNotesNote) {
      if (noteEditing.value.date && noteEditing.value.note) {
        notesBefore.value.push(noteEditing.value)
      }
      notesAfter.value.forEach((n) => notesBefore.value.push(n))
      const idx = notesBefore.value.findIndex((n) => n === note)
      notesAfter.value = notesBefore.value.splice(idx + 1)
      noteEditing.value = notesBefore.value.splice(idx, 1)[0]
      nextTick(() => {
        noteRef.value.focus()
      })
    }

    // Adding a new note.
    function onAdd() {
      if (noteEditing.value.date && noteEditing.value.note) {
        notesBefore.value.push(noteEditing.value)
      }
      notesAfter.value.forEach((n) => notesBefore.value.push(n))
      notesAfter.value = []
      noteEditing.value = {
        date: moment().format('YYYY-MM-DD'),
        note: '',
        visibility: visibilityOptions[0].value,
      }
      nextTick(() => {
        noteRef.value.focus()
      })
    }

    // Save and close.
    const submitting = ref(false)
    async function onSubmit() {
      var body = new FormData()
      notesBefore.value.forEach((n) => {
        body.append('note', n.note)
        body.append('date', n.date)
        body.append('visibility', n.visibility)
      })
      if (noteEditing.value.date && noteEditing.value.note) {
        body.append('note', noteEditing.value.note)
        body.append('date', noteEditing.value.date)
        body.append('visibility', noteEditing.value.visibility)
      }
      notesAfter.value.forEach((n) => {
        body.append('note', n.note)
        body.append('date', n.date)
        body.append('visibility', n.visibility)
      })
      submitting.value = true
      await axios.post(`/api/people/${props.pid}/notes`, body)
      submitting.value = false
      modal.value.close(true)
    }
    function onCancel() {
      modal.value.close(false)
    }

    return {
      editNote,
      loading,
      modal,
      noteEditing,
      noteRef,
      notesAfter,
      notesBefore,
      onAdd,
      onCancel,
      onSubmit,
      show,
      submitting,
      visibilityLabels,
      visibilityOptions,
    }
  },
})
</script>

<style lang="postcss">
#person-edit-note-help {
  margin: 0 0.75rem 1rem;
}
.person-edit-note {
  display: grid;
  grid: auto-flow / min-content 1fr;
  margin: 0 0.75rem 1rem;
  column-gap: 1rem;
}
</style>
