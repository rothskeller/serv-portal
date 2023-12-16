// The <s-password> element displays a pair of password entry fields, a meter
// for displaying password strength, and an area for feedback about the
// password.  It accepts four attributes:
//   name is the name that should be given to the password field when it is part
//     of a form submission.  This is required.
//   value is the value that should be initially populated into the password
//     fields, if any.
//   hints is a comma-separated list of hints to the password strength checker.
//     These are words that should be discouraged in a password, because they
//     are related to the site or the person.
//   override is a boolean attribute; if present, weak passwords will get
//     appropriate feedback but will be accepted anyway.
// This element defines a form field that will be submitted with any form that
// contains the element.  The value of the field will be the supplied password
// if it is valid, or an empty string otherwise.  Valid means that a password
// has been entered identically in both fields, and either it is sufficiently
// strong or the override attribute was set.  This element also exports a
// read-only valid boolean property indicating the same.
class SPassword extends HTMLElement {

  // The zxcvbn script is very large (about 400kB compressed), so we don't load
  // it until we need it, i.e., until the first <s-password> element is
  // constructed.  zxcvbnLoaded is 1 when the loading attempt starts and 2 when
  // loading is complete.
  static zxcvbnLoaded = 0
  static loadZXCVBN() {
    if (SPassword.zxcvbnLoaded) return // don't try loading multiple times
    SPassword.zxcvbnLoaded = 1 // loading in progress
    const script = document.createElement('script')
    script.onload = () => { SPassword.zxcvbnLoaded = 2 } // loading complete
    script.src = '/assets/zxcvbn--62c8cb55.js'
    document.head.appendChild(script)
  }

  constructor() {
    super()
    SPassword.loadZXCVBN()
    this.showMismatch = false
  }

  connectedCallback() {
    // Create the sub-elements.
    this.passwordFields = document.createElement('div')
    this.passwordFields.className = 's-password-fields'
    this.password1 = document.createElement('input')
    this.password1.type = 'password'
    this.password1.className = 'formInput'
    this.password1.autocomplete = 'new-password'
    if (this.hasAttribute('value')) this.password1.value = this.getAttribute('value')
    this.password1.addEventListener('input', this.check.bind(this))
    this.password1.addEventListener('blur', this.blur.bind(this))
    this.passwordFields.appendChild(this.password1)
    this.password2 = document.createElement('input')
    this.password2.type = 'password'
    this.password2.className = 'formInput'
    this.password2.autocomplete = 'new-password'
    this.password2.value = this.password1.value
    this.password2.addEventListener('input', this.check.bind(this))
    this.password2.addEventListener('blur', this.blur.bind(this))
    this.passwordFields.appendChild(this.password2)
    this.appendChild(this.passwordFields)
    this.feedback = document.createElement('div')
    this.feedback.className = 's-password-feedback'
    this.meter = document.createElement('div')
    this.meter.className = 's-password-meter'
    this.meter.style.display = 'none'
    this.feedback.appendChild(this.meter)
    this.message = document.createElement('div')
    this.message.className = 's-password-message'
    this.feedback.appendChild(this.message)
    this.appendChild(this.feedback)
    this.result = document.createElement('input')
    this.result.type = 'hidden'
    this.result.name = this.getAttribute('name')
    this.appendChild(this.result)
    this.hints = (this.getAttribute('hints') || '').split(',')
  }

  // check() is called whenever either password field's content is changed, or
  // when valid() is called.  It verifies the validity of the entry and updates
  // the visual feedback.
  check() {
    // If the two passwords aren't the same, hide the meter and display that
    // error.  However, we don't do this until the blur() handler sets the
    // showMismatch flag.
    if (this.showMismatch && this.password1.value !== this.password2.value) {
      this.meter.style.display = 'none'
      this.message.textContent = document.documentElement.lang === 'es' ? `Las dos contraseñas no son iguales.` : `The two passwords are not the same.`
      this.message.classList.remove('success')
      this.result.value = ''
      return
    }
    // If we don't have a password at all, or the zxcvbn library hasn't finished
    // loading, give no feedback at all.
    if (!this.password1.value || SPassword.zxcvbnLoaded < 2) {
      this.meter.style.display = 'none'
      this.message.textContent = ''
      this.message.classList.remove('success')
      this.result.value = ''
      return
    }
    // Run the zxcvbn analysis on the supplied password and display its quality.
    const analysis = zxcvbn(this.password1.value, this.hints)
    this.meter.style.display = null
    this.meter.textContent = '' // remove all children
    const metercolor = analysis.score < 2 ? 'bad' : analysis.score === 2 ? 'warn' : 'good'
    for (let i = 0; i <= analysis.score; i++) {
      const step = document.createElement('div')
      step.className = `s-password-meter-step ${metercolor}`
      this.meter.appendChild(step)
    }
    if (analysis.score > 2) this.message.classList.add('success')
    else this.message.classList.remove('success')
    this.message.textContent = [
      analysis.feedback.warning,
      ...analysis.feedback.suggestions,
      `This password would take ${analysis.crack_times_display.offline_slow_hashing_1e4_per_second} to crack.`
    ].filter(s => !!s).join('\n')
    // If the password is valid (or the override flag is set), put the password
    // into the form field for submission with the form.  Otherwise, clear the
    // form field.  Note that we check for mismatch between the two fields
    // again here.  Even if a mismatch was ignored above (due to not having set
    // showMismatch), we still don't want to put it into the form field.
    if (this.password1.value === this.password2.value && (analysis.score > 2 || this.hasAttribute('override')))
      this.result.value = this.password1.value
    else this.result.value = ''
  }

  // blur() is called whenever either of the password entry fields loses focus.
  // If both of them have lost focus — i.e., focus has moved away from the whole
  // <s-password> element — we turn on the showMismatch flag and re-check the
  // password.
  blur(evt) {
    if (evt.relatedTarget === this.password1 || evt.relatedTarget === this.password2) return
    this.showMismatch = true
    this.check()
  }

  // The valid property is a read-only boolean indicating whether the field has
  // a valid value for submission.
  get valid() { return this.result.value != '' }
}
customElements.define('s-password', SPassword)
