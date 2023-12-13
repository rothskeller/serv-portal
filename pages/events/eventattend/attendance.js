; (function () {
  function getHours(row) {
    const hoursElm = up.element.get(row.closest('.attendanceRow'), '.attendanceHours')
    return hoursElm ? hoursElm.getAttribute('value') || '0' : '0'
  }
  function setHours(row, value) {
    const hoursElm = up.element.get(row.closest('.attendanceRow'), '.attendanceHours')
    if (hoursElm) hoursElm.setAttribute('value', value)
  }
  function hoursLess(a, b) {
    const a0 = (a.endsWith('½') ? a.substring(0, a.length - 1) : a) || '0'
    const b0 = (b.endsWith('½') ? b.substring(0, b.length - 1) : b) || '0'
    if (parseInt(a) < parseInt(b)) return true
    return a0 === b0 && a === a0 && b !== b0
  }
  function getSignedIn(row) {
    return up.element.get(row.closest('.attendanceRow'), '.attendanceSignedIn').classList.contains('true')
  }
  function toggleSignedIn(row, value) {
    const box = up.element.get(row.closest('.attendanceRow'), '.attendanceSignedIn')
    const signedin = box.classList.toggle('true', value)
    up.element.get(box, 'input').value = signedin
  }
  function getCredited(row) {
    return up.element.get(row.closest('.attendanceRow'), '.attendanceCredited').classList.contains('true')
  }
  function toggleCredited(row, value) {
    const box = up.element.get(row.closest('.attendanceRow'), '.attendanceCredited')
    const credited = box.classList.toggle('true', value)
    up.element.get(box, 's-icon').setAttribute('icon', credited ? 'star-solid' : 'star')
    up.element.get(box, 'input').value = credited
  }
  up.on('change', '.attendanceHours', (evt, elm) => {
    if (elm.timesheet) toggleSignedIn(elm, true)
  })
  up.on('click', '.attendanceSignedIn', (evt, elm) => {
    toggleSignedIn(elm)
  })
  up.on('click', '.attendanceCredited', (evt, elm) => {
    toggleCredited(elm)
  })
  up.on('click', '.attendanceGrid .attendanceName', (evt, elm) => {
    const myHours = getHours(elm), mySignedIn = getSignedIn(elm), myCredited = getCredited(elm)
    const defaults = up.element.get(elm.closest('form'), '.attendanceDefault')
    const defHours = getHours(defaults), defSignedIn = getSignedIn(defaults), defCredited = getCredited(defaults)
    if (myHours === defHours && mySignedIn === defSignedIn && myCredited === defCredited) {
      setHours(elm, '0')
      toggleSignedIn(elm, false)
      toggleCredited(elm, false)
    } else {
      if (hoursLess(myHours, defHours)) setHours(elm, defHours)
      toggleSignedIn(elm, defSignedIn)
      toggleCredited(elm, defCredited)
    }
  })
  up.on('s-change', '.attendanceNew input', (evt, elm) => {
    let pname = elm.value, pkey = elm.getAttribute('s-value')
    if (!pname) return
    elm.setAttribute('s-value', '')
    elm.value = ''
    elm.closest('.attendanceGrid').querySelectorAll('.attendanceName').forEach(name => {
      if (name.dataset.key === pkey) pkey = ''
    })
    if (!pkey) return
    const attendanceNew = elm.parentElement
    const defaults = up.element.get(elm.closest('form'), '.attendanceDefault')
    const defHours = getHours(defaults), defSignedIn = getSignedIn(defaults), defCredited = getCredited(defaults)
    const newRow = document.querySelector('.attendanceTemplate').cloneNode(true).content.firstChild
    const newHours = up.element.get(newRow, '.attendanceHours')
    if (newHours) newHours.setAttribute('name', `hours${pkey.substring(1)}`)
    up.element.get(newRow, '.attendanceSignedIn input').name = `signedin${pkey.substring(1)}`
    up.element.get(newRow, '.attendanceCredited input').name = `credited${pkey.substring(1)}`
    up.element.get(newRow, '.attendanceName').dataset.key = pkey
    up.element.get(newRow, '.attendanceName').textContent = pname
    setHours(newRow, defHours)
    toggleSignedIn(newRow, defSignedIn)
    toggleCredited(newRow, defCredited)
    attendanceNew.parentElement.insertBefore(newRow, attendanceNew)
  })
})()
