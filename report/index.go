package report

import (
	"time"

	"github.com/mailru/easyjson/jwriter"

	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/util"
)

// GetIndex handles GET /api/reports requests.
func GetIndex(r *util.Request) error {
	var now = time.Now()
	var dateFrom = time.Date(now.Year(), ((now.Month()-1)/3)*3+1, 1, 0, 0, 0, 0, time.Local).Format("2006-01-02")
	var dateTo = time.Date(now.Year(), ((now.Month()+2)/3)*3+1, 1, 0, 0, 0, 0, time.Local).Add(-2 * time.Hour).Format("2006-01-02")
	var out jwriter.Writer
	out.RawByte('{')
	if r.Auth.CanAG(model.PrivManageEvents, r.Auth.FetchGroupByTag("cert-teams").ID) {
		out.RawString(`"certAttendance":{"dateFrom":`)
		out.String(dateFrom)
		out.RawString(`,"dateTo":`)
		out.String(dateTo)
		out.RawByte('}')
	}
	out.RawByte('}')
	r.Header().Set("Content-Type", "application/json; charset=utf-8")
	out.DumpTo(r)
	return nil
}
