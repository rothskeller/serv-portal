package static

import (
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/ui"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

// ContactUsPage handles GET /contact requests.
func ContactUsPage(r *request.Request) {
	var user = auth.SessionUser(r, 0, false)
	ui.Page(r, user, ui.PageOpts{Title: r.Loc("Contact Us")}, func(main *htmlb.Element) {
		main = main.A("class=static")
		if r.Language != "es" {
			main.R(`<p>Sunnyvale Emergency Response Volunteers
(SERV) is the volunteer arm of the Sunnyvale Office of Emergency Services,
which is part of the city’s Department of Public Safety.
<blockquote><a href="mailto:oes@sunnyvale.ca.gov" target="_blank">oes@sunnyvale.ca.gov</a><br>
(408) 730–7190 English (messages only)<br>
(408) 730-7294 Spanish (messages only)</blockquote>
<p>Our offices are at
<blockquote>Sunnyvale Public Safety Headquarters<br>
700 All America Way<br>
Sunnyvale, CA  94086</blockquote>
<div style="margin:1.5rem 0"><button class="sbtn sbtn-primary" onclick="history.back()">Back</button></div>`)
		} else {
			main.R(`<p>Voluntarios de Respuesta a Emergencias de
Sunnyvale (SERV, por siglas en inglés) es el brazo voluntario de la Oficina de
Servicios de Emergencia de Sunnyvale, que forma parte del Departamento de
Seguridad Pública de la ciudad.
<blockquote><a href="mailto:oes@sunnyvale.ca.gov" target="_blank">oes@sunnyvale.ca.gov</a><br>
(408) 730–7190 en inglés (mensajes solamente)<br>
(408) 730-7294 en español (mensajes solamente)</blockquote>
<p>Nuestra oficina está en
<blockquote>Jefatura de Seguridad Pública de Sunnyvale<br>
700 All America Way<br>
Sunnyvale, CA  94086</blockquote>
<div style="margin:1.5rem 0"><button class="sbtn sbtn-primary" onclick="history.back()">Regresar</button></div>`)
		}
	})
}
