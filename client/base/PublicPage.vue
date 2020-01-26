<!--
PublicPage is a minimal version of the page structure in MainMenu.vue, used for
public (nonauthenticated) pages.
-->

<template lang="pug">
#page-top(:class="$store.state.touch ? 'touch' : 'mouse'")
  #page-heading
    #page-menu-spacer
    #page-titlebox
      #page-title(v-text="title")
    #page-menu-spacer
  #page-main
    #page-content
      slot
</template>

<script>
export default {
  props: {
    title: String,
  },
}
</script>
