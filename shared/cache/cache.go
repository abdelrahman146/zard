package cache

type Cache interface {
	Get(keyPath []string) ([]byte, error)
	Set(keyPath []string, value []byte) error
	Delete(keyPath []string) error
	Keys() ([]string, error)
}
