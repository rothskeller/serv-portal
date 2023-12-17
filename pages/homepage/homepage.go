package homepage

import (
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/store/folder"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/ui"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

// Serve handles "/" requests.
func Serve(r *request.Request) {
	var user = auth.SessionUser(r, 0, false)
	ui.Page(r, user, ui.PageOpts{}, func(main *htmlb.Element) {
		main.A("class=home")
		homeHeading(r, main)
		homeTopButtons(r, main, user)
		blocks := main.E("div class=homeBlocks")
		homeClasses(r, blocks)
		homePrograms(r, blocks)
		homeLibrary(r, user, blocks)
		homeContact(r, blocks)
		main.E("div class=homeSpacer")
		homeBottomLinks(r, main)
		homeContact2(r, main)
	})
}

func homeHeading(r *request.Request, main *htmlb.Element) {
	main.E("div class=homeHeading").R(r.LangString(
		"Sunnyvale Emergency Response Volunteers",
		"Voluntarios de Respuesta a Emergencias de Sunnyvale",
	))
}

func homeTopButtons(r *request.Request, main *htmlb.Element, user *person.Person) {
	buttons := main.E("div class=homeTButtons")
	if user == nil {
		buttons.E("a href=/login class='sbtn sbtn-primary sbtn-xsmall'").R(r.LangString(
			"Volunteer Login", "Iniciar sesión",
		))
	} else {
		buttons.E("a href=/people/%d class='sbtn sbtn-primary sbtn-xsmall'", user.ID()).R(r.LangString(
			"Profile", "Perfil",
		))
	}
	if r.Language == "es" {
		buttons.E("a href=/en class='sbtn sbtn-primary sbtn-xsmall'>View in English")
	} else {
		buttons.E("a href=/es class='sbtn sbtn-primary sbtn-xsmall'>Vea en español")
	}
}

func homeClasses(r *request.Request, blocks *htmlb.Element) {
	block := blocks.E("div class=homeBlock")
	block.E("div class=homeBlockTitle").R(r.LangString(
		"Classes and Training", "Clases y capacitación",
	))
	classes := block.E("div class=homeClasses")
	pep := classes.E("a href=/pep class=homeClass id=homeClassPEP")
	pep.E("div class=homeClassImg").
		E("img id=homeClassImgPEP src=%s", ui.AssetURL(r.LangString("pep-logo.png", "ppde-logo.png")))
	pep.E("div class=homeClassSlug id=homeClassSlugPEP").R(r.LangString(
		"Preparedness for\nhomes and families",
		"Preparación para su\nfamilia y casa",
	))
	pep.E("div class=homeClassInfoShort id=homeClassInfoShortPEP").R(r.LangString(
		"2 hours\nEnglish and Spanish",
		"2 horas\nespañol e inglés",
	))
	pep.E("div class=homeClassInfoLong id=homeClassInfoLongPEP").R(r.LangString(
		"2 hours\nEnglish Jan. 25\nSpanish Jan. 13",
		"2 horas\nespañol 13 enero\ninglés 25 enero",
	))
	cert := classes.E("a href=/cert-basic class=homeClass id=homeClassCERT")
	cert.E("div class=homeClassImg").
		E("img id=homeClassImgCERT src=%s", ui.AssetURL("cert-logo.png"))
	cert.E("div class=homeClassSlug id=homeClassSlugCERT").R(r.LangString(
		"Helping others safely\nin a disaster",
		"Ayudar a otros\nen un desastre",
	))
	cert.E("div class=homeClassInfoShort id=homeClassInfoShortCERT").R(r.LangString(
		"7 weeks\nEnglish only",
		"7 semanas\ninglés solamente",
	))
	cert.E("div class=homeClassInfoLong id=homeClassInfoLongCERT").R(r.LangString(
		"7 weeks\nEnglish only\nFeb–Mar 2024",
		"7 semanas\ninglés solamente\nfeb–mar 2024",
	))
}

func homePrograms(r *request.Request, blocks *htmlb.Element) {
	block := blocks.E("div class=homeBlock")
	block.E("div class=homeBlockTitle").R(r.LangString(
		"Volunteer Programs", "Programas de voluntariado",
	))
	programs := block.E("div class=homePrograms")
	cert := programs.E("a href=/cert class=homeProgram id=homeProgramCERT")
	cert.E("div class=homeProgramBadge").
		E("img id=homeProgramBadgeCERT src=%s", ui.AssetURL("cert-badge.png"))
	cert.E("div class=homeProgramSlugShort id=homeProgramSlugShortCERT").R(r.LangString(
		"Emergency Response Team",
		"Respuesta en emergencias",
	))
	cert.E("div class=homeProgramSlugLong id=homeProgramSlugLongCERT").R(r.LangString(
		"Community Emergency Response Team",
		"Equipo comunitario de respuesta en emergencias",
	))
	listos := programs.E("a href=/listos class=homeProgram id=homeProgramListos")
	listos.E("div class=homeProgramBadge").
		E("img id=homeProgramBadgeListos src=%s", ui.AssetURL("listos-badge.png"))
	listos.E("div class=homeProgramSlugShort id=homeProgramSlugShortListos").R(r.LangString(
		"Preparedness Education",
		"Educación de preparación",
	))
	listos.E("div class=homeProgramSlugLong id=homeProgramSlugLongListos").R(r.LangString(
		"Listos California: Preparedness Education",
		"Listos California: Educación de preparación",
	))
	sares := programs.E("a href=/sares class=homeProgram id=homeProgramSARES")
	sares.E("div class=homeProgramBadge").
		E("img id=homeProgramBadgeSARES src=%s", ui.AssetURL("sares-badge.png"))
	sares.E("div class=homeProgramSlugShort id=homeProgramSlugShortSARES").R(r.LangString(
		"Emergency Communications",
		"Communicaciones en emergencias",
	))
	sares.E("div class=homeProgramSlugLong id=homeProgramSlugLongSARES").R(r.LangString(
		"Sunnyvale Amateur Radio Emergency Communications Service",
		"Radioaficionados de Sunnyvale:\nCommunicaciones en emergencias",
	))
	snap := programs.E("a href=/snap class=homeProgram id=homeProgramSNAP")
	snap.E("div class=homeProgramBadge").
		E("img id=homeProgramBadgeSNAP src=%s", ui.AssetURL("snap-badge.png"))
	snap.E("div class=homeProgramSlugShort id=homeProgramSlugShortSNAP").R(r.LangString(
		"Neighborhood Preparedness",
		"Preparación del vecindario",
	))
	snap.E("div class=homeProgramSlugLong id=homeProgramSlugLongSNAP").R(r.LangString(
		"Sunnyvale Neighborhoods Actively Prepare",
		"Vecindarios de Sunnyvale se preparan activamente",
	))
}

func homeLibrary(r *request.Request, user *person.Person, blocks *htmlb.Element) {
	const folderFields = folder.FID | folder.FName | folder.FViewer | folder.FParent | folder.FURLName
	blocks.E("div class=homeLibraryShort").
		E("a href=/files class='sbtn sbtn-primary'").R(r.LangString(
		"Information Library", "Archivos y recursos",
	))
	block := blocks.E("div class=homeBlock id=homeBlockLibraryLong")
	block.E("div class=homeBlockTitle").R(r.LangString(
		"Information Library", "Archivos y recursos",
	))
	library := block.E("div class=homeLibraryLong")
	folder.AllWithParent(r, folder.RootID, folderFields, func(f *folder.Folder) {
		if user.HasPrivLevel(f.Viewer()) {
			library.E("div class=folder data-id=%d", f.ID()).
				E("s-icon icon=folder").P().
				E("a href=%s up-target=main", f.Path(r)).
				T(f.Name())
		}
	})
}

func homeContact(r *request.Request, blocks *htmlb.Element) {
	block := blocks.E("div class=homeBlock id=homeBlockContact")
	block.E("div class=homeBlockTitle").R(r.LangString("Contact Us", "Contáctenos"))
	contact := block.E("div class=homeContact")
	contact.R(r.LangString(
		`Office of Emergency Services
Department of Public Safety
City of Sunnyvale

<a href="mailto:oes@sunnyvale.ca.gov">oes@sunnyvale.ca.gov</a>
<a href="tel:+14087307190">(408) 730-7190</a>
(messages only)`,
		`Oficina de Servicios de Emergencia
Departmento de Seguridad Pública
Ciudad de Sunnyvale

<a href="mailto:oes@sunnyvale.ca.gov">oes@sunnyvale.ca.gov</a>
<a href="tel:+14087307294">(408) 730-7294</a>
(mensajes solamente)`,
	))
}

func homeBottomLinks(r *request.Request, main *htmlb.Element) {
	links := main.E("div class=homeLinks")
	links.E("a href=/contact").R(r.LangString("Contact Us", "Contáctenos"))
	links.E("span> • ")
	links.E("a href=/about").R(r.LangString(
		"Web Site Information", "Información del sitio web",
	))
}

func homeContact2(r *request.Request, main *htmlb.Element) {
	contact := main.E("div class=homeContact2")
	contact.E("div").R(r.LangString(
		`Office of Emergency Services
Department of Public Safety
City of Sunnyvale`,
		`Oficina de Servicios de Emergencia
Departmento de Seguridad Pública
Ciudad de Sunnyvale`,
	))
	contact.E("div").R(r.LangString(
		`<a href="mailto:oes@sunnyvale.ca.gov">oes@sunnyvale.ca.gov</a>
<a href="tel:+14087307190">(408) 730-7190</a>
(messages only)`,
		`<a href="mailto:oes@sunnyvale.ca.gov">oes@sunnyvale.ca.gov</a>
<a href="tel:+14087307294">(408) 730-7294</a>
(mensajes solamente)`,
	))
}
