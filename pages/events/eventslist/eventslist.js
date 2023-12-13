up.on('change', '#eventslistYear', (evt, elm) => {
  up.navigate({ url: `/events/list/${elm.getAttribute('value')}` })
})
