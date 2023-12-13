up.on('click', '.listeditRoleEdit', async (evt, elm) => {
  evt.preventDefault()
  const list = elm.dataset.list
  const role = elm.dataset.role || 'NEW'
  up.layer.open({
    url: `/admin/lists/${list}/roleedit/${role}`,
    size: 'grow',
    dismissable: 'key',
    history: false,
    onAccepted: evt => {
      evt.value.roles.forEach(rid => {
        let input = document.getElementById(`listeditRole${rid}`)
        if (!input) {
          input = document.createElement('input')
          input.type = 'hidden'
          input.name = `role${rid}`
          elm.parentElement.appendChild(input)
        }
        input.value = `${evt.value.submodel}:${evt.value.sender}`
      })
      up.validate('form')
    },
  })
})
