<!--
Files displays the Files browser.
-->

<template lang="pug">
div(v-if="!folder")
  b-spinner(small)
#files(v-else @drop="onDrop" @dragover="onDragOver")
  .files-line(v-if="folder.id")
    .files-icon
      svg(xmlns="http://www.w3.org/2000/svg" viewBox="0 0 448 512")
        path(fill="currentColor" d="M34.9 289.5l-22.2-22.2c-9.4-9.4-9.4-24.6 0-33.9L207 39c9.4-9.4 24.6-9.4 33.9 0l194.3 194.3c9.4 9.4 9.4 24.6 0 33.9L413 289.4c-9.5 9.5-25 9.3-34.3-.4L264 168.6V456c0 13.3-10.7 24-24 24h-32c-13.3 0-24-10.7-24-24V168.6L69.2 289.1c-9.3 9.8-24.8 10-34.3.4z")
    b-link.files-name(:to="`/files/${folder.parent ? folder.parent.id : 0}`") {{folder.parent ? folder.parent.name : 'Files'}}
  .files-line(v-for="child in folder.children" :key="`f${child.id}`" @mouseover="onHoverFolder(child.id)" @mouseout="onHoverFolder(0)")
    .files-icon(@click="onEditFolder(child)")
      svg(v-if="folder.canEdit && hoverFolder === child.id" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 576 512")
        path(fill="currentColor" d="M402.3 344.9l32-32c5-5 13.7-1.5 13.7 5.7V464c0 26.5-21.5 48-48 48H48c-26.5 0-48-21.5-48-48V112c0-26.5 21.5-48 48-48h273.5c7.1 0 10.7 8.6 5.7 13.7l-32 32c-1.5 1.5-3.5 2.3-5.7 2.3H48v352h352V350.5c0-2.1.8-4.1 2.3-5.6zm156.6-201.8L296.3 405.7l-90.4 10c-26.2 2.9-48.5-19.2-45.6-45.6l10-90.4L432.9 17.1c22.9-22.9 59.9-22.9 82.7 0l43.2 43.2c22.9 22.9 22.9 60 .1 82.8zM460.1 174L402 115.9 216.2 301.8l-7.3 65.3 65.3-7.3L460.1 174zm64.8-79.7l-43.2-43.2c-4.1-4.1-10.8-4.1-14.8 0L436 82l58.1 58.1 30.9-30.9c4-4.2 4-10.8-.1-14.9z")
      svg(v-else xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512")
        path(fill="currentColor" d="M464 128H272l-54.63-54.63c-6-6-14.14-9.37-22.63-9.37H48C21.49 64 0 85.49 0 112v288c0 26.51 21.49 48 48 48h416c26.51 0 48-21.49 48-48V176c0-26.51-21.49-48-48-48zm0 272H48V112h140.12l54.63 54.63c6 6 14.14 9.37 22.63 9.37H464v224z")
    b-link.files-name(:to="`/files/${child.id}`") {{child.name}}
  .files-line(v-for="doc in folder.documents" :key="`d${doc.id}`" @mouseover="onHoverDocument(doc.id)" @mouseout="onHoverDocument(0)")
    .files-icon(@click="onEditDocument(doc)")
      svg(v-if="folder.canEdit && hoverDocument === doc.id" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 576 512")
        path(fill="currentColor" d="M402.3 344.9l32-32c5-5 13.7-1.5 13.7 5.7V464c0 26.5-21.5 48-48 48H48c-26.5 0-48-21.5-48-48V112c0-26.5 21.5-48 48-48h273.5c7.1 0 10.7 8.6 5.7 13.7l-32 32c-1.5 1.5-3.5 2.3-5.7 2.3H48v352h352V350.5c0-2.1.8-4.1 2.3-5.6zm156.6-201.8L296.3 405.7l-90.4 10c-26.2 2.9-48.5-19.2-45.6-45.6l10-90.4L432.9 17.1c22.9-22.9 59.9-22.9 82.7 0l43.2 43.2c22.9 22.9 22.9 60 .1 82.8zM460.1 174L402 115.9 216.2 301.8l-7.3 65.3 65.3-7.3L460.1 174zm64.8-79.7l-43.2-43.2c-4.1-4.1-10.8-4.1-14.8 0L436 82l58.1 58.1 30.9-30.9c4-4.2 4-10.8-.1-14.9z")
      svg(v-else-if="doc.name.endsWith('.docx') || doc.name.endsWith('.doc')" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 384 512")
        path(fill="currentColor" d="M369.9 97.9L286 14C277 5 264.8-.1 252.1-.1H48C21.5 0 0 21.5 0 48v416c0 26.5 21.5 48 48 48h288c26.5 0 48-21.5 48-48V131.9c0-12.7-5.1-25-14.1-34zM332.1 128H256V51.9l76.1 76.1zM48 464V48h160v104c0 13.3 10.7 24 24 24h104v288H48zm220.1-208c-5.7 0-10.6 4-11.7 9.5-20.6 97.7-20.4 95.4-21 103.5-.2-1.2-.4-2.6-.7-4.3-.8-5.1.3.2-23.6-99.5-1.3-5.4-6.1-9.2-11.7-9.2h-13.3c-5.5 0-10.3 3.8-11.7 9.1-24.4 99-24 96.2-24.8 103.7-.1-1.1-.2-2.5-.5-4.2-.7-5.2-14.1-73.3-19.1-99-1.1-5.6-6-9.7-11.8-9.7h-16.8c-7.8 0-13.5 7.3-11.7 14.8 8 32.6 26.7 109.5 33.2 136 1.3 5.4 6.1 9.1 11.7 9.1h25.2c5.5 0 10.3-3.7 11.6-9.1l17.9-71.4c1.5-6.2 2.5-12 3-17.3l2.9 17.3c.1.4 12.6 50.5 17.9 71.4 1.3 5.3 6.1 9.1 11.6 9.1h24.7c5.5 0 10.3-3.7 11.6-9.1 20.8-81.9 30.2-119 34.5-136 1.9-7.6-3.8-14.9-11.6-14.9h-15.8z")
      svg(v-else-if="doc.name.endsWith('.pdf')" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 384 512")
        path(fill="currentColor" d="M369.9 97.9L286 14C277 5 264.8-.1 252.1-.1H48C21.5 0 0 21.5 0 48v416c0 26.5 21.5 48 48 48h288c26.5 0 48-21.5 48-48V131.9c0-12.7-5.1-25-14.1-34zM332.1 128H256V51.9l76.1 76.1zM48 464V48h160v104c0 13.3 10.7 24 24 24h104v288H48zm250.2-143.7c-12.2-12-47-8.7-64.4-6.5-17.2-10.5-28.7-25-36.8-46.3 3.9-16.1 10.1-40.6 5.4-56-4.2-26.2-37.8-23.6-42.6-5.9-4.4 16.1-.4 38.5 7 67.1-10 23.9-24.9 56-35.4 74.4-20 10.3-47 26.2-51 46.2-3.3 15.8 26 55.2 76.1-31.2 22.4-7.4 46.8-16.5 68.4-20.1 18.9 10.2 41 17 55.8 17 25.5 0 28-28.2 17.5-38.7zm-198.1 77.8c5.1-13.7 24.5-29.5 30.4-35-19 30.3-30.4 35.7-30.4 35zm81.6-190.6c7.4 0 6.7 32.1 1.8 40.8-4.4-13.9-4.3-40.8-1.8-40.8zm-24.4 136.6c9.7-16.9 18-37 24.7-54.7 8.3 15.1 18.9 27.2 30.1 35.5-20.8 4.3-38.9 13.1-54.8 19.2zm131.6-5s-5 6-37.3-7.8c35.1-2.6 40.9 5.4 37.3 7.8z")
      svg(v-else-if="doc.name.endsWith('.ppt') || doc.name.endsWith('.pptx')" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 384 512")
        path(fill="currentColor" d="M369.9 97.9L286 14C277 5 264.8-.1 252.1-.1H48C21.5 0 0 21.5 0 48v416c0 26.5 21.5 48 48 48h288c26.5 0 48-21.5 48-48V131.9c0-12.7-5.1-25-14.1-34zM332.1 128H256V51.9l76.1 76.1zM48 464V48h160v104c0 13.3 10.7 24 24 24h104v288H48zm72-60V236c0-6.6 5.4-12 12-12h69.2c36.7 0 62.8 27 62.8 66.3 0 74.3-68.7 66.5-95.5 66.5V404c0 6.6-5.4 12-12 12H132c-6.6 0-12-5.4-12-12zm48.5-87.4h23c7.9 0 13.9-2.4 18.1-7.2 8.5-9.8 8.4-28.5.1-37.8-4.1-4.6-9.9-7-17.4-7h-23.9v52z")
      svg(v-else-if="doc.name.endsWith('.xls') || doc.name.endsWith('.xlsx')" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 384 512")
        path(fill="currentColor" d="M369.9 97.9L286 14C277 5 264.8-.1 252.1-.1H48C21.5 0 0 21.5 0 48v416c0 26.5 21.5 48 48 48h288c26.5 0 48-21.5 48-48V131.9c0-12.7-5.1-25-14.1-34zM332.1 128H256V51.9l76.1 76.1zM48 464V48h160v104c0 13.3 10.7 24 24 24h104v288H48zm212-240h-28.8c-4.4 0-8.4 2.4-10.5 6.3-18 33.1-22.2 42.4-28.6 57.7-13.9-29.1-6.9-17.3-28.6-57.7-2.1-3.9-6.2-6.3-10.6-6.3H124c-9.3 0-15 10-10.4 18l46.3 78-46.3 78c-4.7 8 1.1 18 10.4 18h28.9c4.4 0 8.4-2.4 10.5-6.3 21.7-40 23-45 28.6-57.7 14.9 30.2 5.9 15.9 28.6 57.7 2.1 3.9 6.2 6.3 10.6 6.3H260c9.3 0 15-10 10.4-18L224 320c.7-1.1 30.3-50.5 46.3-78 4.7-8-1.1-18-10.3-18z")
      svg(v-else-if="doc.name.endsWith('.jpg') || doc.name.endsWith('.jpeg')" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 384 512")
        path(fill="currentColor" d="M369.9 97.9L286 14C277 5 264.8-.1 252.1-.1H48C21.5 0 0 21.5 0 48v416c0 26.5 21.5 48 48 48h288c26.5 0 48-21.5 48-48V131.9c0-12.7-5.1-25-14.1-34zM332.1 128H256V51.9l76.1 76.1zM48 464V48h160v104c0 13.3 10.7 24 24 24h104v288H48zm32-48h224V288l-23.5-23.5c-4.7-4.7-12.3-4.7-17 0L176 352l-39.5-39.5c-4.7-4.7-12.3-4.7-17 0L80 352v64zm48-240c-26.5 0-48 21.5-48 48s21.5 48 48 48 48-21.5 48-48-21.5-48-48-48z")
      svg(v-else xmlns="http://www.w3.org/2000/svg" viewBox="0 0 384 512")
        path(fill="currentColor" d="M369.9 97.9L286 14C277 5 264.8-.1 252.1-.1H48C21.5 0 0 21.5 0 48v416c0 26.5 21.5 48 48 48h288c26.5 0 48-21.5 48-48V131.9c0-12.7-5.1-25-14.1-34zM332.1 128H256V51.9l76.1 76.1zM48 464V48h160v104c0 13.3 10.7 24 24 24h104v288H48z")
    b-link.files-name(v-if="doc.name.endsWith('.pdf')" :href="`/api/folders/${$route.params.id}/${doc.id}`" target="_blank") {{doc.name}}
    b-link.files-name(v-else :href="`/api/folders/${$route.params.id}/${doc.id}`" download) {{doc.name}}
  div.mt-3(v-if="folder.canEdit")
    b-btn(size="sm" variant="primary" @click="onAddFolder") Add Folder
    b-btn.ml-2(size="sm" variant="primary" @click="onAddDocument") Add File
  b-modal#files-edit-folder(:title="editFolder ? 'Edit Folder' : 'Add Folder'" @shown="onEditFolderShown" @ok="onEditFolderOK")
    form(@select.prevent="doEditFolder")
      b-form-group(label="Name" label-for="files-edit-folder-name" :state="editFolderNameError ? false : null" :invalid-feedback="editFolderNameError")
        b-input#files-edit-folder-name(ref="editFolderName" :state="editFolderNameError ? false : null" trim v-model="editFolderName")
      b-form-group(label="Group" label-for="files-edit-folder-group")
        b-select#files-edit-folder-group(:options="folder.allowedGroups" v-model="editFolderGroup" text-field="name" value-field="id")
      b-form-group(v-if="editFolder" label="Parent Folder" label-for="files-edit-folder-parent")
        b-select#files-edit-folder-parent(:options="allowedParents" v-model="editFolderParent" text-field="name" value-field="id")
    template(v-slot:modal-footer="{ok, cancel, hide}")
      b-btn.mr-5(v-if="editFolder" variant="danger" @click="onDeleteFolder(editFolder, hide)") Delete
      b-btn(@click="cancel()") Cancel
      b-btn(variant="primary" @click="ok()") OK
  b-modal#files-edit-doc(:title="editDocument ? 'Edit File' : 'Add File'" @shown="onEditDocumentShown" @ok="onEditDocumentOK")
    form(@select.prevent="doEditDocument")
      b-form-group(v-if="editDocument" label="Name" label-for="files-edit-doc-name" :state="editDocumentNameError ? false : null" :invalid-feedback="editDocumentNameError")
        b-input#files-edit-doc-name(ref="editDocumentName" :state="editDocumentNameError ? false : null" trim v-model="editDocumentName")
      b-form-group(v-else :state="editDocumentFilesError ? false : null" :invalid-feedback="editDocumentFilesError")
        b-form-file#files-edit-doc-files(multiple :state="editDocumentFilesError ? false : null" v-model="editDocumentFiles")
      b-form-group(v-if="editDocument" label="In Folder" label-for="files-edit-doc-folder")
        b-select#files-edit-doc-folder(:options="folder.allowedParents" v-model="editDocumentFolder" text-field="name" value-field="id")
    template(v-slot:modal-footer="{ok, cancel, hide}")
      b-btn.mr-5(v-if="editDocument" variant="danger" @click="onDeleteDocument(editDocument, hide)") Delete
      b-btn(@click="cancel()") Cancel
      b-btn(variant="primary" @click="ok()") OK
</template>

<script>
function extension(path) {
  const dot = path.lastIndexOf('.')
  return dot < 0 ? '' : path.substr(dot)
}

export default {
  data: () => ({
    folder: null,
    editDocument: null,
    editDocumentName: null,
    editDocumentNameError: null,
    editDocumentFiles: null,
    editDocumentFilesError: null,
    editDocumentFolder: null,
    editFolder: null,
    editFolderName: null,
    editFolderNameError: null,
    editFolderGroup: null,
    editFolderParent: null,
    hoverDocument: null,
    hoverFolder: null,
  }),
  computed: {
    allowedParents() {
      let ignoreAbove = null
      return this.folder.allowedParents.filter(p => {
        if (p.id === this.editFolder.id) {
          ignoreAbove = p.indent
          return false
        }
        if (ignoreAbove !== null && p.indent > ignoreAbove) {
          return false
        }
        ignoreAbove = null
        return true
      })
    }
  },
  created() {
    this.loadFolder()
  },
  watch: {
    $route: 'loadFolder',
  },
  methods: {
    async doEditDocument() {
      const body = new FormData
      const docID = this.editDocument ? this.editDocument.id : 'NEW'
      if (this.editDocument) {
        if (!this.editDocumentName) {
          this.editDocumentNameError = 'The file name is required.'
          return false
        }
        if (!extension(this.editDocumentName)) {
          this.editDoumentNameError = 'The file must have an extension indicating its type.'
          return false
        }
        if (extension(this.editDocumentName) !== extension(this.editDocument.name)) {
          this.editDocumentNameError = `The new file name must end with "${extension(this.editDocument.name)}".`
          return false
        }
        if (this.folder.documents.find(f => f !== this.editDocument && f.name === this.editDocumentName)) {
          this.editDocumentNameError = 'This file name is already in use.'
          return false
        }
        body.append('name', this.editDocumentName)
        if (this.editDocument) body.append('folder', this.editDocumentFolder)
      } else {
        if (!this.editDocumentFiles || !this.editDocumentFiles.length) {
          this.editDocumentFilesError = 'Select the file(s) to be added.'
          return false
        }
        const replace = []
        this.editDocumentFiles.forEach(nf => {
          if (this.folder.documents.find(d => d.name === nf.name))
            replace.push(nf.name)
        })
        if (replace.length === 1) {
          if (!await this.$bvModal.msgBoxConfirm(
            `This will replace the existing file named "${replace[0]}".  Are you sure?`,
            {
              title: 'Replace File', headerBgVariant: 'warning', headerTextVariant: 'white',
              okTitle: 'Replace', okVariant: 'warning', cancelTitle: 'Keep',
            })) return false
        } else if (replace.length > 1) {
          if (!await this.$bvModal.msgBoxConfirm(
            `This will replace ${replace.length} existing files with the same names.  Are you sure?`,
            {
              title: 'Replace Files', headerBgVariant: 'warning', headerTextVariant: 'white',
              okTitle: 'Replace', okVariant: 'warning', cancelTitle: 'Keep',
            })) return false
        }
        this.editDocumentFiles.forEach(nf => { body.append('file', nf) })
      }
      this.folder = (await this.$axios.post(`/api/folders/${this.$route.params.id}/${docID}`, body)).data
      return true
    },
    doEditFolder() {
      if (!this.editFolderName) {
        this.editFolderNameError = 'The folder name is required.'
        return false
      }
      if (this.folder.children.find(f => f !== this.editFolder && f.name === this.editFolderName)) {
        this.editFolderNameError = 'This folder name is already in use.'
        return false
      }
      const body = new FormData
      body.append('name', this.editFolderName)
      body.append('group', this.editFolderGroup)
      if (this.editFolder) body.append('parent', this.editFolderParent)
      if (this.editFolder)
        this.$axios.put(`/api/folders/${this.editFolder.id}`, body).then(r => { this.folder = r.data })
      else
        this.$axios.post(`/api/folders/${this.$route.params.id}`, body).then(r => { this.folder = r.data })
      return true
    },
    async loadFolder() {
      this.folder = (await this.$axios.get(`/api/folders/${this.$route.params.id}`)).data
      this.$store.commit('setPage', { title: this.folder.name || 'Files' })
    },
    onAddDocument() {
      this.editDocument = null
      this.editDocumentName = null
      this.editDocumentNameError = null
      this.editDocumentFilesError = null
      this.editDocumentFolder = this.folder.id
      this.$bvModal.show('files-edit-doc')
    },
    onAddFolder() {
      this.editFolder = null
      this.editFolderName = null
      this.editFolderNameError = null
      this.editFolderGroup = this.folder.group
      this.$bvModal.show('files-edit-folder')
    },
    onDragOver(evt) { evt.preventDefault() },
    async onDrop(evt) {
      if (!this.folder.canEdit) return
      evt.preventDefault()
      const body = new FormData
      const replace = []
      for (let i = 0; i < evt.dataTransfer.files.length; i++) {
        const nf = evt.dataTransfer.files[i]
        if (this.folder.documents.find(d => d.name === nf.name))
          replace.push(nf.name)
      }
      if (replace.length === 1) {
        if (!await this.$bvModal.msgBoxConfirm(
          `This will replace the existing file named "${replace[0]}".  Are you sure?`,
          {
            title: 'Replace File', headerBgVariant: 'warning', headerTextVariant: 'white',
            okTitle: 'Replace', okVariant: 'warning', cancelTitle: 'Keep',
          })) return false
      } else if (replace.length > 1) {
        if (!await this.$bvModal.msgBoxConfirm(
          `This will replace ${replace.length} existing files with the same names.  Are you sure?`,
          {
            title: 'Replace Files', headerBgVariant: 'warning', headerTextVariant: 'white',
            okTitle: 'Replace', okVariant: 'warning', cancelTitle: 'Keep',
          })) return false
      }
      for (let i = 0; i < evt.dataTransfer.files.length; i++) {
        body.append('file', evt.dataTransfer.files[i])
      }
      this.folder = (await this.$axios.post(`/api/folders/${this.$route.params.id}/NEW`, body)).data
    },
    onEditDocument(doc) {
      if (!this.folder.canEdit) return
      this.editDocument = doc
      this.editDocumentName = doc.name
      this.editDocumentNameError = null
      this.editDocumentFolder = this.folder.id || 0
      this.$bvModal.show('files-edit-doc')
    },
    onEditFolder(folder) {
      if (!this.folder.canEdit) return
      this.editFolder = folder
      this.editFolderName = folder.name
      this.editFolderNameError = null
      this.editFolderGroup = folder.group
      this.editFolderParent = this.folder.id || 0
      this.$bvModal.show('files-edit-folder')
    },
    onEditDocumentShown() {
      if (this.$refs.editDocumentName) this.$refs.editDocumentName.focus()
    },
    onEditFolderShown() {
      this.$refs.editFolderName.focus()
    },
    async onEditDocumentOK(evt) {
      evt.preventDefault()
      if (await this.doEditDocument())
        this.$bvModal.hide('files-edit-doc')
    },
    onEditFolderOK(evt) {
      if (!this.doEditFolder())
        evt.preventDefault()
    },
    async onDeleteDocument(doc, hide) {
      hide()
      if (!await this.$bvModal.msgBoxConfirm(
        `Are you sure you want to delete the file "${doc.name}"?  It cannot be restored.`,
        {
          title: 'Delete File', headerBgVariant: 'danger', headerTextVariant: 'white',
          okTitle: 'Delete', okVariant: 'danger', cancelTitle: 'Keep',
        }))
        return
      this.folder = (await this.$axios.delete(`/api/folders/${this.$route.params.id}/${doc.id}`)).data
    },
    async onDeleteFolder(child, hide) {
      hide()
      if (!await this.$bvModal.msgBoxConfirm(
        `Are you sure you want to delete the folder "${child.name}", along with all of its files and sub-folders?  It cannot be restored.`,
        {
          title: 'Delete Folder', headerBgVariant: 'danger', headerTextVariant: 'white',
          okTitle: 'Delete', okVariant: 'danger', cancelTitle: 'Keep',
        }))
        return
      this.folder = (await this.$axios.delete(`/api/folders/${child.id}`)).data
    },
    onHoverDocument(id) { this.hoverDocument = id },
    onHoverFolder(id) { this.hoverFolder = id },
  },
}
</script>

<style lang="stylus">
#files
  min-height calc(100vh - 40px)
  .mouse &
    margin 1.5rem 0.75rem
.files-line
  display flex
  align-items center
  .touch &
    padding 0.25rem 0.5rem
    min-height 40px
    border-bottom 1px solid #ccc
.files-icon
  display flex
  flex none
  justify-content center
  align-items center
  margin-right 0.5rem
  padding 0
  width 1rem
  height 1rem
  border none
  background-color white
  color black
  &:hover, &:active, &:focus
    background-color white !important
    color black !important
  .touch &
    width 1.5rem
    height 1.5rem
  svg
    width 100%
    height 100%
.files-name
  flex 1 1 auto
  overflow hidden
  min-width 0
  text-overflow ellipsis
  white-space nowrap
  .touch &
    white-space normal
    line-height 1.2
</style>
