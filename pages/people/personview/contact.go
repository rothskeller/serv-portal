package personview

import (
	"strings"

	"sunnyvaleserv.org/portal/store/enum"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

const contactPersonFields = person.FEmail | person.FEmail2 | person.FCellPhone | person.FHomePhone | person.FWorkPhone | person.FAddresses | person.FEmContacts

func showContact(r *request.Request, main *htmlb.Element, user, p *person.Person, viewLevel person.ViewLevel) {
	editable := user.ID() == p.ID() || user.HasPrivLevel(0, enum.PrivLeader)
	if !editable && viewLevel < person.ViewWorkContact {
		return
	}
	if p.Email() == "" && p.Email2() == "" && p.CellPhone() == "" && p.HomePhone() == "" && p.WorkPhone() == "" &&
		p.Addresses().Home == nil && p.Addresses().Work == nil && p.Addresses().Mail == nil && !editable {
		return // nothing to show
	}
	if viewLevel == person.ViewWorkContact && (p.Email() == "" && p.Email2() == "" && p.WorkPhone() == "" && p.Addresses().Work == nil) && !editable {
		return // nothing to show
	}
	section := main.E("div class='personviewSection personviewContact'")
	sheader := section.E("div class=personviewSectionHeader")
	sheader.E("div class=personviewSectionHeaderText").R(r.LangString("Contact Information", "Información del contacto"))
	if editable {
		sheader.E("div class=personviewSectionHeaderEdit").
			E("a href=/people/%d/edcontact up-layer=new up-size=grow up-dismissable=key up-history=false class='sbtn sbtn-small sbtn-primary'", p.ID()).R(r.LangString("Edit", "Editar"))
	}
	if p.Email() != "" || p.Email2() != "" {
		ediv := section.E("div class=personviewContactEmails")
		if p.Email() != "" {
			ediv.E("div").E("a href=mailto:%s target=_blank>%s", p.Email(), p.Email())
		}
		if p.Email2() != "" {
			ediv.E("div").E("a href=mailto:%s target=_blank>%s", p.Email2(), p.Email2())
		}
	}
	if viewLevel == person.ViewFull && p.CellPhone() != "" {
		section.E("div class=personviewContactPhone").
			E("a href=tel:%s target=_blank>%s", p.CellPhone(), p.CellPhone()).
			P().E("span class=personviewContactPhoneLabel").R(r.LangString("(Cell)", "(Móvil)"))
	}
	if viewLevel == person.ViewFull && p.HomePhone() != "" {
		section.E("div class=personviewContactPhone").
			E("a href=tel:%s target=_blank>%s", p.HomePhone(), p.HomePhone()).
			P().E("span class=personviewContactPhoneLabel").R(r.LangString("(Home)", "(Casa)"))
	}
	if p.WorkPhone() != "" {
		section.E("div class=personviewContactPhone").
			E("a href=tel:%s target=_blank>%s", p.WorkPhone(), p.WorkPhone()).
			P().E("span class=personviewContactPhoneLabel").R(r.LangString("(Work)", "(Trabajo)"))
	}
	if p.Addresses().Work != nil && p.Addresses().Work.SameAsHome {
		showAddress(r, section, p.Addresses().Home, r.LangString("Home Address (all day)", "Direccion de casa (todo el día)"), true)
	} else if viewLevel == person.ViewFull {
		showAddress(r, section, p.Addresses().Home, r.LangString("Home Address", "Direccion de casa"), true)
	}
	showAddress(r, section, p.Addresses().Work, r.LangString("Work Address", "Direccion de trabajo"), true)
	if viewLevel == person.ViewFull {
		showAddress(r, section, p.Addresses().Mail, r.LangString("Mailing Address", "Direccion de correos"), false)
	}
	if editable {
		switch len(p.EmContacts()) {
		case 0:
			section.E("div class=personviewContactEmerg").R(r.LangString("No emergency contacts on file.", "No hay contactos de emergencia registrados."))
		case 1:
			section.E("div class=personviewContactEmerg").R(r.LangString("1 emergency contact on file.", "1 contacto de emergencia registrado."))
		default:
			section.E("div class=personviewContactEmerg>%d", len(p.EmContacts())).R(r.LangString(" emergency contacts on file.", " contactos de emergencia registrados."))
		}
	}
}

func showAddress(r *request.Request, section *htmlb.Element, address *person.Address, label string, showMap bool) {
	if address == nil || address.SameAsHome {
		return
	}
	div := section.E("div class=personviewContactAddress")
	labeldiv := div.E("div>%s:", label)
	if showMap {
		labeldiv.E("a href=https://www.google.com/maps/search/?api=1&query=%s class=personviewContactAddressMap target=_blank", address.Address).R(r.LangString("Map", "Mapa"))
	}
	parts := strings.SplitN(address.Address, ",", 2)
	div.E("div>%s", parts[0])
	if len(parts) == 2 {
		div.E("div>%s", parts[1])
	}
	if address.FireDistrict != 0 {
		div.E(r.LangString("div>Sunnyvale Fire District %d", "div>Distrito de bomberos %d de Sunnyvale"), address.FireDistrict)
	}
}
