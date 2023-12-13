package personview

import (
	"sunnyvaleserv.org/portal/store/enum"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

const notesPersonFields = person.FNotes

func showNotes(r *request.Request, main *htmlb.Element, user, p *person.Person, canViewContactInfo bool) {
	var (
		section  *htmlb.Element
		notes    = p.Notes()
		editable = user.HasPrivLevel(0, enum.PrivLeader)
	)
	for i := len(notes) - 1; i >= 0; i-- {
		n := notes[i]
		switch n.Visibility {
		case person.NoteVisibleToWebmaster:
			if !user.IsWebmaster() {
				continue
			}
		case person.NoteVisibleToAdmins:
			if !user.IsAdminLeader() {
				continue
			}
		case person.NoteVisibleToLeaders:
			if !user.HasPrivLevel(0, enum.PrivLeader) {
				continue
			}
		case person.NoteVisibleWithContact:
			if !canViewContactInfo {
				continue
			}
		}
		section = startNotes(main, section, p, editable)
		ndiv := section
		if editable {
			ndiv = ndiv.E("a href=/people/%d/ednote/%d up-layer=new up-size=grow up-dismissable=key up-history=false", p.ID(), i)
		}
		ndiv = ndiv.E("div class=personviewNote")
		ndiv.E("div class=personviewNoteDate").R(n.Date.Format("2006-01-02"))
		ndiv.E("div class=personviewNoteText>%s", n.Note)
	}
	if editable {
		if section == nil {
			startNotes(main, nil, p, true).E("div>No notes on file.")
		} else {
			section.E("div class=personviewNotesHelp>(Click a note to edit or remove it.)")
		}
	}
}

func startNotes(main *htmlb.Element, section *htmlb.Element, p *person.Person, editable bool) *htmlb.Element {
	if section == nil {
		section = main.E("div class=personviewSection")
		sheader := section.E("div class=personviewSectionHeader")
		sheader.E("div class=personviewSectionHeaderText>Notes")
		if editable {
			sheader.E("div class=personviewSectionHeaderEdit").
				E("a href=/people/%d/ednote up-layer=new up-size=grow up-dismissable=key up-history=false class='sbtn sbtn-small sbtn-primary'>Add", p.ID())
		}
		section = section.E("div class=personviewNotes")
	}
	return section
}
