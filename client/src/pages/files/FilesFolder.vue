<!--
FilesFolder
-->

<template lang="pug">
.files-folder(
  :style='[indentStyle, dragOverStyle]',
  :draggable='folder.canEdit && folder.url !== ""',
  @dragstart='onDragStart',
  @dragenter='onDragEnter',
  @dragover='onDragOver',
  @dragleave='onDragLeave',
  @dragend='onDragEnd',
  @drop='onDrop',
  @mouseover='onMouseOver',
  @mouseout='onMouseOut'
)
  SIcon.files-folder-icon(:icon='icon', :style='iconStyle', @click='onEdit')
  router-link.files-folder-link(:to='url', draggable='false', v-text='folderName')
  MessageBox(ref='errorModal', title='Error', cancelLabel='', variant='danger')
    | {{ errorText }}
  MessageBox(
    ref='replaceModal',
    :title='replaceTitle',
    cancelLabel='Keep',
    okLabel='Replace',
    variant='danger'
  )
    | {{ replaceText }}
  FilesEditFolder(ref='editFolderModal')
</template>

<script lang="ts">
import { computed, defineComponent, inject, PropType, ref, Ref } from 'vue'
import axios from '../../plugins/axios'
import { LoginData } from '../../plugins/login'
import { MessageBox, SIcon } from '../../base'
import type { GetFolderFolder } from '../Files.vue'
import FilesEditFolder from './FilesEditFolder.vue'

interface PostFolderDrop {
  error?: string
  replace?: Array<string>
  urlName?: string
}
interface folderDragData {
  name: string
  url: string
  visibility: string
}

export default defineComponent({
  components: { FilesEditFolder, MessageBox, SIcon },
  props: {
    folder: { type: Object as PropType<GetFolderFolder>, required: true },
    indent: { type: Number, required: true },
  },
  emits: ['progress', 'reload', 'showTrash'],
  setup(props, { emit }) {
    const me = inject<Ref<LoginData>>('me')!
    const folderName = computed(() =>
      props.folder.url ? props.folder.name : me.value ? 'Files' : 'Sunnyvale SERV'
    )
    const icon = computed(() => {
      if (hovering.value && props.folder.canEdit && props.folder.url !== '') return 'edit'
      if (props.folder.url || me.value) return 'folder'
      return 'home'
    })
    const iconStyle = computed(() =>
      hovering.value && props.folder.canEdit && props.folder.url !== '' ? { color: 'blue' } : null
    )
    const indentStyle = computed(() => ({ marginLeft: `${props.indent * 1.5}em` }))
    const url = computed(() =>
      props.folder.url ? `/files${props.folder.url}` : me.value ? '/files' : '/'
    )
    const dragOverStyle = computed(() => (draggingOver.value ? { backgroundColor: '#ccc' } : ''))

    // Drop support.
    const draggingOver = ref(0)
    const replace = ref([] as Array<string>)
    const errorText = ref('')
    const errorModal = ref(null as any)
    function dropSupported(evt: DragEvent): boolean {
      if (!props.folder || !props.folder.canEdit) return false
      if (evt.dataTransfer!.types.includes('x-serv-folder')) return true
      if (props.folder.url === '') return false // no documents at top level
      if (evt.dataTransfer!.types.includes('x-serv-document')) return true
      if (evt.dataTransfer!.types.includes('Files')) return true
      if (evt.dataTransfer!.types.includes('text/uri-list')) return true
      return false
    }
    function onDragEnter(evt: DragEvent) {
      if (!dropSupported(evt)) return
      evt.preventDefault()
      draggingOver.value++
    }
    function onDragOver(evt: DragEvent) {
      if (!dropSupported(evt)) return
      evt.preventDefault()
    }
    function onDragLeave(evt: DragEvent) {
      if (!dropSupported(evt)) return
      evt.preventDefault()
      if (draggingOver.value) draggingOver.value--
    }
    async function onDrop(evt: DragEvent) {
      if (!dropSupported(evt)) return
      evt.preventDefault()
      draggingOver.value = 0
      if (evt.dataTransfer!.types.includes('Files')) onDropFiles(evt.dataTransfer!)
      else if (evt.dataTransfer!.types.includes('x-serv-folder')) onDropFolder(evt.dataTransfer!)
      else if (evt.dataTransfer!.types.includes('x-serv-document'))
        onDropDocument(evt.dataTransfer!)
      else if (evt.dataTransfer!.types.includes('text/uri-list')) onDropLink(evt.dataTransfer!)
    }

    async function onDropFiles(dt: DataTransfer) {
      // We have to build both request bodies up front because dt will be lost
      // as soon as this function returns, i.e., at the first async call.
      const checkBody = new FormData()
      const uploadBody = new FormData()
      checkBody.append('replace', 'true')
      for (let i = 0; i < dt.files.length; i++) {
        const f = dt.files[i]
        if (!f.type) {
          errorText.value = 'Only regular files can be added to the folder.'
          await errorModal.value.show()
          return
        }
        checkBody.append('name', f.name)
        uploadBody.append('file', f)
      }
      try {
        await axios.post(`/api/folders${props.folder.url}?op=checkNames`, checkBody)
        // No conflicts, so:
        await axios.post(`/api/folders${props.folder.url}?op=newFiles`, uploadBody, {
          onUploadProgress,
        })
        emit('reload')
      } catch (err) {
        if (!err.response || err.response.status !== 409) throw err
        if (!err.response.headers['x-can-replace']) {
          errorText.value = err.response.data
          await errorModal.value.show()
          return
        }
        replaceTitle.value = 'Replace Files'
        replaceText.value = err.response.data
        if (!(await replaceModal.value.show())) return
        await axios.post(`/api/folders${props.folder.url}?op=newFiles`, uploadBody, {
          onUploadProgress,
        })
        emit('reload')
      } finally {
        emit('progress', 0)
      }
    }
    function onUploadProgress(evt: ProgressEvent) {
      if (evt.lengthComputable && evt.total) emit('progress', evt.loaded / evt.total)
    }

    // Receive a dropped SERV folder.
    async function onDropFolder(dt: DataTransfer) {
      const fdd: { name: string; url: string } = JSON.parse(dt.getData('x-serv-folder'))
      if (fdd.url === props.folder.url) return
      try {
        const newURL = (
          await axios.post<string>(`/api/folders${fdd.url}?op=move`, null, {
            params: { parent: props.folder.url },
          })
        ).data
        if (newURL != props.folder.url)
          emit('reload', function (u: string): string {
            if (u === props.folder.url) return newURL
            if (u.startsWith(props.folder.url) + '/')
              return newURL + u.substr(props.folder.url.length)
            return u
          })
        else emit('reload')
      } catch (err) {
        if (err.response && err.response.status === 409) {
          errorText.value = err.response.data
          await errorModal.value.show()
          return
        }
        throw err
      }
    }

    // Receive a dropped SERV document.
    async function onDropDocument(dt: DataTransfer) {
      const ddd: { name: string; url: string; type: string } = JSON.parse(
        dt.getData('x-serv-document')
      )
      if (ddd.url.replace(/\/[^/]*$/, '') === props.folder.url) return
      try {
        await axios.post<string>(`/api/document${ddd.url}?op=move`, null, {
          params: { parent: props.folder.url },
        })
        emit('reload')
      } catch (err) {
        if (!err.response || err.response.status !== 409) throw err
        if (err.response.headers['x-can-replace']) {
          replaceTitle.value = ddd.type === 'link' ? 'Replace Link' : 'Replace File'
          replaceText.value = err.response.data
          if (!(await replaceModal.value.show())) return
        } else {
          errorText.value = err.response.data
          await errorModal.value.show()
          return
        }
        await axios.post(`/api/document${ddd.url}?op=move&replace=true`, null, {
          params: { parent: props.folder.url },
        })
        emit('reload')
      }
    }

    // Receive a dropped URL.
    async function onDropLink(dt: DataTransfer) {
      const url = dt.getData('text/uri-list')!.replace(/\r\n.*/, '')
      if (!url.startsWith('http://') && !url.startsWith('https://')) {
        errorText.value = 'Only links starting with http:// or https:// are supported.'
        await errorModal.value.show()
        return
      }
      try {
        await axios.post(`/api/folders${props.folder.url}?op=newLink`, null, { params: { url } })
        emit('reload')
      } catch (err) {
        if (!err.response || err.response.status !== 409) throw err
        if (err.response.headers['x-can-replace']) {
          replaceTitle.value = 'Replace Link'
          replaceText.value = err.response.data
          if (!(await replaceModal.value.show())) return
        } else {
          errorText.value = err.response.data
          await errorModal.value.show()
          return
        }
        await axios.post(`/api/folders${props.folder.url}?op=newLink&replace=true`, null, {
          params: { url },
        })
        emit('reload')
      }
    }

    // Replacements.
    const replaceModal = ref(null as any)
    const replaceTitle = ref('')
    const replaceText = ref('')

    // Drag support.
    function onDragStart(evt: DragEvent) {
      if (!props.folder.canEdit || props.folder.url === '') return
      const fdd = JSON.stringify({ name: props.folder.name, url: props.folder.url })
      evt.dataTransfer!.setData('x-serv-folder', fdd)
      evt.dataTransfer!.setData(
        'text/uri-list',
        `https://sunnyvaleserv.org/${url.value}\r\n# ${props.folder.name}`
      )
      evt.dataTransfer!.setData('text/plain', `https://sunnyvaleserv.org/${url.value}`)
      emit('showTrash', true)
    }
    function onDragEnd() {
      emit('showTrash', false)
    }

    // Edit button support.
    const hovering = ref(false)
    const editFolderModal = ref(null as any)
    function onMouseOver() {
      hovering.value = true
    }
    function onMouseOut() {
      hovering.value = false
    }
    async function onEdit() {
      if (!props.folder.canEdit || props.folder.url === '') return
      const newURL: string = await editFolderModal.value.showEdit(props.folder.url)
      if (newURL && newURL != props.folder.url)
        emit('reload', function (u: string): string {
          if (u === props.folder.url) return newURL
          if (u.startsWith(props.folder.url) + '/')
            return newURL + u.substr(props.folder.url.length)
          return u
        })
      else emit('reload')
    }

    return {
      dragOverStyle,
      editFolderModal,
      errorModal,
      errorText,
      folderName,
      icon,
      iconStyle,
      indentStyle,
      onDragEnd,
      onDragEnter,
      onDragLeave,
      onDragOver,
      onDragStart,
      onDrop,
      onEdit,
      onMouseOver,
      onMouseOut,
      replaceModal,
      replaceText,
      replaceTitle,
      url,
    }
  },
})
</script>

<style lang="postcss">
.files-folder {
  display: flex;
  align-items: center;
  min-height: 0;
  min-width: 0;
}
.files-folder-icon {
  flex: none;
  height: 1rem;
  width: 1.5rem;
  margin-right: 0.25rem;
  /* This margin optimizes for the case that the next line down will be
   * indented, and will have a folder icon.  The folder icon is 1rem wide,
   * centered in a 1.5rem svg, so it has 0.25rem of whitespace on the left.
   * By adding this margin, the text after this line's icon will line up with
   * the left side of the folder icon on the line below.
   */
}
.files-folder-link {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
