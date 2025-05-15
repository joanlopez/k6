// Package noop offers an implementation of the storage interface that does nothing.
// Predominantly built for testing purposes.
package noop

import (
	"go.k6.io/k6/storage"
)

func init() {
	storage.RegisterExtension("noop", newNoOpStorageFromParams)
}

func newNoOpStorageFromParams(_ storage.Params) (storage.Storage, error) {
	return NewNoOpStorage(), nil
}

// NewNoOpStorage returns a new storage implementation that does nothing.
func NewNoOpStorage() storage.Storage {
	return &noOpStorage{}
}

type noOpStorage struct{}

func (*noOpStorage) Description() string {
	return "this is a no-op storage, built for testing purposes"
}
