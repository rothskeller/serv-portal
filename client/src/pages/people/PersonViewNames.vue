<!--
PersonViewNames displays the names part of the person view page.
-->

<template lang="pug">
#person-view-names
  #person-view-names-ifc
    #person-view-names-ic
      #person-view-names-i(v-text='person.informalName')
      #person-view-names-c(v-if='person.callSign', v-text='person.callSign')
    #person-view-names-f(
      v-if='person.formalName !== person.informalName',
      v-text='`(${person.formalName})`'
    )
  #person-view-names-edit(v-if='person.canEdit')
    SButton(variant='primary', small, @click='onEditNames') Edit
  PersonEditNames(v-if='person.canEdit', ref='editNamesModal', :pid='person.id')
</template>

<script lang="ts">
import { defineComponent, PropType, ref } from 'vue'
import { SButton } from '../../base'
import PersonEditNames from './PersonEditNames.vue'
import { GetPerson } from './PersonView.vue'

export default defineComponent({
  components: { PersonEditNames, SButton },
  props: {
    person: { type: Object as PropType<GetPerson>, required: true },
  },
  emits: ['reload'],
  setup(props, { emit }) {
    const editNamesModal = ref(null as any)
    async function onEditNames() {
      if (!(await editNamesModal.value.show())) return
      emit('reload')
    }
    return {
      editNamesModal,
      onEditNames,
    }
  },
})
</script>

<style lang="postcss">
#person-view-names {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}
#person-view-names-ifc {
  margin-bottom: 1.5rem;
}
#person-view-names-ic {
  font-size: 1.25rem;
  line-height: 1.2;
  color: black;
}
#person-view-names-i {
  display: inline;
  font-weight: bold;
}
#person-view-names-c {
  display: inline;
  margin-left: 0.5rem;
}
#person-view-names-f {
  color: #888;
  line-height: 1;
}
#person-view-names-edit {
  margin-left: 0.5rem;
}
</style>
