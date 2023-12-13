package ui

import (
	"math"
	"regexp"
	"strconv"
	"strings"
)

// MinutesToHours renders a number of minutes as an hours string suitable for an
// s-hours control.
func MinutesToHours(minutes uint) string {
	minutes += 19
	h, m := minutes/60, minutes%60
	if m >= 30 && h == 0 {
		return "½"
	} else if m >= 30 {
		return strconv.Itoa(int(h)) + "½"
	} else {
		return strconv.Itoa(int(h))
	}
}

var floatHoursRE = regexp.MustCompile(`^\d+\.\d*$`)
var timeRangeRE = regexp.MustCompile(`^(?:[01]\d?|2[0-3]?):?([0-5]\d)[- ]?(?:[01]\d?|2[0-3]?):?([0-5]\d)$`)

// SHoursValue parses the value returned from an s-hours control.  It allows for
// the possibility that the Javascript was disabled and the control contains a
// raw value.
func SHoursValue(s string) (minutes uint, ok bool) {
	// Special cases.
	if s == "" {
		return 0, true
	}
	if s == "½" {
		return 30, true
	}
	// The next simplest input is an integer number of hours.
	if h, err := strconv.Atoi(s); err == nil && h >= 0 {
		return uint(h * 60), true
	}
	// If the string ends with "½", the rest could be an integer number of
	// hours.
	if strings.HasSuffix(s, "½") {
		if h, err := strconv.Atoi(s[:len(s)-2]); err == nil && h >= 0 {
			return uint(h*60) + 30, true
		}
	}
	// If the string looks like a float, parse it that way.
	if floatHoursRE.MatchString(s) {
		h, _ := strconv.ParseFloat(s, 64)
		m := uint(math.Round(h * 60))
		if rem := m % 30; rem <= 10 {
			m -= rem
		} else {
			m += 30 - rem
		}
		return m, true
	}
	// If the string looks like a timesheet range, parse it that way.
	if match := timeRangeRE.FindStringSubmatch(s); match != nil {
		sh, _ := strconv.Atoi(match[1])
		sm, _ := strconv.Atoi(match[2])
		eh, _ := strconv.Atoi(match[3])
		em, _ := strconv.Atoi(match[4])
		start := sh*60 + sm
		end := eh*60 + em
		diff := end - start
		if diff >= 0 {
			m := uint(diff)
			if rem := m % 30; rem <= 10 {
				m -= rem
			} else {
				m += 30 - rem
			}
			return m, true
		}
	}
	return 0, false
}
