<!--
Files displays the Files browser.
-->

<template lang="pug">
div(v-if='!folder')
  SSpinner
#files(
  v-else,
  :class='dragging ? "dragging" : null',
  @drop='onDrop',
  @dragover='onDragOver',
  @dragleave='onDragLeave'
)
  SProgress#files-progress(v-if='uploadProgress !== null', :value='uploadProgress')
  .files-line(v-if='folder.id')
    .files-icon: SIcon(icon='up')
    router-link.files-name(:to='`/files/${folder.parent ? folder.parent.id : 0}`') {{ folder.parent ? folder.parent.name : "Files" }}
  FilesFolder(
    v-for='child in folder.children',
    :key='`f${child.id}`',
    :folder='child',
    :canEdit='folder.canEdit',
    @edit='onEditFolder(child)'
  )
  FilesDocument(
    v-for='doc in folder.documents',
    :key='`d${doc.id}`',
    :folderID='folder.id',
    :doc='doc',
    :canEdit='folder.canEdit',
    @edit='onEditDocument(doc)'
  )
  #files-add(v-if='folder.canEdit || folder.canAdd')
    SButton(v-if='folder.canEdit', variant='primary', @click='onAddFolder') Add Folder
    SButton(variant='primary', @click='onAddDocument') Add File
  MessageBox(
    ref='replaceConfirmRef',
    :title='replace.length > 1 ? "Replace Files" : "Replace File"',
    cancelLabel='Keep',
    okLabel='Replace',
    variant='warning'
  ) This will replace {{ replace.length > 1 ? `${replace.length} existing files with the same names` : `the existing file named "${replace[0]}"` }}. Are you sure?
  FilesAddDocument(ref='addDocumentRef', :parentID='folder.id', :siblings='folder.documents')
  FilesEditDocument(
    ref='editDocumentRef',
    :doc='editDocument',
    :allowedParents='folder.allowedParents',
    :parentID='folder.id',
    :siblings='folder.documents'
  )
  FilesEditFolder(
    ref='editFolderRef',
    :allowedGroups='folder.allowedGroups',
    :allowedParents='folder.allowedParents',
    :editFolder='editFolder',
    :parentID='folder.id',
    :siblings='folder.children'
  )
</template>

<script lang="ts">
import { defineComponent, ref, watch, watchEffect } from 'vue'
import { useRoute } from 'vue-router'
import axios from '../plugins/axios'
import setPage from '../plugins/page'
import { MessageBox, SButton, SIcon, SProgress, SSpinner } from '../base'
import FilesAddDocument from './files/FilesAddDocument.vue'
import FilesDocument from './files/FilesDocument.vue'
import FilesEditDocument from './files/FilesEditDocument.vue'
import FilesEditFolder from './files/FilesEditFolder.vue'
import FilesFolder from './files/FilesFolder.vue'

type GetFolderParent = {
  id: number
  name: string
  url: string
}
export type GetFolderDocument = {
  id: number
  name: string
  needsApproval?: true
}
export type GetFolderChild = {
  id: number
  name: string
  url: string
  group: number
  approvals?: number
}
interface GetFolderBase {
  id: number
  parent?: GetFolderParent
  group: number
  name: string
  url: string
  documents: Array<GetFolderDocument>
  children?: Array<GetFolderChild>
  canEdit: boolean
  canAdd: boolean
}
export type GetFolderEditAllowedGroup = {
  id: number
  name: string
}
export type GetFolderEditAllowedParent = {
  id: number
  name: string
  indent: number
  disabled?: true
}
export interface GetFolderEdit extends GetFolderBase {
  canEdit: true
  allowedGroups: Array<GetFolderEditAllowedGroup>
  allowedParents: Array<GetFolderEditAllowedParent>
}
type GetFolder = GetFolderBase | GetFolderEdit

export default defineComponent({
  components: {
    FilesAddDocument,
    FilesDocument,
    FilesEditDocument,
    FilesEditFolder,
    FilesFolder,
    MessageBox,
    SButton,
    SIcon,
    SProgress,
    SSpinner,
  },
  setup() {
    setPage({ title: 'Files' })

    // First, load the folder data.
    const route = useRoute()
    const folder = ref(null as null | GetFolder)
    watch(folder, () => {
      if (!folder.value!.children) folder.value!.children = []
      if (!folder.value!.documents) folder.value!.documents = []
      setPage({ title: folder.value!.name || 'Files' })
    })
    watchEffect(async () => {
      folder.value = (await axios.get<GetFolder>(`/api/folders/${route.params.id}`)).data
    })

    // Handle drag and drop, and the message box for replace confirmation.
    const dragging = ref(false)
    const uploadProgress = ref(null as null | number)
    const replace = ref([] as Array<string>)
    const replaceConfirmRef = ref(null as any)
    function onUploadProgress(evt: ProgressEvent) {
      if (evt.lengthComputable && evt.total) uploadProgress.value = evt.loaded / evt.total
    }
    function onDragOver(evt: DragEvent) {
      evt.preventDefault()
      dragging.value = true
    }
    function onDragLeave() {
      dragging.value = false
    }
    async function onDrop(evt: DragEvent) {
      if (!folder.value!.canAdd && !folder.value!.canEdit) return
      evt.preventDefault()
      dragging.value = false
      replace.value = []
      const droppedFiles = evt.dataTransfer!.files
      for (let i = 0; i < droppedFiles.length; i++) {
        const nf = droppedFiles[i]
        if (folder.value!.documents.find((d) => d.name === nf.name)) replace.value.push(nf.name)
      }
      if (replace.value.length > 0) if (!(await replaceConfirmRef.value.show())) return
      const body = new FormData()
      for (let i = 0; i < droppedFiles.length; i++) {
        body.append('file', droppedFiles[i])
      }
      folder.value = (
        await axios.post<GetFolder>(`/api/folders/${route.params.id}/NEW`, body, {
          onUploadProgress,
        })
      ).data
      uploadProgress.value = null
    }

    // Editing of folders.
    const editFolder = ref(null as null | GetFolderChild)
    const editFolderRef = ref(null as any)
    async function onEditFolder(cf: GetFolderChild) {
      if (!folder.value!.canEdit) return
      editFolder.value = cf
      const updated: GetFolderEdit = await editFolderRef.value.show()
      if (updated) folder.value = updated
      editFolder.value = null
    }
    async function onAddFolder() {
      if (!folder.value!.canEdit) return
      editFolder.value = null
      const updated: GetFolderEdit = await editFolderRef.value.show()
      if (updated) folder.value = updated
    }

    // Editing of documents.
    const editDocument = ref(null as null | GetFolderDocument)
    const editDocumentRef = ref(null as any)
    async function onEditDocument(doc: GetFolderDocument) {
      if (!folder.value!.canEdit) return
      editDocument.value = doc
      const updated: GetFolderEdit = await editDocumentRef.value.show()
      if (updated) folder.value = updated
      editDocument.value = null
    }
    const addDocumentRef = ref(null as any)
    async function onAddDocument() {
      if (!folder.value!.canEdit) return
      const updated: GetFolderEdit = await addDocumentRef.value.show()
      if (updated) folder.value = updated
    }

    return {
      addDocumentRef,
      dragging,
      editDocument,
      editDocumentRef,
      editFolder,
      editFolderRef,
      folder,
      onAddDocument,
      onAddFolder,
      onDragLeave,
      onDragOver,
      onDrop,
      onEditDocument,
      onEditFolder,
      replace,
      replaceConfirmRef,
      uploadProgress,
    }
  },
  /*
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
        if (this.folder.documents.find(f => f !== this.editDocument && f.name === this.editDocumentName && !this.editDocument.needsApproval)) {
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
      this.folder = (await this.$axios.post(`/api/folders/${this.$route.params.id}/${docID}`, body, {
        onUploadProgress: this.onUploadProgress,
      })).data
      this.uploadProgress = null
      this.$bvModal.hide('files-edit-doc')
      return true
    },
    onAddDocument() {
      this.editDocument = null
      this.editDocumentName = null
      this.editDocumentNameError = null
      this.editDocumentFilesError = null
      this.editDocumentFolder = this.folder.id
      this.$bvModal.show('files-edit-doc')
    },
    onEditDocument(doc) {
      if (!this.folder.canEdit) return
      this.editDocument = doc
      this.editDocumentName = doc.name
      this.editDocumentNameError = null
      this.editDocumentFolder = this.folder.id || 0
      this.$bvModal.show('files-edit-doc')
    },
  */
})
</script>

<style lang="postcss">
#files {
  position: relative;
  height: 100%;
  .mouse & {
    padding: 1.5rem 0.75rem;
  }
  &.dragging {
    box-shadow: inset 0 0 0 0.25rem #006600;
  }
}
#files-progress {
  position: absolute;
  top: 0;
  right: 0;
  left: 0;
}
.files-line {
  display: flex;
  align-items: center;
  .touch & {
    padding: 0.25rem 0.5rem;
    min-height: 40px;
    border-bottom: 1px solid #ccc;
  }
}
.files-icon {
  display: flex;
  flex: none;
  justify-content: center;
  align-items: center;
  margin-right: 0.5rem;
  padding: 0;
  width: 1rem;
  height: 1rem;
  border: none;
  background-color: white;
  color: black;
  &:hover,
  &:active,
  &:focus {
    background-color: white !important;
    color: black !important;
  }
  .touch & {
    width: 1.5rem;
    height: 1.5rem;
  }
  svg {
    width: 100%;
    height: 100%;
  }
}
.files-name {
  flex: 1 1 auto;
  overflow: hidden;
  min-width: 0;
  text-overflow: ellipsis;
  white-space: nowrap;
  .touch & {
    white-space: normal;
    line-height: 1.2;
  }
}
.files-pending {
  color: red;
}
#files-add {
  margin-top: 1rem;
  & .sbtn {
    margin-right: 0.5rem;
  }
}
</style>
