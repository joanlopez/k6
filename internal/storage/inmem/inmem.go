// Package inmem offers an implementation of the storage interface that holds data in memory.
// Predominantly built for testing purposes.
package inmem

import (
	"go.k6.io/k6/storage"
)

func init() {
	storage.RegisterExtension("inmem", newInMemoryStorageFromParams)
}

func newInMemoryStorageFromParams(_ storage.Params) (storage.Storage, error) {
	return NewInMemStorage(), nil
}

// NewInMemStorage returns a new storage implementation that holds data in memory.
func NewInMemStorage() storage.Storage {
	return &inMemStorage{}
}

type inMemStorage struct{}

func (*inMemStorage) Description() string {
	return "this is a in-memory storage, built for testing purposes"
}
