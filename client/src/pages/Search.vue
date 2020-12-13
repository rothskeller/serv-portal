<!--
Search displays the search page.
-->

<template lang="pug">
#search
  form#search-form(@submit.prevent='onSubmit')
    #search-query-row
      SInput#search-query(autofocus, v-model='query')
      SButton(type='submit', variant='primary') Search
  #search-error(v-if='error', v-text='error')
  #search-results
    template(v-if='roles.length')
      .search-result-type Roles
      .search-result(v-for='r in roles')
        router-link(:to='`/people/list?role=${r.id}`') {{ r.name }}
    template(v-if='people.length')
      .search-result-type People
      .search-result(v-for='p in people')
        router-link(:to='`/people/${p.id}`') {{ p.informalName }}
    template(v-if='events.length')
      .search-result-type Events
      .search-result(v-for='e in events')
        router-link(:to='`/events/${e.id}`') {{ e.date }} {{ e.name }}
    template(v-if='folders.length')
      .search-result-type Folders
      .search-result(v-for='f in folders')
        router-link(:to='`/files${f.url}`') {{ f.name }}
        span.search-result-path(v-text='resultPath(f)')
    template(v-if='documents.length')
      .search-result-type Files
      .search-result(v-for='d in documents')
        a(:href='d.url', :target='d.newtab ? "_blank" : null') {{ d.name }}
        span.search-result-path(v-text='resultPath(d)')
    template(v-if='textMessages.length')
      .search-result-type Text Messages
      .search-result(v-for='tm in textMessages')
        router-link(:to='`/texts/${tm.id}`') From {{ tm.sender }} on {{ tm.timestamp.substr(0, 10) }}
    .search-result(
      v-if='submitted && !error && !roles.length && !people.length && !events.length && !folders.length && !documents.length && !textMessages.length'
    ) No results found.
</template>

<script lang="ts">
import { defineComponent, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from '../plugins/axios'
import setPage from '../plugins/page'
import { SButton, SInput } from '../base'

interface GetSearchResultDoc {
  type: 'document'
  name: string
  url: string
  path: Array<string>
  newtab: boolean
}
interface GetSearchResultEvent {
  type: 'event'
  id: number
  date: string
  name: string
}
interface GetSearchResultFolder {
  type: 'folder'
  name: string
  url: string
  path: Array<string>
}
interface GetSearchResultPerson {
  type: 'person'
  id: number
  informalName: string
}
interface GetSearchResultRole {
  type: 'role'
  id: number
  name: string
}
interface GetSearchResultText {
  type: 'textMessage'
  sender: string
  timestamp: string
}
type GetSearchResult =
  | GetSearchResultDoc
  | GetSearchResultEvent
  | GetSearchResultFolder
  | GetSearchResultPerson
  | GetSearchResultRole
  | GetSearchResultText
interface GetSearch {
  results: Array<GetSearchResult>
  error?: string
}

export default defineComponent({
  components: { SButton, SInput },
  setup() {
    const route = useRoute()
    const router = useRouter()
    setPage({ title: 'Search' })

    const query = ref('')
    if (route.query.q) {
      query.value = route.query.q as string
      onSubmit()
    }

    const documents = ref([] as Array<GetSearchResultDoc>)
    const events = ref([] as Array<GetSearchResultEvent>)
    const folders = ref([] as Array<GetSearchResultFolder>)
    const people = ref([] as Array<GetSearchResultPerson>)
    const roles = ref([] as Array<GetSearchResultRole>)
    const textMessages = ref([] as Array<GetSearchResultText>)
    const error = ref('')
    const submitted = ref(false)
    async function onSubmit() {
      query.value = query.value.trim()
      if (!query.value) {
        documents.value = events.value = folders.value = roles.value = people.value = textMessages.value = []
        error.value = ''
        return
      }
      if (query.value !== route.query.q)
        router.replace({ path: '/search', query: { q: query.value } })
      const resp = (await axios.get<GetSearch>('/api/search', { params: { q: query.value } })).data
      error.value = resp.error || ''
      documents.value = resp.results.filter(
        (r) => r.type === 'document'
      ) as Array<GetSearchResultDoc>
      events.value = resp.results.filter((r) => r.type === 'event') as Array<GetSearchResultEvent>
      folders.value = resp.results.filter(
        (r) => r.type === 'folder'
      ) as Array<GetSearchResultFolder>
      people.value = resp.results.filter((r) => r.type === 'person') as Array<GetSearchResultPerson>
      roles.value = resp.results.filter((r) => r.type === 'role') as Array<GetSearchResultRole>
      textMessages.value = resp.results.filter(
        (r) => r.type === 'textMessage'
      ) as Array<GetSearchResultText>
      submitted.value = true
    }

    function resultPath(result: GetSearchResultDoc | GetSearchResultFolder) {
      if (!result.path.length) return ''
      return 'in ' + result.path.join(' > ')
    }

    return {
      documents,
      error,
      events,
      folders,
      onSubmit,
      people,
      query,
      resultPath,
      roles,
      submitted,
      textMessages,
    }
  },
})
</script>

<style lang="postcss">
#search {
  margin: 1.5rem 0.75rem;
}
#search-query-row {
  text-align: center;
}
#search-query {
  display: inline;
  margin-right: 0.25rem;
  width: 10rem;
  vertical-align: bottom;
}
#search-error {
  margin-top: 0.75rem;
  color: red;
  text-align: center;
}
#search-results {
  margin-top: 0.75rem;
}
.search-result-type {
  padding-top: 0.75rem;
  text-decoration: underline;
}
.search-result {
  margin-left: 2rem;
  text-indent: -2rem;
  line-height: 1.2;
}
.search-result-path {
  padding-left: 1rem;
  font-style: italic;
  font-size: 0.75rem;
}
</style>
