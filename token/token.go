package token

import (
	"crypto/rand"
	"encoding/base64"
	"math/big"
)

// RandomBytes generate a 32 random bytes, encode using base64 URL encoding
// and return the string
func RandomBytes() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

var runes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQUVWXYZ0123456789")

// RandomToken returns a secure random string with length n
func RandomToken(n int) (*string, error) {
	b := make([]rune, n)
	for i := range b {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(runes))))
		if err != nil {
			return nil, err
		}
		b[i] = runes[num.Int64()]
	}
	result := string(b)
	return &result, nil
}
