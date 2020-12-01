<!--
Files displays a folder's contents.  It is used both in public and logged-in
contexts.
-->

<template lang="pug">
NotFound(v-if='notFound')
#files(v-else-if='!folder.ancestors')
  SSpinner
#files(v-else)
  SProgress.files-progress(v-if='progress', :value='progress')
  transition-group(name='files')
    FilesFolder.files-item(
      v-for='(p, i) in folder.ancestors',
      :folder='p',
      :key='p.url',
      :indent='i * ancestorsIndent',
      @progress='onProgress',
      @reload='onReload',
      @showTrash='onShowTrash'
    )
    FilesFolder.files-item(
      v-for='c in folder.children',
      :folder='c',
      :key='c.url',
      :indent='contentsIndent',
      @progress='onProgress',
      @reload='onReload',
      @showTrash='onShowTrash'
    )
    FilesDocument.files-item(
      v-for='d in folder.documents',
      :folder='folder.ancestors[folder.ancestors.length - 1]',
      :doc='d',
      :key='d.name',
      :indent='contentsIndent',
      @progress='onProgress',
      @reload='onReload',
      @showTrash='onShowTrash'
    )
    div(
      v-if='!folder.children.length && !folder.documents.length',
      key='.empty',
      :style='{ marginLeft: `${1.5 * contentsIndent + 0.25}rem` }'
    ) (folder is empty)
  FilesTrash(v-if='showTrash > 0', @reload='onReload', @showTrash='onShowTrash')
  #files-buttons(v-else-if='canEdit')
    SButton(variant='primary', small, @click='onAddFolder') Add Folder
    SButton(variant='primary', small, @click='onAddFiles') Add File
    SButton(variant='primary', small, @click='onAddURL') Add Link
  FilesEditFolder(v-if='canEdit', ref='editFolderModal')
  FilesEditLink(v-if='canEdit', ref='editLinkModal')
  FilesAddFiles(v-if='canEdit', ref='addFilesModal')
</template>

<script lang="ts">
import { computed, defineComponent, inject, Ref, ref, watchEffect } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from '../plugins/axios'
import { LoginData } from '../plugins/login'
import setPage from '../plugins/page'
import { Size } from '../plugins/size'
import { SButton, SProgress, SSpinner } from '../base'
import NotFound from './NotFound.vue'
import FilesAddFiles from './files/FilesAddFiles.vue'
import FilesDocument from './files/FilesDocument.vue'
import FilesEditFolder from './files/FilesEditFolder.vue'
import FilesEditLink from './files/FilesEditLink.vue'
import FilesFolder from './files/FilesFolder.vue'
import FilesTrash from './files/FilesTrash.vue'

export interface GetFolderFolder {
  name: string
  url: string
  canEdit: boolean
}
export interface GetFolderDocument {
  name: string
  url: string
  newtab: boolean
}
interface GetFolder {
  ancestors: Array<GetFolderFolder>
  children: Array<GetFolderFolder>
  documents: Array<GetFolderDocument>
}

export default defineComponent({
  components: {
    FilesAddFiles,
    FilesDocument,
    FilesEditFolder,
    FilesEditLink,
    FilesFolder,
    FilesTrash,
    NotFound,
    SButton,
    SProgress,
    SSpinner,
  },
  setup() {
    const route = useRoute()
    const router = useRouter()
    const me = inject<Ref<LoginData>>('me')!
    const size = inject<Size>('containerSize')!
    setPage({ title: me.value ? 'Files' : '' })

    // Load the requested folder.
    const folder = ref({} as GetFolder)
    const notFound = ref(false)
    const canEdit = ref(false)
    async function loadFolder() {
      try {
        const url = route.params.path
          ? `/api/folders/${route.params.path}?op=browse`
          : '/api/folders?op=browse'
        const resp = await axios.get<GetFolder>(url)
        notFound.value = false
        folder.value = resp.data
        canEdit.value = folder.value.ancestors[folder.value.ancestors.length - 1].canEdit
      } catch (err) {
        if (err.response && err.response.status === 404) notFound.value = true
        else throw err
      }
    }
    watchEffect(loadFolder)

    // Reload the folder when requested by a sub-component.  The optional
    // parameter is a URL mapper function.  If provided, it tells how to map the
    // URL we've been viewing to the correct URL reflecting whatever change was
    // just made.
    function onReload(urlMapper?: (u: string) => string) {
      if (urlMapper) {
        const oldURL = route.params.path ? `/${route.params.path}` : ''
        let newURL = urlMapper(oldURL)
        if (newURL !== oldURL && newURL !== '') router.replace(`/files${newURL}`)
        else if (newURL !== oldURL) router.replace('/files')
        else loadFolder()
      } else loadFolder()
    }

    // The indentation style depends on the size of the window.
    const ancestorsIndent = computed(() => (size.w >= 24 ? 1 : 0))
    const contentsIndent = computed(() =>
      size.w >= 24 && folder.value.ancestors ? folder.value.ancestors.length : 1
    )

    // Showing or hiding the trash can.  It is shown while an item is being
    // dragged, and also after an item has been dropped on the trash can and the
    // deletion is still in progress.
    const showTrash = ref(0)
    function onShowTrash(flag: boolean) {
      if (flag) showTrash.value++
      else showTrash.value--
    }

    // Adding folders, files, and links.
    const editFolderModal = ref(null as any)
    const editLinkModal = ref(null as any)
    const addFilesModal = ref(null as any)
    async function onAddFolder() {
      const parent = folder.value.ancestors[folder.value.ancestors.length - 1].url
      if (await editFolderModal.value.showAdd(parent)) loadFolder()
    }
    async function onAddFiles() {
      const parent = folder.value.ancestors[folder.value.ancestors.length - 1].url
      if (await addFilesModal.value.show(parent)) loadFolder()
    }
    async function onAddURL() {
      const parent = folder.value.ancestors[folder.value.ancestors.length - 1].url
      if (await editLinkModal.value.showAdd(parent)) loadFolder()
    }

    // Displaying upload progress.
    const progress = ref(0)
    function onProgress(p: number) {
      progress.value = p
    }

    return {
      addFilesModal,
      ancestorsIndent,
      canEdit,
      contentsIndent,
      editFolderModal,
      editLinkModal,
      folder,
      onReload,
      notFound,
      onAddFiles,
      onAddFolder,
      onAddURL,
      onProgress,
      onShowTrash,
      progress,
      showTrash,
    }
  },
})
</script>

<style lang="postcss">
#files {
  margin: 1.5rem 0.75rem;
  display: grid;
  grid: auto-flow min-content / fit-content(100%);
  position: relative;
}
.files-progress {
  position: absolute;
  top: -1.5rem;
  left: -0.75rem;
  right: -0.75rem;
}
.files-item {
  transition: opacity 0.5s ease, height 0.5s ease, margin-left 0.5s ease;
  height: 1.5rem;
}
.files-enter-from,
.files-leave-to {
  opacity: 0;
  height: 0;
}
#files-buttons {
  margin-top: 1.5rem;
  & .sbtn {
    margin-right: 0.5rem;
  }
}
</style>
