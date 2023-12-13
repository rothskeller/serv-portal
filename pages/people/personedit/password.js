up.on('up:form:submit', '.personeditPassword', evt => {
  if (!document.getElementById('personeditPasswordNew').valid) evt.preventDefault()
  const old = document.getElementById('personeditPasswordOld')
  if (old && !old.value) evt.preventDefault()
})
