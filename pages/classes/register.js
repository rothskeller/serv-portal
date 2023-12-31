up.on('click', '.classregClear', (evt, elm) => {
  const row = elm.dataset.row
  document.getElementById(`classregFirstname${row}`).value = ''
  document.getElementById(`classregLastname${row}`).value = ''
  document.getElementById(`classregEmail${row}`).value = ''
  document.getElementById(`classregCellPhone${row}`).value = ''
  document.getElementById(`classregFirstname${row}`).focus()
})
up.on('input', '.classregFirstname', (evt, elm) => {
  if (!elm.value) return
  const row = parseInt(elm.id.substring(17)) + 1
  const max = elm.closest('form').dataset.max
  if (max && row >= parseInt(max)) return
  if (document.getElementById(`classregFirstname${row}`)) return
  const insertBefore = elm.parentElement.parentElement.nextElementSibling.nextElementSibling.nextElementSibling
  const divider = elm.parentElement.parentElement.previousElementSibling.cloneNode(true)
  divider.classList.remove('first')
  divider.firstElementChild.textContent = divider.firstElementChild.textContent.replace(/ .*/, ` ${row + 1}`)
  divider.lastElementChild.dataset.row = row
  insertBefore.parentElement.insertBefore(divider, insertBefore)
  const names = elm.parentElement.parentElement.cloneNode(true)
  names.firstElementChild.htmlFor = names.firstElementChild.htmlFor.substring(0, 17) + row
  names.lastElementChild.firstElementChild.id = names.lastElementChild.firstElementChild.id.substring(0, 17) + row
  names.lastElementChild.firstElementChild.value = ''
  names.lastElementChild.lastElementChild.id = names.lastElementChild.firstElementChild.id.substring(0, 16) + row
  names.lastElementChild.lastElementChild.value = ''
  insertBefore.parentElement.insertBefore(names, insertBefore)
  const email = elm.parentElement.parentElement.nextElementSibling.cloneNode(true)
  email.firstElementChild.htmlFor = email.firstElementChild.htmlFor.substring(0, 13) + row
  email.lastElementChild.id = email.lastElementChild.id.substring(0, 13) + row
  email.lastElementChild.value = ''
  insertBefore.parentElement.insertBefore(email, insertBefore)
  const cellPhone = elm.parentElement.parentElement.nextElementSibling.nextElementSibling.cloneNode(true)
  cellPhone.firstElementChild.htmlFor = cellPhone.firstElementChild.htmlFor.substring(0, 17) + row
  cellPhone.lastElementChild.id = cellPhone.lastElementChild.id.substring(0, 17) + row
  cellPhone.lastElementChild.value = ''
  insertBefore.parentElement.insertBefore(cellPhone, insertBefore)
})
