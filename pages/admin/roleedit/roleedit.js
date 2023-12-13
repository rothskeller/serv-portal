up.on('input', '#roleeditName', (evt, elm) => {
  const title = document.getElementById('roleeditTitle')
  if (title.dataset.match) {
    title.value = elm.value
  }
})
up.on('input', '#roleeditTitle', (evt, elm) => {
  const name = document.getElementById('roleeditName')
  if (name.value === elm.value) elm.dataset.match = true
  else elm.dataset.match = ''
})
