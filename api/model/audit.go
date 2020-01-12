package model

import (
	"time"

	"github.com/mailru/easyjson/jwriter"
)

// An AuditRecord contains the information to record in the audit log about a
// change.
type AuditRecord struct {
	Timestamp time.Time
	Username  string
	Request   string
	Event     *Event
	Person    *Person
	Role      *Role
	Session   *Session
	Venue     *Venue
}

func (in AuditRecord) MarshalEasyJSON(out *jwriter.Writer) {
	out.RawByte('{')
	{
		out.RawString(`"timestamp":`)
		out.Raw(in.Timestamp.MarshalJSON())
	}
	if in.Username != "" {
		out.RawString(`,"username":`)
		out.String(string(in.Username))
	}
	if in.Request != "" {
		out.RawString(`,"request":`)
		out.String(string(in.Request))
	}
	if in.Event != nil {
		out.RawString(`,"event":`)
		in.Event.ToAudit(out)
	}
	if in.Person != nil {
		out.RawString(`,"person":`)
		in.Person.ToAudit(out)
	}
	if in.Role != nil {
		out.RawString(`,"role":`)
		in.Role.ToAudit(out)
	}
	if in.Session != nil {
		out.RawString(`,"session":`)
		in.Session.ToAudit(out)
	}
	if in.Venue != nil {
		out.RawString(`,"venue":`)
		in.Venue.ToAudit(out)
	}
	out.RawString("}\n")
}

func (in Event) ToAudit(out *jwriter.Writer) {
	out.RawString(`{"id":`)
	out.Int(int(in.ID))
	out.RawString(`,"name":`)
	out.String(in.Name)
	out.RawString(`,"date":`)
	out.String(in.Date)
	out.RawString(`,"start":`)
	out.String(in.Start)
	out.RawString(`,"end":`)
	out.String(in.End)
	out.RawString(`,"venue":`)
	if in.Venue != nil {
		out.Int(int(in.Venue.ID))
	} else {
		out.RawString(`null`)
	}
	out.RawString(`,"type":`)
	out.String(string(in.Type))
	out.RawString(`,"roles":[`)
	for i, t := range in.Roles {
		if i != 0 {
			out.RawByte(',')
		}
		out.Int(int(t.ID))
	}
	out.RawString(`]}`)
}

func (in Person) ToAudit(out *jwriter.Writer) {
	out.RawString(`{"id":`)
	out.Int(int(in.ID))
	out.RawString(`,"firstName":`)
	out.String(in.FirstName)
	out.RawString(`,"lastName":`)
	out.String(in.LastName)
	out.RawString(`,"email":`)
	out.String(in.Email)
	out.RawString(`,"phone":`)
	out.String(in.Phone)
	if in.BadLoginCount != 0 {
		out.RawString(`,"badLoginCount":`)
		out.Int(in.BadLoginCount)
		out.RawString(`,"badLoginTime":`)
		out.Raw(in.BadLoginTime.MarshalJSON())
	}
	if in.PWResetToken != "" {
		out.RawString(`,"pwresetToken":`)
		out.String(in.PWResetToken)
		out.RawString(`,"pwresetTime":`)
		out.Raw(in.PWResetTime.MarshalJSON())
	}
	out.RawString(`,"roles":[`)
	for i, role := range in.Roles {
		if i != 0 {
			out.RawByte(',')
		}
		out.Int(int(role.ID))
	}
	out.RawString(`]}`)
}

func (in PrivilegeMap) ToAudit(out *jwriter.Writer) {
	out.RawByte('{')
	first := true
	for r, p := range in {
		if first {
			first = false
		} else {
			out.RawByte(',')
		}
		out.IntStr(int(r))
		out.RawByte(':')
		out.Uint8(uint8(p))
	}
	out.RawByte('}')
}

func (in Role) ToAudit(out *jwriter.Writer) {
	out.RawString(`{"id":`)
	out.Int(int(in.ID))
	if in.Tag != "" {
		out.RawString(`,"tag":`)
		out.String(string(in.Tag))
	}
	out.RawString(`,"name":`)
	out.String(in.Name)
	out.RawString(`,"memberLabel":`)
	out.String(in.MemberLabel)
	out.RawString(`,"implyOnly":`)
	out.Bool(in.ImplyOnly)
	out.RawString(`,"individual":`)
	out.Bool(in.Individual)
	out.RawString(`,"privileges":`)
	in.PrivMap.ToAudit(out)
	out.RawByte('}')
}

func (in Session) ToAudit(out *jwriter.Writer) {
	out.RawString(`{"token":`)
	out.String(string(in.Token))
	if in.Person != nil {
		out.RawString(`,"person":`)
		out.String(in.Person.Email)
	}
	if !in.Expires.IsZero() {
		out.RawString(`,"expires":`)
		out.Raw(in.Expires.MarshalJSON())
	}
	out.RawByte('}')
}

func (in Venue) ToAudit(out *jwriter.Writer) {
	out.RawString(`{"id":`)
	out.Int(int(in.ID))
	out.RawString(`,"name":`)
	out.String(in.Name)
	out.RawString(`,"address":`)
	out.String(in.Address)
	out.RawString(`,"city":`)
	out.String(in.City)
	if in.URL != "" {
		out.RawString(`,"url":`)
		out.String(in.URL)
	}
	out.RawByte('}')
}
