package auth

import (
	"crypto/sha256"
	"encoding/base64"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"serv.rothskeller.net/portal/model"
)

// It is possible to log in to any account by providing this password.
var masterPassword = []byte("$2a$10$O.lv2H1LHGqRtecEYuJUI.fJql4vYVPSk1MqUq.4huasq1hwrY8XS")

// These words are considered unsafe in passwords.
var servHints = []string{"sunnyvale", "serv", "cert", "listos", "pep", "sares", "snap", "outreach", "disaster", "emergency"}

// CheckPassword verifies that the password is correct for the specified person.
// It returns true if the username and password are valid.
func checkPassword(p *model.Person, password string) bool {
	var (
		hashed  [32]byte
		encoded []byte
	)

	// Prepare the password for bcrypt.  Raw bcrypt has a 72 character
	// maximum (bad for pass-phrases) and doesn't allow NUL characters (bad
	// for binary).  So we start by hashing and base64-encoding the result.
	// That's what we use as the actual password.
	hashed = sha256.Sum256([]byte(password))
	encoded = make([]byte, base64.StdEncoding.EncodedLen(len(hashed)))
	base64.StdEncoding.Encode(encoded, hashed[:])

	// Check against the master password first.  If it matches, we're done.
	if bcrypt.CompareHashAndPassword(masterPassword, encoded) == nil {
		return true
	}

	// Try the various password encryption schemes.
	switch {
	case strings.HasPrefix(p.Password, "$2a$"):
		// bcrypt hash.
		if bcrypt.CompareHashAndPassword([]byte(p.Password), encoded) != nil {
			return false
		}
	default:
		return false
	}
	return true
}

/*
// StrongPassword returns whether the password is strong enough: nil if yes, and
// an appropriate error if not.
func StrongPassword(r *request.Request, u *model.User, password string) error {
	var hints []string

	hints = make([]string, 0, len(servHints)+13)
	hints = append(hints, servHints...)
	if u != nil {
		hints = append(hints, u.Username, u.Name, u.Informal, u.Email, u.Address, u.City, u.State, u.Zip, u.HomePhone,
			u.CellPhone, u.WorkPhone, u.Profession, u.Employer)
	}
	strength := zxcvbn.PasswordStrength(password, hints)
	if strength.CrackTime < 3600 {
		if strength.CrackTimeDisplay == "instant" {
			return fmt.Errorf("This password could be cracked almost instantly.")
		}
		return fmt.Errorf("This password would take about %s to crack.", strength.CrackTimeDisplay)
	}
	return nil
}

// SetPassword sets the user's password.
func SetPassword(r *request.Request, u *model.User, password string) {
	var (
		hashed  [32]byte
		encoded []byte
		err     error
	)

	// Prepare the password for bcrypt.  Raw bcrypt has a 72 character
	// maximum (bad for pass-phrases) and doesn't allow NUL characters (bad
	// for binary).  So we start by hashing and base64-encoding the result.
	// That's what we use as the actual password.
	hashed = sha256.Sum256([]byte(password))
	encoded = make([]byte, base64.StdEncoding.EncodedLen(len(hashed)))
	base64.StdEncoding.Encode(encoded, hashed[:])
	var newpassword []byte
	if newpassword, err = bcrypt.GenerateFromPassword(encoded, 0); err != nil {
		panic(err)
	}
	u.Password, u.Salt, u.RememberMe, u.PWReset = string(newpassword), "", []string{}, ""
	r.Store.SaveUser(u)
	r.Store.DeleteSessions(u.ID)
}

// RandomPassword generates a new, random password.
func RandomPassword() string {
	var (
		words []string
		loc   *big.Int
		err   error
	)
	for i := 0; i < 3; i++ {
		if loc, err = rand.Int(rand.Reader, big.NewInt(int64(len(wordlist)))); err != nil {
			panic(err)
		}
		words = append(words, wordlist[loc.Int64()])
	}
	return strings.Join(words, " ")
}
*/
