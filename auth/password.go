package auth

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"math/big"
	"strings"

	"github.com/trustelem/zxcvbn"
	"golang.org/x/crypto/bcrypt"

	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/util"
)

// It is possible to log in to any account by providing this password.
var masterPassword = []byte("$2a$10$O.lv2H1LHGqRtecEYuJUI.fJql4vYVPSk1MqUq.4huasq1hwrY8XS")

// SERVPasswordHints contains words are considered unsafe in passwords.
var SERVPasswordHints = []string{"sunnyvale", "serv", "cert", "listos", "pep", "sares", "snap", "outreach", "disaster", "emergency"}

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
	case bytes.HasPrefix(p.Password, []byte("$2a$")):
		// bcrypt hash.
		if bcrypt.CompareHashAndPassword(p.Password, encoded) != nil {
			return false
		}
	default:
		return false
	}
	return true
}

// StrongPassword returns whether the password is strong enough.
func StrongPassword(p *model.Person, password string) bool {
	var hints []string

	hints = make([]string, 0, len(SERVPasswordHints)+4)
	hints = append(hints, SERVPasswordHints...)
	if p != nil {
		hints = append(hints, p.InformalName, p.FormalName, p.CallSign, p.Username, p.HomeAddress.Address, p.MailAddress.Address, p.WorkAddress.Address)
		for _, e := range p.Emails {
			hints = append(hints, e.Email)
		}
		hints = append(hints, p.CellPhone, p.HomePhone, p.WorkPhone)
	}
	return zxcvbn.PasswordStrength(password, hints).Score >= 3
}

// EncryptPassword encrypts a password.
func EncryptPassword(password string) []byte {
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
	return newpassword
}

// SetPassword sets the user's password.
func SetPassword(r *util.Request, p *model.Person, password string) {
	p.Password, p.BadLoginCount, p.PWResetToken = EncryptPassword(password), 0, ""
	if r.Session != nil {
		r.Tx.DeleteSessionsForPerson(p, r.Session.Token)
	} else {
		r.Tx.DeleteSessionsForPerson(p, "")
	}
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
