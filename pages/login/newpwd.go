package login

import (
	"fmt"
	"math"
	"regexp"
	"strings"

	"github.com/nbutton23/zxcvbn-go"
	"github.com/nbutton23/zxcvbn-go/match"
	"github.com/nbutton23/zxcvbn-go/scoring"

	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/ui/form"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

// NewPasswordPairRow is a form row containing a pair of entry fields for a new
// password.
type NewPasswordPairRow struct {
	form.LabeledRow
	FocusID  string
	Name     string
	ValueP1  *string
	ValueP2  *string
	Validate string
	Override bool
	Person   *person.Person
	Score    int
	Messages []string
}

func (nppr *NewPasswordPairRow) Get() {
	nppr.Score = -1
}

func (nppr *NewPasswordPairRow) Read(r *request.Request) bool {
	*nppr.ValueP1 = r.FormValue(nppr.Name + "-1")
	*nppr.ValueP2 = r.FormValue(nppr.Name + "-2")
	// If the two passwords don't match, and both are non-empty, raise that
	// error first.
	if *nppr.ValueP1 != "" && *nppr.ValueP2 != "" && *nppr.ValueP1 != *nppr.ValueP2 {
		nppr.Score = -1
		nppr.Messages = []string{r.Loc("The two passwords are not the same.")}
		return false
	}
	// If the two passwords don't match, and we're submitting rather than
	// validating, raise that error.
	validating := r.Request.Header.Get("X-Up-Validate") != ""
	if !validating && *nppr.ValueP1 != *nppr.ValueP2 {
		nppr.Score = -1
		nppr.Messages = []string{r.Loc("The two passwords are not the same.")}
		return false
	}
	// If they haven't specified a new password at all, flag that.
	if !validating && *nppr.ValueP1 == "" {
		nppr.Score = -1
		nppr.Messages = []string{r.Loc("Please specify a new password, twice.")}
		return false
	}
	// Either the new passwords match, or we're validating and the second
	// one is empty.  Determine the strength of the password and record
	// score and issues.
	hints := auth.StrongPasswordHints(nppr.Person)
	analysis := zxcvbn.PasswordStrength(*nppr.ValueP1, hints)
	nppr.Score = analysis.Score
	nppr.Messages = analysisMessages(r, analysis)
	return nppr.Score > 2 || nppr.Override
}

func (nppr *NewPasswordPairRow) Emit(r *request.Request, parent *htmlb.Element, focus bool) {
	if nppr.FocusID == "" && nppr.RowID != "" {
		nppr.FocusID = nppr.RowID + "-in"
	}
	row := nppr.EmitPrefix(r, parent, nppr.FocusID)
	box := row.E("div class=formInput")
	box.E("input type=password name=%s-1 class=formInput value=%s autocomplete=new-password up-watch-event=input up-watch-delay=100 up-keep",
		nppr.Name, *nppr.ValueP1,
		nppr.FocusID != "", "id=%s", nppr.FocusID,
		focus, "autofocus",
		nppr.Validate == "" && nppr.RowID != "", "up-validate=#%s", nppr.RowID,
		nppr.Validate != "" && nppr.Validate != form.NoValidate, "up-validate=%s", nppr.Validate)
	box.E("input type=password name=%s-2 class=formInput value=%s autocomplete=new-password up-watch-event=input up-watch-delay=100 up-keep",
		nppr.Name, *nppr.ValueP2,
		nppr.Validate == "" && nppr.RowID != "", "up-validate=#%s", nppr.RowID,
		nppr.Validate != "" && nppr.Validate != form.NoValidate, "up-validate=%s", nppr.Validate)
	if nppr.Score < 0 && len(nppr.Messages) == 0 {
		return
	}
	feedback := row.E("div class='formError newpwdFeedback'")
	var scoreColor = "bad"
	if nppr.Score > 2 {
		scoreColor = "good"
	} else if nppr.Score == 2 {
		scoreColor = "warn"
	}
	if nppr.Score >= 0 {
		meter := feedback.E("div class=newpwdFeedbackMeter")
		for i := 0; i <= nppr.Score; i++ {
			meter.E("div class='newpwdFeedbackMeterStep %s'", scoreColor)
		}
	}
	if len(nppr.Messages) != 0 {
		feedback.E("div class=%s", scoreColor).T(strings.Join(nppr.Messages, "  "))
	}
}

func analysisMessages(r *request.Request, a scoring.MinEntropyMatch) (messages []string) {
	messages = append(messages, crackTimeMessage(r, a.CrackTime))
	if len(a.MatchSequence) == 0 {
		return append(messages,
			r.Loc("Use a few words.  Avoid common phrases."),
			r.Loc("No need for symbols, digits, or uppercase letters."))
	}
	if a.Score > 2 {
		return messages
	}
	longest := a.MatchSequence[0]
	for _, seq := range a.MatchSequence {
		if len(seq.Token) > len(longest.Token) {
			longest = seq
		}
	}
	messages = append(messages, matchMessages(r, longest, len(a.MatchSequence) == 1)...)
	return messages
}

func crackTimeMessage(r *request.Request, seconds float64) (message string) {
	const minute = 60.0
	const hour = minute * 60
	const day = hour * 24
	const month = day * 31
	const year = month * 12
	const century = year * 100
	var count int

	message = "This password would take %d "
	if seconds < minute {
		return r.Loc("This password would take less than a minute to crack.")
	} else if seconds < hour {
		count = int(math.Round(seconds / minute))
		message += "minute"
	} else if seconds < day {
		count = int(math.Round(seconds / hour))
		message += "hour"
	} else if seconds < month {
		count = int(math.Round(seconds / day))
		message += "day"
	} else if seconds < year {
		count = int(math.Round(seconds / month))
		message += "month"
	} else if seconds < century {
		count = int(math.Round(seconds / year))
		message += "year"
	} else {
		return r.Loc("This password would take centuries to crack.")
	}
	if count != 1 {
		message += "s"
	}
	message += " to crack."
	return fmt.Sprintf(r.Loc(message), count)
}

func matchMessages(r *request.Request, m match.Match, only bool) (messages []string) {
	switch m.Pattern {
	case "dictionary":
		messages = dictionaryMatchMessages(r, m, only)
	case "spatial":
		messages = append(messages,
			r.Loc("Short keyboard patterns are easy to guess."),
			r.Loc("Use a longer keyboard pattern with more turns."))
	case "repeat":
		if len(m.DictionaryName) == 1 {
			messages = append(messages, r.Loc("Repeats like “aaa” are easy to guess."))
		} else {
			messages = append(messages, r.Loc("Repeats like “abcabcabc” are only slightly harder to guess than “abc”."))
		}
		messages = append(messages, r.Loc("Avoid repeated words and characters."))
	case "sequence":
		messages = append(messages,
			r.Loc("Sequences like “abc” or “6543” are easy to guess."),
			r.Loc("Avoid sequences."))
	case "":
		if m.DictionaryName != "date_match" {
			break
		}
		fallthrough
	case "date":
		messages = append(messages,
			r.Loc("Dates are often easy to guess."),
			r.Loc("Avoid dates and years that are associated with you."))
	}
	messages = append(messages,
		r.Loc("Add another word or two.  Uncommon words are better."))
	return messages
}

var startUpperRE = regexp.MustCompile(`^[A-Z][^A-Z]+$`)
var allUpperRE = regexp.MustCompile(`^[^a-z]+$`)

func dictionaryMatchMessages(r *request.Request, m match.Match, only bool) (messages []string) {
	switch m.DictionaryName {
	case "passwords", "passwords_3117":
		messages = append(messages, r.Loc("This is similar to a commonly used password."))
	case "english_wikipedia", "english_wikipedia_3117":
		if only {
			messages = append(messages, r.Loc("A word by itself is easy to guess."))
		}
	case "surnames", "surnames_3117", "male_names", "male_names_3117", "female_names", "female_names_3117":
		messages = append(messages, r.Loc("Common names and surnames are easy to guess."))
	}
	if startUpperRE.MatchString(m.Token) {
		messages = append(messages, r.Loc("Capitalization doesn’t help very much."))
	} else if allUpperRE.MatchString(m.Token) && strings.ToLower(m.Token) != m.Token {
		messages = append(messages, r.Loc("All upper case is almost as easy to guess as all lower case."))
	}
	if strings.HasSuffix(m.DictionaryName, "_3117") {
		messages = append(messages, r.Loc("Predictable substitutions like “@” instead of “a” don’t help very much."))
	}
	return messages
}
