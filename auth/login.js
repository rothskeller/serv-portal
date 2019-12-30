window.addEventListener('load', function () {
  const lf = document.getElementById('login-form')
  if (lf) lf.addEventListener('submit', function (evt) {
    if (evt.target['login-email'].value.trim() === '' || evt.target['login-password'].value.trim() === '')
      evt.preventDefault()
  })
})
