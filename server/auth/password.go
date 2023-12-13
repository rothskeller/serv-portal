package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"math/big"
	"strings"
	"time"

	"github.com/nbutton23/zxcvbn-go"
	"golang.org/x/crypto/bcrypt"

	"sunnyvaleserv.org/portal/store"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/store/session"
	"sunnyvaleserv.org/portal/util/request"
)

// It is possible to log in to any account by providing this password.
var masterPassword = []byte("$2a$10$O.lv2H1LHGqRtecEYuJUI.fJql4vYVPSk1MqUq.4huasq1hwrY8XS")

// SERVPasswordHints contains words are considered unsafe in passwords.
var SERVPasswordHints = []string{"sunnyvale", "serv", "cert", "listos", "pep", "sares", "snap", "outreach", "disaster", "emergency"}

// CheckPassword verifies that the password is correct for the specified user.
// (As a side effect, if the password is stored in any outdated encryption
// method, it is re-encrypted in a modern one.)  It returns true if the username
// and password are valid.
func CheckPassword(storer store.Storer, p *person.Person, password string) bool {
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
	case strings.HasPrefix(p.Password(), "$2a$"):
		// bcrypt hash.
		if bcrypt.CompareHashAndPassword([]byte(p.Password()), encoded) != nil {
			return false
		}
	default:
		return false
	}
	return true
}

// StrongPasswordPersonFields is a bitmask of the fields that the caller of
// StrongPassword must request when fetching the person whose password is being
// verified.
const StrongPasswordPersonFields = person.FInformalName | person.FFormalName | person.FCallSign | person.FEmail | person.FEmail2 | person.FAddresses | person.FCellPhone | person.FHomePhone | person.FWorkPhone

// StrongPasswordHints returns the set of hints to give to zxcvbn when checking
// password strength.
func StrongPasswordHints(p *person.Person) []string {
	var hints = make([]string, 0, len(SERVPasswordHints)+11)
	hints = append(hints, SERVPasswordHints...)
	if p != nil {
		hints = append(hints, p.InformalName(), p.FormalName(), p.CallSign(), p.Email(), p.Email2(), p.CellPhone(), p.HomePhone(), p.WorkPhone())
		if p.Addresses().Home != nil {
			hints = append(hints, p.Addresses().Home.Address)
		}
		if p.Addresses().Work != nil {
			hints = append(hints, p.Addresses().Work.Address)
		}
		if p.Addresses().Mail != nil {
			hints = append(hints, p.Addresses().Mail.Address)
		}
	}
	return hints
}

// StrongPassword returns whether the password is strong enough.
func StrongPassword(p *person.Person, password string) bool {
	var hints = StrongPasswordHints(p)
	return zxcvbn.PasswordStrength(password, hints).Score >= 3
}

// EncryptPassword encrypts a new password.
func EncryptPassword(password string) string {
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
	return string(newpassword)
}

// SetPassword sets the user's password.
func SetPassword(r *request.Request, p *person.Person, password string) {
	r.Transaction(func() {
		up := p.Updater()
		up.Password = EncryptPassword(password)
		up.BadLoginCount = 0
		up.BadLoginTime = time.Time{}
		up.PWResetToken = ""
		up.PWResetTime = time.Time{}
		p.Update(r, up, person.FPassword|person.FBadLoginCount|person.FBadLoginTime|person.FPWResetToken|person.FPWResetTime)
		session.DeleteForPerson(r, p, r.SessionToken)
	})
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
