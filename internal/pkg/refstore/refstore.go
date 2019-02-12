package refstore

// RefStore interface
type RefStore interface {
	Put(key, value string)
	Exists(key string) bool
}
