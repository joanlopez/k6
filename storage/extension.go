package storage

import (
	"github.com/sirupsen/logrus"
	"go.k6.io/k6/ext"
)

// Constructor returns an instance of a secret source extension module.
// This should return an instance of [Source] given the parameters.
// The Secret Source should not log its secrets and any returned secret will be cached and redacted
// by the [Manager]. No additional work needs to be done by the Secret source apart from retrieving
// the secret.
type Constructor func(Params) (Storage, error)

// Storage is the interface a storage needs to implement.
type Storage interface {
	// Description returns a human-readable description to be printed in the command-line
	// output to let users know which storage extension(s) are in use.
	Description() string
}

// Params contains all possible constructor parameters a storage may need.
type Params struct {
	Logger logrus.FieldLogger
	//Environment map[string]string
	//FS          fsext.Fs
	//Usage       *usage.Usage
}

// RegisterExtension registers the given storage extension constructor.
// This function panics if a module with the same name is already registered.
func RegisterExtension(name string, c Constructor) {
	ext.Register(name, ext.StorageExtension, c)
}
