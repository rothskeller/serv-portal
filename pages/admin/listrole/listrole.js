up.on('submit', '.listeditRoleForm', (evt, elm) => {
  evt.preventDefault()
  const data = {
    roles: elm.elements.roles.getAttribute('value').split(/ /),
    submodel: elm.elements.submodel.value,
    sender: elm.elements.sender.checked,
  }
  console.log(data)
  up.layer.accept(data)
})
