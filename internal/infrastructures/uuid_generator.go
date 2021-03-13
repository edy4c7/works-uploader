package infrastructures

type UUIDGeneratorImpl struct{}

var DefaultUUIDGenerator *UUIDGeneratorImpl = &UUIDGeneratorImpl{}

func (r *UUIDGeneratorImpl) Generate() string {
	return ""
}
