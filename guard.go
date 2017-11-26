package telebot

import (
	"crypto/sha1"
	"fmt"
)

// Guard is a callback guard, it performs some sort
// of encryption for callbacks and callback data.
type Guard interface {
	Encrypt(text, secret string) string
	Decrypt(ciphertext, secret string) string
}

// XorGuard implements simple XOR encryption.
type XorGuard struct{}

func sha1str(data string) string {
	h := sha1.New()
	h.Write([]byte(data))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

func (x *XorGuard) Encrypt(text, secret string) string {
	// and hope for the best
	key := sha1str(secret)

	ctext := make([]byte, len(text))

	for i, _ := range text {
		ctext[i] = text[i] ^ key[i%len(key)]
	}

	return string(ctext)
}

func (x *XorGuard) Decrypt(ctext, secret string) string {
	return x.Encrypt(ctext, secret)
}
