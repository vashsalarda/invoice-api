package middleware

import (
	"crypto/aes"
	"fmt"
)

func Encrypt(secretKey string, text string) string {
	aes, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		fmt.Println("err", err)
		panic(err)
	}

	// Make a buffer the same length as plaintext
	cipherText := make([]byte, len(text))
	aes.Encrypt(cipherText, []byte(text))

	return string(cipherText)
}
