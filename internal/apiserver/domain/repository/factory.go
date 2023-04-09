package repository

var client Factory

// Factory defines the repository interface.
type Factory interface {
	UserRepository() UserRepository
	Close() error
}

// Client return the repository client instance.
func Client() Factory {
	return client
}

// SetClient set the repository client.
func SetClient(factory Factory) {
	client = factory
}
