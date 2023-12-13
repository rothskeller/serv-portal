// The <s-hours> element displays an entry field for entering volunteer hours.
// It accepts three attributes:
//   disabled is a boolean value indicating that the control should be disabled.
//   name is the name that should be given to the hours field when it is part of
//     a form submission.  This is required.
//   value is the value that should be initially populated into the field.  It
//     must be either empty, "½", or an integer possibly followed by "½".
// This works similarly to a regular entry field, except:
//   - Hours can be entered with decimals (e.g. 3.5), but are converted to the
//     above (3½) on blur.
//   - The up and down arrow buttons, and +/- keys, raise and lower the value by
//     ½ hour.
//   - If the user enters more than two digits (without a decimal point), the
//     control switches to timesheet mode.  In this mode, it expects eight
//     digits to be entered; it treats that as a start and end time, and uses
//     the difference between them as the value.
// Callers can query the "timesheet" Boolean property to determine whether the
// value was entered in timesheet mode.
class SHours extends HTMLElement {
  static get observedAttributes() { return ['value'] };

  constructor() {
    super()
    this._in = this._in2 = this._frame = null
    this.timesheet = false
  }

  connectedCallback() {
    // Create the main input field.
    this._in = document.createElement('input')
    this._in.className = 'formInput SHoursInput1'
    this._in.autocomplete = 'off'
    this._in.disabled = this.hasAttribute('disabled')
    this._in.name = this.getAttribute('name')
    this._in.value = this.getAttribute('value')
    this._prev = this._in.value
    if (!this._in.disabled) {
      this._in.addEventListener('keydown', this.onKeyDown.bind(this))
      this._in.addEventListener('input', this.onInput.bind(this))
      this._in.addEventListener('focus', this.onFocus.bind(this))
      this._in.addEventListener('blur', this.onBlur.bind(this))
    }
    this.appendChild(this._in)
  }

  attributeChangedCallback(name, oldValue, newValue) {
    if (this._in) this._in.value = newValue
  }

  onKeyDown(evt) {
    const val = this._in.value
    // Nothing special unless the entry matches a valid number of hours.
    if (!val.match(/^\d{0,2}(?:½|\.\d*)?$/)) return
    // Nothing special unless the key is an arrow.
    let offset
    switch (evt.key) {
      case 'ArrowUp':
      case 'Up':
      case '+':
        offset = 0.5
        break
      case 'ArrowDown':
      case 'Down':
      case '-':
        offset = -0.5
        break
    }
    if (!offset) return
    evt.preventDefault()
    // Convert the entry into a number of hours.
    let hours
    if (val.match(/\./)) {
      hours = parseFloat(val) || 0
    } else if (val.match(/½$/)) {
      hours = (parseInt(val.substring(0, val.length - 1)) || 0) + 0.5
    } else {
      hours = parseInt(val) || 0
    }
    // Apply the offset, with corrections.
    hours += offset
    if (hours < 0) hours = 0
    if (hours >= 100) hours = 99.5
    // Export the new value.
    let fraction = hours - Math.floor(hours)
    if (hours <= 10 / 60) this._in.value = '0'
    else if (hours < 40 / 60) this._in.value = '½'
    else if (fraction <= 10 / 60) this._in.value = `${Math.floor(hours)}`
    else if (fraction > 40 / 60) this._in.value = `${Math.ceil(hours)}`
    else this._in.value = `${Math.floor(hours)}½`
    this._prev = this._in.value
    // We prevented the default action above (don't want to treat the key
    // literally), so the input will not emit an input event.  We should do that
    // ourselves in case users of s-hours are watching for it.
    this._in.dispatchEvent(new InputEvent('input', { bubbles: true }))
  }

  onInput(evt) {
    const val = this._in.value
    // If this entry matches a valid number of hours, we're golden.
    if (val.match(/^\d{0,2}(?:½|\.\d*)?$/)) {
      this._prev = val
      this.closeTimesheet(true)
      return
    }
    // If it's got 3 digits, insert a colon between the second and third.
    if (val.match(/^(?:[01]\d|2[0-3])[0-5]$/)) {
      this._prev = this._in.value = `${val.substring(0, 2)}:${val.substring(2)}`
      this.openTimesheet()
      return
    }
    // If it's got xx:xx, focus the timesheet entry.
    if (val.match(/^(?:[01]\d?|2[0-3]?):[0-5]\d$/)) {
      this._prev = val
      this.openTimesheet()
      this._in2.focus()
      return
    }
    // If it's got xx:x, that's OK.
    if (val.match(/^(?:[01]\d?|2[0-3]?):[0-5]$/)) {
      this._prev = val
      this.openTimesheet()
      return
    }
    // If it's got xx:, remove the colon and close the timesheet.
    if (val.match(/^(?:[01]\d?|2[0-3]?):$/)) {
      this._prev = this._in.value = val.substring(0, 2)
      this.closeTimesheet(true)
      return
    }
    // Anything else is bogus, revert the entry.
    this._in.value = this._prev
  }

  onFocus(evt) {
    // Ignore if focus coming from second input.
    if (this._in2 && evt.relatedTarget === this._in2) return
    this._in.select()
  }

  onBlur(evt) {
    // Ignore if focus moving to second input.
    if (this._in2 && evt.relatedTarget === this._in2) return
    // If the input value contains a decimal, change it to be integer with maybe
    // a ½.
    if (this._in.value.match(/\./)) {
      let hours = parseFloat(this._in.value) || 0
      let fraction = hours - Math.floor(hours)
      if (hours <= 10 / 60) this._in.value = '0'
      else if (hours < 40 / 60) this._in.value = '½'
      else if (fraction <= 10 / 60) this._in.value = `${Math.floor(hours)}`
      else if (fraction > 40 / 60) this._in.value = `${Math.ceil(hours)}`
      else this._in.value = `${Math.floor(hours)}½`
    }
    this.timesheet = false
    this.setAttribute('value', this._in.value)
    this.dispatchEvent(new Event('change', { bubbles: true }))
  }

  openTimesheet(text) {
    if (!this._in2) {
      this._in2 = document.createElement('input')
      this._in2.className = 'formInput SHoursInput2'
      this._in2.autocomplete = 'off'
      this._prev2 = ''
      this._in2.addEventListener('keydown', this.onKeyDown2.bind(this))
      this._in2.addEventListener('input', this.onInput2.bind(this))
      this._in2.addEventListener('blur', this.onBlur2.bind(this))
      this.appendChild(this._in2)
      this._frame = document.createElement('div')
      this._frame.className = 'SHoursFrame'
      this._frame.textContent = 'Timesheet Entry'
      const ti = document.createElement('div')
      ti.className = 'SHoursTI'
      ti.textContent = 'Time In'
      this._frame.appendChild(ti)
      const to = document.createElement('div')
      to.className = 'SHoursTO'
      to.textContent = 'Time Out'
      this._frame.appendChild(to)
      this.insertBefore(this._frame, this._in)
    }
  }

  closeTimesheet(refocus) {
    if (this._in2) {
      this.removeChild(this._in2)
      this._in2 = null
      this.removeChild(this._frame)
      this._frame = null
    }
    if (refocus) this._in.focus()
  }

  onKeyDown2(evt) {
    if (evt.key === 'Backspace' && this._in2.value === '') {
      evt.preventDefault()
      this._prev = this._in.value = this._in.value.substring(0, this._in.value.length - 1)
      this._in.focus()
    }
  }

  onInput2(evt) {
    const val = this._in2.value
    // Empty: return focus to first input.
    if (val === '') {
      this.prev2 = val
      this._in.focus()
      return
    }
    // x, xx, xx:x, xx:xx -- accept unchanged.
    if (val.match(/^[012]$/) || val.match(/^(?:[01]\d|2[0-3])$/) || val.match(/^(?:[01]\d?|2[0-3]?):[0-5]$/) || val.match(/^(?:[01]\d?|2[0-3]?):[0-5]\d$/)) {
      this._prev2 = val
      return
    }
    // If it's got 3 digits, insert a colon between the second and third.
    if (val.match(/^(?:[01]\d|2[0-3])[0-5]$/)) {
      this._prev2 = this._in2.value = `${val.substring(0, 2)}:${val.substring(2)}`
      return
    }
    // If it's got xx:, remove the colon.
    if (val.match(/^(?:[01]\d?|2[0-3]?):$/)) {
      this._prev2 = this._in2.value = val.substring(0, 2)
      return
    }
    // Anything else is bogus, revert the entry.
    this._in2.value = this._prev2
  }

  onBlur2(evt) {
    // Ignore if focus moving to first input.
    if (evt.relatedTarget === this._in) return
    // Test for correct formatting.
    const val1 = this._in.value
    const val2 = this._in2.value
    this.closeTimesheet(false)
    if (!val1.match(/^(?:[01]\d|2[0-3]):[0-5]\d$/) || !val2.match(/^(?:[01]\d|2[0-3]):[0-5]\d$/)) {
      this._in.value = ''
      return
    }
    const t1 = parseInt(val1.substring(0, 2)) * 60 + parseInt(val1.substring(3))
    const t2 = parseInt(val2.substring(0, 2)) * 60 + parseInt(val2.substring(3))
    const diff = t2 - t1
    if (diff < 0) {
      this._in.value = ''
      return
    }
    const hours = diff / 60
    let fraction = hours - Math.floor(hours)
    if (hours <= 10 / 60) this._in.value = '0'
    else if (hours < 40 / 60) this._in.value = '½'
    else if (fraction <= 10 / 60) this._in.value = `${Math.floor(hours)}`
    else if (fraction > 40 / 60) this._in.value = `${Math.ceil(hours)}`
    else this._in.value = `${Math.floor(hours)}½`
    this.timesheet = true
    this.setAttribute('value', this._in.value)
    this.dispatchEvent(new Event('change', { bubbles: true }))
  }
}
customElements.define('s-hours', SHours)
