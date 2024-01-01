package l10n

import (
	"fmt"
	"time"
)

// spanish maps English phrases used in the UI into Spanish phrases.
// Note:  the files in pages/static also have Spanish text in them.
var spanish = map[string]string{

	// common
	"and":                                "y",
	"Cancel":                             "Cancelar", // Button label
	"Edit":                               "Editar",   // Button label
	"Save":                               "Guardar",  // Button label
	"Details":                            "Detalles",
	"List":                               "Lista",
	"%q is not a valid YYYY-MM-DD date.": "%q no es una fecha válida AAAA-MM-DD.",

	// classes/common.go
	"Sign Up":               "Inscribirse",
	"This session is full.": "Esta sesión está llena.",
	"No sessions of this class are currently scheduled.": "No hay sesiones programadas de esta clase actualmente.",
	"This class is presented by Sunnyvale Emergency Response Volunteers (SERV), the volunteer arm of the Sunnyvale Office of Emergency Services.": "Esta clase es presentada por Voluntarios de Respuesta a Emergencias de Sunnyvale (SERV, en inglés), el brazo voluntario de la Oficina de Servicios de Emergencia de Sunnyvale.",

	// classes/cert.go
	"CERT Basic Training":                         "Capacitación básica del CERT",
	"How to help your community after a disaster": "Cómo ayudar a su comunidad después de un desastre",
	"<p>In a disaster, professional emergency responders will be overwhelmed, and people will have to rely on their neighbors for help.  If you want to be one of the helpers, this class is for you.  It teaches basic emergency response skills, and how to use them safely.  Topics include:</p><ul><li>Disaster Preparedness<li>The CERT Organization<li>Usage of Personal Protective Equipment (PPE)<li>Disaster Medical Operations<li>Triaging, Assessing, and Treating Patients<li>Disaster Psychology<li>Fire Safety and Utility Control<li>Extinguishing Small Fires<li>Light Search and Rescue<li>Terrorism and CERT<li>Disaster Simulation Exercise</ul><p>This class meets for seven weekday evenings and one full Saturday (see dates below).  On successful completion of the class, you will be invited to join the Sunnyvale CERT Deployment Team, which supports the professional responders in Sunnyvale's Department of Public Safety.</p><p>IMPORTANT:  Space in this class is limited.  Please do not sign up unless you fully expect to attend all of the sessions.  This class is open to anyone aged 18 or over, but preference will be given to Sunnyvale residents.  High school students under age 18 are welcome if their parent or other responsible adult is also in the class.</p>": "<p>En un desastre, los servicios de emergencia profesionales se verán abrumados y los residentes tendrán que depender de la ayuda de sus vecinos.  Si quiere ser uno de los ayudantes, esta clase es para usted.  Enseña habilidades básicas de respuesta a emergencias y cómo usarlas de manera segura.  Los temas incluyen:</p><ul><li>Preparación para desastres<li>La organización CERT<li>Uso de equipo de protección personal<li>Operaciones médicas en casos de desastre<li>Selección, evaluación y tratamiento de pacientes<li >Psicología de desastres<li>Seguridad contra incendios y control de servicios públicos<li>Extinción de pequeños incendios<li>Búsqueda y rescate ligeros<li>Terrorismo y CERT<li>Ejercicio de simulación de desastres</ul><p>Esta clase se reúne durante siete tardes entre semana y un sábado completo (ver fechas a continuación).  Al completar exitosamente la clase, se le invitará a unirse al equipo de despliegue de Sunnyvale CERT, que apoya a los socorristas profesionales del Departamento de Seguridad Pública de Sunnyvale.</p><p><b>IMPORTANTE:</b>  El espacio en esta clase es limitado.  No se registre a menos que espere asistir a todas las sesiones.  Esta clase está abierta a cualquier persona mayor de 18 años, pero se dará preferencia a los residentes de Sunnyvale.  Los estudiantes de secundaria menores de 18 años son bienvenidos si sus padres u otro adulto responsable también están en la clase.</p><p><b>IMPORTANTE:</b> Esta clase se imparte únicamenta en inglés.  Sin embargo, los materiales impresos están disponibles en español.</p>",

	// classes/pep.go
	"Personal Emergency Preparedness":   "Preparación para desastres y emergencias",
	"Are you prepared\nfor a disaster?": "¿Está preparado\npara un desastre?",
	"Earthquakes, fires, floods, pandemics, power outages, chemical spills ... these are just some of the disasters than can strike our area without warning.  After a disaster strikes, professional emergency services may not be available to help you for several days.  Are you fully prepared to take care of yourself and your family if the need arises?\n\nOur <b>Personal Emergency Preparedness</b> class can help you prepare for disasters.  It will teach you about the various disasters you might face, what preparations you can make for them, and how to prioritize.": "Terremotos, incendios, inundaciones, pandemias, cortes de energía, derrames químicos ... estos son solo algunos de los desastres que pueden afectarnos sin aviso.  Después de un desastre, es posible que los servicios de emergencia profesionales no estén disponibles durante varios días.  ¿Está completamente preparado para cuidar de usted y de su familia si se necesita?\n\nNuestra clase puede ayudarle a prepararse para desastres.  Enseñaremos sobre los diversos desastres que podría enfrentar, qué preparativos puede hacer para ellos y cómo establecer prioridades.",

	// classes/register.go
	"Class Registration":                     "Inscripción de clase",
	"First":                                  "nombre de pila",
	"Last":                                   "apellido(s)",
	"Student %d":                             "Estudioso %d",
	"Clear":                                  "Vaciar",
	"How did you find out about this class?": "¿Cómo se enteró de esta clase?",
	"(select one)":                           "(elija uno)",
	"This class is now full.":                "Esta clase ahora está llena.",
	"Both first and last name are required. ":        "Se requieren tanto el nombre como el apellido. ",
	"Each student must have a different name. ":      "Cada estudioso debe tener un nombre diferente. ",
	"The email address is not valid. ":               "La dirección de correo electrónico no es válida. ",
	"The cell phone number is not valid. ":           "El número de teléfono móvil no es válido. ",
	"The class does not have this many spaces left.": "A la clase no le quedan tantos espacios.",
	"Greetings, %s,": "Saludos, %s:",
	"Thank you for your interest in our “%s” class:": "Gracias por su interés en nuestra clase “%s”.",
	"We confirm the registration of:":                "Confirmamos la inscripción de:",
	"We confirm the registrations of:":               "Confirmamos las inscripciones de:",
	"You have canceled the registration of:":         "Ha cancelado la inscripción de:",
	"You have canceled the registrations of:":        "Ha cancelado las inscripciones de:",
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
	"A confirmation message has been sent to %s.":                                                               "Se ha enviado un mensaje de confirmación a %s.",
	"A confirmation message has been sent to %s. If you don’t receive it promptly, look for it in your Junk Mail folder. Move it to your inbox so that future messages from us about the class are not marked as Junk Mail.": "Se ha enviado un mensaje de confirmación a %s. Si no lo recibe rápidamente, búsquelo en su carpeta de correo no deseado. Muévalo a su bandeja de entrada para que futuros mensajes nuestros sobre la clase no se marquen como correo no deseado.",
	"If you need to withdraw from the class, please return to this website and remove your registration.  You may also send email to serv@sunnyvale.ca.gov.":                                                                 "Si necesita retirarse de la clase, regrese a este sitio web y vacie su inscripción. También puede enviar un correo electrónico a serv@sunnyvale.ca.gov.",

	// classes/reglogin.go
	"To register for this class, please enter your email address.": "Para inscribirse en esta clase, introduzca su dirección de correo electrónico.",
	"Submit":                                     "Enviar",
	"Your email address is required.":            "Se requiere su dirección de correo electrónico.",
	"This is not a valid email address.":         "Esta dirección de correo electrónico no es válida.",
	"To register for this class, please log in.": "Para inscribirse en esta clase, inicie sesión.",
	"Login":                              "Iniciar sesión",
	"Your password is required.":         "Se requiere su contraseña.",
	"Login incorrect. Please try again.": "Acceso incorrecto. Por favor, inténtelo de nuevo.",

	// errpage/errpage.go
	"No Such Page": "No existe esa página",
	"Sorry, the page you asked for doesn’t exist.  But we have plenty of other good ones!  You can <a href=\"javascript:history.back()\">go back</a> to where you were, or return to <a href=\"/\">the home page</a>.  Look around; you’re sure to find a page you like.": "Lo sentimos, la página que solicitó no existe.  ¡Pero tenemos muchas otras buenas! Puede <a href=\"javascript:history.back()\">volver</a> a donde estaba o regrese a la <a href=\"/\">la página de inicio</a>.  Mire alrededor; Seguro que encontrará una página que le gusta.",
	"Forbidden": "Prohibido",
	"Sorry, but your account doesn’t have permissions for the operation you requested.  If you think you should have permissions, contact <a href=mailto:admin@sunnyvaleserv.org>admin@SunnyvaleSERV.org</a> for assistance.": "Lo sentimos, pero su cuenta no tiene permisos para la operación que usted solicitó.  Si cree que debería tener permisos, póngase en contacto conmailto:admin@sunnyvaleserv.org>admin@SunnyvaleSERV.org</a> para asistencia.",
	"Web Site Error": "Error del sitio web",
	"We’re sorry, but this web site isn’t working correctly right now.  This problem has been reported to the site administrator.  We’ll get it fixed as soon as possible.": "Lo sentimos, pero este sitio web no funciona correctamente en este momento.  Este problema ha sido informado al administrador del sitio.  Lo solucionaremos lo antes posible.",

	// events/*
	"Events":   "Eventos",
	"Calendar": "Calendario",
	"Signups":  "Inscripciones",

	// events/eventscal/eventscal.go
	"SMTWTFS": "DLMMJVS", // initials of days of the week

	// events/eventslist/eventslist.go
	"Date":     "Fecha",
	"Event":    "Evento",
	"Location": "Sitio",
	"TBD":      "Por determinar", // location of event unknown

	// events/eventview/details.go
	"at %s":         "a las %s",   // time
	"from %s to %s": "de %s a %s", // time range
	"Location TBD":  "Sitio por determinar",

	// events/eventview/task.go
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
	"You spent %s volunteer hour.":            "Pasó %s hora de voluntariado.", // %s is "1" or "½"
	"You spent %s volunteer hours.":           "Pasó %s horas de voluntariado.",
	"You did not record volunteer hours.":     "No registró horas de voluntariado.",

	// events/langdate/langdate.go
	"Sunday":    "Domingo",
	"Monday":    "Lunes",
	"Tuesday":   "Martes",
	"Wednesday": "Miércoles",
	"Thursday":  "Jueves",
	"Friday":    "Vieres",
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

	// events/signups/shared.go
	// Describing the number of people signed up and needed for a task:
	"Have %d,": "Tenemos %d,",
	"need %d":  "necesitamos %d",
	"limit %d": "límite %d",
	"no limit": "no hay límite",
	// Reasons for not being able to sign up:
	"Already signed up for a conflicting shift.": "Ya se inscribió a un turno conflictivo.",
	"Signups are closed.":                        "Las inscripciones están cerradas.",
	"Not eligible to sign up.":                   "No es elegible para registrarse.",
	"DSW registration is required.":              "Se requiere registro DSW.",
	"A background check is required.":            "Se requiere una verificación de antecedentes.",
	"The shift has ended.":                       "El turno ha terminado.",
	"The shift has already started.":             "El turno ya ha comenzado.",
	"The shift is full.":                         "El turno está completo.",
	"No person selected.":                        "Ninguna persona seleccionada.",

	// events/signups/signups.go
	"Event Signups": "Inscripciones para eventos",
	"There are no upcoming events with signups.": "No hay eventos próximos con inscripciones.",

	// homepage/homepage.go
	"Sunnyvale Emergency Response Volunteers": "Voluntarios de Respuesta a Emergencias de Sunnyvale",
	"Volunteer Login": "Iniciar sesión",
	// "Profile": "Perfil",
	"Classes and Training":                                     "Clases y capacitación",
	"Preparedness for\nhomes and families":                     "Preparación para su\nfamilia y casa",
	"2 hours\nEnglish and Spanish":                             "2 horas\nespañol e inglés",
	"2 hours\nEnglish Jan. 25\nSpanish Jan. 13":                "2 horas\nespañol 13 enero\ninglés 25 enero",
	"Helping others safely\nin a disaster":                     "Ayudar a otros\nen un desastre",
	"7 weeks\nEnglish only":                                    "7 semanas\ninglés solamente",
	"7 weeks\nEnglish only\nFeb–Mar 2024":                      "7 semanas\ninglés solamente\nfeb–mar 2024",
	"Volunteer Programs":                                       "Programas de voluntariado",
	"Emergency Response Team":                                  "Respuesta en emergencias",
	"Community Emergency Response Team":                        "Equipo comunitario de respuesta a emergencias",
	"Preparedness Education":                                   "Educación de preparación",
	"Listos California: Preparedness Education":                "Listos California: Educación de preparación",
	"Emergency Communications":                                 "Communicaciones en emergencias",
	"Sunnyvale Amateur Radio Emergency Communications Service": "Radioaficionados de Sunnyvale:\nCommunicaciones en emergencias",
	"Neighborhood Preparedness":                                "Preparación del vecindario",
	"Sunnyvale Neighborhoods Actively Prepare":                 "Vecindarios de Sunnyvale se preparan activamente",
	"Information Library":                                      "Archivos y recursos",
	"Contact Us":                                               "Contáctenos",
	"Office of Emergency Services\nDepartment of Public Safety\nCity of Sunnyvale":                                                         "Oficina de Servicios de Emergencia\nDepartamento de Seguridad Pública\nCiudad de Sunnyvale",
	"<a href=\"mailto:serv@sunnyvale.ca.gov\">serv@sunnyvale.ca.gov</a>\n<a href=\"tel:+14087307190\">(408) 730-7190</a>\n(messages only)": "<a href=\"mailto:serv@sunnyvale.ca.gov\">serv@sunnyvale.ca.gov</a>\n<a href=\"tel:+14087307294\">(408) 730-7294</a>\n(mensajes solamente)",
	"Web Site Information": "Información del sitio web",
	// Asset URL for logo for PEP class:
	"pep-logo.png": "ppde-logo.png",

	// people/*
	"(all)":    "(todos)", // as in, all roles, everyone
	"Activity": "Actividad",
	"Map":      "Mapa",
	"Name":     "Nombre",
	"People":   "Personas",

	// people/activity/activity.go:
	"Volunteer Activity": "Actividad de voluntariado",
	"%s Activity":        "Actividad de %s", // %s is a person's name
	"Signed In":          "Registrado",
	"Credited":           "Acreditado",
	"Other %s Hours":     "Otras horas para %s", // %s is a SERV organization name
	"No activity.":       "No hay actividad.",

	// people/peoplelist/peoplelist.go:
	"Sort":              "Ordenar", // button label
	"cell":              "móvil",   // after a phone number
	"home":              "casa",    // after a phone number
	"work":              "trabajo", // after a phone number
	"1 person listed.":  "1 persona en la lista.",
	"%d people listed.": "%d personas en la lista.",

	// people/peoplemap/peoplemap.go:
	"Home":             "En casa",
	"Business":         "A trabajo",
	"(Business Hours)": "(Horas de trabajo)",

	// people/person{edit,view}/contact.go
	"Contact Information":      "Información de contacto",
	"Edit Contact Information": "Editar información de contacto",
	// Email addresses:
	"Email":      "Email",
	"Alt. Email": "Otro email",
	"This is the email address you log in with.":        "Esta es la dirección de correo electrónico con la que inicia sesión.",
	"%q is not a valid email address.":                  "%q no es una dirección de correo electrónico válida.",
	"The email address %q is in use by another person.": "La dirección de correo electrónico %q está siendo utilizada por otra persona.",
	// Phone numbers:
	"Cell Phone":                      "Tel. móvil",
	"(Cell)":                          "(Móvil)", // after the number
	"Home Phone":                      "Tel. de casa",
	"(Home)":                          "(Casa)", // after the number
	"Work Phone":                      "Tel. de trabajo",
	"(Work)":                          "(Trabajo)", // after the number
	"%q is not a valid phone number.": "%q no es un número de teléfono válido.",
	"%q is not a valid 10-digit phone number.":     "%q no es un número de teléfono válido de 10 dígitos.",
	"Another person has the cell phone number %q.": "Otra persona tiene el número de teléfono móvil %q.",
	// Addresses:
	"Home Address (all day)":     "Dirección de casa (todo el día)",
	"Home Address":               "Dirección de casa",
	"Work Address":               "Dirección de trabajo",
	"Mailing Address":            "Dirección de correos",
	"Same as home address":       "Igual que la de casa",
	"Sunnyvale Fire District %d": "Distrito de bomberos %d de Sunnyvale",
	"This address cannot be marked “same as home” when there is no home address.":                       "Esta dirección no se puede marcar como “igual que la de casa” cuando no hay una dirección de casa.",
	"Address changes cannot be accepted right now because the address verification service is offline.": "No se pueden aceptar cambios de dirección en este momento porque el servicio de verificación de dirección está fuera de línea.",
	"This is not a valid address.": "Esta no es una dirección válida.",
	// Emergency contacts:
	"Emergency Contact":                                   "Contacto de emergencias",
	"No emergency contacts on file.":                      "No hay contactos de emergencia registrados.",
	"1 emergency contact on file.":                        "1 contacto de emergencia registrado.",
	"%d emergency contacts on file.":                      "%d contactos de emergencia registrados.",
	"A phone number may not be specified without a name.": "No se puede especificar un número de teléfono sin un nombre.",
	"At least one phone number is required.":              "Se requiere al menos un número de teléfono.",
	"Relationship":                                        "Relación",
	"(select relationship)":                               "(seleccione una relación)",
	"A relationship may not be specified without a name.": "No se puede especificar una relación sin un nombre.",
	"The relationship is required.":                       "Se requiere la relación.",
	"%q is not one of the relationship choices.":          "%q no es una de las opciones de relación.",
	"Co-worker":  "Compañero de trabajo",
	"Daughter":   "Hija",
	"Father":     "Padre",
	"Friend":     "Amigo",
	"Mother":     "Madre",
	"Neighbor":   "Vecino",
	"Other":      "Otro",
	"Relative":   "Pariente",
	"Son":        "Hijo",
	"Spouse":     "Cónyuge",
	"Supervisor": "Supervisor",

	// people/personedit/names.go
	"Edit Names":            "Editar nombres",
	"The name is required.": "Se requiere el nombre.",
	"What you like to be called, e.g. “Joe Banks”": "Cómo le gusta que le llamen, p.e. “Paco García”",
	"Formal name":                  "Nombre formal",
	"The formal name is required.": "Se requiere el nombre formal.",
	"For formal documents, e.g. “Joseph A. Banks, Jr.”": "Para documentos formales, p.e. “Francisco García Ramírez”",
	"Sort name": "Nombre ordenado",
	"For appearance in sorted lists, e.g. “Banks, Joe”": "Para aparecer en listas ordenadas, p.e. “García, Paco”",
	"The sort name is required.":                        "Se requiere el nombre ordenado.",
	"Another person has the sort name %q.":              "Otra persona tiene el nombre ordenado %q.",
	"Call sign":                                         "Indicativo",
	"FCC amateur radio license (if any)":                "Indicativo de licencia de radioaficionado de la FCC (si corresponde)",
	"%q is not a valid FCC amateur radio call sign.":    "%q no es un indicativo válido para radioaficionados de la FCC.",
	"Another person has the call sign %q.":              "Otra persona tiene el indicativo %q.",
	"Birthdate":                                         "Fecha de nacimiento",
	"Pronouns":                                          "Pronumbres",
	"he/him/his":                                        "él/lo",
	"she/her/hers":                                      "ella/la",
	"they/them/theirs":                                  "elle/le",

	// people/personview/notes.go
	"Notes": "Notas",

	// people/personview/password.go
	"Change Password":                       "Cambiar contraseña", // Button label
	"Password":                              "Contraseña",
	"Password Change":                       "Cambiar de contraseña",
	"Old Password":                          "Contraseña anterior",
	"Please specify your old password.":     "Por favor ingrese su contraseña anterior.",
	"This is not the correct old password.": "Esta no es la contraseña anterior correcta.",
	"New Password":                          "Contraseña nueva",
	"Please specify a valid new password.":  "Por favor ingrese una nueva contraseña válida.",
	"The new password is too weak.":         "La nueva contraseña es demasiado débil.",

	// people/personedit/pwreset.go
	// Email sent when a password is reset by an administrator.  "\n" is a newline.
	"From: %s\nSubject: SunnyvaleSERV.org Password Reset\nContent-Type: text/plain; charset=utf8\n\nHello, %s,\n\n%s has reset the password for your account on SunnyvaleSERV.org.  Your new login information is:\n\n    Email:    %s\n    Password: %s\n\nThis password is three words chosen randomly from a dictionary — a method that generally produces a very secure and often memorable password.  If the resulting phrase has any meaning, it’s unintentional coincidence.\n\nYou can change this password by logging into SunnyvaleSERV.org and clicking the “Change Password” button on your Profile page.  If you have any questions, just reply to this email.\n\nRegards,\nSunnyvaleSERV.org\n": "From: %s\nSubject: SunnyvaleSERV.org restablecimiento de contraseña\nContent-Type: text/plain; charset=utf8\n\nHola, %s,\n\n%s ha restablecido la contraseña de su cuenta en SunnyvaleSERV.org.  Su información nueva es:\n\n    Email:      %s\n    Contraseña: %s\n\nEsta contraseña es tres palabras elegidas al azar de un diccionario inglés.  Este es un método que generalmente produce una contraseña muy segura y, a menudo, fácil de recordar.  Si la frase resultante tiene algún significado, es una coincidencia involuntaria.\n\nPuede cambiar esta contraseña iniciando sesión en SunnyvaleSERV.org y haciendo clic en el botón “Cambiar contraseña” en su página de perfil.  Si tiene alguna pregunta, simplemente responda a este mensaje.\n\nSaludos,\nSunnyvaleSERV.org\n",

	// people/personview/roles.go
	"No current role in any SERV org.": "No tiene ningún papel actual en ninguna organization de SERV.",
	"SERV Role":                        "Papel en SERV",
	"SERV Roles":                       "Papeles en SERV",

	// people/personview/status.go
	"Volunteer Status": "Estado del voluntario",
	// Volunteer registration:
	"City volunteer":       "Voluntario de la ciudad",
	"Registration pending": "Registro pendiente",
	// DSW registration:
	"Not registered":            "No está registrado",
	"Registered %s":             "Registrado el %s",                 // %s is date
	"Registered %s, expires %s": "Registrado el %s, caducará el %s", // %s is date
	"Expired on %s":             "Caducó el %s",                     // %s is date
	// Background check:
	"Background check": "Verificación de antecedentes",
	"Cleared":          "Aprobada",
	"Needed":           "Necesaria",

	// people/person{edit,view}/subscriptions.go
	"Subscriptions":           "Suscripciones",
	"Edit List Subscriptions": "Editar suscripciones",
	"Unsubscribe All":         "Desuscribirse a todos", // button label
	"Not subscribed to any email or text messaging.": "No está suscrito a ningún correo electrónico o mensaje de texto.",
	"Unsubscribed from all email.":                   "Se ha desuscribido de todos los correos electrónicos.",
	"Unsubscribed from all text messaging.":          "Se ha desuscribido de todos los mensajes de texto.",
	// The next three are used when there are 1, 2, and 3-or-more roles affected:
	"Messages sent to %s are considered required for the “%s” role.  Unsubscribing from it may cause you to lose that role.":              "Los mensajes enviados a %s se consideran obligatorios para el papel “%s”.  Desuscribirse puede hacer que pierda ese papel.",
	"Messages sent to %s are considered required for the “%s” and “%s” roles.  Unsubscribing from it may cause you to lose those roles.":  "Los mensajes enviados a %s se consideran obligatorios para los papeles “%s” y “%s”.  Desuscribirse puede hacer que pierda esos papeles.",
	"Messages sent to %s are considered required for the “%s”, and “%s” roles.  Unsubscribing from it may cause you to lose those roles.": "Los mensajes enviados a %s se consideran obligatorios para los papeles “%s”, y “%s”.  Desuscribirse puede hacer que pierda esos papeles.",

	// people/personedit/vregister.go
	"Register as a City Volunteer": "Registrarse como voluntario de ciudad",
	"Register":                     "Registrarse", // button label
	"Thank you for your interest in volunteering with the City of Sunnyvale, Office of Emergency Services.  Please complete this form to register as a City of Sunnyvale Volunteer.  Once we receive your registration (which usually takes a few days) we will contact you to schedule an appointment for your fingerprinting.  (Please note: registering as a city volunteer is not required for taking one of our classes.  It is only required when joining one of our volunteer groups.)": "Gracias por su interés en ser voluntario en la Oficina de Servicios de Emergencia de la ciudad de Sunnyvale.  Complete este formulario para registrarse como voluntario de la ciudad de Sunnyvale.  Una vez que recibamos su registro (lo que generalmente demora unos días), nos comunicaremos con usted para programar una cita para su toma de huellas digitales.  (Tenga en cuenta: no es necesario registrarse como voluntario de la ciudad para tomar una de nuestras clases.  Solo es necesario cuando se une a uno de nuestros grupos de voluntarios).",
	"Interests":                             "Intereses",
	"CERT Deployment Team":                  "Equipo de despliegue CERT",
	"Community Outreach":                    "Alcance comunitario",
	"Amateur Radio (SARES)":                 "Radioaficionados (SARES)",
	"Neighborhood Preparedness Facilitator": "Facilitador de preparación vecinal",
	"Preparedness Class Instructor":         "Instructor de clases de preparación",
	"CERT Basic Training Instructor":        "Instructor de clases CERT",
	"By submitting this application, I certify that all statements I have made on this application are true and correct and I hereby authorize the City of Sunnyvale to investigate the accuracy of this information.  I am aware that fingerprinting and a criminal records search is required for volunteers 18 years of age or older.  I understand that I am working at all times on a voluntary basis, without monetary compensation or benefits, and not as a paid employee.  I give the City of Sunnyvale permission to use any photographs or videos taken of me during my service without obligation or compensation to me.  I understand that the City of Sunnyvale reserves the right to terminate a volunteer's service at any time.  I understand that volunteers are covered under the City of Sunnyvale's Worker's Compensation Program for an injury or accident occurring while on duty.": "Al enviar esta solicitud, certifico que todas las declaraciones que he hecho en esta solicitud son verdaderas y correctas y por la presente autorizo a la ciudad de Sunnyvale a investigar la exactitud de esta información.  Soy consciente de que se requieren huellas dactilares y una búsqueda de antecedentes penales para los voluntarios mayores de 18 años.  Entiendo que estoy trabajando en todo momento de forma voluntaria, sin compensación monetaria ni beneficios, y no como empleado remunerado.  Doy permiso a la ciudad de Sunnyvale para utilizar fotografías o videos tomados de mí durante mi servicio sin obligación ni compensación para mí.  Entiendo que la ciudad de Sunnyvale se reserva el derecho de cancelar el servicio de un voluntario en cualquier momento.  Entiendo que los voluntarios están cubiertos por el Programa de Compensación para Trabajadores de la Ciudad de Sunnyvale por una lesión o accidente que ocurra mientras están de servicio.",
	"I agree":                     "Estoy de acuerdo",
	"Your birthdate is required.": "Se requiere su fecha de nacimiento.",
	"A cell or home phone number is required.":                                   "Se requiere un número de teléfono móvil o a casa.",
	"Your home address is required.":                                             "Se requiere su dirección de casa.",
	"The emergency contacts are required.":                                       "Se requieren los contactos de emergencia.",
	"Please check that you agree with the above statement in order to register.": "Por favor, marque la casilla para mostrar que está de acuerdo con la declaración anterior para poder registrarse.",

	// search/search.go
	"Search":                       "Buscar", // button label
	"Documents":                    "Archivos",
	"Folders":                      "Carpetas",
	"in folder":                    "en la carpeta",
	"Venues":                       "Sitios",
	"Nothing matched your search.": "No se encontró nada en su búsqueda.",

	// static/calendar.go
	"SERV Calendar Subscription": "Suscripción al calendario de SERV",

	// static/emaillists.go
	"SERV Email Lists": "Listas de correo electrónico de SERV",

	// store/class/referral.go
	"Word of mouth":                 "Boca a boca",
	"Information table at an event": "Mesa informativa en un evento",
	"Printed advertisement":         "Publicidad impresa",
	"Online advertisement":          "Publicidad en línea",

	// ui/orgdot/orgdot.go
	"CERT Deployment": "Despliegue de CERT",
	"CERT Training":   "Capacitación CERT",

	// ui/page.go
	// Main menu items:
	"Welcome":         "Bienvenido",
	"View in English": "Vea en español",
	// "Events":  "Eventos",
	// "People": "Personas",
	"Files":         "Archivos",
	"Profile":       "Perfil",
	"Logout":        "Cerrar",
	"Web Site Info": "Info. del sitio",
}

func spanishDate(day time.Time) string {
	return fmt.Sprintf("%s, %d de %s de %d", spanish[day.Weekday().String()], day.Day(), spanish[day.Month().String()], day.Year())
}
