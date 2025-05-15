package storage

// Manager manages storage...
type Manager struct {
	storages map[string]Storage
}

// NewManager returns a new NewManager...
func NewManager(storages map[string]Storage) (*Manager, error) {
	sm := &Manager{
		storages: storages,
	}
	return sm, nil
}
