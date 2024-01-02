package static

import (
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/ui"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

// ListosPage handles GET /listos requests.
func ListosPage(r *request.Request) {
	var user = auth.SessionUser(r, 0, false)
	ui.Page(r, user, ui.PageOpts{Title: r.Loc("Listos California")}, func(main *htmlb.Element) {
		main = main.A("class=static")
		if r.Language != "es" {
			main.R(`<p><b>Listos California</b>
<p>Listos California is a state program, managed by the California Office of
Emergency Services (CalOES), focusing on disaster preparedness education for
California residents.  Under their umbrella, the Sunnyvale Listos program
provides disaster preparedness education in Sunnyvale.
<p>Our flagship offering is our
<a href=/pep up-target=main>Personal Emergency Preparedness</a> class.  This is
a two-hour class that teaches home and family preparedness.  We offer this
class to the general public every 2–3 months, in both English and Spanish.  We
also offer it to neighborhood associations, businesses, etc. when requested.
<p>Our disaster preparedness education efforts also include Outreach booths and
tables at public events (the Arts and Wine Festival, the Diwali Festival, the
Firefighters Pancake Breakfast, neighborhood block parties, etc.).  At these
events, we set up tables and distribute disaster preparedness information to
participants.
<p>For more information about Listos California or our disaster preparedness
education programs, write us at
<a href="mailto:listos@sunnyvale.ca.gov" target="_blank">listos@sunnyvale.ca.gov</a>.
Also write to us if you want to arrange a private preparedness class for your
neighborhood or group, or have a preparedness table at your event.
<div style="margin:1.5rem 0"><button class="sbtn sbtn-primary" onclick="history.back()">Back</button></div>`)
		} else {
			main.R(`<p><b>Listos California</b>
<p>Listos California es un programa estatal, gestionado por la Ofinica de
Servicios de Emergencia de California (CalOES), centrándose en la educación de
preparación para desastres para los residentes de California.  Bajo su escudo,
el programa Listos Sunnyvale proporciona educación de preparación para desastres
en Sunnyvale.
<p>Nuestra oferta estrella es nuestra clase
<a href=/pep up-target=main>Preparación para desastres y emergencias</a>.
Se trata de una clase de dos horas que enseña la preparación del hogar y la
familia.  Ofrecemos esta clase al público en general cada 2-3 meses, tanto en
inglés como en español.  También la ofrecemos a asociaciones de vecinos,
empresas, etc. cuando lo solicitan.
<p>Nuestra labor de educación sobre la preparación ante desastres también
incluye puestos y mesas de divulgación en eventos públicos (el Festival de las
Artes y el Vino, el Festival Diwali, el Desayuno de Panqueques de los Bomberos,
fiestas vecinales, etc.).  En estos eventos, instalamos mesas y distribuimos
información sobre preparación ante desastres a los participantes.
<p>Para más información sobre Listos California o nuestros programas educativos
de preparación ante desastres, escríbanos a
<a href="mailto:listos@sunnyvale.ca.gov" target="_blank">listos@sunnyvale.ca.gov</a>.
También escríbanos si desea organizar una clase privada de preparación para su
vecindario o grupo, o tener una mesa de preparación en su evento.
<div style="margin:1.5rem 0"><button class="sbtn sbtn-primary" onclick="history.back()">Regresar</button></div>`)
		}
	})
}
