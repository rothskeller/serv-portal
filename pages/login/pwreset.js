up.on('up:form:submit', '.pwResetForm', evt => {
  const newpwd = document.getElementById('pwresetNewPassword')
  if (!newpwd.valid) evt.preventDefault()
})
