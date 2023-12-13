package util

import (
	"crypto/rand"
	"encoding/base64"
	"strconv"
)

// ParseID parses a string as an integer, returning -1 if it's not valid.
func ParseID(s string) int {
	if val, err := strconv.Atoi(s); err == nil {
		return val
	}
	return -1
}

// RandomToken returns a random token string, used for various purposes.
func RandomToken() string {
	var (
		tokenb [24]byte
		err    error
	)
	if _, err = rand.Read(tokenb[:]); err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(tokenb[:])
}
