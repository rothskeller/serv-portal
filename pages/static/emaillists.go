package static

import (
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/ui"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

// EmailListsPage handles GET /email-lists requests.
func EmailListsPage(r *request.Request) {
	var user *person.Person

	if user = auth.SessionUser(r, 0, true); user == nil {
		return
	}
	ui.Page(r, user, ui.PageOpts{Title: r.LangString("SERV Email Lists", "Listas de correo electrónico de SERV")}, func(main *htmlb.Element) {
		main.A("class=static").R(r.LangString(`<h1>SERV Email Lists</h1>
<p>
  The SunnyvaleSERV.org site offers a number of email distribution lists.
  We have one for each volunteer program, that we give out to the general
  public who might want more information about the program.  Email sent to
  these lists is delivered to designated public contact people for each
  program:
</p>
<ul class="emaillist">
  <li>cert@sunnyvaleserv.org
  <li>listos@sunnyvaleserv.org
  <li>sares@sunnyvaleserv.org
  <li>snap@sunnyvaleserv.org
</ul>
<p>There are also lists for the volunteers on each of our teams:</p>
<ul class="emaillist">
  <li>cert-alpha@sunnyvaleserv.org
  <li>cert-committee@sunnyvaleserv.org
  <li>listos-team@sunnyvaleserv.org
  <li>outreach-team@sunnyvaleserv.org
  <li>sares-active@sunnyvaleserv.org
  <li>sares-leads@sunnyvaleserv.org
  <li>snap-team@sunnyvaleserv.org
</ul>
<p>and for the students in each CERT class:</p>
<ul class="emaillist">
  <li>cert-60@sunnyvaleserv.org
  <li>cert-61@sunnyvaleserv.org
  <li>cert-62@sunnyvaleserv.org
  <li style="font:inherit">etc.
</ul>
<p>Finally, there are some broader lists for special purposes:</p>
<ul class="emaillist">
  <li>serv-all@sunnyvaleserv.org
  <li>volunteer-hours@sunnyvaleserv.org
</ul>
<p>
  All of these email lists have restricted access.  For the team lists, only
  members of the team can send mail to them; for the class lists, only the
  instructors can send mail to them; and for the broader lists, only DPS staff
  can send mail to them.  Any mail sent to any of our lists from someone else
  is held for approval before being routed to the list.  Messages on topics
  unrelated to SERV will generally be rejected.
<p>
  If you are receiving email from one of these lists that you do not want,
  there is an “unsubscribe” link at the bottom of every email.  If you are
  receiving email at the wrong address, you can change your email address in
  the “Profile” section of this web site.
</p>
<div style="margin:1.5rem 0"><button class="sbtn sbtn-primary" onclick="history.back()">Back</button></div>`,

			`<h1>Listas de correo electrónico de SERV</h1>
<p>
  El sitio SunnyvaleSERV.org ofrece varias listas de distribución de correo
  electrónico.  Tenemos uno para cada programa de voluntariado, que entregamos
  al público general que pueda querer más información sobre el programa.  Estas
  listas se entregan a las personas de contacto público designadas para cada
  programa:
</p>
<ul class="emaillist">
  <li>cert@sunnyvaleserv.org
  <li>listos@sunnyvaleserv.org
  <li>sares@sunnyvaleserv.org
  <li>snap@sunnyvaleserv.org
</ul>
<p>También hay listas de los voluntarios de cada uno de nuestros equipos:</p>
<ul class="emaillist">
  <li>cert-alpha@sunnyvaleserv.org
  <li>cert-committee@sunnyvaleserv.org
  <li>listos-team@sunnyvaleserv.org
  <li>outreach-team@sunnyvaleserv.org
  <li>sares-active@sunnyvaleserv.org
  <li>sares-leads@sunnyvaleserv.org
  <li>snap-team@sunnyvaleserv.org
</ul>
<p>y para los estudiantes en cada clase CERT:</p>
<ul class="emaillist">
  <li>cert-60@sunnyvaleserv.org
  <li>cert-61@sunnyvaleserv.org
  <li>cert-62@sunnyvaleserv.org
  <li style="font:inherit">etc.
</ul>
<p>Finalmente, existen algunas listas más amplias para propósitos especiales:</p>
<ul class="emaillist">
  <li>serv-all@sunnyvaleserv.org
  <li>volunteer-hours@sunnyvaleserv.org
</ul>
<p>
  Todas estas listas de correo electrónico tienen acceso restringido.  Para las
  listas de equipos, solo los miembros del equipo pueden enviarles correo; para
  las listas de clases, sólo el los instructores pueden enviarles correo; y para
  las listas más amplias, solo el personal de DSP puedo enviarles correo.
  Cualquier correo enviado a cualquiera de nuestras listas por parte de otra
  persona se retiene para su aprobación antes de ser enviado a la lista.
  Mensajes sobre temas que no estén relacionados con SERV generalmente serán
  rechazados.
<p>
  Si recibe un correo electrónico de una de estas listas que no desea, hay un
  enlace para "cancelar suscripción" en la parte inferior de cada correo
  electrónico. Si recibe un correo electrónico en la dirección incorrecta, puede
  cambiar su dirección de correo electrónico en la sección “Perfil” de este
  sitio web.
</p>
<div style="margin:1.5rem 0"><button class="sbtn sbtn-primary" onclick="history.back()">Regrese</button></div>`))
	})
}
