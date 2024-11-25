// Package id generates random n-digit IDs.
package id

import (
	"crypto/rand"
	"encoding/binary"
)

const (
	AlphanumCaseSensitiveDigits   = "abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ123456789"
	AlphanumCaseInsensitiveDigits = "ABCDEFGHJKLMNPQRSTUVWXYZ123456789"
)

func newDigit(charset string) byte {
	b := make([]byte, 8)
	n, err := rand.Read(b)
	if n != 8 {
		panic(n)
	} else if err != nil {
		panic(err)
	}
	return charset[uint(binary.BigEndian.Uint64(b))%uint(len(charset))]
}

// New returns a randomly created ID string. It is not guaranteed to be unique.
//
// Risk estimation can look like this, assuming six AlphanumCaseInsensitiveDigits and five tries:
//
//	33^6 ≈ 10^9 different combinations
//	/ 10^6 purchases
//	= 10^-3 risk of individual ID being not unique
//	^ 5 because we do five tries
//	= 10^-15 risk of none of five IDs being not unique
//	/ 10^6 purchases
//	= 10^-9 overall risk of failing
func New(length int, charset string) string {
	var result = make([]byte, length)
	for i := range result {
		result[i] = newDigit(charset)
	}
	return string(result)
}
