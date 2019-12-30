window.addEventListener('load', function () {
  var team = document.getElementById('listPeople-team');
  if (team) team.addEventListener('change', function () {
    document.getElementById('listPeople-title').submit();
  });
});
