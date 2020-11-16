package report

import (
	"time"

	"github.com/mailru/easyjson/jwriter"

	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/util"
)

// GetIndex handles GET /api/reports requests.
func GetIndex(r *util.Request) error {
	var out jwriter.Writer
	out.RawByte('{')
	if ao := allowedOrgs(r); len(ao) != 0 {
		now := time.Now()
		out.RawString(`"attendance":{"datePresets":[{"label":"this month","dateFrom":`)
		out.String(time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local).Format("2006-01-02"))
		out.RawString(`,"dateTo":`)
		out.String(time.Date(now.Year(), now.Month()+1, 0, 0, 0, 0, 0, time.Local).Format("2006-01-02"))
		out.RawString(`},{"label":"last month","dateFrom":`)
		out.String(time.Date(now.Year(), now.Month()-1, 1, 0, 0, 0, 0, time.Local).Format("2006-01-02"))
		out.RawString(`,"dateTo":`)
		out.String(time.Date(now.Year(), now.Month(), 0, 0, 0, 0, 0, time.Local).Format("2006-01-02"))
		out.RawString(`},{"label":"this quarter","dateFrom":`)
		out.String(time.Date(now.Year(), (now.Month()-1)/3*3+1, 1, 0, 0, 0, 0, time.Local).Format("2006-01-02"))
		out.RawString(`,"dateTo":`)
		out.String(time.Date(now.Year(), (now.Month()-1)/3*3+4, 0, 0, 0, 0, 0, time.Local).Format("2006-01-02"))
		out.RawString(`},{"label":"last quarter","dateFrom":`)
		out.String(time.Date(now.Year(), (now.Month()-1)/3*3-2, 1, 0, 0, 0, 0, time.Local).Format("2006-01-02"))
		out.RawString(`,"dateTo":`)
		out.String(time.Date(now.Year(), (now.Month()-1)/3*3+1, 0, 0, 0, 0, 0, time.Local).Format("2006-01-02"))
		out.RawString(`},{"label":"this year","dateFrom":`)
		out.String(time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.Local).Format("2006-01-02"))
		out.RawString(`,"dateTo":`)
		out.String(time.Date(now.Year(), 12, 31, 0, 0, 0, 0, time.Local).Format("2006-01-02"))
		out.RawString(`},{"label":"last year","dateFrom":`)
		out.String(time.Date(now.Year()-1, 1, 1, 0, 0, 0, 0, time.Local).Format("2006-01-02"))
		out.RawString(`,"dateTo":`)
		out.String(time.Date(now.Year()-1, 12, 31, 0, 0, 0, 0, time.Local).Format("2006-01-02"))
		out.RawString(`}],"eventTypes":[`)
		for i, et := range model.AllEventTypes {
			if i != 0 {
				out.RawByte(',')
			}
			out.RawString(`{"id":`)
			out.Int(int(et))
			out.RawString(`,"label":`)
			out.String(model.EventTypeNames[et])
			out.RawByte('}')
		}
		out.RawString(`],"organizations":[`)
		var first = true
		for _, o := range model.AllOrganizations {
			if !ao[o] {
				continue
			}
			if first {
				first = false
			} else {
				out.RawByte(',')
			}
			out.RawString(`{"id":`)
			out.Int(int(o))
			out.RawString(`,"label":`)
			out.String(model.OrganizationNames[o])
			out.RawByte('}')
		}
		out.RawString(`],"attendanceTypes":[`)
		for i, at := range model.AllAttendanceTypes {
			if i != 0 {
				out.RawByte(',')
			}
			out.RawString(`{"id":`)
			out.Int(int(at))
			out.RawString(`,"label":`)
			out.String(model.AttendanceTypeNames[at])
			out.RawByte('}')
		}
		out.RawString(`]}`)
	}
	out.RawByte('}')
	r.Header().Set("Content-Type", "application/json; charset=utf-8")
	out.DumpTo(r)
	return nil
}
