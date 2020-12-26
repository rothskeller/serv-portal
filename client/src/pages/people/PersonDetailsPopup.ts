// PersonDetailsPopup is the popup that appears when clicking on the details
// button for a person in the PeopleList.

import { defineComponent, h, PropType } from 'vue'
import { SIcon } from '../../base'
import { GetPeoplePerson } from './api'

const PersonDetailsPopup = defineComponent({
  name: 'PersonDetailsPopup',
  props: {
    person: { type: Object as PropType<GetPeoplePerson>, required: true },
  },
  setup(props) {
    return () =>
      h('div', { id: 'person-details' }, [
        h('div', { id: 'person-details-name' }, props.person.informalName),
        h(
          'div',
          { id: 'person-details-roles' },
          ...props.person.roles.map((r: string) => h('div', r))
        ),
        props.person.email || props.person.email2
          ? h('div', { class: 'person-details-spacer' })
          : null,
        props.person.email ? renderEmail(props.person.email) : null,
        props.person.email2 ? renderEmail(props.person.email2) : null,
        props.person.cellPhone || props.person.homePhone || props.person.workPhone
          ? h('div', { class: 'person-details-spacer' })
          : null,
        props.person.cellPhone ? renderPhone(props.person.cellPhone, 'cell', true) : null,
        props.person.homePhone ? renderPhone(props.person.homePhone, 'home', false) : null,
        props.person.workPhone ? renderPhone(props.person.workPhone, 'work', false) : null,
      ])
  },
})
export default PersonDetailsPopup

function renderEmail(email: string) {
  return [
    h('div', email),
    h('div'),
    h(
      'div',
      {
        class: 'person-details-icon',
        onClick: () => {
          window.open(`mailto:${encodeURIComponent(email)}`, '_blank')
        },
      },
      h(
        'a',
        { href: `mailto:${encodeURIComponent(email)}`, target: '_blank' },
        h(SIcon, { icon: 'email' })
      )
    ),
  ]
}

function renderPhone(phone: string, label: string, isCell: boolean) {
  return [
    h('div', `${phone} (${label})`),
    isCell
      ? h(
          'div',
          {
            class: 'person-details-icon',
            onClick: () => {
              window.open(`sms:${encodeURIComponent(phone)}`, '_blank')
            },
          },
          h(
            'a',
            { href: `sms:${encodeURIComponent(phone)}`, target: '_blank' },
            h(SIcon, { icon: 'message' })
          )
        )
      : h('div'),
    h(
      'div',
      {
        class: 'person-details-icon',
        onClick: () => {
          window.open(`tel:${encodeURIComponent(phone)}`, '_blank')
        },
      },
      h(
        'a',
        { href: `tel:${encodeURIComponent(phone)}`, target: '_blank' },
        h(SIcon, { icon: 'phone' })
      )
    ),
  ]
}
