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
	ui.Page(r, user, ui.PageOpts{NoHome: true}, func(main *htmlb.Element) {
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
	main.E("div class=homeHeading").R(r.Loc("Sunnyvale Emergency Response Volunteers"))
}

func homeTopButtons(r *request.Request, main *htmlb.Element, user *person.Person) {
	buttons := main.E("div class=homeTButtons")
	if user == nil {
		buttons.E("a href=/login class='sbtn sbtn-primary sbtn-xsmall'").R(r.Loc("Volunteer Login"))
	} else {
		buttons.E("a href=/people/%d class='sbtn sbtn-primary sbtn-xsmall'", user.ID()).R(r.Loc("Profile"))
	}
	if r.Language == "es" {
		buttons.E("a href=/en class='sbtn sbtn-primary sbtn-xsmall'>View in English")
	} else {
		buttons.E("a href=/es class='sbtn sbtn-primary sbtn-xsmall'>Ver en español")
	}
}

func homeClasses(r *request.Request, blocks *htmlb.Element) {
	block := blocks.E("div class=homeBlock")
	block.E("div class=homeBlockTitle").R(r.Loc("Classes and Training"))
	classes := block.E("div class=homeClasses")
	pep := classes.E("a href=/pep class=homeClass id=homeClassPEP")
	pep.E("div class=homeClassImg").E("img id=homeClassImgPEP src=%s", ui.AssetURL(r.Loc("pep-logo.png")))
	text := pep.E("div class=homeClassText")
	text.E("div class=homeClassTitle").T(r.Loc("Disaster preparedness for homes and families"))
	text.E("div class=homeClassInfo").T(r.Loc("2 hours, English or Spanish"))
	cert := classes.E("a href=/cert-basic class=homeClass id=homeClassCERT")
	cert.E("div class=homeClassImg").E("img id=homeClassImgCERT src=%s", ui.AssetURL("cert-logo.png"))
	text = cert.E("div class=homeClassText")
	text.E("div class=homeClassTitle").T(r.Loc("Helping others safely in a disaster"))
	text.E("div class=homeClassInfo").T(r.Loc("7 weeks, English only"))
}

func homePrograms(r *request.Request, blocks *htmlb.Element) {
	block := blocks.E("div class=homeBlock")
	block.E("div class=homeBlockTitle").R(r.Loc("Volunteer Programs"))
	programs := block.E("div class=homePrograms")
	cert := programs.E("a href=/cert class=homeProgram id=homeProgramCERT")
	cert.E("div class=homeProgramBadge").
		E("img id=homeProgramBadgeCERT src=%s", ui.AssetURL("cert-badge.png"))
	cert.E("div class=homeProgramSlug id=homeProgramSlugCERT").R(r.Loc("Community Emergency Response Team"))
	listos := programs.E("a href=/listos class=homeProgram id=homeProgramListos")
	listos.E("div class=homeProgramBadge").
		E("img id=homeProgramBadgeListos src=%s", ui.AssetURL("listos-badge.png"))
	listos.E("div class=homeProgramSlug id=homeProgramSlugListos").R(r.Loc("Listos California: Preparedness Education"))
	sares := programs.E("a href=/sares class=homeProgram id=homeProgramSARES")
	sares.E("div class=homeProgramBadge").
		E("img id=homeProgramBadgeSARES src=%s", ui.AssetURL("sares-badge.png"))
	sares.E("div class=homeProgramSlug id=homeProgramSlugSARES").R(r.Loc("Sunnyvale Amateur Radio Emergency Communications Service"))
	snap := programs.E("a href=/snap class=homeProgram id=homeProgramSNAP")
	snap.E("div class=homeProgramBadge").
		E("img id=homeProgramBadgeSNAP src=%s", ui.AssetURL("snap-badge.png"))
	snap.E("div class=homeProgramSlug id=homeProgramSlugSNAP").R(r.Loc("Sunnyvale Neighborhoods Actively Prepare"))
}

func homeLibrary(r *request.Request, user *person.Person, blocks *htmlb.Element) {
	const folderFields = folder.FID | folder.FName | folder.FViewer | folder.FParent | folder.FURLName
	block := blocks.E("div class=homeBlock")
	block.E("div class=homeBlockTitle").R(r.Loc("Information Library"))
	library := block.E("div class=homeLibrary")
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
	block.E("div class=homeBlockTitle").R(r.Loc("Contact Us"))
	contact := block.E("div class=homeContact")
	contact.E("img class=homeContactImg src=%s", ui.AssetURL("sny-logo.png"))
	text := contact.E("div class=homeContactText")
	text.R(r.Loc("Office of Emergency Services\nDepartment of Public Safety\nCity of Sunnyvale"))
	text.R("\n\n<a href=\"mailto:serv@sunnyvale.ca.gov\">serv@sunnyvale.ca.gov</a>\n")
	text.R(r.Loc("<a href=\"tel:+14087307190\">(408) 730-7190</a>"))
	text.R("\n")
	text.R(r.Loc("(messages only)"))
}

func homeBottomLinks(r *request.Request, main *htmlb.Element) {
	links := main.E("div class=homeLinks")
	links.E("a href=/contact").R(r.Loc("Contact Us"))
	links.E("span> • ")
	links.E("a href=/about").R(r.Loc("Web Site Information"))
}

func homeContact2(r *request.Request, main *htmlb.Element) {
	contact := main.E("div class=homeContact2")
	contact.E("img class=homeContactImg2 src=%s", ui.AssetURL("sny-logo.png"))
	contact.E("div").R(r.Loc("Office of Emergency Services\nDepartment of Public Safety\nCity of Sunnyvale"))
	div := contact.E("div")
	div.R("<a href=\"mailto:serv@sunnyvale.ca.gov\">serv@sunnyvale.ca.gov</a>\n")
	div.R(r.Loc("<a href=\"tel:+14087307190\">(408) 730-7190</a>"))
	div.R("\n")
	div.R(r.Loc("(messages only)"))
}
