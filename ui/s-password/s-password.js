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
  static crackTimeMessage(ct) {
    const lang = document.documentElement.lang
    let message = lang == 'es' ? 'Esta contraseña tardaría ' : 'This password would take '
    if (ct < 1) {
      message += lang == 'es' ? 'menos que un segundo' : 'less than a second'
    } else if (ct < 60) {
      ct = Math.round(ct)
      if (ct === 1) message += lang == 'es' ? '1 segundo' : '1 second'
      else message += lang == 'es' ? `${ct} segundos` : `${ct} seconds`
    } else if (ct < 60 * 60) {
      ct = Math.round(ct / 60)
      if (ct === 1) message += lang == 'es' ? '1 minuto' : '1 minute'
      else message += lang == 'es' ? `${ct} minutos` : `${ct} minutes`
    } else if (ct < 60 * 60 * 24) {
      ct = Math.round(ct / (60 * 60))
      if (ct === 1) message += lang == 'es' ? '1 hora' : '1 hour'
      else message += lang == 'es' ? `${ct} horas` : `${ct} hours`
    } else if (ct < 60 * 60 * 24 * 31) {
      ct = Math.round(ct / (60 * 60 * 24))
      if (ct === 1) message += lang == 'es' ? '1 día' : '1 day'
      else message += lang == 'es' ? `${ct} días` : `${ct} days`
    } else if (ct < 60 * 60 * 24 * 31 * 12) {
      ct = Math.round(ct / (60 * 60 * 24 * 31))
      if (ct === 1) message += lang == 'es' ? '1 mes' : '1 month'
      else message += lang == 'es' ? `${ct} meses` : `${ct} months`
    } else if (ct < 60 * 60 * 24 * 31 * 12 * 100) {
      ct = Math.round(ct / (60 * 60 * 24 * 31 * 12))
      if (ct === 1) message += lang == 'es' ? '1 año' : '1 year'
      else message += lang == 'es' ? `${ct} años` : `${ct} years`
    } else {
      ct = Math.round(ct / (60 * 60 * 24 * 31 * 12 * 100))
      if (ct === 1) message += lang == 'es' ? '1 siglo' : '1 century'
      else message += lang == 'es' ? `${ct} siglos` : `${ct} centuries`
    }
    message += lang == 'es' ? ' en descifrarse.' : ' to crack.'
    return message
  }
  static translateFeedback(m) {
    if (document.documentElement.lang != 'es') return m
    const s = {
      'A word by itself is easy to guess': 'Una palabra por sí sola es fácil de adivinar',
      'Add another word or two. Uncommon words are better.': 'Añada una o dos palabras más. Las palabras poco comunes son mejores.',
      'All-uppercase is almost as easy to guess as all-lowercase': 'Todas las mayúsculas son casi tan fáciles de adivinar como todas las minúsculas',
      'Avoid dates and years that are associated with you': 'Evite fechas y años que estén asociados con usted',
      'Avoid recent years': 'Evite años recientes',
      'Avoid repeated words and characters': 'Evite palabras y caracteres repetidos',
      'Avoid sequences': 'Evite secuencias',
      'Avoid years that are associated with you': 'Evite años que se asocien con usted',
      'Capitalization doesn\'t help very much': 'Las mayúsculas no ayudan mucho',
      'Common names and surnames are easy to guess': 'Los nombres y apellidos comunes son fáciles de adivinar',
      'Dates are often easy to guess': 'Las fechas suelen ser fáciles de adivinar',
      'Names and surnames by themselves are easy to guess': 'Los nombres y apellidos por sí solos son fáciles de adivinar',
      'No need for symbols, digits, or uppercase letters': 'No necesita símbolos, dígitos ni mayúsculas',
      'Predictable substitutions like \'@\' instead of \'a\' don\'t help very much': 'Las sustituciones predecibles como "@" en lugar de "a" no ayudan mucho.',
      'Recent years are easy to guess': 'Los años recientes son fáciles de adivinar',
      'Repeats like "aaa" are easy to guess': 'Las repeticiones como "aaa" son fáciles de adivinar',
      'Repeats like "abcabcabc" are only slightly harder to guess than "abc"': 'Las repeticiones como "abcabcabc" son un poco más difíciles de adivinar que "abc".',
      'Reversed words aren\'t much harder to guess': 'Las palabras invertidas no son mucho más difíciles de adivinar',
      'Sequences like abc or 6543 are easy to guess': 'Secuencias como "abc" o "6543" son fáciles de adivinar.',
      'Short keyboard patterns are easy to guess': 'Los patrones de teclado cortos son fáciles de adivinar',
      'Straight rows of keys are easy to guess': 'Las filas rectas de teclas son fáciles de adivinar',
      'This is a top-10 common password': 'Esta es una de las 10 contraseñas más comunes',
      'This is a top-100 common password': 'Esta es una de las 100 contraseñas más comunes',
      'This is a very common password': 'Esta es una contraseña muy común',
      'This is similar to a commonly used password': 'Es similar a una contraseña común',
      'Use a few words, avoid common phrases': 'Use pocas palabras, evite frases comunes',
      'Use a longer keyboard pattern with more turns': 'Utilice un patrón de teclado más largo y con más vueltas.',
    }[m]
    return s || m
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
      SPassword.crackTimeMessage(analysis.crack_times_seconds.offline_slow_hashing_1e4_per_second)
    ].filter(s => !!s).map(s => SPassword.translateFeedback(s)).join('\n')
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
