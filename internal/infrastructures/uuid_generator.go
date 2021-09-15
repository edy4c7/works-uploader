package infrastructures

import "github.com/google/uuid"

type UUIDGeneratorImpl struct{}

var DefaultUUIDGenerator *UUIDGeneratorImpl = &UUIDGeneratorImpl{}

func (r *UUIDGeneratorImpl) Generate() string {
	return uuid.NewString()
}
