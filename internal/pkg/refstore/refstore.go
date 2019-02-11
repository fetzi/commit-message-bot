package refstore

// RefStore interface
type RefStore interface {
	Put(key string)
	Exists(key string) bool
}
