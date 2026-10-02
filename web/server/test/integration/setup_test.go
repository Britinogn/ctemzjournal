package integration

import "github.com/google/uuid"

// uuidPair holds two seeded user ids for isolation tests.
type uuidPair struct {
	A uuid.UUID
	B uuid.UUID
}
