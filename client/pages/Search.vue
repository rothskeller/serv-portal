<!--
Search displays the search page.
-->

<template lang="pug">
#search
  form#search-form(@submit.prevent="onSubmit")
    #search-query-row
      b-form-input#search-query(autofocus v-model="query")
      b-btn(type="submit" variant="primary") Search
  #search-error(v-if="error" v-text="error")
  #search-results
    template(v-if="groups.length")
      .search-result-type Groups
      .search-result(v-for="g in groups")
        b-link(:to="`/people/list?group=${g.id}`") {{g.name}}
    template(v-if="people.length")
      .search-result-type People
      .search-result(v-for="p in people")
        b-link(:to="`/people/${p.id}`") {{p.informalName}}
    template(v-if="events.length")
      .search-result-type Events
      .search-result(v-for="e in events")
        b-link(:to="`/events/${e.id}`") {{e.date}} {{e.name}}
    template(v-if="folders.length")
      .search-result-type Folders
      .search-result(v-for="f in folders")
        b-link(:to="`/files/${f.id}`") {{f.name}}
        span.search-result-path(v-if="f.path.length" v-text="resultPath(f)")
    template(v-if="documents.length")
      .search-result-type Files
      .search-result(v-for="d in documents")
        b-link(v-if="d.name.endsWith('.pdf')" :href="`/api/folders/${d.folderID}/${d.documentID}`" target="_blank") {{d.name}}
        b-link(v-else :href="`/api/folders/${d.folderID}/${d.documentID}`" download) {{d.name}}
        span.search-result-path(v-text="resultPath(d)")
    template(v-if="textMessages.length")
      .search-result-type Text Messages
      .search-result(v-for="tm in textMessages")
        b-link(:to="`/texts/${tm.id}`") From {{tm.sender}} on {{tm.timestamp.substr(0, 10)}}
    .search-result(v-if="!error && !groups.length && !people.length && !events.length && !folders.length && !documents.length && !textMessages.length") No results found.
</template>

<script>
export default {
  data: () => ({
    query: null,
    error: null,
    documents: [],
    events: [],
    folders: [],
    groups: [],
    people: [],
    textMessages: [],
  }),
  mounted() {
    this.$store.commit('setPage', { title: 'Search' })
    this.query = this.$route.query.q
    if (this.query) this.onSubmit()
  },
  methods: {
    async onSubmit() {
      this.query = this.query.trim()
      if (!this.query) {
        this.documents = this.events = this.folders = this.groups = this.people = this.textMessages = []
        this.error = null
        return
      }
      if (this.query !== this.$route.query.q)
        this.$router.replace({ path: '/search', query: { q: this.query } })
      const resp = (await this.$axios.get('/api/search', { params: { q: this.query } })).data
      this.error = resp.error
      this.documents = resp.results.filter(r => r.type === 'document')
      this.events = resp.results.filter(r => r.type === 'event')
      this.folders = resp.results.filter(r => r.type === 'folder')
      this.groups = resp.results.filter(r => r.type === 'group')
      this.people = resp.results.filter(r => r.type === 'person')
      this.textMessages = resp.results.filter(r => r.type === 'textMessage')
    },
    resultPath(result) {
      return 'in ' + result.path.join(' > ')
    },
  },
}
</script>

<style lang="stylus">
#search-query-row
  text-align center
#search-query
  display inline
  margin-right 0.25rem
  width 10rem
  vertical-align bottom
#search-error
  margin-top 0.75rem
  color red
  text-align center
#search-results
  margin-top 0.75rem
.search-result-type
  padding-top 0.75rem
  text-decoration underline
.search-result
  margin-left 2rem
  text-indent -2rem
  line-height 1.2
.search-result-path
  padding-left 1rem
  font-style italic
  font-size 0.75rem
</style>
