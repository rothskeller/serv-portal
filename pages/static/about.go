package static

import (
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/ui"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

// AboutPage handles GET /about requests.
func AboutPage(r *request.Request) {
	ui.Page(r, auth.SessionUser(r, 0, false), ui.PageOpts{}, func(main *htmlb.Element) {
		main = main.A("class=static")
		if r.Language != "es" {
			main.R(`<h1>Privacy Policy</h1>
<p>
  This web site collects information about people who work for, volunteer for,
  or take classes organized through the Office of Emergency Services (OES) in
  the Sunnyvale Department of Public Safety (DPS).  The information we collect
  includes:
</p>
<ul>
  <li>Basic Information
  <ul>
    <li>name
    <li>amateur radio call sign
    <li>contact information (email addresses, phone numbers, and physical and postal addresses)
    <li>memberships in, and roles held in, SERV volunteer groups
    <li>emergency response classes taken and certificates issued
    <li>credentials that are relevant to SERV operations
    <li>other information voluntarily provided such as skills, languages spoken, available equipment, etc.
  </ul>
  <li>Restricted Information
  <ul>
    <li>attendance at SERV events, and hours spent at them
    <li>Disaster Service Worker registration status
    <li>photo IDs and card access keys issued
    <li>Live Scan fingerprinting success, with date (see note below)
    <li>background check success, with date (see note below)
  </ul>
  <li>Targeted Information
  <ul>
    <li>email messages sent to any SunnyvaleSERV.org address
    <li>text messages sent through this web site
  </ul>
  <li>Private Information
  <ul>
    <li>logs of web site visits and actions taken
    <li>preferred language
  </ul>
</ul>
<p>
  All of the above information is available to the paid and volunteer staff of
  OES and their delegates, including the web site maintainers.  Private
  information is not available to anyone else.
<p>
  If you are a student in an OES-organized class, such as CERT, Listos, or
  PEP, your basic and restricted information may be shared with the class
  instructors as long as the class is in progress.
<p>
  If you are a volunteer in a SERV volunteer group, your basic information may
  be shared with other volunteers in that group, and your restricted
  information may be shared with the leaders of that group.
<p>
  If you are a volunteer in a SERV volunteer group, and you have successfully
  completed Live Scan fingerprinting and/or background checks, that fact (with
  no detail other than the date) may be shared with the leaders of your
  volunteer group.  A negative result will not be shared with them.
<p>
  If you have sent any email or text messages (targeted information) through
  the site, they may be shared with any member of the group(s) to which you
  sent them, including members who join those groups after you send the
  messages.
<p>
  If you volunteer for mutual aid or training with another emergency response
  organization or jurisdiction, we may share your basic and/or restricted
  information with them.
<p>
  The OES staff may share anonymized, aggregate data derived from the above
  information with anyone at their discretion.
</p>
<h1>Cookies</h1>
<p>
  This site uses browser cookies.  While you are logged in, a browser cookie
  contains your session identification; this cookie goes away when you log out
  or your login session expires.  More permanent cookies are used to store
  some of your user interface preferences, such as your preferred language and
  whether you prefer to see the events page in calendar or list form.  No
  personally identifiable information is ever stored in browser cookies.
</p>
<h1>Credits and Copyrights</h1>
<p>
  This site was developed by Steven Roth, as a volunteer for the Sunnyvale
  Department of Public Safety.  The site software is copyrighted © 2020–2021
  by Steven Roth.  Steven Roth has granted the Sunnyvale Department of Public
  Safety a non-exclusive, perpetual, royalty-free, worldwide license to use
  this software.  The Sunnyvale Department of Public Safety owns the
  SunnyvaleSERV.org domain and funds the ongoing usage and maintenance of the
  site.
</p>
<h1>Technologies and Services</h1>
<p>
  The software for this web site is written in
  <a href="https://golang.org" target="_blank">Go</a>, with data storage in a
  <a href="https://sqlite.org" target="_blank">SQLite</a> database.  This web
  site is hosted by
  <a href="https://www.dreamhost.com/" target="_blank">Dreamhost</a>.  It uses
  <a href="https://www.google.com/maps" target="_blank">Google Maps</a> for
  geolocation and mapping,
  <a href="https://www.twilio.com/" target="_blank">Twilio</a> for text
  messaging, and
  <a href="https://www.algolia.com/" target="_blank">Algolia</a> for searching.
</p>
<div style="margin:1.5rem 0"><button class="sbtn sbtn-primary" onclick="history.back()">Back</button></div>
`)
		} else {
			main.R(`<h1>Política de privacidad</h1>
<p>
  Este sitio web recopila información sobre personas que trabajan, son
  voluntarios, o tomar clases organizadas a través de la Oficina de Servicios de
  Emergencia (OSE) en el Departamento de Seguridad Pública de Sunnyvale (DSP).
  La información que recopilamos incluye:
</p>
<ul>
  <li>Información basica:
  <ul>
    <li>nombre
    <li>indicativo de radioaficionado
    <li>información de contacto (direcciones de correo electrónico, números de teléfono y direcciones físicas y postales)
    <li>membresías y roles desempeñados en grupos de voluntarios de SERV
    <li>clases de respuesta a emergencias tomadas y certificados emitidos
    <li>credenciales que son relevantes para las operaciones de SERV
    <li>otra información proporcionada voluntariamente como habilidades, idiomas hablados, equipos disponibles, etc.
  </ul>
  <li>Información restringida
  <ul>
    <li>asistencia a eventos de SERV y horas dedicadas a ellos
    <li>estado de registro como trabajador de servicios de desastre
    <li>identificaciones con fotografía y claves de acceso emitidas
    <li>éxito de la toma de huellas digitales de Live Scan, con fecha (consulte la nota a continuación)
    <li>éxito de la verificación de antecedentes, con fecha (consulte la nota a continuación)
  </ul>
  <li>Información dirigida
  <ul>
    <li>mensajes enviados a cualquier dirección de SunnyvaleSERV.org
    <li>mensajes de texto enviados a través de este sitio web
  </ul>
  <li>Información privada
  <ul>
    <li>registros de visitas al sitio web y acciones realizadas
    <li>idioma preferido
  </ul>
</ul>
<p>
  Toda la información anterior está disponible para el personal remunerado y
  voluntario de OSE y sus delegados, incluidos los mantenedores del sitio web.
  La información privada no está disponible para nadie más.
<p>
  Si es estudiante de una clase organizada por OSE, como CERT, Listos o PPDE, su
  su información básica y restringida puede ser compartida con los instructores
  mientras la clase esté en curso.
<p>
  Si es voluntario en un grupo de voluntarios de SERV, su información básica
  puede ser compartido con otros voluntarios en ese grupo, y su información
  restringida puede ser compartida con los líderes de ese grupo.
<p>
  Si es voluntario en un grupo de voluntarios de SERV y ha logrado completado la
  toma de huellas digitales de Live Scan y/o verificaciones de antecedentes, ese
  hecho (con ningún detalle más que la fecha) puede ser compartido con los
  líderes de su grupo de voluntarios.  Un resultado negativo no será compartido
  con ellos.
<p>
  Si ha enviado algún correo electrónico o mensaje de texto (información
  dirigida) a través de el sitio, pueden ser compartidos con cualquier miembro
  del grupo(s) al cual usted los envió, incluidos los miembros que se unen a
  esos grupos después de que usted envió los mensajes.
<p>
  Si se ofrece como voluntario para ayuda mutua o capacitación con otra
  respuesta de emergencia organización o jurisdicción, podemos compartir su
  información básica y/o restringida información con ellos.
<p>
  El personal de OSE puede compartir datos agregados anonimizados derivados de
  la información anterior con cualquier persona a su discreción.
</p>
<h1>Cookies</h1>
<p>
  Este sitio utiliza cookies del navegador.  Mientras está conectado, una cookie
  del navegador contiene su identificación de sesión; Esta cookie desaparece
  cuando cierra o caduca la sesión.  Se utilizan cookies más permanentes para
  almacenar algunas de sus preferencias de interfaz de usuario, como su idioma
  preferido y si prefiere ver la página de eventos en forma de calendario o de
  lista.  Nunca se almacena información de identificación personal en las
  cookies del navegador.
</p>
<h1>Créditos y derechos de autor</h1>
<p>
  Este sitio fue desarrollado por Steven Roth, como voluntario del Departamento
  de Seguridad Pública de Sunnyvale. El software del sitio tiene derechos de
  autor © 2020–2021 por Steven Roth.  Steven Roth ha concedido al Departamento
  de Seguridad Pública de Sunnyvale una licencia mundial, no exclusiva, perpetua
  y libre de regalías de uso de este software. El Departamento de Seguridad
  Pública de Sunnyvale es propietario del dominio SunnyvaleSERV.org y financia
  el uso y mantenimiento continuo del sitio.
</p>
<h1>Tecnologías y servicios</h1>
<p>
  El software de este sitio web está escrito en
  <a href="https://golang.org" target="_blank">Go</a>, con almacenamiento de
  datos en un base de datos
  <a href="https://sqlite.org" target="_blank">SQLite</a>.  Esta sitio web está
  alojado por
  <a href="https://www.dreamhost.com/" target="_blank">Dreamhost</a>.  Usa
  <a href="https://www.google.com/maps" target="_blank">Google Maps</a> para
  geolocalización y cartografía,
  <a href="https://www.twilio.com/" target="_blank">Twilio</a> para mensajes de
  texto, y
  <a href="https://www.algolia.com/" target="_blank">Algolia</a> para buscar.
</p>
<div style="margin:1.5rem 0"><button class="sbtn sbtn-primary" onclick="history.back()">Regrese</button></div>
`)
		}
	})
}
