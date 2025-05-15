// Package storage registers all the internal storage implementations when imported
package storage

import (
	_ "go.k6.io/k6/internal/storage/inmem" // import them for init
	_ "go.k6.io/k6/internal/storage/noop"  // import them for init
)
