<!--
Policies displays the Policies page.
-->

<template lang="pug">
PublicPage(title="Sunnyvale SERV")
  Policies
  b-btn.my-3(@click="$router.go(-1)") Back
</template>

<script>
import Policies from '@/base/Policies'
import PublicPage from '@/base/PublicPage'

export default {
  components: { Policies, PublicPage },
}
</script>

<style lang="stylus"></style>
