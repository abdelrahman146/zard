package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func (Struct) Encrypt(txt string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(txt))
	encrypted := h.Sum(nil)
	return hex.EncodeToString(encrypted)
}

func (Struct) Decrypt(encrypted string, secret string) (string, error) {
	encryptedBytes, err := hex.DecodeString(encrypted)
	if err != nil {
		return "", err
	}
	h := hmac.New(sha256.New, []byte(secret))
	h.Write(encryptedBytes)
	decrypted := h.Sum(nil)
	return string(decrypted), nil
}

func (a Struct) Compare(encrypted string, txt string, secret string) bool {
	decrypted, err := a.Decrypt(encrypted, secret)
	if err != nil {
		return false
	}
	return decrypted == txt
}
