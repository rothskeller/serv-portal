up.on('change', '#peoplelistRole', (evt, elm) => {
  up.submit(elm.form, { target: 'main' })
})
up.on('click', '.peoplelistPersonNameroles', (evt, elm) => {
  if (getComputedStyle(elm).getPropertyValue('--touch') === '1') {
    up.follow(up.element.get(elm, 'a'))
  }
})
up.on('click', '.peoplelistPersonEmailphone', (evt, elm) => {
  if (getComputedStyle(elm).getPropertyValue('--touch') === '1') {
    const email = up.element.get(elm, 'a')
    window.open(email.href, '_blank')
  }
});
(function () {
  let detailsShown
  function toggleDetails(evt, elm) {
    if (detailsShown) {
      const pus = detailsShown.getElementsByClassName('peoplelistDetails')
      pus[0].style.display = 'none'
      if (detailsShown === elm) {
        detailsShown = null
        return
      }
    }
    detailsShown = elm
    const pus = elm.getElementsByClassName('peoplelistDetails')
    pus[0].style.display = null
  }
  up.on('click', '.peoplelistPersonDetails', toggleDetails)
})()
