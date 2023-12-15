package person

import (
	"strings"
	"time"

	"sunnyvaleserv.org/portal/store/enum"
	"sunnyvaleserv.org/portal/store/internal/phys"
)

const (
	badLoginTimeFormat = "2006-01-02T15:04:05"
	pwresetTimeFormat  = "2006-01-02T15:04:05"
	bgCheckDateFormat  = "2006-01-02"
	dswRegDateFormat   = "2006-01-02"
	noteDateFormat     = "2006-01-02"
)

// ColumnList generates a comma-separated list of column names for the specified
// person fields.  It is used in constructing SQL SELECT statements.
func ColumnList(sb *strings.Builder, fields Fields) {
	sep := phys.NewSeparator(", ")
	if fields&FID != 0 {
		sb.WriteString(sep())
		sb.WriteString("p.id")
	}
	if fields&FVolgisticsID != 0 {
		sb.WriteString(sep())
		sb.WriteString("p.volgistics_id")
	}
	if fields&FInformalName != 0 {
		sb.WriteString(sep())
		sb.WriteString("p.informal_name")
	}
	if fields&FFormalName != 0 {
		sb.WriteString(sep())
		sb.WriteString("p.formal_name")
	}
	if fields&FSortName != 0 {
		sb.WriteString(sep())
		sb.WriteString("p.sort_name")
	}
	if fields&FCallSign != 0 {
		sb.WriteString(sep())
		sb.WriteString("p.call_sign")
	}
	if fields&FPronouns != 0 {
		sb.WriteString(sep())
		sb.WriteString("p.pronouns")
	}
	if fields&FEmail != 0 {
		sb.WriteString(sep())
		sb.WriteString("p.email")
	}
	if fields&FEmail2 != 0 {
		sb.WriteString(sep())
		sb.WriteString("p.email2")
	}
	if fields&FCellPhone != 0 {
		sb.WriteString(sep())
		sb.WriteString("p.cell_phone")
	}
	if fields&FHomePhone != 0 {
		sb.WriteString(sep())
		sb.WriteString("p.home_phone")
	}
	if fields&FWorkPhone != 0 {
		sb.WriteString(sep())
		sb.WriteString("p.work_phone")
	}
	if fields&FPassword != 0 {
		sb.WriteString(sep())
		sb.WriteString("p.password")
	}
	if fields&FBadLoginCount != 0 {
		sb.WriteString(sep())
		sb.WriteString("p.bad_login_count")
	}
	if fields&FBadLoginTime != 0 {
		sb.WriteString(sep())
		sb.WriteString("p.bad_login_time")
	}
	if fields&FPWResetToken != 0 {
		sb.WriteString(sep())
		sb.WriteString("p.pwreset_token")
	}
	if fields&FPWResetTime != 0 {
		sb.WriteString(sep())
		sb.WriteString("p.pwreset_time")
	}
	if fields&FUnsubscribeToken != 0 {
		sb.WriteString(sep())
		sb.WriteString("p.unsubscribe_token")
	}
	if fields&FHoursToken != 0 {
		sb.WriteString(sep())
		sb.WriteString("p.hours_token")
	}
	if fields&FIdentification != 0 {
		sb.WriteString(sep())
		sb.WriteString("p.identification")
	}
	if fields&FBirthdate != 0 {
		sb.WriteString(sep())
		sb.WriteString("p.birthdate")
	}
	if fields&FLanguage != 0 {
		sb.WriteString(sep())
		sb.WriteString("p.language")
	}
	if fields&FFlags != 0 {
		sb.WriteString(sep())
		sb.WriteString("p.flags")
	}
	if fields&FAddresses != 0 {
		panic("FAddresses cannot be fetched with ColumnList/Scan")
	}
	if fields&FBGChecks != 0 {
		panic("FBGChecks cannot be fetched with ColumnList/Scan")
	}
	if fields&FDSWRegistrations != 0 {
		panic("FDSWRegistrations cannot be fetched with ColumnList/Scan")
	}
	if fields&FNotes != 0 {
		panic("FNotes cannot be fetched with ColumnList/Scan")
	}
	if fields&FEmContacts != 0 {
		panic("FEmContacts cannot be fetched with ColumnList/Scan")
	}
	if fields&FPrivLevels != 0 {
		panic("FPrivLevels cannot be fetched with ColumnList/Scan")
	}
}

// Scan reads columns corresponding to the specified fields from the specified
// statement into the receiver.
func (p *Person) Scan(stmt *phys.Stmt, fields Fields) {
	if fields&FID != 0 {
		p.id = ID(stmt.ColumnInt())
	}
	if fields&FVolgisticsID != 0 {
		p.volgisticsID = uint(stmt.ColumnInt())
	}
	if fields&FInformalName != 0 {
		p.informalName = stmt.ColumnText()
	}
	if fields&FFormalName != 0 {
		p.formalName = stmt.ColumnText()
	}
	if fields&FSortName != 0 {
		p.sortName = stmt.ColumnText()
	}
	if fields&FCallSign != 0 {
		p.callSign = stmt.ColumnText()
	}
	if fields&FPronouns != 0 {
		p.pronouns = stmt.ColumnText()
	}
	if fields&FEmail != 0 {
		p.email = stmt.ColumnText()
	}
	if fields&FEmail2 != 0 {
		p.email2 = stmt.ColumnText()
	}
	if fields&FCellPhone != 0 {
		p.cellPhone = stmt.ColumnText()
	}
	if fields&FHomePhone != 0 {
		p.homePhone = stmt.ColumnText()
	}
	if fields&FWorkPhone != 0 {
		p.workPhone = stmt.ColumnText()
	}
	if fields&FPassword != 0 {
		p.password = stmt.ColumnText()
	}
	if fields&FBadLoginCount != 0 {
		p.badLoginCount = uint(stmt.ColumnInt())
	}
	if fields&FBadLoginTime != 0 {
		p.badLoginTime, _ = time.ParseInLocation(badLoginTimeFormat, stmt.ColumnText(), time.Local)
	}
	if fields&FPWResetToken != 0 {
		p.pwresetToken = stmt.ColumnText()
	}
	if fields&FPWResetTime != 0 {
		p.pwresetTime, _ = time.ParseInLocation(pwresetTimeFormat, stmt.ColumnText(), time.Local)
	}
	if fields&FUnsubscribeToken != 0 {
		p.unsubscribeToken = stmt.ColumnText()
	}
	if fields&FHoursToken != 0 {
		p.hoursToken = stmt.ColumnText()
	}
	if fields&FIdentification != 0 {
		p.identification = IdentType(stmt.ColumnInt())
	}
	if fields&FBirthdate != 0 {
		p.birthdate = stmt.ColumnText()
	}
	if fields&FLanguage != 0 {
		p.language = stmt.ColumnText()
	}
	if fields&FFlags != 0 {
		p.flags = Flags(stmt.ColumnHexInt())
	}
	p.fields |= fields &^ (FAddresses | FBGChecks | FDSWRegistrations | FNotes | FEmContacts | FPrivLevels)
}

func (p *Person) readJoins(store phys.Storer, fields Fields) {
	if fields&FAddresses != 0 {
		p.readAddresses(store)
	}
	if fields&FBGChecks != 0 {
		p.readBGChecks(store)
	}
	if fields&FDSWRegistrations != 0 {
		p.readDSWRegistrations(store)
	}
	if fields&FNotes != 0 {
		p.readNotes(store)
	}
	if fields&FEmContacts != 0 {
		p.readEmContacts(store)
	}
	if fields&FPrivLevels != 0 {
		p.readPrivLevels(store)
	}
}

const readAddressesSQL = `SELECT type, same_as_home, address, latitude, longitude, fire_district FROM person_address WHERE person=?`

func (p *Person) readAddresses(store phys.Storer) {
	p.addresses = Addresses{}
	phys.SQL(store, readAddressesSQL, func(stmt *phys.Stmt) {
		stmt.BindInt(int(p.ID()))
		for stmt.Step() {
			var atype = addressType(stmt.ColumnInt())
			var a Address
			a.SameAsHome = stmt.ColumnBool()
			a.Address = stmt.ColumnText()
			a.Latitude = stmt.ColumnFloat()
			a.Longitude = stmt.ColumnFloat()
			a.FireDistrict = uint(stmt.ColumnInt())
			switch atype {
			case addressHome:
				p.addresses.Home = &a
			case addressWork:
				p.addresses.Work = &a
			case addressMail:
				p.addresses.Mail = &a
			default:
				panic("unknown address type in database")
			}
		}
	})
	p.fields |= FAddresses
}

const readBGChecksSQL = `SELECT type, cleared, nli, assumed FROM person_bgcheck WHERE person=?`

func (p *Person) readBGChecks(store phys.Storer) {
	p.bgChecks = BGChecks{}
	phys.SQL(store, readBGChecksSQL, func(stmt *phys.Stmt) {
		stmt.BindInt(int(p.ID()))
		for stmt.Step() {
			var bgtype = bgCheckType(stmt.ColumnInt())
			var bg BGCheck
			bg.Cleared, _ = time.ParseInLocation(bgCheckDateFormat, stmt.ColumnText(), time.Local)
			bg.NLI, _ = time.ParseInLocation(bgCheckDateFormat, stmt.ColumnText(), time.Local)
			bg.Assumed = stmt.ColumnBool()
			switch bgtype {
			case bgCheckDOJ:
				p.bgChecks.DOJ = &bg
			case bgCheckFBI:
				p.bgChecks.FBI = &bg
			case bgCheckPHS:
				p.bgChecks.PHS = &bg
			default:
				panic("unknown bgcheck type in database")
			}
		}
	})
	p.fields |= FBGChecks
}

const readDSWRegistrationsSQL = `SELECT class, registered, expiration FROM person_dswreg WHERE person=?`

func (p *Person) readDSWRegistrations(store phys.Storer) {
	p.dswRegistrations = DSWRegistrations{}
	phys.SQL(store, readDSWRegistrationsSQL, func(stmt *phys.Stmt) {
		stmt.BindInt(int(p.ID()))
		for stmt.Step() {
			var class = dswClass(stmt.ColumnInt())
			var reg DSWRegistration
			reg.Registered, _ = time.ParseInLocation(dswRegDateFormat, stmt.ColumnText(), time.Local)
			reg.Expiration, _ = time.ParseInLocation(dswRegDateFormat, stmt.ColumnText(), time.Local)
			switch class {
			case dswCERT:
				p.dswRegistrations.CERT = &reg
			case dswCommunications:
				p.dswRegistrations.Communications = &reg
			default:
				panic("unknown dsw class in database")
			}
		}
	})
	p.fields |= FDSWRegistrations
}

const readNotesSQL = `SELECT note, date, visibility FROM person_note WHERE person=? ORDER BY date, rowid`

func (p *Person) readNotes(store phys.Storer) {
	p.notes = p.notes[:0]
	phys.SQL(store, readNotesSQL, func(stmt *phys.Stmt) {
		stmt.BindInt(int(p.ID()))
		for stmt.Step() {
			var note Note
			note.Note = stmt.ColumnText()
			note.Date, _ = time.ParseInLocation(noteDateFormat, stmt.ColumnText(), time.Local)
			note.Visibility = NoteVisibility(stmt.ColumnInt())
			p.notes = append(p.notes, &note)
		}
	})
	p.fields |= FNotes
}

const readEmContactsSQL = `SELECT name, home_phone, cell_phone, relationship FROM person_emcontact WHERE person=? ORDER BY rowid`

func (p *Person) readEmContacts(store phys.Storer) {
	p.emContacts = p.emContacts[:0]
	phys.SQL(store, readEmContactsSQL, func(stmt *phys.Stmt) {
		stmt.BindInt(int(p.ID()))
		for stmt.Step() {
			var emContact EmContact
			emContact.Name = stmt.ColumnText()
			emContact.HomePhone = stmt.ColumnText()
			emContact.CellPhone = stmt.ColumnText()
			emContact.Relationship = stmt.ColumnText()
			p.emContacts = append(p.emContacts, &emContact)
		}
	})
	p.fields |= FEmContacts
}

const readPrivLevelsSQL = `SELECT org, privlevel FROM person_privlevel WHERE person=?`

func (p *Person) readPrivLevels(store phys.Storer) {
	p.privLevels = make([]enum.PrivLevel, enum.NumOrgs)
	phys.SQL(store, readPrivLevelsSQL, func(stmt *phys.Stmt) {
		stmt.BindInt(int(p.ID()))
		for stmt.Step() {
			var org = enum.Org(stmt.ColumnInt())
			var privLevel = enum.PrivLevel(stmt.ColumnInt())
			p.privLevels[org] = privLevel
		}
	})
	p.fields |= FPrivLevels
}
