up.on('click', 'button.personeditDSWExpireAdd', (evt, elm) => {
  const expire = elm.previousElementSibling
  if (expire.value === '') {
    let reg = elm.closest('.formRow')
    reg = reg.previousElementSibling
    reg = up.element.get(reg, 'input')
    if (reg.value !== '') {
      expire.value = reg.value.substring(0, 4) + "-12-31"
    }
  } else if (expire.value.match(/-12-31$/)) {
    expire.value = `${parseInt(expire.value.substring(0, 4)) + 1}-12-31`
  } else {
    expire.value = expire.value.substring(0, 4) + "-12-31"
  }
})
