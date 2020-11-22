<!--
ListEdit displays the page for creating and editing lists.
-->

<template lang="pug">
#list
  SSpinner(v-if='loading')
  template(v-else)
    SForm(:submitLabel='submitLabel', :disabled='submitting', @submit='onSubmit')
      SFRadioGroup#list-type(label='Type', inline, :options='typeOptions', v-model='list.type')
      SFInput#list-name(
        :label='nameLabel',
        :help='nameHelp',
        trim,
        v-model='list.name',
        :errorFn='nameError'
      )
      template(v-if='list.id', #extraButtons)
        SButton(@click='onDelete', variant='danger', :disabled='submitting') Delete List
      MessageBox(
        ref='deleteModal',
        title='Delete List',
        cancelLabel='Keep',
        okLabel='Delete',
        variant='danger'
      )
        | Are you sure you want to delete this list? All associated data,
        | including role associations, manual subscribes, and unsubscribes will
        | be permanently lost.
</template>

<script lang="ts">
import { computed, defineComponent, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from '../../plugins/axios'
import setPage from '../../plugins/page'
import { MessageBox, SButton, SForm, SFInput, SFRadioGroup, SSpinner } from '../../base'

interface GetList {
  id: number
  type: string
  name: string
}
interface PostList {
  duplicateName?: true
}

const typeOptions = [
  { value: 'email', label: 'Email' },
  { value: 'sms', label: 'SMS' },
]

export default defineComponent({
  components: { MessageBox, SButton, SForm, SFInput, SFRadioGroup, SSpinner },
  setup() {
    const route = useRoute()
    const router = useRouter()
    setPage({ title: route.params.lid === 'NEW' ? 'Add List' : 'Edit List' })
    const submitLabel = route.params.lid === 'NEW' ? 'Add List' : 'Save List'

    // Load the page data.
    const loading = ref(true)
    const list = ref({} as GetList)
    axios.get<GetList>(`/api/lists/${route.params.lid}`).then((resp) => {
      list.value = resp.data
      loading.value = false
    })

    // The name field.
    const duplicateName = ref('')
    const nameLabel = computed(() => (list.value.type === 'email' ? 'Email address' : 'Name'))
    const nameHelp = computed(() => (list.value.type === 'email' ? '@SunnyvaleSERV.org' : null))
    function nameError(lostFocus: boolean): string {
      if (!lostFocus) return ''
      if (!list.value.name)
        return list.value.type === 'email'
          ? 'The email address is required.'
          : 'The list name is required.'
      if (list.value.type === 'email' && !list.value.name.match(/^[a-z][-a-z0-9]*$/))
        return 'The email address must start with a lowercase letter and consist of lowercase letters and digits.'
      if (duplicateName.value === list.value.name) return 'Another list has this name.'
      return ''
    }

    // Submit the changes.
    const submitting = ref(false)
    async function onSubmit() {
      const body = new FormData()
      body.append('type', list.value.type)
      body.append('name', list.value.name)
      submitting.value = true
      const resp = (await axios.post<PostList>(`/api/lists/${route.params.lid}`, body)).data
      submitting.value = false
      if (resp && resp.duplicateName) {
        duplicateName.value = list.value.name
      } else {
        router.push('/admin/lists')
      }
    }

    // Handle deletions of lists.
    const deleteModal = ref(null as any)
    async function onDelete() {
      if (deleteModal.value) {
        const confirmed: boolean = await deleteModal.value.show()
        if (confirmed) {
          submitting.value = true
          await axios.delete(`/api/lists/${route.params.lid}`)
          submitting.value = false
          router.push('/admin/lists')
        }
      }
    }

    return {
      deleteModal,
      loading,
      list,
      nameError,
      nameHelp,
      nameLabel,
      onDelete,
      onSubmit,
      submitLabel,
      submitting,
      typeOptions,
    }
  },
})
</script>

<style lang="postcss">
#list {
  margin: 1.5rem 0.75rem;
}
</style>
