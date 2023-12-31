package l10n

import (
	"fmt"
	"strings"
	"time"
)

// Localize returns the translation of the specified string into the specified
// language.
func Localize(s, lang string) string {
	if lang == "es" {
		if tr, ok := spanish[s]; ok {
			return tr
		}
	}
	return s
}

// LocalizeDate returns the long form of the date in the specified language.
func LocalizeDate(date time.Time, lang string) string {
	if lang == "es" {
		return spanishDate(date)
	}
	return date.Format("Monday, January 2, 2006")
}

// Conjoin joins the supplied list of strings (which must already be localized)
// with the specified (English) conjunction, translated into the specified
// language.
func Conjoin(list []string, conjunction, lang string) string {
	switch len(list) {
	case 0:
		return ""
	case 1:
		return list[0]
	}
	if lang == "es" {
		last := list[len(list)-1]
		if conjunction == "and" {
			if strings.HasPrefix(last, "i") || (strings.HasPrefix(last, "hi") && !strings.HasPrefix(last, "hie")) {
				conjunction = "e"
			} else {
				conjunction = "y"
			}
		} else if conjunction == "or" {
			if strings.HasPrefix(last, "o") || strings.HasPrefix(last, "ho") {
				conjunction = "u"
			} else {
				conjunction = "o"
			}
		} else {
			conjunction = Localize(conjunction, lang)
		}
	}
	if len(list) == 2 || lang == "es" { // No Oxford comma.
		return fmt.Sprintf("%s %s %s", strings.Join(list[:len(list)-1], ", "), conjunction, list[len(list)-1])
	}
	// Oxford comma.
	return fmt.Sprintf("%s, %s %s", strings.Join(list[:len(list)-1], ", "), conjunction, list[len(list)-1])
}
