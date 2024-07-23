package auth

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

type Struct struct{}

// Create example output: zky_e4d909c290d0fb1ca068ffaddf22cbd0_c29tZUNoZWNrc3Vt
func (Struct) CreateApiKey(prefix string, subject string, secret string) string {
	m := md5.New()
	m.Write([]byte(subject + time.Now().String()))
	hashBytes := m.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(hashString))
	checksum := h.Sum(nil)
	checksumString := hex.EncodeToString(checksum)
	key := fmt.Sprintf("%s_%s_%s", prefix, hashString, checksumString)
	return key
}

func (Struct) ValidateApiKey(key string, secret string) (ok bool) {
	parts := strings.Split(key, "_")
	if len(parts) != 3 {
		return false
	}
	hashString := parts[1]
	checksumString := parts[2]
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(hashString))
	expectedChecksum := h.Sum(nil)
	providedChecksum, err := hex.DecodeString(checksumString)
	if err != nil {
		return false
	}
	return hmac.Equal(expectedChecksum, providedChecksum)
}
