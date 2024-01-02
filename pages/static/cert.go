package static

import (
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/ui"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

// CERTPage handles GET /cert requests.
func CERTPage(r *request.Request) {
	var user = auth.SessionUser(r, 0, false)
	ui.Page(r, user, ui.PageOpts{Title: r.Loc("Sunnyvale CERT")}, func(main *htmlb.Element) {
		main = main.A("class=static")
		if r.Language != "es" {
			main.R(`<p><b>Community Emergency Response Team (CERT)</b>
<p>CERT is a nationwide program, managed by the Federal Emergency Management
Agency (FEMA), that prepares residents to care for themselves and their
communities during and after major disasters.  Its emphasis is on training
residents to be able to respond safely and effectively during an emergency.
<p>The CERT program was created by the Los Angeles Fire Department
after seeing the significant loss of life of volunteer rescuers in the 1985
Mexico City earthquake.  Volunteers are credited with having saved many lives in
the aftermath of that earthquake, but many of the volunteers were killed because
they did not know how to keep themselves safe while doing such work.  LAFD
created the CERT program to ensure that the same thing didn't happen on their
watch.  The 1987 Whittier earthquake near Los Angeles underscored the value of
this program.  It the early 1990s, FEMA expanded the program to cover other
disasters besides earthquakes, and spread it nationwide.
<p>In Sunnyvale, we teach the FEMA-standard
<a href="/cert-basic" up-target=main>CERT Basic Training<a> class, with some
local enhancements, to anyone who wants it.
This is a 30-hour class, taught over seven weeks, covering all aspects of
volunteer disaster response.  For the graduates of that class, we also teach
occasional refresher classes on specific CERT topics to help our volunteers
keep their skills and knowledge fresh.
<p>Sunnyvale also has a “CERT Deployment Team.”  This is a group of CERT-trained
volunteers who have agreed to be on call to assist the professional responders
in the Department of Public Safety when needed.  Our CERT Deployment Team
receives additional, monthly training covering both the CERT topics and
more advanced public safety skills.
<p>For more information about our CERT program, write to
<a href="mailto:cert@sunnyvale.ca.gov" target="_blank">cert@sunnyvale.ca.gov</a>.
<div style="margin:1.5rem 0"><button class="sbtn sbtn-primary" onclick="history.back()">Back</button></div>`)
		} else {
			main.R(`<p><b>Equipo Communitario de Respuesta a Emergencias (CERT)</b>
<p>CERT es un programa nacional, gestionado por la Agencia Federal para la
Gestión de Emergencias (FEMA), que prepara a los residentes para cuidarse a sí
mismos y a sus comunidades durante y después de grandes catástrofes.  Su énfasis
está en formar a los residentes para que sean capaces de responder con seguridad
y eficacia durante una emergencia.
<p>El programa CERT fue creado por el Departamento de Bomberos de Los Ángeles
tras comprobar la grande pérdida de vidas de rescatadores voluntarios en el
terremoto de la Ciudad de México de 1985.  A los voluntarios se les atribuye
haber salvado muchas vidas tras ese terremoto, pero muchos de los voluntarios
murieron porque no sabían cómo mantenerse a salvo mientras realizaban ese
trabajo.  El LAFD creó el programa CERT para asegurarse de que no ocurriera lo
mismo durante su guardia.  El terremoto de Whittier, cerca de Los Ángeles, en
1987, puso de manifiesto el valor de este programa.  A principios de los 90, la
FEMA amplió el programa para cubrir otras catástrofes además de los terremotos,
y lo extendió por todo el país.
<p>En Sunnyvale, impartimos el curso
<a href="/cert-basic" up-target=main>Capacitación básica del CERT</a> estándar
de la FEMA, con algunas mejoras locales, a todo aquel que lo desee. Se trata de
una clase de 30 horas, impartida a lo largo de siete semanas, que cubre todos
los aspectos de la respuesta voluntaria en caso de catástrofe.  Para los
graduados de esa clase, también impartimos clases ocasionales de actualización
sobre temas específicos del CERT para ayudar a nuestros voluntarios a mantener
sus habilidades y conocimientos al día.
<p>Sunnyvale también tiene un “Equipo de Despliegue CERT”.  Este es un grupo de
voluntarios entrenados en CERT que han acordado estar de guardia para ayudar a
los respondedores profesionales en el Departamento de Seguridad Pública cuando
sea necesario.  Nuestro Equipo de Despliegue CERT recibe entrenamiento adicional
mensual que cubre tanto los temas CERT como habilidades más avanzadas de
seguridad pública.
<p>Para más información sobre nuestro programa CERT, escriba a
<a href="mailto:cert@sunnyvale.ca.gov" target="_blank">cert@sunnyvale.ca.gov</a>.
<div style="margin:1.5rem 0"><button class="sbtn sbtn-primary" onclick="history.back()">Regresar</button></div>`)
		}
	})
}
