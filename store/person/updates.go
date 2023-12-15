package person

import (
	"fmt"
	"strings"
	"time"

	"sunnyvaleserv.org/portal/store/internal/phys"
)

// tableFields is the bitmask of fields that are stored in the main person
// table.
const tableFields = FID | FVolgisticsID | FInformalName | FFormalName | FSortName | FCallSign | FPronouns | FEmail | FEmail2 | FCellPhone | FHomePhone | FWorkPhone | FPassword | FBadLoginCount | FBadLoginTime | FPWResetToken | FPWResetTime | FUnsubscribeToken | FHoursToken | FIdentification | FBirthdate | FLanguage | FFlags

// Updater is a structure that can be filled with data for a new or changed
// person, and then later applied.  For creating new people, it can simply be
// instantiated with new().  For updating existing people, either *every* field
// in it must be set, or it should be instantiated with the Updater method of
// the person being changed.
type Updater struct {
	ID               ID
	VolgisticsID     uint
	InformalName     string
	FormalName       string
	SortName         string
	CallSign         string
	Pronouns         string
	Email            string
	Email2           string
	CellPhone        string
	HomePhone        string
	WorkPhone        string
	Password         string
	BadLoginCount    uint
	BadLoginTime     time.Time
	PWResetToken     string
	PWResetTime      time.Time
	UnsubscribeToken string
	HoursToken       string
	Identification   IdentType
	Birthdate        string
	Language         string
	Flags            Flags
	Addresses        Addresses
	BGChecks         BGChecks
	DSWRegistrations DSWRegistrations
	Notes            Notes
	EmContacts       EmContacts
}

// Updater returns a new Updater for the specified person, with its data
// matching the current data for the person.  The person must have fetched FID.
func (p *Person) Updater() *Updater {
	return &Updater{
		ID:               p.ID(),
		VolgisticsID:     p.volgisticsID,
		InformalName:     p.informalName,
		FormalName:       p.formalName,
		SortName:         p.sortName,
		CallSign:         p.callSign,
		Pronouns:         p.pronouns,
		Email:            p.email,
		Email2:           p.email2,
		CellPhone:        p.cellPhone,
		HomePhone:        p.homePhone,
		WorkPhone:        p.workPhone,
		Password:         p.password,
		BadLoginCount:    p.badLoginCount,
		BadLoginTime:     p.badLoginTime,
		PWResetToken:     p.pwresetToken,
		PWResetTime:      p.pwresetTime,
		UnsubscribeToken: p.unsubscribeToken,
		HoursToken:       p.hoursToken,
		Identification:   p.identification,
		Birthdate:        p.birthdate,
		Language:         p.language,
		Flags:            p.flags,
		Addresses:        p.addresses.clone(),
		BGChecks:         p.bgChecks.clone(),
		DSWRegistrations: p.dswRegistrations.clone(),
		Notes:            p.notes.clone(),
		EmContacts:       p.emContacts.clone(),
	}
}

const createSQL = `INSERT INTO person (id, volgistics_id, informal_name, formal_name, sort_name, call_sign, pronouns, email, email2, cell_phone, home_phone, work_phone, password, bad_login_count, bad_login_time, pwreset_token, pwreset_time, unsubscribe_token, hours_token, identification, birthdate, language, flags) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`

// Create creates a new person, with the data in the Updater.
func Create(storer phys.Storer, u *Updater) (p *Person) {
	p = new(Person)
	p.fields = tableFields | FAddresses | FBGChecks | FDSWRegistrations | FNotes
	phys.SQL(storer, createSQL, func(stmt *phys.Stmt) {
		stmt.BindNullInt(int(u.ID))
		stmt.BindNullInt(int(u.VolgisticsID))
		stmt.BindText(u.InformalName)
		stmt.BindText(u.FormalName)
		stmt.BindText(u.SortName)
		stmt.BindNullText(u.CallSign)
		stmt.BindNullText(u.Pronouns)
		stmt.BindNullText(u.Email)
		stmt.BindNullText(u.Email2)
		stmt.BindNullText(u.CellPhone)
		stmt.BindNullText(u.HomePhone)
		stmt.BindNullText(u.WorkPhone)
		stmt.BindNullText(u.Password)
		stmt.BindNullInt(int(u.BadLoginCount))
		if u.BadLoginTime.IsZero() {
			stmt.BindNull()
		} else {
			stmt.BindText(u.BadLoginTime.In(time.Local).Format(badLoginTimeFormat))
		}
		stmt.BindNullText(u.PWResetToken)
		if u.PWResetTime.IsZero() {
			stmt.BindNull()
		} else {
			stmt.BindText(u.PWResetTime.In(time.Local).Format(pwresetTimeFormat))
		}
		stmt.BindNullText(u.UnsubscribeToken)
		stmt.BindNullText(u.HoursToken)
		stmt.BindInt(int(u.Identification))
		stmt.BindNullText(u.Birthdate)
		stmt.BindText(u.Language)
		stmt.BindHexInt(int(u.Flags))
		stmt.Step()
		if u.ID != 0 {
			p.id = u.ID
		} else {
			p.id = ID(phys.LastInsertRowID(storer))
		}
	})
	storeAddresses(storer, p.id, u.Addresses, false)
	storeBGChecks(storer, p.id, u.BGChecks, false)
	storeDSWRegistrations(storer, p.id, u.DSWRegistrations, false)
	storeNotes(storer, p.id, u.Notes, false)
	storeEmContacts(storer, p.id, u.EmContacts, false)
	p.auditAndUpdate(storer, u, p.fields, true)
	phys.Index(storer, p)
	return p
}

var updateSQLCache map[Fields]string

// Update updates the specified fields of the existing person, with the data in
// the Updater.
func (p *Person) Update(storer phys.Storer, u *Updater, fields Fields) {
	if fields&FID != 0 {
		panic("cannot change person ID")
	}
	if fields&FPrivLevels != 0 {
		panic("cannot change privilege levels")
	}
	if tf := fields & tableFields; tf != 0 {
		if updateSQLCache == nil {
			updateSQLCache = make(map[Fields]string)
		}
		if _, ok := updateSQLCache[tf]; !ok {
			var sb strings.Builder
			sb.WriteString("UPDATE person SET ")
			var sep = phys.NewSeparator(", ")
			if tf&FVolgisticsID != 0 {
				sb.WriteString(sep())
				sb.WriteString("volgistics_id=?")
			}
			if tf&FInformalName != 0 {
				sb.WriteString(sep())
				sb.WriteString("informal_name=?")
			}
			if tf&FFormalName != 0 {
				sb.WriteString(sep())
				sb.WriteString("formal_name=?")
			}
			if tf&FSortName != 0 {
				sb.WriteString(sep())
				sb.WriteString("sort_name=?")
			}
			if tf&FCallSign != 0 {
				sb.WriteString(sep())
				sb.WriteString("call_sign=?")
			}
			if tf&FPronouns != 0 {
				sb.WriteString(sep())
				sb.WriteString("pronouns=?")
			}
			if tf&FEmail != 0 {
				sb.WriteString(sep())
				sb.WriteString("email=?")
			}
			if tf&FEmail2 != 0 {
				sb.WriteString(sep())
				sb.WriteString("email2=?")
			}
			if tf&FCellPhone != 0 {
				sb.WriteString(sep())
				sb.WriteString("cell_phone=?")
			}
			if tf&FHomePhone != 0 {
				sb.WriteString(sep())
				sb.WriteString("home_phone=?")
			}
			if tf&FWorkPhone != 0 {
				sb.WriteString(sep())
				sb.WriteString("work_phone=?")
			}
			if tf&FPassword != 0 {
				sb.WriteString(sep())
				sb.WriteString("password=?")
			}
			if tf&FBadLoginCount != 0 {
				sb.WriteString(sep())
				sb.WriteString("bad_login_count=?")
			}
			if tf&FBadLoginTime != 0 {
				sb.WriteString(sep())
				sb.WriteString("bad_login_time=?")
			}
			if tf&FPWResetToken != 0 {
				sb.WriteString(sep())
				sb.WriteString("pwreset_token=?")
			}
			if tf&FPWResetTime != 0 {
				sb.WriteString(sep())
				sb.WriteString("pwreset_time=?")
			}
			if tf&FUnsubscribeToken != 0 {
				sb.WriteString(sep())
				sb.WriteString("unsubscribe_token=?")
			}
			if tf&FHoursToken != 0 {
				sb.WriteString(sep())
				sb.WriteString("hours_token=?")
			}
			if tf&FIdentification != 0 {
				sb.WriteString(sep())
				sb.WriteString("identification=?")
			}
			if tf&FBirthdate != 0 {
				sb.WriteString(sep())
				sb.WriteString("birthdate=?")
			}
			if tf&FLanguage != 0 {
				sb.WriteString(sep())
				sb.WriteString("language=?")
			}
			if tf&FFlags != 0 {
				sb.WriteString(sep())
				sb.WriteString("flags=?")
			}
			sb.WriteString(" WHERE id=?")
			updateSQLCache[tf] = sb.String()
		}
		phys.SQL(storer, updateSQLCache[tf], func(stmt *phys.Stmt) {
			if tf&FVolgisticsID != 0 {
				stmt.BindNullInt(int(u.VolgisticsID))
			}
			if tf&FInformalName != 0 {
				stmt.BindText(u.InformalName)
			}
			if tf&FFormalName != 0 {
				stmt.BindText(u.FormalName)
			}
			if tf&FSortName != 0 {
				stmt.BindText(u.SortName)
			}
			if tf&FCallSign != 0 {
				stmt.BindNullText(u.CallSign)
			}
			if tf&FPronouns != 0 {
				stmt.BindNullText(u.Pronouns)
			}
			if tf&FEmail != 0 {
				stmt.BindNullText(u.Email)
			}
			if tf&FEmail2 != 0 {
				stmt.BindNullText(u.Email2)
			}
			if tf&FCellPhone != 0 {
				stmt.BindNullText(u.CellPhone)
			}
			if tf&FHomePhone != 0 {
				stmt.BindNullText(u.HomePhone)
			}
			if tf&FWorkPhone != 0 {
				stmt.BindNullText(u.WorkPhone)
			}
			if tf&FPassword != 0 {
				stmt.BindNullText(u.Password)
			}
			if tf&FBadLoginCount != 0 {
				stmt.BindNullInt(int(u.BadLoginCount))
			}
			if tf&FBadLoginTime != 0 {
				if u.BadLoginTime.IsZero() {
					stmt.BindNull()
				} else {
					stmt.BindText(u.BadLoginTime.In(time.Local).Format(badLoginTimeFormat))
				}
			}
			if tf&FPWResetToken != 0 {
				stmt.BindNullText(u.PWResetToken)
			}
			if tf&FPWResetTime != 0 {
				if u.PWResetTime.IsZero() {
					stmt.BindNull()
				} else {
					stmt.BindText(u.PWResetTime.In(time.Local).Format(pwresetTimeFormat))
				}
			}
			if tf&FUnsubscribeToken != 0 {
				stmt.BindNullText(u.UnsubscribeToken)
			}
			if tf&FHoursToken != 0 {
				stmt.BindNullText(u.HoursToken)
			}
			if tf&FIdentification != 0 {
				stmt.BindInt(int(u.Identification))
			}
			if tf&FBirthdate != 0 {
				stmt.BindNullText(u.Birthdate)
			}
			if tf&FLanguage != 0 {
				stmt.BindText(u.Language)
			}
			if tf&FFlags != 0 {
				stmt.BindHexInt(int(u.Flags))
			}
			stmt.BindInt(int(p.id))
			stmt.Step()
		})
	}
	if fields&FAddresses != 0 {
		storeAddresses(storer, p.id, u.Addresses, true)
	}
	if fields&FBGChecks != 0 {
		storeBGChecks(storer, p.id, u.BGChecks, true)
	}
	if fields&FDSWRegistrations != 0 {
		storeDSWRegistrations(storer, p.id, u.DSWRegistrations, true)
	}
	if fields&FNotes != 0 {
		storeNotes(storer, p.id, u.Notes, true)
	}
	if fields&FEmContacts != 0 {
		storeEmContacts(storer, p.id, u.EmContacts, true)
	}
	p.auditAndUpdate(storer, u, fields, false)
	phys.Index(storer, p)
}

const deleteAddressesSQL = `DELETE FROM person_address WHERE person=?`
const createAddressSQL = `INSERT INTO person_address (person, type, same_as_home, address, latitude, longitude, fire_district) VALUES (?,?,?,?,?,?,?)`

func storeAddresses(storer phys.Storer, pid ID, addresses Addresses, delfirst bool) {
	if delfirst {
		phys.SQL(storer, deleteAddressesSQL, func(stmt *phys.Stmt) {
			stmt.BindInt(int(pid))
			stmt.Step()
		})
	}
	if addresses.Home != nil {
		storeAddress(storer, pid, addressHome, addresses.Home)
	}
	if addresses.Work != nil {
		storeAddress(storer, pid, addressWork, addresses.Work)
	}
	if addresses.Mail != nil {
		storeAddress(storer, pid, addressMail, addresses.Mail)
	}
}
func storeAddress(storer phys.Storer, pid ID, atype addressType, address *Address) {
	phys.SQL(storer, createAddressSQL, func(stmt *phys.Stmt) {
		stmt.BindInt(int(pid))
		stmt.BindInt(int(atype))
		stmt.BindBool(address.SameAsHome)
		stmt.BindNullText(address.Address)
		stmt.BindNullFloat(address.Latitude)
		stmt.BindNullFloat(address.Longitude)
		stmt.BindNullInt(int(address.FireDistrict))
		stmt.Step()
	})
}

const deleteBGChecksSQL = `DELETE FROM person_bgcheck WHERE person=?`
const createBGCheckSQL = `INSERT INTO person_bgcheck (person, type, cleared, nli, assumed) VALUES (?,?,?,?,?)`

func storeBGChecks(storer phys.Storer, pid ID, bgChecks BGChecks, delfirst bool) {
	if delfirst {
		phys.SQL(storer, deleteBGChecksSQL, func(stmt *phys.Stmt) {
			stmt.BindInt(int(pid))
			stmt.Step()
		})
	}
	if bgChecks.DOJ != nil {
		storeBGCheck(storer, pid, bgCheckDOJ, bgChecks.DOJ)
	}
	if bgChecks.FBI != nil {
		storeBGCheck(storer, pid, bgCheckFBI, bgChecks.FBI)
	}
	if bgChecks.PHS != nil {
		storeBGCheck(storer, pid, bgCheckPHS, bgChecks.PHS)
	}
}
func storeBGCheck(storer phys.Storer, pid ID, bgtype bgCheckType, bgCheck *BGCheck) {
	phys.SQL(storer, createBGCheckSQL, func(stmt *phys.Stmt) {
		stmt.BindInt(int(pid))
		stmt.BindInt(int(bgtype))
		if bgCheck.Cleared.IsZero() {
			stmt.BindNull()
		} else {
			stmt.BindText(bgCheck.Cleared.In(time.Local).Format(bgCheckDateFormat))
		}
		if bgCheck.NLI.IsZero() {
			stmt.BindNull()
		} else {
			stmt.BindText(bgCheck.NLI.In(time.Local).Format(bgCheckDateFormat))
		}
		stmt.BindBool(bgCheck.Assumed)
		stmt.Step()
	})
}

const deleteDSWRegistrationsSQL = `DELETE FROM person_dswreg WHERE person=?`
const createDSWRegistrationSQL = `INSERT INTO person_dswreg (person, class, registered, expiration) VALUES (?,?,?,?)`

func storeDSWRegistrations(storer phys.Storer, pid ID, dswregs DSWRegistrations, delfirst bool) {
	if delfirst {
		phys.SQL(storer, deleteDSWRegistrationsSQL, func(stmt *phys.Stmt) {
			stmt.BindInt(int(pid))
			stmt.Step()
		})
	}
	if dswregs.CERT != nil {
		storeDSWRegistration(storer, pid, dswCERT, dswregs.CERT)
	}
	if dswregs.Communications != nil {
		storeDSWRegistration(storer, pid, dswCommunications, dswregs.Communications)
	}
}
func storeDSWRegistration(storer phys.Storer, pid ID, class dswClass, dswreg *DSWRegistration) {
	phys.SQL(storer, createDSWRegistrationSQL, func(stmt *phys.Stmt) {
		stmt.BindInt(int(pid))
		stmt.BindInt(int(class))
		stmt.BindText(dswreg.Registered.In(time.Local).Format(dswRegDateFormat))
		if dswreg.Expiration.IsZero() {
			stmt.BindNull()
		} else {
			stmt.BindText(dswreg.Expiration.In(time.Local).Format(dswRegDateFormat))
		}
		stmt.Step()
	})
}

const deleteNotesSQL = `DELETE FROM person_note WHERE person=?`
const createNoteSQL = `INSERT INTO person_note (person, note, date, visibility) VALUES (?,?,?,?)`

func storeNotes(storer phys.Storer, pid ID, notes Notes, delfirst bool) {
	if delfirst {
		phys.SQL(storer, deleteNotesSQL, func(stmt *phys.Stmt) {
			stmt.BindInt(int(pid))
			stmt.Step()
		})
	}
	if len(notes) != 0 {
		phys.SQL(storer, createNoteSQL, func(stmt *phys.Stmt) {
			for _, note := range notes {
				stmt.BindInt(int(pid))
				stmt.BindText(note.Note)
				stmt.BindText(note.Date.In(time.Local).Format(noteDateFormat))
				stmt.BindInt(int(note.Visibility))
				stmt.Step()
				stmt.Reset()
			}
		})
	}
}

const deleteEmContactsSQL = `DELETE FROM person_emcontact WHERE person=?`
const createEmContactSQL = `INSERT INTO person_emcontact (person, name, home_phone, cell_phone, relationship) VALUES (?,?,?,?,?)`

func storeEmContacts(storer phys.Storer, pid ID, emContacts EmContacts, delfirst bool) {
	if delfirst {
		phys.SQL(storer, deleteEmContactsSQL, func(stmt *phys.Stmt) {
			stmt.BindInt(int(pid))
			stmt.Step()
		})
	}
	if len(emContacts) != 0 {
		phys.SQL(storer, createEmContactSQL, func(stmt *phys.Stmt) {
			for _, emContact := range emContacts {
				stmt.BindInt(int(pid))
				stmt.BindText(emContact.Name)
				stmt.BindNullText(emContact.HomePhone)
				stmt.BindNullText(emContact.CellPhone)
				stmt.BindText(emContact.Relationship)
				stmt.Step()
				stmt.Reset()
			}
		})
	}
}

func (p *Person) auditAndUpdate(storer phys.Storer, u *Updater, fields Fields, create bool) {
	context := fmt.Sprintf("Person %q [%d]", u.InformalName, p.id)
	if create {
		context = "ADD " + context
	}
	if fields&FVolgisticsID != 0 && u.VolgisticsID != p.volgisticsID {
		phys.Audit(storer, "%s:: volgisticsID = %d", context, u.VolgisticsID)
		p.volgisticsID = u.VolgisticsID
	}
	if fields&FInformalName != 0 && u.InformalName != p.informalName {
		phys.Audit(storer, "%s:: informalName = %q", context, u.InformalName)
		p.informalName = u.InformalName
	}
	if fields&FFormalName != 0 && u.FormalName != p.formalName {
		phys.Audit(storer, "%s:: formalName = %q", context, u.FormalName)
		p.formalName = u.FormalName
	}
	if fields&FSortName != 0 && u.SortName != p.sortName {
		phys.Audit(storer, "%s:: sortName = %q", context, u.SortName)
		p.sortName = u.SortName
	}
	if fields&FCallSign != 0 && u.CallSign != p.callSign {
		phys.Audit(storer, "%s:: callSign = %q", context, u.CallSign)
		p.callSign = u.CallSign
	}
	if fields&FPronouns != 0 && u.Pronouns != p.pronouns {
		phys.Audit(storer, "%s:: pronouns = %q", context, u.Pronouns)
		p.pronouns = u.Pronouns
	}
	if fields&FEmail != 0 && u.Email != p.email {
		phys.Audit(storer, "%s:: email = %q", context, u.Email)
		p.email = u.Email
	}
	if fields&FEmail2 != 0 && u.Email2 != p.email2 {
		phys.Audit(storer, "%s:: email2 = %q", context, u.Email2)
		p.email2 = u.Email2
	}
	if fields&FCellPhone != 0 && u.CellPhone != p.cellPhone {
		phys.Audit(storer, "%s:: cellPhone = %q", context, u.CellPhone)
		p.cellPhone = u.CellPhone
	}
	if fields&FHomePhone != 0 && u.HomePhone != p.homePhone {
		phys.Audit(storer, "%s:: homePhone = %q", context, u.HomePhone)
		p.homePhone = u.HomePhone
	}
	if fields&FWorkPhone != 0 && u.WorkPhone != p.workPhone {
		phys.Audit(storer, "%s:: workPhone = %q", context, u.WorkPhone)
		p.workPhone = u.WorkPhone
	}
	if fields&FPassword != 0 && u.Password != p.password {
		phys.Audit(storer, "%s:: password = %q", context, u.Password)
		p.password = u.Password
	}
	if fields&FBadLoginCount != 0 && u.BadLoginCount != p.badLoginCount {
		phys.Audit(storer, "%s:: badLoginCount = %d", context, u.BadLoginCount)
		p.badLoginCount = u.BadLoginCount
	}
	if fields&FBadLoginTime != 0 && u.BadLoginTime != p.badLoginTime {
		if u.BadLoginTime.IsZero() {
			phys.Audit(storer, "%s:: badLoginTime = (zero)", context)
			p.badLoginTime = u.BadLoginTime
		} else {
			phys.Audit(storer, "%s:: badLoginTime = %s", context, u.BadLoginTime.In(time.Local).Format(badLoginTimeFormat))
			p.badLoginTime = u.BadLoginTime.In(time.Local)
		}
	}
	if fields&FPWResetToken != 0 && u.PWResetToken != p.pwresetToken {
		phys.Audit(storer, "%s:: pwresetToken = %q", context, u.PWResetToken)
		p.pwresetToken = u.PWResetToken
	}
	if fields&FPWResetTime != 0 && u.PWResetTime != p.pwresetTime {
		if u.PWResetTime.IsZero() {
			phys.Audit(storer, "%s:: pwresetTime = (zero)", context)
			p.pwresetTime = u.PWResetTime
		} else {
			phys.Audit(storer, "%s:: pwresetTime = %s", context, u.PWResetTime.In(time.Local).Format(pwresetTimeFormat))
			p.pwresetTime = u.PWResetTime.In(time.Local)
		}
	}
	if fields&FUnsubscribeToken != 0 && u.UnsubscribeToken != p.unsubscribeToken {
		phys.Audit(storer, "%s:: unsubscribeToken = %q", context, u.UnsubscribeToken)
		p.unsubscribeToken = u.UnsubscribeToken
	}
	if fields&FHoursToken != 0 && u.HoursToken != p.hoursToken {
		phys.Audit(storer, "%s:: hoursToken = %q", context, u.HoursToken)
		p.hoursToken = u.HoursToken
	}
	if fields&FIdentification != 0 && u.Identification != p.identification {
		phys.Audit(storer, "%s:: identification = 0x%x", context, u.Identification)
		p.identification = u.Identification
	}
	if fields&FBirthdate != 0 && u.Birthdate != p.birthdate {
		phys.Audit(storer, "%s:: birthdate = %q", context, u.Birthdate)
		p.birthdate = u.Birthdate
	}
	if fields&FLanguage != 0 && u.Language != p.language {
		phys.Audit(storer, "%s:: language = %s", context, u.Language)
		p.language = u.Language
	}
	if fields&FFlags != 0 && u.Flags != p.flags {
		phys.Audit(storer, "%s:: flags = 0x%x", context, u.Flags)
		p.flags = u.Flags
	}
	if fields&FAddresses != 0 {
		auditAndUpdateAddress(storer, &p.addresses.Home, &u.Addresses.Home, context, "addresses.home")
		auditAndUpdateAddress(storer, &p.addresses.Work, &u.Addresses.Work, context, "addresses.work")
		auditAndUpdateAddress(storer, &p.addresses.Mail, &u.Addresses.Mail, context, "addresses.mail")
	}
	if fields&FBGChecks != 0 {
		auditAndUpdateBGCheck(storer, &p.bgChecks.DOJ, &u.BGChecks.DOJ, context, "bgChecks.DOJ")
		auditAndUpdateBGCheck(storer, &p.bgChecks.FBI, &u.BGChecks.FBI, context, "bgChecks.FBI")
		auditAndUpdateBGCheck(storer, &p.bgChecks.PHS, &u.BGChecks.PHS, context, "bgChecks.PHS")
	}
	if fields&FDSWRegistrations != 0 {
		auditAndUpdateDSWRegistration(storer, &p.dswRegistrations.CERT, &u.DSWRegistrations.CERT, context, "dswRegistrations.CERT")
		auditAndUpdateDSWRegistration(storer, &p.dswRegistrations.Communications, &u.DSWRegistrations.Communications, context, "dswRegistrations.Communications")
	}
	if fields&FNotes != 0 {
		var i int
		for i = 0; i < len(p.notes) && i < len(u.Notes); i++ {
			auditAndUpdateNote(storer, p.notes[i], u.Notes[i], context, fmt.Sprintf("notes[%d]", i))
		}
		for ; i < len(p.notes); i++ {
			auditAndUpdateNote(storer, p.notes[i], new(Note), context, fmt.Sprintf("DELETE notes[%d]", i))
		}
		for ; i < len(u.Notes); i++ {
			p.notes = append(p.notes, new(Note))
			auditAndUpdateNote(storer, p.notes[i], u.Notes[i], context, fmt.Sprintf("ADD notes[%d]", i))
		}
		p.notes = p.notes[:len(u.Notes)]
	}
	if fields&FEmContacts != 0 {
		var i int
		for i = 0; i < len(p.emContacts) && i < len(u.EmContacts); i++ {
			auditAndUpdateEmContact(storer, p.emContacts[i], u.EmContacts[i], context, fmt.Sprintf("emContacts[%d]", i))
		}
		for ; i < len(p.emContacts); i++ {
			auditAndUpdateEmContact(storer, p.emContacts[i], new(EmContact), context, fmt.Sprintf("DELETE emContacts[%d]", i))
		}
		for ; i < len(u.EmContacts); i++ {
			p.emContacts = append(p.emContacts, new(EmContact))
			auditAndUpdateEmContact(storer, p.emContacts[i], u.EmContacts[i], context, fmt.Sprintf("ADD emContacts[%d]", i))
		}
		p.emContacts = p.emContacts[:len(u.EmContacts)]
	}
}
func auditAndUpdateAddress(storer phys.Storer, paddr, uaddr **Address, context, addrtype string) {
	if *paddr == nil && *uaddr == nil {
		return
	}
	if *uaddr == nil {
		phys.Audit(storer, "%s:: DELETE %s", context, addrtype)
		*paddr = nil
		return
	}
	if *paddr == nil {
		addrtype = "ADD " + addrtype
		*paddr = new(Address)
	}
	if (*paddr).SameAsHome != (*uaddr).SameAsHome {
		phys.Audit(storer, "%s:: %s:: sameAsHome = %v", context, addrtype, (*uaddr).SameAsHome)
		(*paddr).SameAsHome = (*uaddr).SameAsHome
	}
	if (*paddr).Address != (*uaddr).Address {
		phys.Audit(storer, "%s:: %s:: address = %q", context, addrtype, (*uaddr).Address)
		(*paddr).Address = (*uaddr).Address
	}
	if (*paddr).Latitude != (*uaddr).Latitude {
		phys.Audit(storer, "%s:: %s:: latitude = %f", context, addrtype, (*uaddr).Latitude)
		(*paddr).Latitude = (*uaddr).Latitude
	}
	if (*paddr).Longitude != (*uaddr).Longitude {
		phys.Audit(storer, "%s:: %s:: longitude = %f", context, addrtype, (*uaddr).Longitude)
		(*paddr).Longitude = (*uaddr).Longitude
	}
	if (*paddr).FireDistrict != (*uaddr).FireDistrict {
		phys.Audit(storer, "%s:: %s:: fireDistrict = %d", context, addrtype, (*uaddr).FireDistrict)
		(*paddr).FireDistrict = (*uaddr).FireDistrict
	}
}
func auditAndUpdateBGCheck(storer phys.Storer, paddr, uaddr **BGCheck, context, checktype string) {
	if *paddr == nil && *uaddr == nil {
		return
	}
	if *uaddr == nil {
		phys.Audit(storer, "%s:: DELETE %s", context, checktype)
		*paddr = nil
		return
	}
	if *paddr == nil {
		checktype = "ADD " + checktype
		*paddr = new(BGCheck)
	}
	if (*paddr).Cleared != (*uaddr).Cleared {
		phys.Audit(storer, "%s:: %s:: cleared = %q", context, checktype, (*uaddr).Cleared)
		(*paddr).Cleared = (*uaddr).Cleared
	}
	if (*paddr).NLI != (*uaddr).NLI {
		phys.Audit(storer, "%s:: %s:: nli = %q", context, checktype, (*uaddr).NLI)
		(*paddr).NLI = (*uaddr).NLI
	}
	if (*paddr).Assumed != (*uaddr).Assumed {
		phys.Audit(storer, "%s:: %s:: assumed = %v", context, checktype, (*uaddr).Assumed)
		(*paddr).Assumed = (*uaddr).Assumed
	}
}
func auditAndUpdateDSWRegistration(storer phys.Storer, paddr, uaddr **DSWRegistration, context, class string) {
	if *paddr == nil && *uaddr == nil {
		return
	}
	if *uaddr == nil {
		phys.Audit(storer, "%s:: DELETE %s", context, class)
		*paddr = nil
		return
	}
	if *paddr == nil {
		class = "ADD " + class
		*paddr = new(DSWRegistration)
	}
	if (*paddr).Registered != (*uaddr).Registered {
		phys.Audit(storer, "%s:: %s:: cleared = %q", context, class, (*uaddr).Registered)
		(*paddr).Registered = (*uaddr).Registered
	}
	if (*paddr).Expiration != (*uaddr).Expiration {
		phys.Audit(storer, "%s:: %s:: nli = %q", context, class, (*uaddr).Expiration)
		(*paddr).Expiration = (*uaddr).Expiration
	}
}
func auditAndUpdateNote(storer phys.Storer, pnote, unote *Note, context, label string) {
	if pnote.Note != unote.Note {
		phys.Audit(storer, "%s:: %s:: note = %q", context, label, unote.Note)
		pnote.Note = unote.Note
	}
	if pnote.Date != unote.Date {
		phys.Audit(storer, "%s:: %s:: date = %q", context, label, unote.Date)
		pnote.Date = unote.Date
	}
	if pnote.Visibility != unote.Visibility {
		phys.Audit(storer, "%s:: %s:: visibility = %s [%d]", context, label, unote.Visibility, unote.Visibility)
		pnote.Visibility = unote.Visibility
	}
}
func auditAndUpdateEmContact(storer phys.Storer, pec, uec *EmContact, context, label string) {
	if pec.Name != uec.Name {
		phys.Audit(storer, "%s:: %s:: name = %q", context, label, uec.Name)
		pec.Name = uec.Name
	}
	if pec.HomePhone != uec.HomePhone {
		phys.Audit(storer, "%s:: %s:: homePhone = %q", context, label, uec.HomePhone)
		pec.HomePhone = uec.HomePhone
	}
	if pec.CellPhone != uec.CellPhone {
		phys.Audit(storer, "%s:: %s:: cellPhone = %q", context, label, uec.CellPhone)
		pec.CellPhone = uec.CellPhone
	}
	if pec.Relationship != uec.Relationship {
		phys.Audit(storer, "%s:: %s:: relationship = %q", context, label, uec.Relationship)
		pec.Relationship = uec.Relationship
	}
}

const duplicateCallSignSQL = `SELECT 1 FROM person WHERE id!=? AND call_sign=?`

// DuplicateCallSign returns whether the call sign specified in the Updater
// would be a duplicate if applied.
func (u *Updater) DuplicateCallSign(storer phys.Storer) (found bool) {
	if u.CallSign == "" {
		return false
	}
	phys.SQL(storer, duplicateCallSignSQL, func(stmt *phys.Stmt) {
		stmt.BindInt(int(u.ID))
		stmt.BindText(u.CallSign)
		found = stmt.Step()
	})
	return found
}

const duplicateCellPhoneSQL = `SELECT 1 FROM person WHERE id!=? AND cell_phone=?`

// DuplicateCellPhone returns whether the cell phone number specified in the
// Updater would be a duplicate if applied.
func (u *Updater) DuplicateCellPhone(storer phys.Storer) (found bool) {
	if u.CellPhone == "" {
		return false
	}
	phys.SQL(storer, duplicateCellPhoneSQL, func(stmt *phys.Stmt) {
		stmt.BindInt(int(u.ID))
		stmt.BindText(u.CellPhone)
		found = stmt.Step()
	})
	return found
}

const duplicateEmailSQL = `SELECT 1 FROM person WHERE id!=? AND email=?`

// DuplicateEmail returns whether the primary email address specified in the
// Updater would be a duplicate if applied.
func (u *Updater) DuplicateEmail(storer phys.Storer) (found bool) {
	phys.SQL(storer, duplicateEmailSQL, func(stmt *phys.Stmt) {
		stmt.BindInt(int(u.ID))
		stmt.BindText(u.Email)
		found = stmt.Step()
	})
	return found
}

const duplicateSortNameSQL = `SELECT 1 FROM person WHERE id!=? AND sort_name=?`

// DuplicateSortName returns whether the sort name specified in the Updater
// would be a duplicate if applied.
func (u *Updater) DuplicateSortName(storer phys.Storer) (found bool) {
	phys.SQL(storer, duplicateSortNameSQL, func(stmt *phys.Stmt) {
		stmt.BindInt(int(u.ID))
		stmt.BindText(u.SortName)
		found = stmt.Step()
	})
	return found
}

const duplicateVolgisticsIDSQL = `SELECT 1 FROM person WHERE id!=? AND volgistics_id=?`

// DuplicateVolgisticsID returns whether the Volgistics ID specified in the
// Updater would be a duplicate if applied.
func (u *Updater) DuplicateVolgisticsID(storer phys.Storer) (found bool) {
	if u.VolgisticsID == 0 {
		return false
	}
	phys.SQL(storer, duplicateVolgisticsIDSQL, func(stmt *phys.Stmt) {
		stmt.BindInt(int(u.ID))
		stmt.BindInt(int(u.VolgisticsID))
		found = stmt.Step()
	})
	return found
}

// We intentionally do not have a Delete method for Person objects.  Person
// objects should never be deleted.

// ClearAllHoursReminders turns off the HoursReminder flag on all people.
func ClearAllHoursReminders(storer phys.Storer) {
	phys.SQL(storer, fmt.Sprintf("UPDATE person SET flags=flags-%d WHERE flags&%d", HoursReminder, HoursReminder), func(stmt *phys.Stmt) {
		stmt.Step()
	})
	phys.Audit(storer, "Clear all hours reminder flags")
}
