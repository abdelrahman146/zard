package cache

import "time"

type Cache interface {
	Get(keyPath []string) ([]byte, error)
	Set(keyPath []string, value []byte, ttl time.Duration) error
	Delete(keyPath []string) error
	Keys(keyPath []string) ([]string, error)
}
