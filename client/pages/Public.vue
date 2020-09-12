<!--
Public displays the public page (the root page of the site).  It's the page
shown at / when the user is not logged in.
-->

<template lang="pug">
PublicPage(title='Sunnyvale SERV')
  #public
    #public-head
      img#public-head-logo(:src='servLogo')
      #public-head-text.
        <span style="font-size:1.25rem;font-weight:bold">Sunnyvale Emergency
        Response Volunteers (SERV)</span> is the volunteer arm of the Sunnyvale
        Department of Public Safety, Office of Emergency Services.  SERV
        volunteers teach disaster preparedness classes, assist uniformed
        officers in emergencies, and respond in disasters when professional
        responders are overloaded.
      #public-head-login
        b-btn(to='/login', variant='primary') Volunteer Login
    #public-folders(v-if='root')
      PublicTopFolder(v-for='folder in root.children', :folder='folder')
</template>

<script>
import PublicPage from '@/base/PublicPage'
import PublicTopFolder from './public/PublicTopFolder'
import servLogo from './serv-logo.png'

export default {
  components: { PublicTopFolder, PublicPage },
  data: () => ({
    root: null,
    servLogo
  }),
  created() {
    this.loadRootFolder()
  },
  methods: {
    async loadRootFolder() {
      this.root = (await this.$axios.get(`/api/folders/0`)).data
    },
  }
}
</script>

<style lang="stylus">
#public-head
  margin 0 auto
  max-width 750px
#public-head-logo
  float left
  margin-right 1rem
  max-width 120px
  width 25%
#public-head-login
  text-align right
</style>
