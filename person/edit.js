window.addEventListener('load', function () {
  var table = document.getElementById('editPerson-team-table');
  if (table) {
    var inputs = table.getElementsByTagName("input");
    for (var i = 0; i < inputs.length; i++) {
      var input = inputs.item(i);
      input.addEventListener('change', function (evt) {
        var select = document.getElementById(evt.target.id.replace('-team-', '-role-'));
        if (select) select.style.display = evt.target.checked ? 'inline' : 'none';
      });
    }
  }
});
