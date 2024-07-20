package auth

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/abdelrahman146/zard/shared/cache"
	"github.com/abdelrahman146/zard/shared/config"
	"github.com/abdelrahman146/zard/shared/rpc"
	"github.com/abdelrahman146/zard/shared/rpc/requests"
	"strings"
	"time"
)

type apiKeyAuth struct {
	cache cache.Cache
	rpc   rpc.RPC
	conf  config.Config
}

// Create example output: zky_e4d909c290d0fb1ca068ffaddf22cbd0_c29tZUNoZWNrc3Vt
func (a apiKeyAuth) Create(subject string) (token string, err error) {
	prefix := a.conf.GetString("auth.apikey.prefix")
	secret := a.conf.GetString("auth.apikey.secret")
	m := md5.New()
	m.Write([]byte(subject + time.Now().String()))
	hashBytes := m.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(hashString))
	checksum := h.Sum(nil)
	checksumString := hex.EncodeToString(checksum)
	key := fmt.Sprintf("%s_%s_%s", prefix, hashString, checksumString)
	return key, nil
}

func (a apiKeyAuth) Authenticate(key string) (claims Claims, err error) {
	parts := strings.Split(key, "_")
	if len(parts) != 3 {
		return nil, InvalidTokenError
	}

	// validate checksum
	secret := a.conf.GetString("auth.apikey.secret")
	hashString := parts[1]
	checksumString := parts[2]
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(hashString))
	expectedChecksum := h.Sum(nil)
	providedChecksum, err := hex.DecodeString(checksumString)
	if err != nil {
		return nil, InvalidTokenError
	}
	if !hmac.Equal(expectedChecksum, providedChecksum) {
		return nil, InvalidTokenError
	}

	// get token from cache
	if stored, err := a.cache.Get([]string{"auth", "apikey", key}); err == nil && stored != nil {
		claims = &requests.AuthenticateApiKeyResponse{}
		if err := json.Unmarshal(stored, claims); err != nil {
			return nil, err
		}
		return claims, nil
	}

	// get token from db and set in cache
	claims = &requests.AuthenticateApiKeyResponse{}
	if resp, err := a.rpc.Request(&requests.AuthenticateApiKeyRequest{
		ApiKey: key,
	}); err != nil {
		return nil, UnauthorizedError
	} else {
		ttl := time.Duration(a.conf.GetInt("auth.apikey.ttl"))
		if err := a.cache.Set([]string{"auth", "apikey", key}, resp, time.Second*ttl); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(resp, claims); err != nil {
			return nil, err
		}
	}

	return claims, nil
}

func (a apiKeyAuth) Revoke(key string) error {
	// remove key from cache
	if _, err := a.cache.Get([]string{"auth", "apikey", key}); err != nil {
		return err
	}
	// remove key from db
	if _, err := a.rpc.Request(&requests.RevokeApiKeyRequest{
		ApiKey: key,
	}); err != nil {
		return err
	}
	return nil
}

func NewApiKeyAuth(cache cache.Cache, rpc rpc.RPC, conf config.Config) Auth {
	return &apiKeyAuth{
		cache: cache,
		rpc:   rpc,
		conf:  conf,
	}
}
