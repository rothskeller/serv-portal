// s-seltree is a custom control for selecting one or more items out of a tree
// of items.  (It is primarily used for selecting roles for people and tasks,
// and in the editing of role implications.)
//
// Visually, it appears to be a tree-structured list, with typical expand and
// collapse controls at each level.  On initial display, it is minimally
// expanded to show the initial selections.
//
// Selections can roll either up or down.  When the "rollup" attribute is not
// set, they roll down.  Selecting a node implicitly selects every descendant
// node in the tree.  The descendants are visually pruned so that the selected
// node looks like a leaf.  The selected node's ID is included in the value of
// the control, but the descendant nodes' IDs are not.  In roll down mode,
// deselecting a node restores its descendants to visibility, with none of them
// selected.
//
// When the "rollup" attribute is set, selections role up.  Selecting a node
// automatically selects every ancestor node.  Only the manually selected node's
// ID is included in the value of the control; the ancestors' IDs are not.  In
// this mode, deselecting a node deselects all of its descendants.
//
// Nodes with the same ID and label can appear in multiple places in the tree.
// Selecting or deselecting any one of them is equivalent to selecting or
// deselecting all of them.
//
// The set of nodes available to choose from are encoded in the text content of
// the control, in a line-oriented format.  Each line's indentation determines
// the node hierarchy.  If the first character after the indentation is a '-',
// the node is disabled and the user cannot select it (although it can still be
// selected by rollup).  After the indentation and optional '-', the first word
// on the line is the node ID and the remainder of the line is the node label.
//
// Attributes:
//
// name - parameter name under which the control's value is submitted with the
//     form (required if form submission is expected, otherwise optional)
// value - value of the control: a whitespace-separated, unordered list of
//     selected node IDs
// rollup - determines whether selections roll up or down (see above)
//
// Events:
//
// change - triggered when any node is selected or deselected
//
class SSelTree extends HTMLElement {
  static get observedAttributes() { return ['value'] };
  connectedCallback() {
    if (!this.firstChild) {
      new MutationObserver((list, observer) => {
        this.initialSetup()
        observer.disconnect()
      }).observe(this, { childList: true })
    } else if (this.firstChild && this.firstChild.nodeType === Node.TEXT_NODE) {
      this.initialSetup()
    } else {
      if (this._in) this._in.name = this.getAttribute('name')
      this.setValue(this.getAttribute('value'))
      this._rollup = this.hasAttribute('rollup')
    }
  }
  attributeChangedCallback(name, oldValue, newValue) {
    if (name === 'value') this.setValue(newValue)
  }
  initialSetup() {
    this.parseTree()
    this._tree.forEach(node => { this.createTreeNode(this, node) })
    this._sel = {}
    if (this.hasAttribute('name')) {
      const name = this.getAttribute('name')
      if (name && !this._in) {
        this._in = document.createElement('input')
        this._in.type = 'hidden'
        this._in.name = name
        this.appendChild(this._in)
      }
    }
    this.setValue(this.getAttribute('value'))
    this._rollup = this.hasAttribute('rollup')
  }
  parseTree() {
    const lines = this.firstChild.nodeValue.split(/\r?\n/)
    this.removeChild(this.firstChild)
    this._tree = []
    let stack = []
    lines.forEach(line => {
      const indent = line.search(/[^ ]/)
      if (indent < 0) return
      if (indent > stack.length + 1) throw ('invalid tree data')
      if (indent === stack.length + 1) stack.push(null)
      line = line.substring(indent)
      const disabled = line.startsWith('-')
      if (disabled) line = line.substring(1)
      const sep = line.indexOf(' ')
      if (indent < 0) throw ('invalid tree data')
      const id = line.substring(0, sep)
      const label = line.substring(sep + 1)
      const node = { id, label, disabled, children: [] }
      if (indent === 0) {
        this._tree.push(node)
        stack = [node]
      } else {
        stack[indent - 1].children.push(node)
      }
      stack[indent] = node
      stack.splice(indent + 1)
    })
  }
  createTreeNode(parent, node) {
    const div = document.createElement('div')
    div.className = `s-seltree-node s-seltree--${node.id}`
    if (node.disabled) div.className += ' disabled'
    if (node.children.length) {
      const iconRight = document.createElement('s-icon')
      iconRight.className = 's-seltree-right'
      iconRight.setAttribute('icon', 'chevron-right')
      iconRight.addEventListener('click', this.onIconRight)
      div.appendChild(iconRight)
      const iconDown = document.createElement('s-icon')
      iconDown.className = 's-seltree-down'
      iconDown.setAttribute('icon', 'chevron-down')
      iconDown.addEventListener('click', this.onIconDown)
      div.appendChild(iconDown)
    }
    const label = document.createElement('div')
    label.className = 's-seltree-label'
    const span = document.createElement('span')
    span.textContent = node.label
    label.appendChild(span)
    label.addEventListener('click', evt => { this.onLabelClick(evt, node) })
    div.appendChild(label)
    if (node.children.length) {
      const kidbox = document.createElement('div')
      kidbox.className = 's-seltree-kids'
      node.children.forEach(kid => { this.createTreeNode(kidbox, kid) })
      div.appendChild(kidbox)
    }
    parent.appendChild(div)
  }
  setValue(value) {
    const ids = {}
    if (value) value.split(' ').forEach(v => { ids[v] = true })
    for (let id in this._sel) {
      if (!ids[id]) this.deselect(id)
    }
    for (let id in ids) {
      if (!this._sel[id]) this.select(id)
    }
    if (this._in) this._in.value = value
  }
  select(id) {
    this._sel[id] = true
    this.querySelectorAll(`.s-seltree--${id}`).forEach(div => {
      div.classList.add('selected')
      if (this._rollup) {
        div = div.parentNode.closest('.s-seltree-node')
        while (div) {
          div.classList.add('rolledup', 'open')
          div.setAttribute('rollup-count', (parseInt(div.getAttribute('rollup-count')) || 0) + 1)
          div = div.parentNode.closest('.s-seltree-node')
        }
      } else {
        div.classList.add('hidekids')
        div = div.parentNode.closest('.s-seltree-node')
        while (div) {
          div.classList.add('open')
          div = div.parentNode.closest('.s-seltree-node')
        }
      }
    })
  }
  deselect(id) {
    delete this._sel[id]
    this.querySelectorAll(`.s-seltree--${id}`).forEach(div => {
      div.classList.remove('selected')
      if (this._rollup) {
        div = div.parentNode.closest('.s-seltree-node')
        while (div) {
          let count = parseInt(div.getAttribute('rollup-count')) || 0
          if (count > 0) count--
          if (count === 0) div.classList.remove('rolledup')
          div.setAttribute('rollup-count', count)
          div = div.parentNode.closest('.s-seltree-node')
        }
      } else {
        div.classList.remove('hidekids')
      }
    })
  }
  onIconRight(evt) {
    evt.target.closest('.s-seltree-node').classList.add('open')
  }
  onIconDown(evt) {
    evt.target.closest('.s-seltree-node').classList.remove('open')
  }
  onLabelClick(evt, node) {
    if (node.disabled) return
    let value
    if (this._sel[node.id]) {
      value = Object.keys(this._sel).filter(id => id !== node.id).join(' ')
    } else {
      value = Object.keys(this._sel).join(' ') + ` ${node.id}`
    }
    this.setAttribute('value', value)
    if (this._in) this._in.value = value
    this.dispatchEvent(new Event('change'))
  }
}
customElements.define('s-seltree', SSelTree)
