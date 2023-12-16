package langdate

import (
	"fmt"
	"time"

	"sunnyvaleserv.org/portal/util/request"
)

func LangDate(r *request.Request, day time.Time) string {
	if r.Language != "es" {
		return day.Format("Monday, January 2, 2006")
	}
	return fmt.Sprintf("%s, %d de %s de %d", esWeekdays[day.Weekday()], day.Day(), esMonths[day.Month()], day.Year())
}

var esWeekdays = map[time.Weekday]string{
	time.Sunday:    "Domingo",
	time.Monday:    "Lunes",
	time.Tuesday:   "Martes",
	time.Wednesday: "Miércoles",
	time.Thursday:  "Jueves",
	time.Friday:    "Vieres",
	time.Saturday:  "Sábado",
}
var esMonths = map[time.Month]string{
	time.January:   "enero",
	time.February:  "febrero",
	time.March:     "marzo",
	time.April:     "abril",
	time.May:       "mayo",
	time.June:      "junio",
	time.July:      "julio",
	time.August:    "agosto",
	time.September: "septiembre",
	time.October:   "octubre",
	time.November:  "noviembre",
	time.December:  "diciembre",
}
