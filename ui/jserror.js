window.addEventListener('error', function (evt) {
  if (evt.message.includes('ResizeObserver')) {
    evt.stopImmediatePropagation()
    return
  }
  var details = JSON.stringify({
    msg: evt.message, url: evt.filename, line: evt.lineno, col: evt.colno, err: evt.error,
    stack: evt.error ? evt.error.stack : null,
  })
  var request = new XMLHttpRequest()
  request.open('POST', '/jserror')
  request.send(details)
});
