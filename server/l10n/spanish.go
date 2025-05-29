package l10n

import (
	"fmt"
	"time"
)

// spanish maps English phrases used in the UI into Spanish phrases.
// Note:  the files in pages/static also have Spanish text in them, and so does
// ui/s-password/s-password.js.
var spanish = map[string]string{
	// common:
	"and":                                "y",
	"Cancel":                             "Cancelar",
	"Cell Phone":                         "Tel. móvil",
	"CERT Basic Training":                "Capacitación básica del CERT",
	"Classes and Training":               "Clases y capacitación",
	"Contact Us":                         "Contáctenos",
	"Details":                            "Detalles",
	"Events":                             "Eventos",
	"Files":                              "Archivos",
	"Greetings, %s,":                     "Saludos, %s:",
	"List":                               "Lista",
	"Log in":                             "Iniciar sesión",
	"Login incorrect. Please try again.": "Acceso incorrecto. Por favor, inténtelo de nuevo.",
	"Login":                              "Iniciar sesión",
	"Name":                               "Nombre",
	"New Password":                       "Nueva contraseña",
	"OK":                                 "Aceptar",
	"Password":                           "Contraseña",
	"People":                             "Personas",
	"pep-logo.png":                       "ppde-logo.png",
	"Personal Emergency Preparedness":    "Preparación para desastres y emergencias",
	"Profile":                            "Perfil",
	"Request Information":                "Solicitar información",
	"%q is not a valid YYYY-MM-DD date.": "%q no es una fecha válida AAAA-MM-DD.",
	"Save":                               "Guardar",
	"SunnyvaleSERV.org Password Reset":   "Restablecimiento de contraseña de SunnyvaleSERV.org",
	"Web Site Information":               "Información del sitio web",

	// pages/classes/*:
	"Class Registration":                  "Inscripción de clase",
	"Email":                               "Email",
	"First":                               "nombre de pila",
	"Last":                                "apellido(s)",
	"Sign Up":                             "Inscribirse",
	"Wait List":                           "Lista de espera",
	"Submit":                              "Enviar",
	"The cell phone number is not valid.": "El número de teléfono móvil no es válido.",
	// CERT description:
	"How to help your community after a disaster": "Cómo ayudar a su comunidad después de un desastre",
	"In a disaster, professional emergency responders will be overwhelmed, and people will have to rely on their neighbors for help.  If you want to be one of the helpers, the <b>Community Emergency Response Team (CERT) Basic Training</b> class is for you.  It teaches basic emergency response skills, and how to use them safely.":                                            "En un desastre, los servicios de emergencia profesionales se verán abrumados y los residentes tendrán que depender de la ayuda de sus vecinos.  Si quiere ser uno de los ayudantes, esta clase <b>Capacitación básica del CERT (Equipo comunitario de respuesta a emergencias)</b> es para usted.  Enseña habilidades básicas de respuesta a emergencias y cómo usarlas de manera segura.",
	"Topics include:<ul><li>Disaster Preparedness<li>The CERT Organization<li>Usage of Personal Protective Equipment (PPE)<li>Disaster Medical Operations<li>Triaging, Assessing, and Treating Patients<li>Disaster Psychology<li>Fire Safety and Utility Control<li>Extinguishing Small Fires<li>Light Search and Rescue<li>Terrorism and CERT<li>Disaster Simulation Exercise</ul>": "Los temas incluyen:<ul><li>Preparación para desastres<li>La organización CERT<li>Uso de equipo de protección personal<li>Operaciones médicas en casos de desastre<li>Selección, evaluación y tratamiento de pacientes<li >Psicología de desastres<li>Seguridad contra incendios y control de servicios públicos<li>Extinción de pequeños incendios<li>Búsqueda y rescate ligeros<li>Terrorismo y CERT<li>Ejercicio de simulación de desastres</ul>",
	"This class meets for seven weekday evenings and one full Saturday (see dates below).  On successful completion of the class, you will be invited to join the Sunnyvale CERT Deployment Team, which supports the professional responders in Sunnyvale's Department of Public Safety.":                                                                                             "Esta clase se reúne durante siete tardes entre semana y un sábado completo (ver fechas a continuación).  Al completar exitosamente la clase, se le invitará a unirse al equipo de despliegue de Sunnyvale CERT, que apoya a los socorristas profesionales del Departamento de Seguridad Pública de Sunnyvale.",
	"<b>IMPORTANT:</b>  Space in this class is limited.  Please do not sign up unless you fully expect to attend all of the sessions.  This class is open to anyone aged 18 or over, but preference will be given to Sunnyvale residents.  High school students under age 18 are welcome if their parent or other responsible adult is also in the class.":                            "<b>IMPORTANTE:</b>  El espacio en esta clase es limitado.  No se registre a menos que espere asistir a todas las sesiones.  Esta clase está abierta a cualquier persona mayor de 18 años, pero se dará preferencia a los residentes de Sunnyvale.  Los estudiantes de secundaria menores de 18 años son bienvenidos si sus padres u otro adulto responsable también están en la clase.",
	"<b>IMPORTANT:</b>  This class is taught only in English.  However, the printed materials are available in Spanish.": "<b>IMPORTANTE:</b> Esta clase se imparte únicamenta en inglés.  Sin embargo, los materiales impresos están disponibles en español.",
	// MYN description:
	"Planning for disasters\nwith your neighbors": "Planificar los desastres\ncon sus vecinos",
	"Following a disaster, Sunnyvale residents will need to rely on each other for several days if city and county services are overwhelmed.  The “Map Your Neighborhood” (MYN) program prepares neighbors to organize a timely response and to support each other in a disaster.":                                                                                                                                                               "Tras un desastre, los residentes de Sunnyvale tendrán que depender unos de otros durante varios días si los servicios de la ciudad y el condado se ven desbordados.  El programa MYN (“Mapear su vecindario”, por sus siglas en inglés) prepara a los vecinos para organizar una respuesta oportuna y apoyarse mutuamente en caso de desastre.",
	"In this program, we lead a two-hour meeting of around 15–25 households.  Neighbors learn the 9 Steps to take following a disaster, identify resources and skills available in their neighborhood that will be useful in a disaster response, and “map” any special challenges or people with particular needs.  As part of this model, neighbors get to know each other and are better prepared to work together responding to a disaster.": "En este programa, dirigimos una reunión de dos horas de duración en la que participan entre 15 y 25 hogares.  Los vecinos aprenden los 9 pasos a seguir tras un desastre, identifican los recursos y habilidades disponibles en su vecindario que serán útiles en una respuesta al desastre, y “mapean” cualquier desafío especial o personas con necesidades particulares.  Como parte de este modelo, los vecinos se conocen entre sí y están mejor preparados para trabajar juntos en la respuesta a un desastre.",
	"For more information about this program, or to arrange a MYN meeting for your neighborhood, click the button below and fill out the contact form.  Alternatively, you can write to <a href=mailto:myn@sunnyvale.ca.gov target=_blank>myn@sunnyvale.ca.gov</a>.":                                                                                                                                                                             "Para más información sobre este programa, o para organizar una reunión de MYN para su vecindario, escriba a <a href=mailto:myn@sunnyvale.ca.gov target=_blank>myn@sunnyvale.ca.gov</a>.",
	// PEP description:
	"Are you prepared\nfor a disaster?": "¿Está preparado\npara un desastre?",
	"Earthquakes, fires, floods, pandemics, power outages, chemical spills ... these are just some of the disasters than can strike our area without warning.  After a disaster strikes, professional emergency services may not be available to help you for several days.  Are you fully prepared to take care of yourself and your family if the need arises?": "Terremotos, incendios, inundaciones, pandemias, cortes de energía, derrames químicos ... estos son solo algunos de los desastres que pueden afectarnos sin aviso.  Después de un desastre, es posible que los servicios de emergencia profesionales no estén disponibles durante varios días.  ¿Está completamente preparado para cuidar de usted y de su familia si se necesita?",
	"Our <b>Personal Emergency Preparedness</b> class can help you prepare for disasters.  It will teach you about the various disasters you might face, what preparations you can make for them, and how to prioritize.":                                                                                                                                         "Nuestra clase <b>Preparación para desastres y emergencias</b> puede ayudarle a prepararse para desastres.  Enseñaremos sobre los diversos desastres que podría enfrentar, qué preparativos puede hacer para ellos y cómo establecer prioridades.",
	"We also teach tailored versions of the class for private groups such as apartment complexes, churches, and businesses.  To arrange a class for your group, please contact us at pep@sunnyvaleserv.org.":                                                                                                                                                      "También impartimos versiones adaptadas de la clase para grupos privados, como complejos de apartamentos, iglesias y empresas.  Para organizar una clase para su grupo, póngase en contacto con nosotros en pep@sunnyvaleserv.org.",

	// pages/classes/all.go:
	"View More": "Ver más",

	// pages/classes/common.go:
	"This session is full.": "Esta sesión está llena.",
	"This class is presented by Sunnyvale Emergency Response Volunteers (SERV), the volunteer arm of the Sunnyvale Office of Emergency Services.": "Esta clase es presentada por Voluntarios de Respuesta a Emergencias de Sunnyvale (SERV, en inglés), el brazo voluntario de la Oficina de Servicios de Emergencia de Sunnyvale.",

	// pages/classes/register.go:
	"This class is now full.": "Esta clase ahora está llena.",
	"This class is now full.  You will be placed on a waiting list for the class and will be notified if space becomes available.": "Este curso ahora está llena.  Le incluiremos en una lista de espera para la clase y le notificaremos si hay espacio disponible.",
	"Student %d":                             "Estudioso %d",
	"Clear":                                  "Vaciar",
	"How did you find out about this class?": "¿Cómo se enteró de esta clase?",
	"(select one)":                           "(elija uno)",
	"Both first and last name are required. ":        "Se requieren tanto el nombre como el apellido. ",
	"Each student must have a different name. ":      "Cada estudioso debe tener un nombre diferente. ",
	"The email address is not valid. ":               "La dirección de correo electrónico no es válida. ",
	"The class does not have this many spaces left.": "A la clase no le quedan tantos espacios.",
	"Thank you for your interest in our “%s” class:": "Gracias por su interés en nuestra clase “%s”.",
	"We confirm the registrations of:":               "Confirmamos las inscripciones de:",
	"We confirm the registration of:":                "Confirmamos la inscripción de:",
	"You have canceled the registrations of:":        "Ha cancelado las inscripciones de:",
	"You have canceled the registration of:":         "Ha cancelado la inscripción de:",
	"If you need to withdraw from the class or make other changes, please return to SunnyvaleSERV.org.  You may also reply to this email.": "Si necesita retirarse de la clase o realizar otros cambios, regrese a SunnyvaleSERV.org. También puede responder a este mensaje.",
	"We look forward to seeing you!":                                                                            "¡Esperamos verle!",
	"We hope to be able to accommodate you at some future class.":                                               "Esperamos poder acomodarlo en alguna clase futura.",
	"%s has registered you for our “%s” class:":                                                                 "%s le ha inscribido en nuestra clase “%s”:",
	"If this is incorrect, or you need to withdraw from the class, please reply to this email and let us know.": "Si esto es incorrecto o necesita retirarse de la clase, responda a este mensaje e infórmenos.",
	"%s has canceled your registration for our “%s” class:":                                                     "%s ha cancelado su inscripción en nuestra clase “%s”:",
	"If this is incorrect, please reply to this email and let us know.":                                         "Si esto es incorrecto, responda a este mensaje e infórmenos.",
	"Thank you!  Your class registrations are confirmed.":                                                       "¡Gracias! Sus inscripciones a la clase están confirmadas.",
	"Thank you!  Your class registration is confirmed.":                                                         "¡Gracias! Su inscripción a la clase está confirmada.",
	"Thank you!  Your class registrations are canceled.":                                                        "¡Gracias! Sus inscripciones a la clase están canceladas.",
	"Thank you!  Your class registration is canceled.":                                                          "¡Gracias! Su inscripción a la clase está cancelada.",
	"Thank you!  Your changes have been saved.":                                                                 "¡Gracias! Se han guardado sus cambios.",
	"A confirmation message has been sent to %s. If you don’t receive it promptly, look for it in your Junk Mail folder. Move it to your inbox so that future messages from us about the class are not marked as Junk Mail.": "Se ha enviado un mensaje de confirmación a %s. Si no lo recibe rápidamente, búsquelo en su carpeta de correo no deseado. Muévalo a su bandeja de entrada para que futuros mensajes nuestros sobre la clase no se marquen como correo no deseado.",
	"If you need to withdraw from the class, please return to this website and remove your registration.  You may also send email to serv@sunnyvale.ca.gov.":                                                                 "Si necesita retirarse de la clase, regrese a este sitio web y vacie su inscripción. También puede enviar un correo electrónico a serv@sunnyvale.ca.gov.",
	"A confirmation message has been sent to %s.": "Se ha enviado un mensaje de confirmación a %s.",

	// pages/classes/reglogin.go:
	"To register for this class, please enter your email address.": "Para inscribirse en esta clase, introduzca su dirección de correo electrónico.",
	"Your email address is required.":                              "Se requiere su dirección de correo electrónico.",
	"This is not a valid email address.":                           "Esta dirección de correo electrónico no es válida.",
	"To register for this class, please log in.":                   "Para inscribirse en esta clase, inicie sesión.",
	"Your password is required.":                                   "Se requiere su contraseña.",
	"Your name is required.":                                       "Se requiere su nombre y apellido.",
	"To register for this class, please create an account.":        "Para inscribirse en esta clase, cree una cuenta.",
	"We do not have an account with this email address.  To create a new account, please provide the following information.": "No tenemos una cuenta con esta dirección de correo electrónico.  Para crear una cuenta nueva, facilite la siguiente información.",
	"The cell phone is used only for urgent notifications, such as last-minute cancellation of a class.  It is optional.":    "El teléfono móvil sólo se utiliza para notificaciones urgentes, como la cancelación de una clase en el último momento.  Es opcional.",
	"Create Account": "Crear cuenta",

	// pages/errpage/errpage.go:
	"No Such Page": "No existe esa página",
	"Sorry, the page you asked for doesn’t exist.  But we have plenty of other good ones!  You can <a href=\"javascript:history.back()\">go back</a> to where you were, or return to <a href=\"/\">the home page</a>.  Look around; you’re sure to find a page you like.": "Lo sentimos, la página que solicitó no existe.  ¡Pero tenemos muchas otras buenas! Puede <a href=\"javascript:history.back()\">volver</a> a donde estaba o regrese a la <a href=\"/\">la página de inicio</a>.  Mire alrededor; Seguro que encontrará una página que le gusta.",
	"Forbidden": "Prohibido",
	"Sorry, but your account doesn’t have permissions for the operation you requested.  If you think you should have permissions, contact <a href=mailto:admin@sunnyvaleserv.org>admin@SunnyvaleSERV.org</a> for assistance.": "Lo sentimos, pero su cuenta no tiene permisos para la operación que usted solicitó.  Si cree que debería tener permisos, póngase en contacto conmailto:admin@sunnyvaleserv.org>admin@SunnyvaleSERV.org</a> para asistencia.",
	"Web Site Error": "Error del sitio web",
	"We’re sorry, but this web site isn’t working correctly right now.  This problem has been reported to the site administrator.  We’ll get it fixed as soon as possible.": "Lo sentimos, pero este sitio web no funciona correctamente en este momento.  Este problema ha sido informado al administrador del sitio.  Lo solucionaremos lo antes posible.",

	// pages/events/*:
	"Calendar":     "Calendario",
	"Location TBD": "Sitio por determinar",
	"Signups":      "Inscripciones",

	// pages/events/eventscal/eventscal.go:
	"SMTWTFS": "DLMMJVS",

	// pages/events/eventslist/eventslist.go:
	"Date":     "Fecha",
	"Event":    "Evento",
	"Location": "Sitio",
	"TBD":      "Por determinar",

	// pages/events/eventview/details.go:
	"from %s to %s": "de %s a %s",
	"at %s":         "a las %s",

	// pages/events/eventview/task.go:
	"No one can sign up right now.":                               "Nadie puede inscribirse en este momento.",
	"Only %s can sign up.":                                        "Sólo %s pueden inscribirse.",
	"Signups for this task require a completed background check.": "Las inscripciones para esta tarea requieren una verificación de antecedentes completa.",
	"Signups for this task require current DSW registration.":     "Las inscripciones para esta tarea requieren un registro DSW actualizado.",
	"Attendance":                              "Asistencia",
	"You signed in.":                          "Se registró.",
	"You did not sign in.":                    "No se registró.",
	"You were credited for this session.":     "Se le acreditó por esta sesión.",
	"You were not credited for this session.": "No se le acreditó por esta sesión.",
	"Volunteer hours":                         "Horas de voluntariado",
	"You spent %s volunteer hours.":           "Pasó %s horas de voluntariado.",
	"You spent %s volunteer hour.":            "Pasó %s hora de voluntariado.",
	"You did not record volunteer hours.":     "No registró horas de voluntariado.",

	// pages/events/signups/shared.go:
	"Have %d,": "Tenemos %d,",
	"need %d":  "necesitamos %d",
	"limit %d": "límite %d",
	"no limit": "no hay límite",

	// pages/events/signups/signups.go:
	"Event Signups": "Inscripciones para eventos",
	"There are no upcoming events with signups.": "No hay eventos próximos con inscripciones.",

	// pages/homepage/homepage.go:
	"Sunnyvale Emergency Response Volunteers": "Voluntarios de Respuesta a Emergencias de Sunnyvale",
	"Volunteer Login": "Iniciar sesión",
	"Disaster preparedness for homes and families":                                 "Cómo preparar su hogar y su familia para un desastre",
	"2 hours, English or Spanish":                                                  "2 horas, español o inglés",
	"Helping others safely in a disaster":                                          "Ayudar a otros de forma segura en un desastre",
	"7 weeks, English only":                                                        "7 semanas, sólo en inglés",
	"Planning for disasters with your neighbors":                                   "Planificar los desastres con sus vecinos",
	"Volunteer Programs":                                                           "Programas de voluntariado",
	"Community Emergency Response Team":                                            "Equipo comunitario de respuesta a emergencias",
	"Listos California: Preparedness Education":                                    "Listos California: Educación de preparación",
	"Sunnyvale Amateur Radio Emergency Communications Service":                     "Radioaficionados de Sunnyvale:\nCommunicaciones en emergencias",
	"Sunnyvale Neighborhoods Actively Prepare":                                     "Vecindarios de Sunnyvale se preparan activamente",
	"Information Library":                                                          "Archivos y recursos",
	"Office of Emergency Services\nDepartment of Public Safety\nCity of Sunnyvale": "Oficina de Servicios de Emergencia\nDepartamento de Seguridad Pública\nCiudad de Sunnyvale",
	"<a href=\"tel:+14087307190\">(408) 730-7190</a>":                              "<a href=\"tel:+14087307294\">(408) 730-7294</a>",
	"(messages only)":                                                              "(mensajes solamente)",

	// pages/login/*:
	"Email address": "Email",

	// pages/login/login.go:
	"Please log in.": "Por favor, inicie sesión.",
	"Your browser is out of date and lacks features needed by this web site. The site may not look or behave correctly.": "Su navegador no está actualizado y carece de las funciones necesarias para este sitio web. El sitio puede no verse o comportarse correctamente.",
	"Remember me":       "Recuérdeme",
	"Reset my password": "Restablecer contraseña",

	// pages/login/newpwd.go:
	"The two passwords are not the same.":                                     "",
	"Please specify a new password, twice.":                                   "",
	"This password would take less than a minute to crack.":                   "Esta contraseña tardaría menos que un minuto en descifrarse.",
	"This password would take %d minute to crack.":                            "Esta contraseña tardaría %d minuto en descifrarse.",
	"This password would take %d minutes to crack.":                           "Esta contraseña tardaría %d minutos en descifrarse.",
	"This password would take %d hour to crack.":                              "Esta contraseña tardaría %d hora en descifrarse.",
	"This password would take %d hours to crack.":                             "Esta contraseña tardaría %d horas en descifrarse.",
	"This password would take %d day to crack.":                               "Esta contraseña tardaría %d día en descifrarse.",
	"This password would take %d days to crack.":                              "Esta contraseña tardaría %d días en descifrarse.",
	"This password would take %d month to crack.":                             "Esta contraseña tardaría %d mes en descifrarse.",
	"This password would take %d months to crack.":                            "Esta contraseña tardaría %d meses en descifrarse.",
	"This password would take %d year to crack.":                              "Esta contraseña tardaría %d año en descifrarse.",
	"This password would take %d years to crack.":                             "Esta contraseña tardaría %d años en descifrarse.",
	"This password would take centuries to crack.":                            "Esta contraseña tardaría siglos en descifrarse.",
	"A word by itself is easy to guess.":                                      "Una palabra por sí sola es fácil de adivinar.",
	"Add another word or two.  Uncommon words are better.":                    "Añada una o dos palabras más.  Las palabras poco comunes son mejores.",
	"All upper case is almost as easy to guess as all lower case.":            "Todas mayúsculas son casi tan fáciles de adivinar como todas minúsculas.",
	"Avoid dates and years that are associated with you.":                     "Evite fechas y años que estén asociados con usted.",
	"Avoid repeated words and characters.":                                    "Evite palabras y caracteres repetidos.",
	"Avoid sequences.":                                                        "Evite secuencias.",
	"Capitalization doesn’t help very much.":                                  "Las mayúsculas no ayudan mucho.",
	"Common names and surnames are easy to guess.":                            "Los nombres y apellidos comunes son fáciles de adivinar.",
	"Dates are often easy to guess.":                                          "Las fechas suelen ser fáciles de adivinar.",
	"No need for symbols, digits, or uppercase letters.":                      "No necesita símbolos, dígitos ni mayúsculas.",
	"Predictable substitutions like “@” instead of “a” don’t help very much.": "Las sustituciones predecibles como “@” en lugar de “a” no ayudan mucho.",
	"Repeats like “aaa” are easy to guess.":                                   "Las repeticiones como “aaa” son fáciles de adivinar.",
	"Repeats like “abcabcabc” are only slightly harder to guess than “abc”.":  "Las repeticiones como “abcabcabc” son un poco más difíciles de adivinar que “abc”.",
	"Sequences like “abc” or “6543” are easy to guess.":                       "Secuencias como “abc” o “6543” son fáciles de adivinar.",
	"Short keyboard patterns are easy to guess.":                              "Los patrones de teclado cortos son fáciles de adivinar.",
	"This is similar to a commonly used password.":                            "Es similar a una contraseña común.",
	"Use a few words.  Avoid common phrases":                                  "Use pocas palabras.  Evite frases comunes.",
	"Use a longer keyboard pattern with more turns.":                          "Utilice un patrón de teclado más largo y con más vueltas.",

	// pages/login/pwreset.go:
	"Password Reset": "Restablecer contraseña",
	"To reset your password, please enter your email address.  If it’s one we have on file, we’ll send a password reset link to it.": "Para restablecer su contraseña, introduzca su dirección de correo electrónico.  Si es una de las que tenemos archivadas, le enviaremos un enlace para restablecer la contraseña.",
	"Reset Password": "Restablecer contraseña",
	"To reset your password on SunnyvaleSERV.org, click this link:":                                                                                                                                                                                   "Para restablecer su contraseña de SunnyvaleSERV.org, haga clic en este enlace:",
	"If you have any problems, reply to this email. If you did not request a password reset, you can safely ignore this email.":                                                                                                                       "Si tiene algún problema, responda a este mensaje.  Si no ha solicitado un restablecimiento de contraseña, puede ignorar este mensaje.",
	"We have sent a password reset link to the email address you provided. It is valid for one hour. Please check your email and follow the link we sent to reset your password.":                                                                     "Hemos enviado un enlace para restablecer la contraseña a la dirección de correo electrónico que nos ha facilitado. Es válido durante una hora. Compruebe su correo electrónico y siga el enlace que le hemos enviado para restablecer su contraseña.",
	"If you do not receive an email with a password reset link, it may be that the email address you provided is not the one we have on file for you. Contact <a href=\"mailto:admin@sunnyvaleserv.org\">admin@SunnyvaleSERV.org</a> for assistance.": "Si no recibe un mensaje con un enlace para restablecer la contraseña, es posible que la dirección de correo electrónico que nos ha facilitado no sea la que tenemos registrada. Póngase en contacto con <a href=\"mailto:admin@sunnyvaleserv.org\">admin@SunnyvaleSERV.org</a> para obtener ayuda.",
	"This password reset link is invalid or has expired.": "Este enlace para restablecer la contraseña no es válido o ha caducado.",
	"Try Again": "Intentárlo de nuevo",

	// pages/people/*:
	"(all)":           "(todos)",
	"Edit":            "Editar",
	"Home Address":    "Dirección de casa",
	"Mailing Address": "Dirección de correos",
	"Map":             "Mapa",
	"Work Address":    "Dirección de trabajo",

	// pages/people/activity/activity.go:
	"Activity":           "Actividad",
	"Volunteer Activity": "Actividad de voluntariado",
	"%s Activity":        "Actividad de %s",
	"Signed In":          "Registrado",
	"Credited":           "Acreditado",
	"Other %s Hours":     "Otras horas para %s",
	"No activity.":       "No hay actividad.",

	// pages/people/peoplelist/peoplelist.go:
	"1 person listed.":  "1 persona en la lista.",
	"%d people listed.": "%d personas en la lista.",
	"Sort":              "Ordenar",
	"cell":              "móvil",
	"home":              "casa",
	"work":              "trabajo",

	// pages/people/peoplemap/peoplemap.go:
	"(Business Hours)": "(Horas de trabajo)",
	"Home[ADDR]":       "En casa",
	"Business":         "A trabajo",

	// pages/people/personedit/contact.go:
	"Edit Contact Information":                          "Editar información de contacto",
	"%q is not a valid email address.":                  "%q no es una dirección de correo electrónico válida.",
	"The email address %q is in use by another person.": "La dirección de correo electrónico %q está siendo utilizada por otra persona.",
	"This is the email address you log in with.":        "Esta es la dirección de correo electrónico con la que inicia sesión.",
	"Alt. Email": "Otro email",
	"%q is not a valid 10-digit phone number.":     "%q no es un número de teléfono válido de 10 dígitos.",
	"Another person has the cell phone number %q.": "Otra persona tiene el número de teléfono móvil %q.",
	"Home Phone":                      "Tel. de casa",
	"%q is not a valid phone number.": "%q no es un número de teléfono válido.",
	"Work Phone":                      "Tel. de trabajo",
	"This address cannot be marked “same as home” when there is no home address.":                       "Esta dirección no se puede marcar como “igual que la de casa” cuando no hay una dirección de casa.",
	"Address changes cannot be accepted right now because the address verification service is offline.": "No se pueden aceptar cambios de dirección en este momento porque el servicio de verificación de dirección está fuera de línea.",
	"This is not a valid address.":                        "Esta no es una dirección válida.",
	"Same as home address":                                "Igual que la de casa",
	"A phone number may not be specified without a name.": "No se puede especificar un número de teléfono sin un nombre.",
	"At least one phone number is required.":              "Se requiere al menos un número de teléfono.",
	"A relationship may not be specified without a name.": "No se puede especificar una relación sin un nombre.",
	"The relationship is required.":                       "Se requiere la relación.",
	"%q is not one of the relationship choices.":          "%q no es una de las opciones de relación.",
	"Emergency Contact":                                   "Contacto de emergencias",
	"Relationship":                                        "Relación",
	"(select relationship)":                               "(seleccione una relación)",
	"Co-worker":                                           "Compañero de trabajo",
	"Daughter":                                            "Hija",
	"Father":                                              "Padre",
	"Friend":                                              "Amigo",
	"Mother":                                              "Madre",
	"Neighbor":                                            "Vecino",
	"Other":                                               "Otro",
	"Relative":                                            "Pariente",
	"Son":                                                 "Hijo",
	"Spouse":                                              "Cónyuge",
	"Supervisor":                                          "Supervisor",

	// pages/people/personedit/names.go:
	"Edit Names":            "Editar nombres",
	"The name is required.": "Se requiere el nombre.",
	"What you like to be called, e.g. “Joe Banks”":      "Cómo le gusta que le llamen, p.e. “Paco García”",
	"The formal name is required.":                      "Se requiere el nombre formal.",
	"Formal name":                                       "Nombre formal",
	"For formal documents, e.g. “Joseph A. Banks, Jr.”": "Para documentos formales, p.e. “Francisco García Ramírez”",
	"The sort name is required.":                        "Se requiere el nombre ordenado.",
	"Another person has the sort name %q.":              "Otra persona tiene el nombre ordenado %q.",
	"Sort name":                                         "Nombre ordenado",
	"For appearance in sorted lists, e.g. “Banks, Joe”": "Para aparecer en listas ordenadas, p.e. “García, Paco”",
	"%q is not a valid FCC amateur radio call sign.":    "%q no es un indicativo válido para radioaficionados de la FCC.",
	"Another person has the call sign %q.":              "Otra persona tiene el indicativo %q.",
	"Call sign":                                         "Indicativo",
	"FCC amateur radio license (if any)":                "Indicativo de licencia de radioaficionado de la FCC (si corresponde)",
	"Birthdate":                                         "Fecha de nacimiento",
	"Pronouns":                                          "Pronumbres",
	"he/him/his":                                        "él/lo",
	"she/her/hers":                                      "ella/la",
	"they/them/theirs":                                  "elle/le",

	// pages/people/personedit/password.go:
	"Password Change":                       "Cambiar de contraseña",
	"Please specify your old password.":     "Por favor ingrese su contraseña anterior.",
	"This is not the correct old password.": "Esta no es la contraseña anterior correcta.",
	"Old Password":                          "Contraseña anterior",
	"Please specify a valid new password.":  "Por favor ingrese una nueva contraseña válida.",
	"The new password is too weak.":         "La nueva contraseña es demasiado débil.",

	// pages/people/personedit/pwreset.go:
	"%s has reset the password for your account on SunnyvaleSERV.org.  Your new login information is:": "%s ha restablecido la contraseña de su cuenta en SunnyvaleSERV.org.  Su nueva información de acceso es:",
	"Email:    %s": "Email:      %s",
	"Password: %s": "Contraseña: %s",
	"This password is three words chosen randomly from a dictionary — a method that generally produces a very secure and often memorable password.  If the resulting phrase has any meaning, it’s unintentional coincidence.": "Esta contraseña consiste en tres palabras elegidas al azar de un diccionario inglés — un método que generalmente produce una contraseña muy segura y a menudo memorable.  Si la frase resultante tiene algún significado, se trata de una coincidencia involuntaria.",
	"You can change this password by logging into SunnyvaleSERV.org and clicking the “Change Password” button on your Profile page.  If you have any questions, just reply to this email.":                                    "Puede cambiar esta contraseña accediendo a SunnyvaleSERV.org y haciendo clic en el botón “Cambiar contraseña” en su página de perfil.  Si usted tiene alguna pregunta, simplemente responder a este mensaje.",
	"Regards,": "Saludos,",

	// pages/people/personedit/subscriptions.go:
	"Edit List Subscriptions": "Editar suscripciones",
	"Messages sent to %s are considered required for the %s role.  Unsubscribing from it may cause you to lose that role.":    "Los mensajes enviados a %s se consideran obligatorios para el papel “%s”.  Desuscribirse puede hacer que pierda ese papel.",
	"Messages sent to %s are considered required for the %s roles.  Unsubscribing from it may cause you to lose those roles.": "Los mensajes enviados a %s se consideran obligatorios para los papeles “%s” y “%s”.  Desuscribirse puede hacer que pierda esos papeles.",
	"Unsubscribe All": "Desuscribirse a todos",

	// pages/people/personedit/vregister.go:
	"Register as a City Volunteer": "Registrarse como voluntario de ciudad",
	"Thank you for your interest in volunteering with the City of Sunnyvale, Office of Emergency Services.  Please complete this form to register as a City of Sunnyvale Volunteer.  (Please note: registering as a city volunteer is not required for taking one of our classes.  It is only required when joining one of our volunteer groups.)":                                                                                                                                             "Gracias por su interés en ser voluntario en la Oficina de Servicios de Emergencia de la ciudad de Sunnyvale.  Complete este formulario para registrarse como voluntario de la ciudad de Sunnyvale.  Una vez que recibamos su registro (lo que generalmente demora unos días), nos comunicaremos con usted para programar una cita para su toma de huellas digitales.  (Tenga en cuenta: no es necesario registrarse como voluntario de la ciudad para tomar una de nuestras clases.  Solo es necesario cuando se une a uno de nuestros grupos de voluntarios).",
	"Thank you for your interest in volunteering with the City of Sunnyvale, Office of Emergency Services.  Please complete this form to register as a City of Sunnyvale Volunteer.  Once we receive your registration (which usually takes a few days) we will contact you to schedule an appointment for your fingerprinting.  (Please note: registering as a city volunteer is not required for taking one of our classes.  It is only required when joining one of our volunteer groups.)": "Gracias por su interés en ser voluntario en la Oficina de Servicios de Emergencia de la ciudad de Sunnyvale.  Complete este formulario para registrarse como voluntario de la ciudad de Sunnyvale.  (Tenga en cuenta: no es necesario registrarse como voluntario de la ciudad para tomar una de nuestras clases.  Solo es necesario cuando se une a uno de nuestros grupos de voluntarios).",
	"A cell or home phone number is required.": "Se requiere un número de teléfono móvil o a casa.",
	"Your home address is required.":           "Se requiere su dirección de casa.",
	"Interests":                                "Intereses",
	"CERT Deployment Team":                     "Equipo de despliegue CERT",
	"Community Outreach":                       "Alcance comunitario",
	"Amateur Radio (SARES)":                    "Radioaficionados (SARES)",
	"Neighborhood Preparedness Facilitator":    "Facilitador de preparación vecinal",
	"Preparedness Class Instructor":            "Instructor de clases de preparación",
	"CERT Basic Training Instructor":           "Instructor de clases CERT",
	"Please check that you agree with the above statement in order to register.": "Por favor, marque la casilla para mostrar que está de acuerdo con la declaración anterior para poder registrarse.",
	"By submitting this application, I certify that all statements I have made on this application are true and correct and I hereby authorize the City of Sunnyvale to investigate the accuracy of this information.  I am aware that fingerprinting and a criminal records search is required for volunteers 18 years of age or older.  I understand that I am working at all times on a voluntary basis, without monetary compensation or benefits, and not as a paid employee.  I give the City of Sunnyvale permission to use any photographs or videos taken of me during my service without obligation or compensation to me.  I understand that the City of Sunnyvale reserves the right to terminate a volunteer's service at any time.  I understand that volunteers are covered under the City of Sunnyvale's Worker's Compensation Program for an injury or accident occurring while on duty.": "Al enviar esta solicitud, certifico que todas las declaraciones que he hecho en esta solicitud son verdaderas y correctas y por la presente autorizo a la ciudad de Sunnyvale a investigar la exactitud de esta información.  Soy consciente de que se requieren huellas dactilares y una búsqueda de antecedentes penales para los voluntarios mayores de 18 años.  Entiendo que estoy trabajando en todo momento de forma voluntaria, sin compensación monetaria ni beneficios, y no como empleado remunerado.  Doy permiso a la ciudad de Sunnyvale para utilizar fotografías o videos tomados de mí durante mi servicio sin obligación ni compensación para mí.  Entiendo que la ciudad de Sunnyvale se reserva el derecho de cancelar el servicio de un voluntario en cualquier momento.  Entiendo que los voluntarios están cubiertos por el Programa de Compensación para Trabajadores de la Ciudad de Sunnyvale por una lesión o accidente que ocurra mientras están de servicio.",
	"I agree":  "Estoy de acuerdo",
	"Register": "Registrarse",
	"Thank you for volunteering with the City of Sunnyvale, Office of Emergency Services.  One of our staff will contact you to schedule a fingerprinting appointment.  (Criminal history checks are required by city policy for all public-facing volunteers.)  If you have not heard from us within a few days, please email us at oes@sunnyvale.ca.gov to follow up.  We look forward to working with you!": "Gracias por ser voluntario con la Ciudad de Sunnyvale, Oficina de Servicios de Emergencia.  Uno de nuestro personal se pondrá en contacto con usted para programar una cita para la toma de huellas dactilares.  (Las comprobaciones de antecedentes penales son requeridas por la política de la ciudad para todos los voluntarios de cara al público).  Si usted no ha oído hablar de nosotros dentro de unos días, por favor envíenos un correo electrónico a oes@sunnyvale.ca.gov para hacer un seguimiento.  Estamos deseando trabajar con usted.",

	// pages/people/personview/contact.go:
	"Contact Information":            "Información de contacto",
	"(Cell)":                         "(Móvil)",
	"(Home)":                         "(Casa)",
	"(Work)":                         "(Trabajo)",
	"Home Address (all day)":         "Dirección de casa (todo el día)",
	"No emergency contacts on file.": "No hay contactos de emergencia registrados.",
	"1 emergency contact on file.":   "1 contacto de emergencia registrado.",
	"%d emergency contacts on file.": "%d contactos de emergencia registrados.",
	"Sunnyvale Fire District %d":     "Distrito de bomberos %d de Sunnyvale",

	// pages/people/personview/notes.go:
	"Notes": "Notas",

	// pages/people/personview/password.go:
	"Change Password": "Cambiar contraseña",

	// pages/people/personview/roles.go:
	"SERV Roles":                       "Papeles en SERV",
	"SERV Role":                        "Papel en SERV",
	"No current role in any SERV org.": "No tiene ningún papel actual en ninguna organization de SERV.",

	// pages/people/personview/status.go:
	"Volunteer Status":          "Estado del voluntario",
	"City volunteer":            "Voluntario de la ciudad",
	"Registration pending":      "Registro pendiente",
	"Registered %s":             "Registrado el %s",
	"Registered %s, expires %s": "Registrado el %s, caducará el %s",
	"Not registered":            "No está registrado",
	"Background check":          "Verificación de antecedentes",
	"Cleared":                   "Aprobada",
	"Needed":                    "Necesaria",

	// pages/people/personview/subscriptions.go:
	"Unsubscribed from all email.":                   "Se ha desuscribido de todos los correos electrónicos.",
	"Unsubscribed from all text messaging.":          "Se ha desuscribido de todos los mensajes de texto.",
	"Not subscribed to any email or text messaging.": "No está suscrito a ningún correo electrónico o mensaje de texto.",
	"Subscriptions": "Suscripciones",

	// pages/search/search.go:
	"Search":                       "Buscar",
	"Folders":                      "Carpetas",
	"Documents":                    "Archivos",
	"in folder":                    "en la carpeta",
	"Venues":                       "Sitios",
	"Nothing matched your search.": "No se encontró nada en su búsqueda.",

	// pages/static/*:
	"Back": "Regresar",

	// pages/static/calendar.go:
	"SERV Calendar Subscription": "Suscripción al calendario de SERV",
	"You can subscribe to the SERV calendar so that SERV events will automatically appear in the calendar app on your phone, or in your desktop calendar software. Please see the instructions for your phone or software below.": "Puede suscribirse al calendario de SERV para que los eventos de SERV aparezcan automáticamente en la aplicación de calendario de su teléfono o en el software de calendario de su computadora.  Consulte las instrucciones para su teléfono o software a continuación.",
	"iPhone or iPad Calendar App":       "Aplicación de calendario de iPhone o iPad",
	"Open the Settings app.":            "Abra la aplicación Configuración.",
	"Go to “Calendar”.":                 "Vaya a “Calendario”.",
	"Go to “Accounts”.":                 "Vaya a “Cuentas”.",
	"Go to “Add Account”.":              "Vaya a “Agregar cuenta”.",
	"Tap on “Other”.":                   "Toque “Otro”.",
	"Tap on “Add Subscribed Calendar”.": "Toque “Agregar calendario suscrito”.",
	"In the “Server” field, enter <code>https://sunnyvaleserv.org/calendar.ics</code>.": "En el campo “Servidor”, ingrese <code>https://sunnyvaleserv.org/calendar.ics</code>.",
	"Tap “Next”.": "Toque “Siguiente”.",
	"Optional: change the “Description” field to a name that’s meaningful to you, such as “SERV Calendar”.": "Opcional: cambie el campo “Descripción” por un nombre que sea significativo para usted, como “Calendario de SERV”.",
	"Tap “Save”.": "Toque “Guardar”.",
	"Google Calendar (including Android Phones)":                                                               "Google Calendar (incluidos los teléfonos Android)",
	"In a web browser, go to Google Calendar (<code>https://calendar.google.com</code>). Log in if necessary.": "En un navegador web, vaya a Google Calendar (<code>https://calendar.google.com</code>). Inicie sesión si es necesario.",
	"In the left sidebar, click the large “+” sign next to “Other Calendars”.":                                 "En la barra lateral izquierda, haga clic en el signo grande “+” junto a “Otros calendarios”.",
	"Click “From URL”.": "Haga clic en “Desde URL”.",
	"In the “URL of calendar” field, enter <code>https://sunnyvaleserv.org/calendar.ics</code>.": "En el campo “URL del calendario”, ingrese <code>https://sunnyvaleserv.org/calendar.ics</code>.",
	"Click “Add calendar”.":    "Haga clic en “Agregar calendario”.",
	"Microsoft Outlook":        "Microsoft Outlook",
	"Open Microsoft Outlook.":  "Abra Microsoft Outlook.",
	"Go to the calendar page.": "Vaya a la página del calendario.",
	"In the Home ribbon, click on “Open Calendar”, then “From Internet”.": "En la cinta Inicio, haga clic en “Abrir calendario” y luego en “Desde Internet”.",
	"Enter <code>https://sunnyvaleserv.org/calendar.ics</code>.":          "Ingrese <code>https://sunnyvaleserv.org/calendar.ics</code>.",
	"Click “Yes”.": "Haga clic en “Sí”.",
	"In the left sidebar, under “Other Calendars”, right-click on “Untitled” and choose “Rename Calendar”.": "En la barra lateral izquierda, en “Otros calendarios”, haga clic derecho en “Untitled” y seleccione “Cambiar nombre de calendario”.",
	"Give the calendar a name meaningful to you, such as “SERV Calendar”.":                                  "Asigne al calendario un nombre significativo para usted, como “Calendario de SERV”.",
	"Mac Calendar App":       "Aplicación de calendario para Mac",
	"Open the Calendar app.": "Abra la aplicación Calendario.",
	"From the menu, choose File → New Calendar Subscription.":  "En el menú, elija Archivo → Nueva suscripción a calendario.",
	"Click “Subscribe”.":                                       "Haga clic en “Suscribir”.",
	"Set the options to suit your preferences and click “OK”.": "Ajuste las opciones a sus preferencias y haga clic en “Aceptar”.",
	"Other Software": "Otro software",
	"Most calendar software has the ability to subscribe to Internet calendars. Consult the documentation for your software to find out how. The address of the SERV calendar is <code>https://sunnyvaleserv.org/calendar.ics</code>.": "La mayoría del software de calendario tiene la capacidad de suscribirse a calendarios de Internet. Consulta la documentación de tu software para descubrir cómo. La dirección del calendario de SERV es <code>https://sunnyvaleserv.org/calendar.ics</code>.",

	// pages/static/cert.go:
	"Sunnyvale CERT": "CERT de Sunnyvale",
	"Community Emergency Response Team (CERT)": "Equipo Communitario de Respuesta a Emergencias (CERT)",
	"CERT is a nationwide program, managed by the Federal Emergency Management Agency (FEMA), that prepares residents to care for themselves and their communities during and after major disasters.  Its emphasis is on training residents to be able to respond safely and effectively during an emergency.":                                                                                                                                                                                                                                                                                                                                                                                   "CERT es un programa nacional, gestionado por la Agencia Federal para la Gestión de Emergencias (FEMA), que prepara a los residentes para cuidarse a sí mismos y a sus comunidades durante y después de grandes catástrofes.  Su énfasis está en formar a los residentes para que sean capaces de responder con seguridad y eficacia durante una emergencia.",
	"The CERT program was created by the Los Angeles Fire Department after seeing the significant loss of life of volunteer rescuers in the 1985 Mexico City earthquake.  Volunteers are credited with having saved many lives in the aftermath of that earthquake, but many of the volunteers were killed because they did not know how to keep themselves safe while doing such work.  LAFD created the CERT program to ensure that the same thing didn't happen on their watch.  The 1987 Whittier earthquake near Los Angeles underscored the value of this program.  It the early 1990s, FEMA expanded the program to cover other disasters besides earthquakes, and spread it nationwide.": "El programa CERT fue creado por el Departamento de Bomberos de Los Ángeles tras comprobar la grande pérdida de vidas de rescatadores voluntarios en el terremoto de la Ciudad de México de 1985.  A los voluntarios se les atribuye haber salvado muchas vidas tras ese terremoto, pero muchos de los voluntarios murieron porque no sabían cómo mantenerse a salvo mientras realizaban ese trabajo.  El LAFD creó el programa CERT para asegurarse de que no ocurriera lo mismo durante su guardia.  El terremoto de Whittier, cerca de Los Ángeles, en 1987, puso de manifiesto el valor de este programa.  A principios de los 90, la FEMA amplió el programa para cubrir otras catástrofes además de los terremotos, y lo extendió por todo el país.",
	"In Sunnyvale, we teach the FEMA-standard <a href=/cert-basic up-target=main>CERT Basic Training<a> class, with some local enhancements, to anyone who wants it. This is a 30-hour class, taught over seven weeks, covering all aspects of volunteer disaster response.  For the graduates of that class, we also teach occasional refresher classes on specific CERT topics to help our volunteers keep their skills and knowledge fresh.":                                                                                                                                                                                                                                                  "En Sunnyvale, impartimos el curso <a href=/cert-basic up-target=main>Capacitación básica del CERT</a> estándar de la FEMA, con algunas mejoras locales, a todo aquel que lo desee. Se trata de una clase de 30 horas, impartida a lo largo de siete semanas, que cubre todos los aspectos de la respuesta voluntaria en caso de catástrofe.  Para los graduados de esa clase, también impartimos clases ocasionales de actualización sobre temas específicos del CERT para ayudar a nuestros voluntarios a mantener sus habilidades y conocimientos al día.",
	"Sunnyvale also has a “CERT Deployment Team.”  This is a group of CERT-trained volunteers who have agreed to be on call to assist the professional responders in the Department of Public Safety when needed.  Our CERT Deployment Team receives additional, monthly training covering both the CERT topics and more advanced public safety skills.":                                                                                                                                                                                                                                                                                                                                         "Sunnyvale también tiene un “Equipo de Despliegue CERT”.  Este es un grupo de voluntarios entrenados en CERT que han acordado estar de guardia para ayudar a los respondedores profesionales en el Departamento de Seguridad Pública cuando sea necesario.  Nuestro Equipo de Despliegue CERT recibe entrenamiento adicional mensual que cubre tanto los temas CERT como habilidades más avanzadas de seguridad pública.",
	"For more information about our CERT program, write to <a href=mailto:cert@sunnyvale.ca.gov target=_blank>cert@sunnyvale.ca.gov</a>.": "Para más información sobre nuestro programa CERT, escriba a <a href=mailto:cert@sunnyvale.ca.gov target=_blank>cert@sunnyvale.ca.gov</a>.",

	// pages/static/contact.go:
	"Sunnyvale Emergency Response Volunteers (SERV) is the volunteer arm of the Sunnyvale Office of Emergency Services, which is part of the city’s Department of Public Safety.": "Voluntarios de Respuesta a Emergencias de Sunnyvale (SERV, por siglas en inglés) es el brazo voluntario de la Oficina de Servicios de Emergencia de Sunnyvale, que forma parte del Departamento de Seguridad Pública de la ciudad.",
	"(408) 730–7190 English (messages only)": "(408) 730–7190 en inglés (mensajes solamente)",
	"(408) 730-7294 Spanish (messages only)": "(408) 730-7294 en español (mensajes solamente)",
	"Our offices are at":                     "Nuestra oficina está en",
	"Sunnyvale Public Safety Headquarters":   "Jefatura de Seguridad Pública de Sunnyvale",

	// pages/static/credits.go (and also about.go):
	"Credits and Copyrights": "Créditos y derechos de autor",
	"This site was developed by Steven Roth, as a volunteer for the Sunnyvale Department of Public Safety.  The site software is copyrighted © 2020–2021 by Steven Roth.  Steven Roth has granted the Sunnyvale Department of Public Safety a non-exclusive, perpetual, royalty-free, worldwide license to use this software.  The Sunnyvale Department of Public Safety owns the SunnyvaleSERV.org domain and funds the ongoing usage and maintenance of the site.": "Este sitio fue desarrollado por Steven Roth, como voluntario del Departamento de Seguridad Pública de Sunnyvale. El software del sitio tiene derechos de autor © 2020–2021 por Steven Roth.  Steven Roth ha concedido al Departamento de Seguridad Pública de Sunnyvale una licencia mundial, no exclusiva, perpetua y libre de regalías de uso de este software. El Departamento de Seguridad Pública de Sunnyvale es propietario del dominio SunnyvaleSERV.org y financia el uso y mantenimiento continuo del sitio.",
	"Technologies and Services": "Tecnologías y servicios",
	"The software for this web site is written in <a href=https://golang.org target=_blank>Go</a>, with data storage in a <a href=https://sqlite.org target=_blank>SQLite</a> database.  This web site is hosted by <a href=https://www.dreamhost.com/ target=_blank>Dreamhost</a>.  It uses <a href=https://www.google.com/maps target=_blank>Google Maps</a> for geolocation and mapping, <a href=https://www.twilio.com/ target=_blank>Twilio</a> for text messaging, and <a href=https://www.algolia.com/ target=_blank>Algolia</a> for searching.": "El software de este sitio web está escrito en <a href=https://golang.org target=_blank>Go</a>, con almacenamiento de datos en un base de datos <a href=https://sqlite.org target=_blank>SQLite</a>.  Esta sitio web está alojado por <a href=https://www.dreamhost.com/ target=_blank>Dreamhost</a>.  Usa <a href=https://www.google.com/maps target=_blank>Google Maps</a> para geolocalización y cartografía, <a href=https://www.twilio.com/ target=_blank>Twilio</a> para mensajes de texto, y <a href=https://www.algolia.com/ target=_blank>Algolia</a> para buscar.",

	// pages/static/emaillists.go:
	"SERV Email Lists": "Listas de correo electrónico de SERV",
	"The SunnyvaleSERV.org site offers a number of email distribution lists. We have one for each volunteer program, that we give out to the general public who might want more information about the program.  Email sent to these lists is delivered to designated public contact people for each program:": "El sitio SunnyvaleSERV.org ofrece varias listas de distribución de correo electrónico.  Tenemos uno para cada programa de voluntariado, que entregamos al público general que pueda querer más información sobre el programa.  Estas listas se entregan a las personas de contacto público designadas para cada programa:",
	"There are also lists for the volunteers on each of our teams:": "También hay listas de los voluntarios de cada uno de nuestros equipos:",
	"and for the students in each CERT class:":                      "y para los estudiantes en cada clase CERT:",
	"Finally, there are some broader lists for special purposes:":   "Finalmente, existen algunas listas más amplias para propósitos especiales:",
	"All of these email lists have restricted access.  For the team lists, only members of the team can send mail to them; for the class lists, only the instructors can send mail to them; and for the broader lists, only DPS staff can send mail to them.  Any mail sent to any of our lists from someone else is held for approval before being routed to the list.  Messages on topics unrelated to SERV will generally be rejected.": "Todas estas listas de correo electrónico tienen acceso restringido.  Para las listas de equipos, solo los miembros del equipo pueden enviarles correo; para las listas de clases, sólo el los instructores pueden enviarles correo; y para las listas más amplias, solo el personal de DSP puedo enviarles correo. Cualquier correo enviado a cualquiera de nuestras listas por parte de otra persona se retiene para su aprobación antes de ser enviado a la lista. Mensajes sobre temas que no estén relacionados con SERV generalmente serán rechazados.",
	"If you are receiving email from one of these lists that you do not want, there is an “unsubscribe” link at the bottom of every email.  If you are receiving email at the wrong address, you can change your email address in the “Profile” section of this web site.":                                                                                                                                                                 "Si recibe un correo electrónico de una de estas listas que no desea, hay un enlace para “Unsubscribe” (cancelar suscripción) en la parte inferior de cada correo electrónico. Si recibe un correo electrónico en la dirección incorrecta, puede cambiar su dirección de correo electrónico en la sección “Perfil” de este sitio web.",

	// pages/static/listos.go:
	"Listos California is a state program, managed by the California Office of Emergency Services (CalOES), focusing on disaster preparedness education for California residents.  Under their umbrella, the Sunnyvale Listos program provides disaster preparedness education in Sunnyvale.":                                                             "Listos California es un programa estatal, gestionado por la Ofinica de Servicios de Emergencia de California (CalOES), centrándose en la educación de preparación para desastres para los residentes de California.  Bajo su escudo, el programa Listos Sunnyvale proporciona educación de preparación para desastres en Sunnyvale.",
	"Our flagship offering is our <a href=/pep up-target=main>Personal Emergency Preparedness</a> class.  This is a two-hour class that teaches home and family preparedness.  We offer this class to the general public every 2–3 months, in both English and Spanish.  We also offer it to neighborhood associations, businesses, etc. when requested.": "Nuestra oferta principal es nuestra clase <a href=/pep up-target=main>Preparación para desastres y emergencias</a>. Se trata de una clase de dos horas que enseña la preparación del hogar y la familia.  Ofrecemos esta clase al público en general cada 2-3 meses, tanto en inglés como en español.  También la ofrecemos a asociaciones de vecinos, empresas, etc. cuando lo solicitan.",
	"Our disaster preparedness education efforts also include Outreach booths and tables at public events (the Arts and Wine Festival, the Diwali Festival, the Firefighters Pancake Breakfast, neighborhood block parties, etc.).  At these events, we set up tables and distribute disaster preparedness information to participants.":                  "Nuestra labor de educación sobre la preparación ante desastres también incluye puestos y mesas de divulgación en eventos públicos (el Festival de las Artes y el Vino, el Festival Diwali, el Desayuno de Panqueques de los Bomberos, fiestas vecinales, etc.).  En estos eventos, instalamos mesas y distribuimos información sobre preparación ante desastres a los participantes.",
	"For more information about Listos California or our disaster preparedness education programs, write us at <a href=mailto:listos@sunnyvale.ca.gov target=_blank>listos@sunnyvale.ca.gov</a>. Also write to us if you want to arrange a private preparedness class for your neighborhood or group, or have a preparedness table at your event.":        "Para más información sobre Listos California o nuestros programas educativos de preparación ante desastres, escríbanos a <a href=mailto:listos@sunnyvale.ca.gov target=_blank>listos@sunnyvale.ca.gov</a>. También escríbanos si desea organizar una clase privada de preparación para su vecindario o grupo, o tener una mesa de preparación en su evento.",

	// pages/static/privacy.go (and also about.go):
	"Privacy Policy": "Política de privacidad",
	"This web site collects information about people who work for, volunteer for, or take classes organized through the Office of Emergency Services (OES) in the Sunnyvale Department of Public Safety (DPS).  The information we collect includes:": "Este sitio web recopila información sobre personas que trabajan, son voluntarios, o tomar clases organizadas a través de la Oficina de Servicios de Emergencia (OSE) en el Departamento de Seguridad Pública de Sunnyvale (DSP). La información que recopilamos incluye:",
	"Basic Information":       "Información basica",
	"name":                    "nombre",
	"amateur radio call sign": "indicativo de radioaficionado",
	"contact information (email addresses, phone numbers, and physical and postal addresses)":            "información de contacto (direcciones de correo electrónico, números de teléfono y direcciones físicas y postales)",
	"memberships in, and roles held in, SERV volunteer groups":                                           "membresías y papeles desempeñados en grupos de voluntarios de SERV",
	"emergency response classes taken and certificates issued":                                           "clases de respuesta a emergencias tomadas y certificados emitidos",
	"credentials that are relevant to SERV operations":                                                   "credenciales que son relevantes para las operaciones de SERV",
	"other information voluntarily provided such as skills, languages spoken, available equipment, etc.": "otra información proporcionada voluntariamente como habilidades, idiomas hablados, equipos disponibles, etc.",
	"Restricted Information":                                       "Información restringida",
	"attendance at SERV events, and hours spent at them":           "asistencia a eventos de SERV y horas dedicadas a ellos",
	"Disaster Service Worker registration status":                  "estado de registro como trabajador de servicios de desastre",
	"photo IDs and card access keys issued":                        "identificaciones con fotografía y claves de acceso emitidas",
	"Live Scan fingerprinting success, with date (see note below)": "éxito de la toma de huellas digitales de Live Scan, con fecha (consulte la nota a continuación)",
	"background check success, with date (see note below)":         "éxito de la verificación de antecedentes, con fecha (consulte la nota a continuación)",
	"Targeted Information":                                         "Información dirigida",
	"email messages sent to any SunnyvaleSERV.org address":         "mensajes enviados a cualquier dirección de SunnyvaleSERV.org",
	"text messages sent through this web site":                     "mensajes de texto enviados a través de este sitio web",
	"Private Information":                                          "Información privada",
	"logs of web site visits and actions taken":                    "registros de visitas al sitio web y acciones realizadas",
	"All of the above information is available to the paid and volunteer staff of OES and their delegates, including the web site maintainers.  Private information is not available to anyone else.":                                                                                                  "Toda la información anterior está disponible para el personal remunerado y voluntario de OSE y sus delegados, incluidos los mantenedores del sitio web. La información privada no está disponible para nadie más.",
	"If you are a student in an OES-organized class, such as CERT, Listos, or PEP, your basic and restricted information may be shared with the class instructors as long as the class is in progress.":                                                                                                "Si es estudiante de una clase organizada por OSE, como CERT, Listos o PPDE, su su información básica y restringida puede ser compartida con los instructores mientras la clase esté en curso.",
	"If you are a volunteer in a SERV volunteer group, your basic information may be shared with other volunteers in that group, and your restricted information may be shared with the leaders of that group.":                                                                                        "Si es voluntario en un grupo de voluntarios de SERV, su información básica puede ser compartido con otros voluntarios en ese grupo, y su información restringida puede ser compartida con los líderes de ese grupo.",
	"If you are a volunteer in a SERV volunteer group, and you have successfully completed Live Scan fingerprinting and/or background checks, that fact (with no detail other than the date) may be shared with the leaders of your volunteer group.  A negative result will not be shared with them.": "Si es voluntario en un grupo de voluntarios de SERV y ha logrado completado la toma de huellas digitales de Live Scan y/o verificaciones de antecedentes, ese hecho (con ningún detalle más que la fecha) puede ser compartido con los líderes de su grupo de voluntarios.  Un resultado negativo no será compartido con ellos.",
	"If you have sent any email or text messages (targeted information) through the site, they may be shared with any member of the group(s) to which you sent them, including members who join those groups after you send the messages.":                                                             "Si ha enviado algún correo electrónico o mensaje de texto (información dirigida) a través de el sitio, pueden ser compartidos con cualquier miembro del grupo(s) al cual usted los envió, incluidos los miembros que se unen a esos grupos después de que usted envió los mensajes.",
	"If you volunteer for mutual aid or training with another emergency response organization or jurisdiction, we may share your basic and/or restricted information with them.":                                                                                                                       "Si se ofrece como voluntario para ayuda mutua o capacitación con otra respuesta de emergencia organización o jurisdicción, podemos compartir su información básica y/o restringida información con ellos.",
	"The OES staff may share anonymized, aggregate data derived from the above information with anyone at their discretion.":                                                                                                                                                                           "El personal de OSE puede compartir datos agregados anonimizados derivados de la información anterior con cualquier persona a su discreción.",
	"Cookies": "Cookies",
	"This site uses browser cookies.  While you are logged in, a browser cookie contains your session identification; this cookie goes away when you log out or your login session expires.  More permanent cookies are used to store some of your user interface preferences, such as your preferred language and whether you prefer to see the events page in calendar or list form.  No personally identifiable information is ever stored in browser cookies.": "Este sitio utiliza cookies del navegador.  Mientras está conectado, una cookie del navegador contiene su identificación de sesión; Esta cookie desaparece cuando cierra o caduca la sesión.  Se utilizan cookies más permanentes para almacenar algunas de sus preferencias de interfaz de usuario, como su idioma preferido y si prefiere ver la página de eventos en forma de calendario o de lista.  Nunca se almacena información de identificación personal en las cookies del navegador.",

	// pages/static/sares.go:
	"Sunnyvale Amateur Radio Emergency Service": "Servicio de Emergencias de Radioaficionados de Sunnyvale",
	"The Sunnyvale Amateur Radio Emergency Service (SARES) is the local chapter of the nationwide Amateur Radio Emergency Service operated by the Amateur Radio Relay League (ARRL).  During times of emergency, it also operates as a local branch of the federal Radio Amateur Civil Emergency Service (RACES).  SARES provides emergency communications services, usually but not always using amateur radio, when regular communications methods are unavailable or saturated.":                           "El Servicio de Emergencia de Radioaficionados de Sunnyvale (SARES, por siglas en inglés) es el capítulo local del Servicio de Emergencia de Radioaficionados (ARES) a nivel nacional operado por la Liga de Radioaficionados (ARRL).  En situaciones de emergencia, también funciona como una rama local del Servicio de Emergencia Civil de Radioaficionados (RACES).  El SARES proporciona servicios de comunicaciones en emergencias, normalmente pero no siempre utilizando radioafición, cuando los métodos de comunicación habituales no están disponibles o están saturados.",
	"In a disaster, telephones and the Internet will likely be down.  Or, if they are working, they will be unable to keep up with demand.  Radio communications serve as an effective backup because they do not rely on massive, fragile infrastructure.  SARES operators can provide essential emergency communications when no other methods are working.  Outside of emergencies, SARES operators provide ongoing community service by supplying communications assistance at public events on request.": "En caso de desastre, es probable que los teléfonos e Internet no funcionen. O, si funcionan, serán incapaces de satisfacer la demanda.  Las comunicaciones por radio son un medio de reserva eficaz porque no dependen de infraestructuras masivas y frágiles.  Los operadores de SARES pueden proporcionar comunicaciones de emergencia esenciales cuando no funcionan otros métodos.  Fuera de las emergencias, los operadores de SARES prestan un servicio continuo a la comunidad proporcionando asistencia de comunicaciones en actos públicos a pedido.",
	"Membership in SARES requires a current FCC amateur radio license.  If you are interested in emergency communications but do not have a license, SARES members will connect you with resources to help you get one.":                                                                                                                                                                                                                                                                                      "Para ser miembro de SARES se requiere una licencia de radioaficionado de la FCC en vigor.  Si se interesan las comunicaciones de emergencia pero no tiene licencia, los miembros de SARES le pondrán en contacto con recursos que le ayudarán a conseguirla.",
	"For more information about SARES or amateur radio, write to <a href=mailto:sares@sunnyvale.ca.gov target=_blank>sares@sunnyvale.ca.gov</a>.": "Para más información sobre SARES or radioafición, escriba a <a href=mailto:sares@sunnyvale.ca.gov target=_blank>sares@sunnyvale.ca.gov</a>.",

	// pages/static/snap.go:
	"Sunnyvale Neighborhoods Actively Prepare (SNAP)": "Vecindarios de Sunnyvale Se Preparan Activamente",
	"SNAP is our neighborhood disaster preparedness program.  Following a disaster, Sunnyvale residents will need to rely on each other for several days if city and county services are overwhelmed.  While our <a href=/listos up-target=main>Listos</a> program teaches preparedness for individuals and families, SNAP prepares neighbors to organize a timely response and to support each other in a disaster.":                                                                                                                           "El programa SNAP (“Vecindarios de Sunnyvale se preparan activamente”, por sus siglas en inglés) es nuestro programa vecinal de preparación ante desastres.  Tras un desastre, los residentes de Sunnyvale tendrán que depender unos de otros durante varios días si los servicios de la ciudad y el condado se ven desbordados.  Mientras que nuestro programa <a href=/listos up-target=main>Listos</a> enseña preparación a individuos y familias, SNAP prepara a los vecinos para organizar una respuesta oportuna y apoyarse mutuamente en caso de desastre.",
	"Using the “Map Your Neighborhood” (MYN) program provided by the Washington State Emergency Management Division, we lead a two-hour meeting of around 15–25 households.  Neighbors learn the 9 Steps to take following a disaster, identify resources and skills available in their neighborhood that will be useful in a disaster response, and “map” any special challenges or people with particular needs.  As part of this model, neighbors get to know each other and are better prepared to work together responding to a disaster.": "Utilizando el programa MYN (“Mapear su vecindario”, por sus siglas en inglés) proporcionado por la División de Gestión de Emergencias del Estado de Washington, dirigimos una reunión de dos horas de duración en la que participan entre 15 y 25 hogares.  Los vecinos aprenden los 9 pasos a seguir tras un desastre, identifican los recursos y habilidades disponibles en su vecindario que serán útiles en una respuesta al desastre, y “mapean” cualquier desafío especial o personas con necesidades particulares.  Como parte de este modelo, los vecinos se conocen entre sí y están mejor preparados para trabajar juntos en la respuesta a un desastre.",
	"For more information about SNAP, or to arrange a MYN meeting for your neighborhood, write to <a href=mailto:snap@sunnyvale.ca.gov target=_blank>snap@sunnyvale.ca.gov</a>.": "Para más información sobre SNAP, o para organizar una reunión de MYN para su vecindario, escriba a <a href=mailto:snap@sunnyvale.ca.gov target=_blank>snap@sunnyvale.ca.gov</a>.",

	// store/class/referral.go:
	"Word of mouth":                 "Boca a boca",
	"Information table at an event": "Mesa informativa en un evento",
	"Printed advertisement":         "Publicidad impresa",
	"Online advertisement":          "Publicidad en línea",

	// store/shiftperson/eligibility.go:
	"Already signed up for a conflicting shift.": "Ya se inscribió a un turno conflictivo.",
	"Signups are closed.":                        "Las inscripciones están cerradas.",
	"Not eligible to sign up.":                   "No es elegible para registrarse.",
	"DSW registration is required.":              "Se requiere registro DSW.",
	"A background check is required.":            "Se requiere una verificación de antecedentes.",
	"The shift has ended.":                       "El turno ha terminado.",
	"The shift has already started.":             "El turno ya ha comenzado.",
	"The shift is full.":                         "El turno está completo.",
	"No person selected.":                        "Ninguna persona seleccionada.",

	// ui/form/formrow.go:
	"%q is not a valid number.":       "%q no es un número válido.",
	"%q is not a valid value for %s.": "%q no es un valor válido para %s.",

	// ui/page.go:
	"Welcome":       "Bienvenido",
	"Home[PAGE]":    "Inicio",
	"Classes":       "Clases",
	"Logout":        "Cerrar sesión",
	"Web Site Info": "Información del sitio",

	// ui/orgdot/orgdot.go:
	"CERT Deployment": "Despliegue de CERT",
	"CERT Training":   "Capacitación CERT",

	// spanishDate(), below:
	"Sunday":    "Domingo",
	"Monday":    "Lunes",
	"Tuesday":   "Martes",
	"Wednesday": "Miércoles",
	"Thursday":  "Jueves",
	"Friday":    "Viernes",
	"Saturday":  "Sábado",
	"January":   "enero",
	"February":  "febrero",
	"March":     "marzo",
	"April":     "abril",
	"May":       "mayo",
	"June":      "junio",
	"July":      "julio",
	"August":    "agosto",
	"September": "septiembre",
	"October":   "octubre",
	"November":  "noviembre",
	"December":  "diciembre",
}

func spanishDate(day time.Time) string {
	return fmt.Sprintf("%s, el %d de %s de %d", spanish[day.Weekday().String()], day.Day(), spanish[day.Month().String()], day.Year())
}
