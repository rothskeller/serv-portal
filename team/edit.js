window.addEventListener('load', function () {
  var table = document.getElementById('editTeam-privileges-table');
  if (table) {
    var inputs = table.getElementsByTagName('input');
    for (var i = 0; i < inputs.length; i++) {
      inputs.item(i).addEventListener('change', function (evt) {
        var split = evt.target.name.split('-');
        if (evt.target.checked) {
          switch (split[0]) {
            case 'admin':
              evt.target.form['viewer-' + split[1]].checked = true;
              break;
            case 'manager':
              evt.target.form['viewer-' + split[1]].checked = true;
              evt.target.form['admin-' + split[1]].checked = true;
              break;
          }
        } else {
          switch (split[0]) {
            case 'viewer':
              evt.target.form['admin-' + split[1]].checked = false;
              evt.target.form['manager-' + split[1]].checked = false;
              break;
            case 'admin':
              evt.target.form['manager-' + split[1]].checked = false;
              break;
          }
        }
      });
    }
  }
  table = document.getElementById('editTeam-inPrivs-table');
  if (table) {
    var inputs = table.getElementsByTagName('input');
    for (var i = 0; i < inputs.length; i++) {
      inputs.item(i).addEventListener('change', function (evt) {
        var split = evt.target.name.split('-');
        if (evt.target.checked) {
          switch (split[0]) {
            case 'inadmin':
              evt.target.form['inviewer-' + split[1]].checked = true;
              break;
            case 'inmanager':
              evt.target.form['inviewer-' + split[1]].checked = true;
              evt.target.form['inadmin-' + split[1]].checked = true;
              break;
          }
        } else {
          switch (split[0]) {
            case 'inviewer':
              evt.target.form['inadmin-' + split[1]].checked = false;
              evt.target.form['inmanager-' + split[1]].checked = false;
              break;
            case 'inadmin':
              evt.target.form['inmanager-' + split[1]].checked = false;
              break;
          }
        }
      });
    }
  }
});
