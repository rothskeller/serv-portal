<!--
PersonEditRoles is the dialog box for editing a person's roles.
-->

<template lang="pug">
Modal(ref='modal')
  SForm(
    dialog,
    variant='primary',
    title='Edit Roles',
    submitLabel='Save',
    :disabled='submitting',
    @submit='onSubmit',
    @cancel='onCancel'
  )
    SSpinner(v-if='!orgs.length')
    template(v-else)
      SFSelect#person-roles-org(
        v-if='orgs.length > 1',
        label='Org',
        :options='orgs',
        valueKey='org',
        labelKey='fmtOrg',
        v-model='org'
      )
      SFCheckGroup#person-roles-roles(
        :label='orgs.length > 1 ? "Roles" : ""',
        :options='orgRoleOptions',
        disabledKey='disabled',
        v-model='orgRoles'
      )
</template>

<script lang="ts">
import { defineComponent, PropType, ref, watch } from 'vue'
import axios from '../../plugins/axios'
import { Modal, SForm, SFCheckGroup, SFSelect, SSpinner } from '../../base'

interface GetPersonRolesOrgRole {
  id: number
  name: string
  held: boolean
  direct: boolean
  implicitOnly?: boolean
  implies: Array<number>
}
interface GetPersonRolesOrg {
  org: string
  roles: Array<GetPersonRolesOrgRole>
  fmtOrg: string // added locally
}
interface RoleOption {
  value: number
  label: string
  disabled: boolean
}

const fmtOrg: Record<string, string> = {
  admin: 'Admin',
  'cert-d': 'CERT Deployment',
  'cert-t': 'CERT Training',
  listos: 'Listos',
  sares: 'SARES',
  snap: 'SNAP',
}

function resetImpliedRoles(orgs: Array<GetPersonRolesOrg>) {
  const toSet = new Set<number>()
  orgs.forEach((o) => {
    o.roles.forEach((r) => {
      if (r.held && r.direct) r.implies.forEach((irid) => toSet.add(irid))
      else if (r.held) r.held = false
    })
  })
  orgs.forEach((o) => {
    o.roles.forEach((r) => {
      if (toSet.has(r.id)) r.held = true
    })
  })
}

export default defineComponent({
  components: { Modal, SFCheckGroup, SForm, SFSelect, SSpinner },
  props: {
    pid: { type: Number, required: true },
  },
  setup(props) {
    const modal = ref(null as any)
    function show() {
      return modal.value.show()
    }

    // Load the form data.
    const orgs = ref([] as Array<GetPersonRolesOrg>)
    const org = ref('')
    axios.get<Array<GetPersonRolesOrg>>(`/api/people/${props.pid}/roles`).then((resp) => {
      orgs.value = resp.data
      orgs.value.forEach((o) => {
        o.fmtOrg = fmtOrg[o.org]
      })
      org.value = orgs.value[0].org
    })

    // Handle the changing of the orgs.
    const orgRoleOptions = ref([] as Array<RoleOption>)
    const orgRoles = ref(new Set<number>())
    watch(org, () => {
      const roles = orgs.value.find((o) => o.org === org.value)!.roles
      orgRoles.value = new Set(roles.filter((r) => r.held).map((r) => r.id))
    })
    watch([orgRoles, org], (n, o) => {
      // To avoid an infinite loop, we need to make sure the two sets' contents
      // are actually different.
      let different = n[1] != o[1]
      ;(n[0] as Set<number>).forEach((r) => {
        if (!(o[0] as Set<number>).has(r)) different = true
      })
      ;(o[0] as Set<number>).forEach((r) => {
        if (!(n[0] as Set<number>).has(r)) different = true
      })
      if (!different) return
      // The roles that should be set directly are:
      // - those that already are set directly in a different org;
      // - those that were set directly in this org and still are set
      // - those that were not set at all in this org and are now
      const toSetDirectly = new Set<number>()
      const toSetIndirectly = new Set<number>()
      orgs.value.forEach((o) => {
        if (o.org === org.value) {
          o.roles
            .filter((r) => (!r.held || r.direct) && orgRoles.value.has(r.id))
            .forEach((r) => {
              toSetDirectly.add(r.id)
              r.implies.forEach((ir) => toSetIndirectly.add(ir))
            })
        } else {
          o.roles
            .filter((r) => r.held && r.direct)
            .forEach((r) => {
              toSetDirectly.add(r.id)
              r.implies.forEach((ir) => toSetIndirectly.add(ir))
            })
        }
      })
      console.log('toSetDirectly', toSetDirectly)
      console.log('toSetIndirectly', toSetIndirectly)
      // Reset the roles appropriately.
      orgs.value.forEach((o) => {
        o.roles.forEach((r) => {
          if (toSetDirectly.has(r.id)) {
            r.held = r.direct = true
          } else if (toSetIndirectly.has(r.id)) {
            r.held = true
            r.direct = false
          } else {
            r.held = r.direct = false
          }
        })
      })
      // Update orgRoles to reflect what's now held.
      const roles = orgs.value.find((o) => o.org === org.value)!.roles
      orgRoles.value = new Set<number>(roles.filter((r) => r.held).map((r) => r.id))
      // Update orgRoleOptions to reflect what's disabled.
      orgRoleOptions.value = roles.map((r) => ({
        value: r.id,
        label: r.name,
        disabled: r.implicitOnly || (r.held && !r.direct),
      }))
    })

    // Save and close.
    const submitting = ref(false)
    async function onSubmit() {
      var body = new FormData()
      orgs.value.forEach((o) => {
        o.roles
          .filter((r) => r.held && r.direct)
          .forEach((r) => {
            body.append('role', r.id.toString())
          })
      })
      submitting.value = true
      await axios.post(`/api/people/${props.pid}/roles`, body)
      submitting.value = false
      modal.value.close(true)
    }
    function onCancel() {
      modal.value.close(false)
    }

    return { modal, onCancel, onSubmit, org, orgRoleOptions, orgRoles, orgs, show, submitting }
  },
})
</script>

<style lang="postcss">
</style>
