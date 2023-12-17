package l10n

import "time"

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
