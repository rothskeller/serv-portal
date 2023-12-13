; (function () {
  function makeFooter(elm, clicked) {
    const footer = document.getElementById('eventscalFooter')
    footer.classList.toggle('eventscalFooter-clicked', clicked)
    const heading = document.createElement('div')
    heading.className = 'eventscalFooterDate'
    heading.textContent = elm.dataset.date
    let events = elm.querySelector('.eventscalEvents')
    if (events) events = events.cloneNode(true)
    else {
      events = document.createElement('div')
      events.textContent = 'No events scheduled.'
    }
    footer.replaceChildren(heading, events)
  }
  up.on('mouseover', '.eventscalDay', (evt, elm) => {
    const footer = document.getElementById('eventscalFooter')
    if (footer.classList.contains('eventscalFooter-clicked')) return
    const date = elm.dataset.date
    if (!date) {
      if (footer.children.length) footer.replaceChildren()
      return
    }
    if (!footer.children.length || footer.children[0].textContent !== date) makeFooter(elm, false)
  })
  up.on('mouseout', '.eventscalDay', () => {
    const footer = document.getElementById('eventscalFooter')
    if (footer.classList.contains('eventscalFooter-clicked')) return
    if (footer.children.length) footer.replaceChildren()
  })
  function onClick(evt, elm) {
    const footer = document.getElementById('eventscalFooter')
    const date = elm.dataset.date
    if (date && footer.children.length && footer.children[0].textContent === date) {
      if (footer.classList.contains('eventscalFooter-clicked')) {
        footer.classList.remove('eventscalFooter-clicked')
        elm.classList.remove('eventscalDay-clicked')
      } else {
        footer.classList.add('eventscalFooter-clicked')
        elm.classList.add('eventscalDay-clicked')
      }
      return
    }
    if (date) makeFooter(elm, true)
    else footer.replaceChildren()
    document.querySelectorAll('.eventscalDay-clicked').forEach(elm => {
      elm.classList.remove('eventscalDay-clicked')
    })
    if (date) elm.classList.add('eventscalDay-clicked')
  }
  up.on('click', '.eventscalDay', onClick)
  up.on('click', '.eventscalDay .eventscalEventLink', (evt, elm) => {
    // On a touch display, the stylesheet makes the links in the grid look like
    // plain text.  We don't want them to act like links, so we'll prevent their
    // click action.  However, we do still want them to toggle the click state
    // of the containing date in the grid, so we'll directly call our own click
    // handler.
    if (getComputedStyle(elm).cursor === 'default') {
      evt.preventDefault()
      onClick(evt, elm.closest('.eventscalDay'))
    }
  })
  up.on('change', '#eventscalMonth', (evt, elm) => {
    up.navigate({ url: `/events/calendar/${elm.getAttribute('value')}` })
  })
})()
