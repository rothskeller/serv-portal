<!--
PublicTopFolder displays one folder on the public page.
-->

<template lang="pug">
.public-folder
  .public-folder-logo-div
    img.public-folder-logo(v-if='folderLogo', :src='folderLogo')
    span(v-else) &nbsp;
  .public-folder-text
    router-link.public-folder-link(:to='folderURL', v-text='folder.name')
    div(v-if='folderData.description', v-html='folderData.description')
</template>

<script>
export default {
  props: {
    folder: Object,
  },
  data: () => ({ folderData: null }),
  computed: {
    folderLogo() {
      if (!this.folderData || !this.folderData.documents) return null
      const logo = this.folderData.documents.find(d => d.name === 'folder.png')
      if (!logo) return null
      return `/api/folders/${this.folder.id}/${logo.id}`
    },
    folderURL() {
      return `/public/${this.folder.id}`
    },
  },
  created() {
    this.loadFolder()
  },
  methods: {
    async loadFolder() {
      this.folderData = (await this.$axios.get(`/api/folders/${this.folder.id}`)).data
    },
  }
}
</script>

<style lang="stylus">
.public-folder
  clear both
  margin 1.5rem auto 0
  max-width 750px
.public-folder-logo-div
  float left
  margin-right 1rem
  max-width 120px
  max-height 120px
  width 25%
  text-align center
.public-folder-logo
  max-width 120px
  max-height 120px
.public-folder-link
  font-weight bold
  font-size 1.25rem
</style>
