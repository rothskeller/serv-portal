// When the mouse enters an editable folder item, change its icon to the edit
// icon.
up.on('mouseover', '.folder[editable],.document[editable]', (evt, elm) => {
  const icon = up.element.get(elm, 's-icon')
  if (!icon.hasAttribute('icon-save')) icon.setAttribute('icon-save', icon.getAttribute('icon'))
  icon.setAttribute('icon', 'edit')
})
// When the mouse leaves an editable folder item, change its icon back to its
// original icon.
up.on('mouseout', '.folder[editable],.document[editable]', (evt, elm) => {
  const icon = up.element.get(elm, 's-icon')
  if (icon.hasAttribute('icon-save')) icon.setAttribute('icon', icon.getAttribute('icon-save'))
})
// When someone clicks on the icon of an editable folder, open the folder edit
// dialog in a layer.
up.on('click', '.folder[editable] svg', (evt, elm) => {
  const folder = elm.closest('.folder')
  up.layer.open({
    url: `/folderedit/${folder.dataset.id}`,
    history: false,
    layer: 'new',
    size: 'grow',
    dismissable: 'key',
  })
})
// When someone clicks on the icon of an editable document, open the document
// edit dialog in a layer.
up.on('click', '.document[editable] s-icon', (evt, elm) => {
  const doc = elm.closest('.document')
  const folder = doc.closest('.folderContents').previousElementSibling
  up.layer.open({
    url: `/docedit/${folder.dataset.id}/${doc.dataset.id}`,
    history: false,
    layer: 'new',
    size: 'grow',
    dismissable: 'key',
  })
})
// When someone selects a file in the document edit dialog, and the file name
// field is empty, put the name of the selected file into the file name field.
up.on('input', '#doceditFile', (evt, elm) => {
  let filename = elm.value
  let idx = filename.lastIndexOf('/')
  if (idx >= 0) filename = filename.substr(idx + 1)
  idx = filename.lastIndexOf('\\')
  if (idx >= 0) filename = filename.substr(idx + 1)
  if (!filename) return
  const namefield = document.getElementById('doceditName')
  if (!namefield.value) namefield.value = filename
})
// When a folder item is added with the "folderItem-new" class, which means it
// was just added to the folder, pause a moment and then remove that class.
// This causes the new item background color to fade away.
up.compiler('.folderItem-new', { batch: true }, elms => {
  setTimeout(() => {
    elms.forEach(elm => { elm.classList.remove('folderItem-new') })
  }, 0)
})
  ; (function () {
    // containingFolder returns the .folder element for the folder that contains
    // the supplied element (a .folder or .document).
    function containingFolder(elm) {
      elm = elm.closest('.folderContents')
      if (elm) elm = elm.previousElementSibling
      return elm
    }
    // canDrop returns whether the toparent folder can receive the item that is
    // being dragged (as identified in the supplied event).
    function canDrop(evt, toparent) {
      return evt.dataTransfer.types.some(t => {
        // If the item that is being dragged is a URL or one or more files, any
        // folder can receive those.
        if (t === 'Files' || t === 'text/uri-list') return true
        // If the item being dragged is a document, it can be received as long
        // as it's not being dragged to the folder it's already in.  The folder
        // it's currently in is encoded in the data transfer type.
        const docmatch = t.match(/^x-serv-document;folder=(\d+)$/)
        if (docmatch && docmatch[1] !== toparent.dataset.id) return true
        // If the item being dragged is a folder, it can be received as long as
        // it's not being dragged to the folder it's already in, to itself, or
        // to any descendant of itself.  Its ID and its parent ID are encoded in
        // the data transfer type.
        const foldermatch = t.match(/^x-serv-folder;folder=(\d+);parent=(\d+)$/)
        if (foldermatch) {
          let movingID = foldermatch[1]
          let fromparentID = foldermatch[2]
          // If the destination folder is the parent of the folder that's moving
          // — i.e., nothing's changing — deny the move.
          if (fromparentID === toparent.dataset.id) return false
          // If the destination folder is, or is descended from, the folder
          // that's moving, deny the move.
          let f = toparent
          while (f) {
            if (movingID === f.dataset.id) return false
            f = containingFolder(f)
          }
          return true
        }
        return false
      })
    }
    // canDelete returns whether the trash target can receive the item that is
    // being dragged (as identified in the supplied event).
    function canDelete(evt) {
      return evt.dataTransfer.types.some(t => t.startsWith('x-serv-document;') || t.startsWith('x-serv-folder;'))
    }
    // When we start dragging a folder, set the data transfer type, embedding
    // its ID and its parent's ID.  Also hide the Add buttons and, if the item
    // is deletable, show the trash drop target.
    up.on('dragstart', '.folder[draggable]', (evt, elm) => {
      evt.dataTransfer.setData(`x-serv-folder;folder=${elm.dataset.id};parent=${containingFolder(elm).dataset.id}`, elm.dataset.id)
      document.getElementById('folderButtons').style.display = 'none'
      if (elm.hasAttribute('deletable')) document.getElementById('folderDelete').style.display = null
    })
    // When we start dragging a document, set the data transfer type, embedding
    // its containing folder's ID.  Also hide the Add buttons and show the trash
    // drop target.
    up.on('dragstart', '.document[draggable]', (evt, elm) => {
      evt.dataTransfer.setData(`x-serv-document;folder=${containingFolder(elm).dataset.id}`, elm.dataset.id)
      document.getElementById('folderButtons').style.display = 'none'
      document.getElementById('folderDelete').style.display = null
    })
    // When we stop dragging something, hide the trash drop target and show the
    // Add buttons.
    up.on('dragend', '.folder[draggable], .document[draggable]', () => {
      document.getElementById('folderButtons').style.display = null
      document.getElementById('folderDelete').style.display = 'none'
    })
    // When something is being dragged over a folder that can receive it, tell
    // the browser that by preventing the dragover event.
    up.on('dragover', '.folder[editable]', (evt, elm) => {
      if (canDrop(evt, elm)) evt.preventDefault()
    })
    // When something deletable is being dragged over the trash drop target,
    // tell the browser that by preventing the dragover event.
    up.on('dragover', '#folderDelete', (evt, elm) => {
      if (canDelete(evt)) evt.preventDefault()
    })
    // When dragging into or out of an editable folder (or the icon or link
    // within the .folder element), add or remove the highlight on that folder.
    up.on('dragenter dragleave', '.folder[editable], .folder[editable] *', (evt, elm) => {
      folder = elm.closest('.folder')
      // Don't remove the highlight if the drag location is a different element
      // within the same .folder elemen.t
      if (evt.type === 'dragleave' && evt.relatedTarget && evt.relatedTarget.closest('.folder') === folder) return
      // Don't add the highlight if the target folder can't receive the item
      // being dragged.
      if (evt.type === 'dragenter' && !canDrop(evt, folder)) return
      // Toggle the highlight.
      folder.classList.toggle('folder-dragging', evt.type === 'dragenter')
    })
    // When dragging into or out of the trash drop target (or the icon within
    // it), add or remove the highlight on that target.
    up.on('dragenter dragleave', '#folderDelete, #folderDelete s-icon', (evt, elm) => {
      target = document.getElementById('folderDelete')
      // Don't add the highlight if the target folder can't receive the item
      // being dragged.
      if (evt.type === 'dragenter' && !canDelete(evt)) return
      // Toggle the highlight.
      target.classList.toggle('folder-dragging', evt.type === 'dragenter')
    })
    // Handle a drop of an item onto a folder.
    up.on('drop', '.folder[editable]', (evt, elm) => {
      if (!canDrop(evt, elm)) return
      let form
      evt.dataTransfer.types.some(t => {
        if (t.startsWith('x-serv-document')) {
          // Dropping a document.  Fill in that field of the drop form.
          const field = document.getElementById('folderDropDoc')
          field.value = evt.dataTransfer.getData(t)
          form = field.form
        }
        if (t.startsWith('x-serv-folder')) {
          // Dropping a folder.  Fill in that field of the drop form.
          const field = document.getElementById('folderDropFolder')
          field.value = evt.dataTransfer.getData(t)
          form = field.form
        }
        if (t === 'text/uri-list') {
          // Dropping a URL.  Fill in that field of the drop form.
          const url = evt.dataTransfer.getData('text/uri-list').replace(/\r\n.*/, '')
          const field = document.getElementById('folderDropURL')
          field.value = url
          form = field.form
        }
        if (t === 'Files') {
          // Dropping a set of files.  Fill in that field of the drop form.
          const field = document.getElementById('folderDropFiles')
          field.files = evt.dataTransfer.files
          form = field.form
        }
      })
      if (form) {
        // We have something to drop.  Set the form action based on the folder
        // being dropped, and submit the form.  The submission counts as a
        // browser history change only if it is a new URL.
        evt.preventDefault()
        form.action = elm.dataset.path
        up.submit(form, { history: form.action !== document.getElementById('folderpath').value })
      }
    })
    // Handle a drop of an item onto the trash target.
    up.on('drop', '#folderDelete', (evt, elm) => {
      if (!canDelete(evt)) return
      let form, history
      let action = document.getElementById('folderpath').value
      evt.dataTransfer.types.some(t => {
        if (t.startsWith('x-serv-document')) {
          // Dropping a document.  Fill in that field of the drop form.
          const field = document.getElementById('folderTrashDoc')
          field.value = evt.dataTransfer.getData(t)
          form = field.form
          history = false
        }
        if (t.startsWith('x-serv-folder')) {
          // Dropping a folder.  Fill in that field of the drop form.
          const field = document.getElementById('folderTrashFolder')
          field.value = evt.dataTransfer.getData(t)
          form = field.form
          if (evt.dataTransfer.getData(t) === document.getElementById('folderid').value) {
            action = action.replace(/\/[^/]*$/, '')
            history = true
          }
        }
      })
      if (form) {
        // We have something to drop.  Set the form action based on the item
        // being dropped, and submit the form.  The submission counts as a
        // browser history change only if it is a new URL.
        evt.preventDefault()
        form.action = action
        up.submit(form, { history })
      }
    })
  })()
