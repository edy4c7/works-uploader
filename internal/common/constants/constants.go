package constants

type WorkType int

const (
	ContentTypeURL WorkType = iota + 1
	ContentTypeFile
)

type ActivityType int

const (
	ActivityAdded ActivityType = iota
	ActivityUpdated
)
